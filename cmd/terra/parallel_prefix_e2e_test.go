//go:build integration

package main_test

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// writeFakeCLI writes an executable shell script that stands in for a real terragrunt or
// terraform binary. It only echoes recognizable lines — including a trailing line with no
// newline, to exercise the writer's Flush path — and never touches real infrastructure.
func writeFakeCLI(t *testing.T, path string) {
	t.Helper()
	const script = `#!/bin/sh
case "$1" in
  --version|-v|version)
    echo "v0.99.0"
    ;;
  *)
    echo "FAKE-CLI stdout line one"
    echo "FAKE-CLI stdout line two"
    echo "FAKE-CLI stderr line" 1>&2
    printf "FAKE-CLI stdout tail-without-newline"
    ;;
esac
`
	require.NoError(t, os.WriteFile(path, []byte(script), 0o755))
}

// hashFile returns the hex-encoded SHA-256 of the file at path.
func hashFile(t *testing.T, path string) string {
	t.Helper()
	file, err := os.Open(path)
	require.NoError(t, err)
	defer func() { _ = file.Close() }()

	hasher := sha256.New()
	_, err = io.Copy(hasher, file)
	require.NoError(t, err)

	return hex.EncodeToString(hasher.Sum(nil))
}

// snapshotRealBinaries records the SHA-256 of the real terragrunt/terraform on the test
// process's PATH so the test can prove afterward that their contents are byte-for-byte
// unchanged by the run.
func snapshotRealBinaries(t *testing.T) map[string]string {
	t.Helper()
	snapshot := make(map[string]string)
	for _, name := range []string{"terragrunt", "terraform"} {
		resolved, err := exec.LookPath(name)
		if err != nil {
			continue
		}
		snapshot[resolved] = hashFile(t, resolved)
	}
	return snapshot
}

// TestParallelOutputPrefix_E2E runs the real terra binary end to end against fake
// terragrunt/terraform stubs and asserts that --parallel prefixes each module's output.
//
// Safety: HOME, PATH, the install path, and the cache dirs are set ON THE SUBPROCESS ONLY
// to directories under t.TempDir(). The terra run therefore cannot read or overwrite the
// developer's real ~/.local/bin binaries or ~/.cache, and the test additionally asserts the
// real terragrunt/terraform binaries are byte-for-byte unchanged.
func TestParallelOutputPrefix_E2E(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()

	realBefore := snapshotRealBinaries(t)

	// 1. Build the terra binary into the sandbox. Without ldflags the version stays "dev",
	//    so terra skips its self-update check and makes no network calls.
	terraBin := filepath.Join(tmp, "terra")
	build := exec.Command("go", "build", "-o", terraBin, "github.com/rios0rios0/terra/cmd/terra")
	build.Stderr = os.Stderr
	require.NoError(t, build.Run(), "failed to build terra binary")

	// 2. Fake terragrunt + terraform on an isolated PATH bin directory.
	fakeBin := filepath.Join(tmp, "bin")
	// nosemgrep: go.lang.correctness.permissions.file_permission.incorrect-default-permission
	require.NoError(t, os.MkdirAll(fakeBin, 0o700))
	writeFakeCLI(t, filepath.Join(fakeBin, "terragrunt"))
	writeFakeCLI(t, filepath.Join(fakeBin, "terraform"))

	// 3. A tree of two modules for terra to process in parallel.
	modulesDir := filepath.Join(tmp, "modules")
	for _, module := range []string{"module-a", "module-b"} {
		dir := filepath.Join(modulesDir, module)
		// nosemgrep: go.lang.correctness.permissions.file_permission.incorrect-default-permission
		require.NoError(t, os.MkdirAll(dir, 0o700))
		require.NoError(t, os.WriteFile(
			filepath.Join(dir, "terragrunt.hcl"), []byte("# fake\n"), 0o644))
	}

	// 4. Fully isolated environment for the terra subprocess. A sandbox HOME means even the
	//    ~/.local/bin and ~/.cache fallbacks resolve inside t.TempDir().
	sandboxHome := filepath.Join(tmp, "home")
	// nosemgrep: go.lang.correctness.permissions.file_permission.incorrect-default-permission
	require.NoError(t, os.MkdirAll(sandboxHome, 0o700))
	env := append(os.Environ(),
		"HOME="+sandboxHome,
		"PATH="+fakeBin+string(os.PathListSeparator)+os.Getenv("PATH"),
		"TERRA_INSTALL_PATH="+filepath.Join(tmp, "install"),
		"TERRA_MODULE_CACHE_DIR="+filepath.Join(tmp, "cache", "modules"),
		"TERRA_PROVIDER_CACHE_DIR="+filepath.Join(tmp, "cache", "providers"),
	)

	// 5. Any command carrying --parallel=N goes through terra's worker pool, which prefixes
	//    each module's output. plan never prompts, so no confirmation flag is needed.
	var stdout, stderr bytes.Buffer
	run := exec.Command(terraBin, "plan", "--parallel=2", modulesDir)
	run.Dir = tmp
	run.Env = env
	run.Stdout = &stdout
	run.Stderr = &stderr
	require.NoErrorf(t, run.Run(),
		"terra run failed\n--- stdout ---\n%s\n--- stderr ---\n%s", stdout.String(), stderr.String())

	out := stdout.String()

	// Each module's fake-terragrunt stdout is streamed through terra and tagged with the
	// module's directory name.
	assert.Containsf(t, out, "[module-a] FAKE-CLI stdout line one",
		"module-a output should be prefixed; full stdout:\n%s", out)
	assert.Containsf(t, out, "[module-b] FAKE-CLI stdout line one",
		"module-b output should be prefixed; full stdout:\n%s", out)

	// The trailing line the stub prints without a newline is flushed with its prefix.
	assert.Contains(t, out, "[module-a] FAKE-CLI stdout tail-without-newline")
	assert.Contains(t, out, "[module-b] FAKE-CLI stdout tail-without-newline")

	// Output is a pipe (not a TTY), so the prefix carries no ANSI color codes.
	assert.NotContains(t, out, "\x1b[", "prefix must be plain when stdout is not a terminal")

	// Safety guard: the developer's real terragrunt/terraform binaries must be byte-for-byte
	// unchanged, i.e. their SHA-256 must match the pre-run snapshot.
	for binaryPath, before := range realBefore {
		assert.Equalf(t, before, hashFile(t, binaryPath),
			"e2e test must not modify the real binary at %s (SHA-256 changed)", binaryPath)
	}
}

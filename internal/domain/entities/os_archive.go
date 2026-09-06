package entities

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	// maxExtractedBytes caps the total number of bytes written while
	// unpacking a single archive. Released `terraform` / `terragrunt`
	// archives are well below 1 GiB, so this ceiling is only reached by a
	// decompression bomb: a small archive that expands into an arbitrarily
	// large payload and fills the installation disk.
	maxExtractedBytes int64 = 1 << 30

	// extractedDirPerm is applied to the directories created while
	// unpacking an archive. Archives are unpacked into a private
	// `os.MkdirTemp` directory, so the owner bits are the only ones any
	// step of the installation ever needs: group and other get nothing.
	extractedDirPerm os.FileMode = 0o700

	// extractedFilePermMask keeps the executable bit carried by an archive
	// entry (release binaries ship as 0o755) while dropping group and other
	// write permissions, mirroring what `unzip` does under a normal umask.
	extractedFilePermMask os.FileMode = 0o755
)

// errArchiveTooLarge signals that an archive expands past maxExtractedBytes.
var errArchiveTooLarge = errors.New("archive expands beyond the allowed size limit")

// extractZipArchive unpacks a ZIP archive into destPath using the Go
// standard library. Doing the work in-process removes the need to resolve
// an external `unzip` / `powershell` binary through PATH, so the archive
// handling no longer depends on what happens to be installed on the host.
func extractZipArchive(archivePath, destPath string) error {
	reader, err := zip.OpenReader(archivePath)
	if err != nil {
		return fmt.Errorf("failed to perform decompressing of %q: %w", archivePath, err)
	}
	defer reader.Close()

	var written int64
	for _, entry := range reader.File {
		entryWritten, entryErr := extractZipEntry(entry, destPath, maxExtractedBytes-written)
		if entryErr != nil {
			return entryErr
		}
		written += entryWritten
	}

	return nil
}

// escapingEntryError reports an archive entry that would be written outside
// the extraction directory.
func escapingEntryError(name string) error {
	return fmt.Errorf(
		"failed to perform decompressing: entry %q escapes the destination directory",
		name,
	)
}

// extractZipEntry writes a single archive entry below destPath and reports
// how many bytes it produced. Entries whose name would resolve outside
// destPath ("Zip Slip") are rejected, and an entry that would push the
// extraction past the remaining budget aborts the operation.
//
// Both guards are spelled out inline rather than delegated to helpers on
// purpose: a taint analyser only accepts an archive entry as sanitized when
// the check dominates the filesystem call in the same function, so hiding
// the comparison behind an `isWithinDir(...)` wrapper reads as an
// unsanitized name flowing into `os.MkdirAll` and `os.OpenFile`.
func extractZipEntry(entry *zip.File, destPath string, budget int64) (int64, error) {
	// First guard: reject a name that is absolute or that climbs out with
	// `..` — the "Zip Slip" attack, e.g. `../../etc/cron.d/payload`.
	relativeName := filepath.Clean(filepath.FromSlash(entry.Name))
	if filepath.IsAbs(relativeName) || relativeName == ".." ||
		strings.HasPrefix(relativeName, ".."+string(os.PathSeparator)) {
		return 0, escapingEntryError(entry.Name)
	}
	// An entry naming the destination itself (`./`) carries nothing to unpack.
	if relativeName == "." {
		return 0, nil
	}

	// Second guard: even a name that cleared the first one must resolve to a
	// path under destPath before anything touches the filesystem.
	targetPath := filepath.Join(destPath, relativeName)
	if !strings.HasPrefix(targetPath, filepath.Clean(destPath)+string(os.PathSeparator)) {
		return 0, escapingEntryError(entry.Name)
	}

	if entry.FileInfo().IsDir() {
		// nosemgrep: go.lang.correctness.permissions.file_permission.incorrect-default-permission
		if err := os.MkdirAll(targetPath, extractedDirPerm); err != nil {
			return 0, fmt.Errorf("failed to perform decompressing of %q: %w", entry.Name, err)
		}
		return 0, nil
	}

	// nosemgrep: go.lang.correctness.permissions.file_permission.incorrect-default-permission
	if err := os.MkdirAll(filepath.Dir(targetPath), extractedDirPerm); err != nil {
		return 0, fmt.Errorf("failed to perform decompressing of %q: %w", entry.Name, err)
	}

	source, err := entry.Open()
	if err != nil {
		return 0, fmt.Errorf("failed to perform decompressing of %q: %w", entry.Name, err)
	}
	defer source.Close()

	// nosemgrep: go.lang.correctness.permissions.file_permission.incorrect-default-permission
	target, err := os.OpenFile(
		targetPath,
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
		entry.Mode().Perm()&extractedFilePermMask,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to perform decompressing of %q: %w", entry.Name, err)
	}
	defer target.Close()

	written, err := copyWithinBudget(target, source, budget)
	if err != nil {
		return written, fmt.Errorf("failed to perform decompressing of %q: %w", entry.Name, err)
	}

	// Close explicitly for the same reason as in copyFile: a write-back
	// failure is reported by close, so a deferred close would hand a
	// silently truncated binary to findBinaryInArchive. The deferred Close
	// above then becomes a no-op.
	if err = target.Close(); err != nil {
		return written, fmt.Errorf("failed to perform decompressing of %q: %w", entry.Name, err)
	}

	return written, nil
}

// copyWithinBudget streams src into dst and refuses to produce more than
// budget bytes, which is what stops a decompression bomb from filling the
// disk. Reading one byte past the budget tells an oversized entry apart
// from one that merely exhausts it exactly.
func copyWithinBudget(dst io.Writer, src io.Reader, budget int64) (int64, error) {
	written, err := io.Copy(dst, io.LimitReader(src, budget+1))
	if err != nil {
		return written, err
	}
	if written > budget {
		return written, errArchiveTooLarge
	}

	return written, nil
}

// moveFile relocates srcPath onto destPath through the Go standard library
// rather than the platform's `mv` / `move` command. [os.Rename] cannot cross
// filesystem boundaries — the download lands in the temp directory while
// the installation directory frequently sits on another mount — so a
// copy-then-delete fallback preserves the behaviour of `mv`.
func moveFile(srcPath, destPath string) error {
	if err := os.Rename(srcPath, destPath); err == nil {
		return nil
	}

	if err := copyFile(srcPath, destPath); err != nil {
		return fmt.Errorf("failed to perform moving of %q to %q: %w", srcPath, destPath, err)
	}

	if err := os.Remove(srcPath); err != nil {
		return fmt.Errorf("failed to perform moving of %q to %q: %w", srcPath, destPath, err)
	}

	return nil
}

// copyFile duplicates srcPath onto destPath, preserving the source's
// permission bits. It backs moveFile's cross-filesystem fallback.
func copyFile(srcPath, destPath string) error {
	source, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("failed to open %q: %w", srcPath, err)
	}
	defer source.Close()

	info, err := source.Stat()
	if err != nil {
		return fmt.Errorf("failed to inspect %q: %w", srcPath, err)
	}

	// nosemgrep: go.lang.correctness.permissions.file_permission.incorrect-default-permission
	target, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, info.Mode().Perm())
	if err != nil {
		return fmt.Errorf("failed to create %q: %w", destPath, err)
	}
	defer target.Close()

	if _, err = io.Copy(target, source); err != nil {
		return fmt.Errorf("failed to copy %q to %q: %w", srcPath, destPath, err)
	}

	// Close explicitly: write-back errors (ENOSPC, EDQUOT, an NFS commit
	// failure) surface at close, and moveFile deletes the source once this
	// returns nil. The deferred Close above then becomes a no-op.
	if err = target.Close(); err != nil {
		return fmt.Errorf("failed to finalize %q: %w", destPath, err)
	}

	return nil
}

// removeFile deletes filePath through the Go standard library instead of
// the platform's `rm` / `del` command.
func removeFile(filePath string) error {
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to perform deleting of %q: %w", filePath, err)
	}

	return nil
}

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
	// unpacking an archive.
	extractedDirPerm os.FileMode = 0o750

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

// extractZipEntry writes a single archive entry below destPath and reports
// how many bytes it produced. Entries whose name would resolve outside
// destPath ("Zip Slip") are rejected, and an entry that would push the
// extraction past the remaining budget aborts the operation.
func extractZipEntry(entry *zip.File, destPath string, budget int64) (int64, error) {
	relativeName, err := safeArchivePath(entry.Name)
	if err != nil {
		return 0, err
	}

	targetPath := filepath.Join(destPath, relativeName)
	// Second layer of defence: even for a name that passed the checks above,
	// refuse to write anything that does not land under destPath.
	if !isWithinDir(targetPath, destPath) {
		return 0, fmt.Errorf(
			"failed to perform decompressing: entry %q escapes the destination directory",
			entry.Name,
		)
	}

	if entry.FileInfo().IsDir() {
		if err = os.MkdirAll(targetPath, extractedDirPerm); err != nil {
			return 0, fmt.Errorf("failed to perform decompressing of %q: %w", entry.Name, err)
		}
		return 0, nil
	}

	if err = os.MkdirAll(filepath.Dir(targetPath), extractedDirPerm); err != nil {
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

	// Reading one byte past the budget tells an oversized entry apart from
	// one that merely exhausts it exactly.
	written, err := io.Copy(target, io.LimitReader(source, budget+1))
	if err != nil {
		return written, fmt.Errorf("failed to perform decompressing of %q: %w", entry.Name, err)
	}
	if written > budget {
		return written, fmt.Errorf(
			"failed to perform decompressing of %q: %w", entry.Name, errArchiveTooLarge,
		)
	}

	return written, nil
}

// safeArchivePath converts an archive entry name into a relative path that
// is guaranteed to stay inside the extraction directory. Entries that are
// absolute or that climb out with `..` — the "Zip Slip" attack, e.g.
// `../../etc/cron.d/payload` — are rejected outright.
func safeArchivePath(name string) (string, error) {
	cleaned := filepath.Clean(filepath.FromSlash(name))

	escapes := filepath.IsAbs(cleaned) ||
		cleaned == ".." ||
		strings.HasPrefix(cleaned, ".."+string(os.PathSeparator))
	if escapes {
		return "", fmt.Errorf(
			"failed to perform decompressing: entry %q escapes the destination directory",
			name,
		)
	}

	return cleaned, nil
}

// isWithinDir reports whether target resolves to a path inside dir. It is
// the guard against archive entries such as `../../etc/cron.d/payload`
// that would otherwise be written outside the extraction directory.
func isWithinDir(target, dir string) bool {
	cleanDir := filepath.Clean(dir)
	cleanTarget := filepath.Clean(target)
	if cleanTarget == cleanDir {
		return true
	}

	return strings.HasPrefix(cleanTarget, cleanDir+string(os.PathSeparator))
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

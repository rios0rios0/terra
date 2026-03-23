package commands

// ExtractArchiveForTest exposes the private extractArchive method for use in external test packages.
func (it *SelfUpdateCommand) ExtractArchiveForTest(archivePath, destDir string) error {
	return it.extractArchive(archivePath, destDir)
}

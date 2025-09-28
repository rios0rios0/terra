package repositories

// ShellRepository is not totally necessary, but it is rather a good example for other applications.
type ShellRepository interface {
	ExecuteCommand(command string, arguments []string, directory string) error
}

// ShellRepositoryWithUpgrade extends ShellRepository with auto-upgrade detection capabilities.
type ShellRepositoryWithUpgrade interface {
	ShellRepository
	ExecuteCommandWithUpgradeDetection(command string, arguments []string, directory string) error
}

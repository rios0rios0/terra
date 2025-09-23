package repositories

// InteractiveShellRepository defines the interface for interactive shell operations with auto-answering.
type InteractiveShellRepository interface {
	ExecuteCommand(command string, arguments []string, directory string) error
	ExecuteCommandWithAnswer(command string, arguments []string, directory string, autoAnswer string) error
}

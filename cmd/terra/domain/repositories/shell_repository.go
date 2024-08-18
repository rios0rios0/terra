package repositories

type ShellRepository interface {
	ExecuteCommand(command string, arguments []string, directory string) error
}

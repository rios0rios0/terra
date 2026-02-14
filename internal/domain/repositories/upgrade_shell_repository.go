package repositories

// UpgradeShellRepository extends ShellRepository with automatic upgrade detection.
// When a command fails with output indicating that terraform/terragrunt needs
// initialization or upgrade, it runs "init --upgrade" and retries the original command.
type UpgradeShellRepository interface {
	ExecuteCommandWithUpgrade(command string, arguments []string, directory string) error
}

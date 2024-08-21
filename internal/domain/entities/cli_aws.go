package entities

type CLIAws struct {
	settings *Settings
}

func NewCLIAws(settings *Settings) *CLIAws {
	return &CLIAws{settings: settings}
}

func (it *CLIAws) GetName() string {
	return "aws"
}

func (it *CLIAws) CanChangeAccount() bool {
	return it.settings.TerraAwsRoleArn != ""
}

func (it *CLIAws) GetCommandChangeAccount() []string {
	return []string{"sts", "assume-role", "--role-arn", it.settings.TerraAwsRoleArn, "--role-session-name", "session1"}
}

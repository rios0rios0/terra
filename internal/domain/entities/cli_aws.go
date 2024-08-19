package entities

import "os"

type CLIAws struct{}

func (it CLIAws) GetCLIName() string {
	return "aws"
}

func (it CLIAws) GetCommandChangeAccount() []string {
	return []string{"sts", "assume-role", "--role-arn", getRoleArn(), "--role-session-name", "session1"}
}

func (it CLIAws) CanChangeAccount() bool {
	return getRoleArn() != ""
}

func getRoleArn() string {
	return os.Getenv("TERRA_AWS_ROLE_ARN")
}

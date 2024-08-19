package entities

import (
	"os"
)

type CLIAzm struct{}

func (it CLIAzm) GetCLIName() string {
	return "az"
}

func (it CLIAzm) GetCommandChangeAccount() []string {
	return []string{"account", "set", "--subscription", getSubscriptionID()}
}

func (it CLIAzm) CanChangeAccount() bool {
	return getSubscriptionID() != ""
}

func getSubscriptionID() string {
	return os.Getenv("TERRA_AZURE_SUBSCRIPTION_ID")
}

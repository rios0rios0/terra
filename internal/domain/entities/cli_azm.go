package entities

type CLIAzm struct {
	settings *Settings
}

func NewCLIAzm(settings *Settings) *CLIAzm {
	return &CLIAzm{settings: settings}
}

func (it *CLIAzm) GetName() string {
	return "az"
}

func (it *CLIAzm) CanChangeAccount() bool {
	return it.settings.TerraAzureSubscriptionID != ""
}

func (it *CLIAzm) GetCommandChangeAccount() []string {
	return []string{"account", "set", "--subscription", it.settings.TerraAzureSubscriptionID}
}

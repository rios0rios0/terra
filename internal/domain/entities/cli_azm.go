package entities

import (
	logger "github.com/sirupsen/logrus"
	"os"
)

type CLIAzm struct{}

func (it CLIAzm) ChangeAccount(accountName string) error {
	subscriptionID := ""
	acceptedEnvs := []string{"TERRA_AZURE_SUBSCRIPTION_ID"}
	for _, env := range acceptedEnvs {
		subscriptionID = os.Getenv(env)
		if subscriptionID != "" {
			break
		}
	}

	if subscriptionID != "" {
		err := runInDir("az", []string{"account", "set", "--subscription", subscriptionID}, dir)
		if err != nil {
			logger.Fatalf("Error changing subscription: %s", err)
		}
	}
}

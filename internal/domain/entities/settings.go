package entities

import (
	"github.com/go-playground/validator/v10"
	"github.com/kelseyhightower/envconfig"
	logger "github.com/sirupsen/logrus"
)

type Settings struct {
	TerraCloud               string `envconfig:"TERRA_CLOUD"                 required:"false" validate:"oneof=aws azure"`
	TerraTerraformWorkspace  string `envconfig:"TERRA_WORKSPACE"             required:"false"`
	TerraAwsRoleArn          string `envconfig:"TERRA_AWS_ROLE_ARN"          required:"false"`
	TerraAzureSubscriptionID string `envconfig:"TERRA_AZURE_SUBSCRIPTION_ID" required:"false"`
}

func NewSettings() *Settings {
	var settings Settings
	err := envconfig.Process("", &settings)
	if err != nil {
		logger.Fatalf("Failed to process environment variables: %s", err)
	}

	validate := validator.New()
	err = validate.Struct(settings)
	if err != nil {
		logger.Fatalf("Environment variables validation error: %v", err)
	}

	return &settings
}

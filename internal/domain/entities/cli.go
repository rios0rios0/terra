package entities

import (
	logger "github.com/sirupsen/logrus"
)

type CLI interface {
	GetName() string
	CanChangeAccount() bool
	GetCommandChangeAccount() []string
}

// NewCLI selects the cloud-specific CLI adapter for account-switching commands.
//
// Selection precedence:
//
//  1. Explicit `TERRA_CLOUD` ("aws" | "azure") wins -- backwards-compatible
//     for consumers that already set it. Settings validation (oneof=aws azure)
//     guarantees the value is either one of these or empty.
//  2. Auto-detection from the cloud-specific credential variable: a non-empty
//     `TERRA_AZURE_SUBSCRIPTION_ID` selects the Azure adapter; a non-empty
//     `TERRA_AWS_ROLE_ARN` selects the AWS adapter. This lets consumers wire
//     a single variable in their pipeline / `.env` instead of repeating the
//     cloud name -- the cloud is already implied by which credential they
//     set.
//  3. If both credential variables are populated and `TERRA_CLOUD` is empty,
//     emit a warning and return nil rather than guessing -- the operator is
//     ambiguous, ask them to be explicit.
//  4. Nothing matches -> nil. Downstream call sites already guard with
//     `it.cli != nil && it.cli.CanChangeAccount()` so nil is the no-op
//     value.
func NewCLI(settings *Settings) CLI {
	mapping := map[string]CLI{
		"aws":   NewCLIAws(settings),
		"azure": NewCLIAzm(settings),
	}

	if value, ok := mapping[settings.TerraCloud]; ok {
		return value
	}

	azureSet := settings.TerraAzureSubscriptionID != ""
	awsSet := settings.TerraAwsRoleArn != ""

	switch {
	case azureSet && awsSet:
		logger.Warn(
			"Both TERRA_AZURE_SUBSCRIPTION_ID and TERRA_AWS_ROLE_ARN are set but TERRA_CLOUD is empty; " +
				"set TERRA_CLOUD=azure or TERRA_CLOUD=aws to disambiguate -- " +
				"account-switch commands will be skipped",
		)
		return nil
	case azureSet:
		logger.Debugf("Auto-detected Azure cloud from TERRA_AZURE_SUBSCRIPTION_ID")
		return mapping["azure"]
	case awsSet:
		logger.Debugf("Auto-detected AWS cloud from TERRA_AWS_ROLE_ARN")
		return mapping["aws"]
	default:
		logger.Debugf("No cloud CLI found, avoiding to execute customized commands...")
		return nil
	}
}

//go:build unit

package entities_test

import (
	"testing"

	"github.com/rios0rios0/terra/internal/domain/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCLI(t *testing.T) {
	t.Run("should return nil when no cloud and no credential variables are set", func(t *testing.T) {
		// given
		settings := &entities.Settings{}

		// when
		cli := entities.NewCLI(settings)

		// then
		assert.Nil(t, cli, "without TERRA_CLOUD or any credential variable set, NewCLI must return nil so downstream `it.cli != nil && it.cli.CanChangeAccount()` skips the account-switch command")
	})

	t.Run("should return Azure adapter when TERRA_CLOUD is azure", func(t *testing.T) {
		// given
		settings := &entities.Settings{TerraCloud: "azure"}

		// when
		cli := entities.NewCLI(settings)

		// then
		require.NotNil(t, cli)
		assert.Equal(t, "az", cli.GetName(), "explicit TERRA_CLOUD=azure must yield the Azure adapter")
	})

	t.Run("should return AWS adapter when TERRA_CLOUD is aws", func(t *testing.T) {
		// given
		settings := &entities.Settings{TerraCloud: "aws"}

		// when
		cli := entities.NewCLI(settings)

		// then
		require.NotNil(t, cli)
		assert.Equal(t, "aws", cli.GetName(), "explicit TERRA_CLOUD=aws must yield the AWS adapter")
	})

	t.Run("should auto-detect Azure when only TERRA_AZURE_SUBSCRIPTION_ID is set", func(t *testing.T) {
		// given
		settings := &entities.Settings{
			TerraAzureSubscriptionID: "12345678-1234-1234-1234-123456789012",
		}

		// when
		cli := entities.NewCLI(settings)

		// then
		require.NotNil(t, cli, "TERRA_AZURE_SUBSCRIPTION_ID must be sufficient to select the Azure adapter without TERRA_CLOUD")
		assert.Equal(t, "az", cli.GetName())
		assert.True(t, cli.CanChangeAccount(), "auto-detected Azure CLI should report it can change account when subscription is set")
	})

	t.Run("should auto-detect AWS when only TERRA_AWS_ROLE_ARN is set", func(t *testing.T) {
		// given
		settings := &entities.Settings{
			TerraAwsRoleArn: "arn:aws:iam::123456789012:role/terraform-role",
		}

		// when
		cli := entities.NewCLI(settings)

		// then
		require.NotNil(t, cli, "TERRA_AWS_ROLE_ARN must be sufficient to select the AWS adapter without TERRA_CLOUD")
		assert.Equal(t, "aws", cli.GetName())
		assert.True(t, cli.CanChangeAccount(), "auto-detected AWS CLI should report it can change account when role ARN is set")
	})

	t.Run("should prefer explicit TERRA_CLOUD over conflicting credential variables", func(t *testing.T) {
		// given
		settings := &entities.Settings{
			TerraCloud:               "azure",
			TerraAwsRoleArn:          "arn:aws:iam::123456789012:role/terraform-role",
			TerraAzureSubscriptionID: "12345678-1234-1234-1234-123456789012",
		}

		// when
		cli := entities.NewCLI(settings)

		// then
		require.NotNil(t, cli)
		assert.Equal(t, "az", cli.GetName(), "explicit TERRA_CLOUD must win over auto-detection so existing consumers don't break")
	})

	t.Run("should return nil and warn when both credential variables are set without explicit TERRA_CLOUD", func(t *testing.T) {
		// given
		settings := &entities.Settings{
			TerraAwsRoleArn:          "arn:aws:iam::123456789012:role/terraform-role",
			TerraAzureSubscriptionID: "12345678-1234-1234-1234-123456789012",
		}

		// when
		cli := entities.NewCLI(settings)

		// then
		assert.Nil(t, cli, "ambiguous configuration (both credentials set, no TERRA_CLOUD) should NOT silently pick one -- operator must disambiguate via TERRA_CLOUD")
	})

	t.Run("should fall through to credential-based auto-detection when TERRA_CLOUD has unrecognised value", func(t *testing.T) {
		// given
		settings := &entities.Settings{
			TerraCloud:               "gcp", // unrecognised by the mapping
			TerraAzureSubscriptionID: "12345678-1234-1234-1234-123456789012",
		}

		// when
		cli := entities.NewCLI(settings)

		// then
		require.NotNil(t, cli, "unrecognised TERRA_CLOUD must not block auto-detection -- the credential variable still wins")
		assert.Equal(t, "az", cli.GetName())
	})
}

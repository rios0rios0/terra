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
		// GIVEN:
		settings := &entities.Settings{}

		// WHEN:
		cli := entities.NewCLI(settings)

		// THEN:
		assert.Nil(t, cli, "without TERRA_CLOUD or any credential variable set, NewCLI must return nil so downstream `it.cli != nil && it.cli.CanChangeAccount()` skips the account-switch command")
	})

	t.Run("should return Azure adapter when TERRA_CLOUD is azure", func(t *testing.T) {
		// GIVEN:
		settings := &entities.Settings{TerraCloud: "azure"}

		// WHEN:
		cli := entities.NewCLI(settings)

		// THEN:
		require.NotNil(t, cli)
		assert.Equal(t, "az", cli.GetName(), "explicit TERRA_CLOUD=azure must yield the Azure adapter")
	})

	t.Run("should return AWS adapter when TERRA_CLOUD is aws", func(t *testing.T) {
		// GIVEN:
		settings := &entities.Settings{TerraCloud: "aws"}

		// WHEN:
		cli := entities.NewCLI(settings)

		// THEN:
		require.NotNil(t, cli)
		assert.Equal(t, "aws", cli.GetName(), "explicit TERRA_CLOUD=aws must yield the AWS adapter")
	})

	t.Run("should auto-detect Azure when only TERRA_AZURE_SUBSCRIPTION_ID is set", func(t *testing.T) {
		// GIVEN:
		settings := &entities.Settings{
			TerraAzureSubscriptionID: "12345678-1234-1234-1234-123456789012",
		}

		// WHEN:
		cli := entities.NewCLI(settings)

		// THEN:
		require.NotNil(t, cli, "TERRA_AZURE_SUBSCRIPTION_ID must be sufficient to select the Azure adapter without TERRA_CLOUD")
		assert.Equal(t, "az", cli.GetName())
		assert.True(t, cli.CanChangeAccount(), "auto-detected Azure CLI should report it can change account when subscription is set")
	})

	t.Run("should auto-detect AWS when only TERRA_AWS_ROLE_ARN is set", func(t *testing.T) {
		// GIVEN:
		settings := &entities.Settings{
			TerraAwsRoleArn: "arn:aws:iam::123456789012:role/terraform-role",
		}

		// WHEN:
		cli := entities.NewCLI(settings)

		// THEN:
		require.NotNil(t, cli, "TERRA_AWS_ROLE_ARN must be sufficient to select the AWS adapter without TERRA_CLOUD")
		assert.Equal(t, "aws", cli.GetName())
		assert.True(t, cli.CanChangeAccount(), "auto-detected AWS CLI should report it can change account when role ARN is set")
	})

	t.Run("should prefer explicit TERRA_CLOUD over conflicting credential variables", func(t *testing.T) {
		// GIVEN:
		settings := &entities.Settings{
			TerraCloud:               "azure",
			TerraAwsRoleArn:          "arn:aws:iam::123456789012:role/terraform-role",
			TerraAzureSubscriptionID: "12345678-1234-1234-1234-123456789012",
		}

		// WHEN:
		cli := entities.NewCLI(settings)

		// THEN:
		require.NotNil(t, cli)
		assert.Equal(t, "az", cli.GetName(), "explicit TERRA_CLOUD must win over auto-detection so existing consumers don't break")
	})

	t.Run("should return nil and warn when both credential variables are set without explicit TERRA_CLOUD", func(t *testing.T) {
		// GIVEN:
		settings := &entities.Settings{
			TerraAwsRoleArn:          "arn:aws:iam::123456789012:role/terraform-role",
			TerraAzureSubscriptionID: "12345678-1234-1234-1234-123456789012",
		}

		// WHEN:
		cli := entities.NewCLI(settings)

		// THEN:
		assert.Nil(t, cli, "ambiguous configuration (both credentials set, no TERRA_CLOUD) should NOT silently pick one -- operator must disambiguate via TERRA_CLOUD")
	})
}

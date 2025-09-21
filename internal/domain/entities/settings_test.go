package entities_test

import (
	"testing"

	"github.com/rios0rios0/terra/internal/domain/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSettings_ShouldCreateValidInstance_WhenEmptyTerraCloudProvided(t *testing.T) {
	// GIVEN: Empty TERRA_CLOUD environment variable
	// Note: No environment variable set means empty string

	// WHEN: Creating settings
	settings := entities.NewSettings()

	// THEN: Should create valid settings with empty TerraCloud
	require.NotNil(t, settings)
	assert.Empty(t, settings.TerraCloud)
}

func TestNewSettings_ShouldCreateValidInstance_WhenValidAwsCloudProvided(t *testing.T) {
	// GIVEN: Valid AWS cloud environment variable
	t.Setenv("TERRA_CLOUD", "aws")

	// WHEN: Creating settings
	settings := entities.NewSettings()

	// THEN: Should create valid settings with AWS cloud
	require.NotNil(t, settings)
	assert.Equal(t, "aws", settings.TerraCloud)
}

func TestNewSettings_ShouldCreateValidInstance_WhenValidAzureCloudProvided(t *testing.T) {
	// GIVEN: Valid Azure cloud environment variable
	t.Setenv("TERRA_CLOUD", "azure")

	// WHEN: Creating settings
	settings := entities.NewSettings()

	// THEN: Should create valid settings with Azure cloud
	require.NotNil(t, settings)
	assert.Equal(t, "azure", settings.TerraCloud)
}

func TestNewCLI_ShouldReturnNil_WhenEmptyCloudProviderProvided(t *testing.T) {
	// GIVEN: Settings with empty cloud provider
	settings := entities.NewSettings()

	// WHEN: Creating CLI interface
	cli := entities.NewCLI(settings)

	// THEN: Should return nil for empty cloud provider
	assert.Nil(t, cli)
}

func TestNewCLI_ShouldReturnAwsCLI_WhenAwsCloudProviderProvided(t *testing.T) {
	// GIVEN: Settings with AWS cloud provider
	t.Setenv("TERRA_CLOUD", "aws")
	settings := entities.NewSettings()

	// WHEN: Creating CLI interface
	cli := entities.NewCLI(settings)

	// THEN: Should return AWS CLI implementation
	require.NotNil(t, cli)
	assert.Equal(t, "aws", cli.GetName())
}

func TestNewCLI_ShouldReturnAzureCLI_WhenAzureCloudProviderProvided(t *testing.T) {
	// GIVEN: Settings with Azure cloud provider
	t.Setenv("TERRA_CLOUD", "azure")
	settings := entities.NewSettings()

	// WHEN: Creating CLI interface
	cli := entities.NewCLI(settings)

	// THEN: Should return Azure CLI implementation
	require.NotNil(t, cli)
	assert.Equal(t, "az", cli.GetName())
}

func TestNewCLIAws_ShouldReturnValidInstance_WhenCalled(t *testing.T) {
	// GIVEN: Settings instance
	settings := entities.NewSettings()

	// WHEN: Creating AWS CLI
	cli := entities.NewCLIAws(settings)

	// THEN: Should return valid AWS CLI instance
	require.NotNil(t, cli)
	assert.Equal(t, "aws", cli.GetName())
}

func TestCLIAws_ShouldNotAllowAccountChange_WhenNoRoleArnProvided(t *testing.T) {
	// GIVEN: AWS CLI with settings without role ARN
	settings := entities.NewSettings()
	cli := entities.NewCLIAws(settings)

	// WHEN: Checking if account change is allowed
	canChange := cli.CanChangeAccount()

	// THEN: Should not allow account change
	assert.False(t, canChange)
}

func TestCLIAws_ShouldAllowAccountChange_WhenValidRoleArnProvided(t *testing.T) {
	// GIVEN: AWS CLI with settings containing valid role ARN
	t.Setenv("TERRA_AWS_ROLE_ARN", "arn:aws:iam::123456789012:role/terraform-role")
	settings := entities.NewSettings()
	cli := entities.NewCLIAws(settings)

	// WHEN: Checking if account change is allowed
	canChange := cli.CanChangeAccount()

	// THEN: Should allow account change
	assert.True(t, canChange)
}

func TestCLIAws_ShouldReturnCorrectCommand_WhenGetCommandChangeAccountCalled(t *testing.T) {
	// GIVEN: AWS CLI with valid role ARN
	roleArn := "arn:aws:iam::123456789012:role/terraform-role"
	t.Setenv("TERRA_AWS_ROLE_ARN", roleArn)
	settings := entities.NewSettings()
	cli := entities.NewCLIAws(settings)

	// WHEN: Getting the change account command
	command := cli.GetCommandChangeAccount()

	// THEN: Should return correct AWS STS assume-role command
	require.Len(t, command, 6)
	assert.Equal(t, "sts", command[0])
	assert.Equal(t, "assume-role", command[1])
	assert.Equal(t, "--role-arn", command[2])
	assert.Equal(t, roleArn, command[3])
	assert.Equal(t, "--role-session-name", command[4])
	assert.Equal(t, "session1", command[5])
}

func TestNewCLIAzm_ShouldReturnValidInstance_WhenCalled(t *testing.T) {
	// GIVEN: Settings instance
	settings := entities.NewSettings()

	// WHEN: Creating Azure CLI
	cli := entities.NewCLIAzm(settings)

	// THEN: Should return valid Azure CLI instance
	require.NotNil(t, cli)
	assert.Equal(t, "az", cli.GetName())
}

func TestCLIAzm_ShouldNotAllowAccountChange_WhenNoSubscriptionIdProvided(t *testing.T) {
	// GIVEN: Azure CLI with settings without subscription ID
	settings := entities.NewSettings()
	cli := entities.NewCLIAzm(settings)

	// WHEN: Checking if account change is allowed
	canChange := cli.CanChangeAccount()

	// THEN: Should not allow account change
	assert.False(t, canChange)
}

func TestCLIAzm_ShouldAllowAccountChange_WhenValidSubscriptionIdProvided(t *testing.T) {
	// GIVEN: Azure CLI with settings containing valid subscription ID
	t.Setenv("TERRA_AZURE_SUBSCRIPTION_ID", "12345678-1234-1234-1234-123456789012")
	settings := entities.NewSettings()
	cli := entities.NewCLIAzm(settings)

	// WHEN: Checking if account change is allowed
	canChange := cli.CanChangeAccount()

	// THEN: Should allow account change
	assert.True(t, canChange)
}

func TestCLIAzm_ShouldReturnCorrectCommand_WhenGetCommandChangeAccountCalled(t *testing.T) {
	// GIVEN: Azure CLI with valid subscription ID
	subscriptionID := "12345678-1234-1234-1234-123456789012"
	t.Setenv("TERRA_AZURE_SUBSCRIPTION_ID", subscriptionID)
	settings := entities.NewSettings()
	cli := entities.NewCLIAzm(settings)

	// WHEN: Getting the change account command
	command := cli.GetCommandChangeAccount()

	// THEN: Should return correct Azure account set command
	expectedCommand := []string{"account", "set", "--subscription", subscriptionID}
	assert.Equal(t, expectedCommand, command)
}

// Note: Additional tests that were using table-driven tests with loops have been
// removed in accordance with the contributing guidelines that state:
// "NEVER use loops (for range) to create test cases inside a test method."
// Each test scenario is now a separate, focused test function.

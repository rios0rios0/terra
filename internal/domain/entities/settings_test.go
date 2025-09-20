package entities

import (
	"os"
	"testing"
)

func TestSettings_OptionalEnvironmentVariables(t *testing.T) {
	tests := []struct {
		name       string
		terraCloud string
	}{
		{
			name:       "Empty TERRA_CLOUD should be valid",
			terraCloud: "",
		},
		{
			name:       "Valid AWS cloud should be valid",
			terraCloud: "aws",
		},
		{
			name:       "Valid Azure cloud should be valid",
			terraCloud: "azure",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear any existing environment variables
			os.Unsetenv("TERRA_CLOUD")
			os.Unsetenv("TERRA_WORKSPACE")
			os.Unsetenv("TERRA_AWS_ROLE_ARN")
			os.Unsetenv("TERRA_AZURE_SUBSCRIPTION_ID")

			// Set the test value if not empty
			if tt.terraCloud != "" {
				os.Setenv("TERRA_CLOUD", tt.terraCloud)
			}

			settings := NewSettings()

			// Verify the setting was loaded correctly
			if settings.TerraCloud != tt.terraCloud {
				t.Errorf("Expected TerraCloud to be %q, got %q", tt.terraCloud, settings.TerraCloud)
			}
		})
	}
}

func TestCLI_OptionalCloudProvider(t *testing.T) {
	tests := []struct {
		name         string
		terraCloud   string
		expectNilCLI bool
		expectedName string
	}{
		{
			name:         "Empty TERRA_CLOUD should return nil CLI",
			terraCloud:   "",
			expectNilCLI: true,
		},
		{
			name:         "AWS cloud should return AWS CLI",
			terraCloud:   "aws",
			expectNilCLI: false,
			expectedName: "aws",
		},
		{
			name:         "Azure cloud should return Azure CLI",
			terraCloud:   "azure",
			expectNilCLI: false,
			expectedName: "az",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear environment variables
			os.Unsetenv("TERRA_CLOUD")
			os.Unsetenv("TERRA_AWS_ROLE_ARN")
			os.Unsetenv("TERRA_AZURE_SUBSCRIPTION_ID")

			if tt.terraCloud != "" {
				os.Setenv("TERRA_CLOUD", tt.terraCloud)
			}

			settings := NewSettings()
			cli := NewCLI(settings)

			if tt.expectNilCLI {
				if cli != nil {
					t.Error("Expected nil CLI but got non-nil")
				}
			} else {
				if cli == nil {
					t.Error("Expected non-nil CLI but got nil")
				} else if cli.GetName() != tt.expectedName {
					t.Errorf("Expected CLI name %q, got %q", tt.expectedName, cli.GetName())
				}
			}
		})
	}
}

func TestCLIAws_CanChangeAccount(t *testing.T) {
	tests := []struct {
		name               string
		terraAwsRoleArn    string
		expectedCanChange  bool
	}{
		{
			name:              "Empty role ARN should return false",
			terraAwsRoleArn:   "",
			expectedCanChange: false,
		},
		{
			name:              "Valid role ARN should return true",
			terraAwsRoleArn:   "arn:aws:iam::123456789012:role/terraform-role",
			expectedCanChange: true,
		},
		{
			name:              "Invalid role ARN should still return true",
			terraAwsRoleArn:   "invalid-arn",
			expectedCanChange: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear environment variables
			os.Unsetenv("TERRA_AWS_ROLE_ARN")

			if tt.terraAwsRoleArn != "" {
				os.Setenv("TERRA_AWS_ROLE_ARN", tt.terraAwsRoleArn)
			}

			settings := NewSettings()
			cli := NewCLIAws(settings)

			canChange := cli.CanChangeAccount()
			if canChange != tt.expectedCanChange {
				t.Errorf("Expected CanChangeAccount to be %v, got %v", tt.expectedCanChange, canChange)
			}
		})
	}
}

func TestCLIAws_GetCommandChangeAccount(t *testing.T) {
	roleArn := "arn:aws:iam::123456789012:role/terraform-role"
	
	// Clear and set environment variable
	os.Unsetenv("TERRA_AWS_ROLE_ARN")
	os.Setenv("TERRA_AWS_ROLE_ARN", roleArn)

	settings := NewSettings()
	cli := NewCLIAws(settings)

	command := cli.GetCommandChangeAccount()
	expectedCommand := []string{"sts", "assume-role", "--role-arn", roleArn, "--role-session-name", "session1"}

	if len(command) != len(expectedCommand) {
		t.Errorf("Expected command length %d, got %d", len(expectedCommand), len(command))
	}

	for i, expected := range expectedCommand {
		if i < len(command) && command[i] != expected {
			t.Errorf("Expected command[%d] to be %q, got %q", i, expected, command[i])
		}
	}
}

func TestCLIAzm_CanChangeAccount(t *testing.T) {
	tests := []struct {
		name                       string
		terraAzureSubscriptionID   string
		expectedCanChange          bool
	}{
		{
			name:                     "Empty subscription ID should return false",
			terraAzureSubscriptionID: "",
			expectedCanChange:        false,
		},
		{
			name:                     "Valid subscription ID should return true",
			terraAzureSubscriptionID: "12345678-1234-1234-1234-123456789012",
			expectedCanChange:        true,
		},
		{
			name:                     "Invalid subscription ID should still return true",
			terraAzureSubscriptionID: "invalid-subscription-id",
			expectedCanChange:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear environment variables
			os.Unsetenv("TERRA_AZURE_SUBSCRIPTION_ID")

			if tt.terraAzureSubscriptionID != "" {
				os.Setenv("TERRA_AZURE_SUBSCRIPTION_ID", tt.terraAzureSubscriptionID)
			}

			settings := NewSettings()
			cli := NewCLIAzm(settings)

			canChange := cli.CanChangeAccount()
			if canChange != tt.expectedCanChange {
				t.Errorf("Expected CanChangeAccount to be %v, got %v", tt.expectedCanChange, canChange)
			}
		})
	}
}

func TestCLIAzm_GetCommandChangeAccount(t *testing.T) {
	subscriptionID := "12345678-1234-1234-1234-123456789012"
	
	// Clear and set environment variable
	os.Unsetenv("TERRA_AZURE_SUBSCRIPTION_ID")
	os.Setenv("TERRA_AZURE_SUBSCRIPTION_ID", subscriptionID)

	settings := NewSettings()
	cli := NewCLIAzm(settings)

	command := cli.GetCommandChangeAccount()
	expectedCommand := []string{"account", "set", "--subscription", subscriptionID}

	if len(command) != len(expectedCommand) {
		t.Errorf("Expected command length %d, got %d", len(expectedCommand), len(command))
	}

	for i, expected := range expectedCommand {
		if i < len(command) && command[i] != expected {
			t.Errorf("Expected command[%d] to be %q, got %q", i, expected, command[i])
		}
	}
}

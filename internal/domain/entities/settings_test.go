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

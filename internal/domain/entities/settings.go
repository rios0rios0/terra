package entities

import (
	"fmt"
	"os"
	"path/filepath"

	validator "github.com/go-playground/validator/v10"
	"github.com/kelseyhightower/envconfig"
	logger "github.com/sirupsen/logrus"
)

type Settings struct {
	TerraCloud               string `envconfig:"TERRA_CLOUD"                 required:"false" validate:"omitempty,oneof=aws azure"`
	TerraTerraformWorkspace  string `envconfig:"TERRA_WORKSPACE"             required:"false"`
	TerraAwsRoleArn          string `envconfig:"TERRA_AWS_ROLE_ARN"          required:"false"`
	TerraAzureSubscriptionID string `envconfig:"TERRA_AZURE_SUBSCRIPTION_ID" required:"false"`
	TerraModuleCacheDir      string `envconfig:"TERRA_MODULE_CACHE_DIR"      required:"false"`
	TerraProviderCacheDir    string `envconfig:"TERRA_PROVIDER_CACHE_DIR"    required:"false"`
	TerraNoCAS               bool   `envconfig:"TERRA_NO_CAS"                required:"false"`
}

// GetModuleCacheDir returns the module cache directory path.
// It uses the configured value or falls back to ~/.cache/terra/modules.
func (s *Settings) GetModuleCacheDir() (string, error) {
	if s.TerraModuleCacheDir != "" {
		return s.TerraModuleCacheDir, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to determine home directory: %w", err)
	}

	return filepath.Join(home, ".cache", "terra", "modules"), nil
}

// GetProviderCacheDir returns the provider cache directory path.
// It uses the configured value or falls back to ~/.cache/terra/providers.
func (s *Settings) GetProviderCacheDir() (string, error) {
	if s.TerraProviderCacheDir != "" {
		return s.TerraProviderCacheDir, nil
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to determine home directory: %w", err)
	}

	return filepath.Join(home, ".cache", "terra", "providers"), nil
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

//nolint:testpackage // Testing private functions and fields requires same package
package controllers

import (
	"testing"

	"github.com/rios0rios0/terra/internal/domain/entities"
	"github.com/spf13/cobra"
)

// MockInstallDependenciesCommand is a mock implementation of the InstallDependencies interface
type MockInstallDependenciesCommand struct {
	ExecuteCallCount int
	LastDependencies []entities.Dependency
}

func (m *MockInstallDependenciesCommand) Execute(dependencies []entities.Dependency) {
	m.ExecuteCallCount++
	m.LastDependencies = dependencies
}

func TestNewInstallDependenciesController(t *testing.T) {
	mockCommand := &MockInstallDependenciesCommand{}
	dependencies := []entities.Dependency{
		{
			Name:              "Test Dependency",
			CLI:               "test",
			BinaryURL:         "https://example.com/test",
			VersionURL:        "https://example.com/version",
			RegexVersion:      `"version":"([^"]+)"`,
			FormattingCommand: []string{"format"},
		},
	}

	controller := NewInstallDependenciesController(mockCommand, dependencies)

	if controller == nil {
		t.Fatal("NewInstallDependenciesController returned nil")
	}

	if controller.command != mockCommand {
		t.Error("Controller command was not set correctly")
	}

	if len(controller.dependencies) != len(dependencies) {
		t.Errorf(
			"Expected %d dependencies, got %d",
			len(dependencies),
			len(controller.dependencies),
		)
	}

	if controller.dependencies[0].Name != dependencies[0].Name {
		t.Errorf(
			"Expected dependency name %s, got %s",
			dependencies[0].Name,
			controller.dependencies[0].Name,
		)
	}
}

func TestInstallDependenciesController_GetBind(t *testing.T) {
	mockCommand := &MockInstallDependenciesCommand{}
	dependencies := []entities.Dependency{}

	controller := NewInstallDependenciesController(mockCommand, dependencies)
	bind := controller.GetBind()

	expectedUse := "install"
	expectedShort := "Install or update Terraform and Terragrunt to the latest versions"
	expectedLong := "Install all the dependencies required to run Terra, or update them if newer versions are available. Dependencies are installed to ~/.local/bin on Linux."

	if bind.Use != expectedUse {
		t.Errorf("Expected Use to be %q, got %q", expectedUse, bind.Use)
	}

	if bind.Short != expectedShort {
		t.Errorf("Expected Short to be %q, got %q", expectedShort, bind.Short)
	}

	if bind.Long != expectedLong {
		t.Errorf("Expected Long to be %q, got %q", expectedLong, bind.Long)
	}
}

func TestInstallDependenciesController_Execute(t *testing.T) {
	mockCommand := &MockInstallDependenciesCommand{}
	dependencies := []entities.Dependency{
		{
			Name: "Test Dependency",
			CLI:  "test",
		},
		{
			Name: "Another Dependency",
			CLI:  "another",
		},
	}

	controller := NewInstallDependenciesController(mockCommand, dependencies)

	// Create a mock cobra command and args
	cmd := &cobra.Command{}
	args := []string{}

	// Execute the controller
	controller.Execute(cmd, args)

	// Verify that the command was called
	if mockCommand.ExecuteCallCount != 1 {
		t.Errorf("Expected Execute to be called once, got %d calls", mockCommand.ExecuteCallCount)
	}

	// Verify that the correct dependencies were passed
	if len(mockCommand.LastDependencies) != len(dependencies) {
		t.Errorf("Expected %d dependencies passed to Execute, got %d",
			len(dependencies), len(mockCommand.LastDependencies))
	}

	if mockCommand.LastDependencies[0].Name != dependencies[0].Name {
		t.Errorf("Expected first dependency name %s, got %s",
			dependencies[0].Name, mockCommand.LastDependencies[0].Name)
	}

	if mockCommand.LastDependencies[1].Name != dependencies[1].Name {
		t.Errorf("Expected second dependency name %s, got %s",
			dependencies[1].Name, mockCommand.LastDependencies[1].Name)
	}
}

func TestInstallDependenciesController_ExecuteMultipleCalls(t *testing.T) {
	mockCommand := &MockInstallDependenciesCommand{}
	dependencies := []entities.Dependency{
		{Name: "Test", CLI: "test"},
	}

	controller := NewInstallDependenciesController(mockCommand, dependencies)
	cmd := &cobra.Command{}
	args := []string{}

	// Execute multiple times
	controller.Execute(cmd, args)
	controller.Execute(cmd, args)
	controller.Execute(cmd, args)

	// Verify that the command was called the correct number of times
	if mockCommand.ExecuteCallCount != 3 {
		t.Errorf(
			"Expected Execute to be called 3 times, got %d calls",
			mockCommand.ExecuteCallCount,
		)
	}
}

package controllers

import (
	"testing"

	"github.com/rios0rios0/terra/internal/domain/entities"
	"github.com/spf13/cobra"
)

// MockRunFromRootCommand is a mock implementation of the RunFromRoot interface
type MockRunFromRootCommand struct {
	ExecuteCallCount int
	LastTargetPath   string
	LastArguments    []string
	LastDependencies []entities.Dependency
}

func (m *MockRunFromRootCommand) Execute(
	targetPath string,
	arguments []string,
	dependencies []entities.Dependency,
) {
	m.ExecuteCallCount++
	m.LastTargetPath = targetPath
	m.LastArguments = arguments
	m.LastDependencies = dependencies
}

func TestNewRunFromRootController(t *testing.T) {
	mockCommand := &MockRunFromRootCommand{}
	dependencies := []entities.Dependency{
		{
			Name: "Test Dependency",
			CLI:  "test",
		},
	}

	controller := NewRunFromRootController(mockCommand, dependencies)

	if controller == nil {
		t.Fatal("NewRunFromRootController returned nil")
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

func TestRunFromRootController_GetBind(t *testing.T) {
	mockCommand := &MockRunFromRootCommand{}
	dependencies := []entities.Dependency{}

	controller := NewRunFromRootController(mockCommand, dependencies)
	bind := controller.GetBind()

	expectedUse := "terra [flags] [terragrunt command] [directory]"
	expectedShort := "Terra is a CLI wrapper for Terragrunt"
	expectedLong := "Terra is a CLI wrapper for Terragrunt that allows changing directory before executing commands. It also allows changing the account/subscription and workspace for AWS and Azure."

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

func TestRunFromRootController_Execute(t *testing.T) {
	mockCommand := &MockRunFromRootCommand{}
	dependencies := []entities.Dependency{
		{
			Name: "Terraform",
			CLI:  "terraform",
		},
	}

	controller := NewRunFromRootController(mockCommand, dependencies)

	// Create a mock cobra command and args - use current directory "." to avoid path issues
	cmd := &cobra.Command{}
	args := []string{"apply", "."}

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
		t.Errorf("Expected dependency name %s, got %s",
			dependencies[0].Name, mockCommand.LastDependencies[0].Name)
	}

	// Note: The actual path handling and argument filtering is done by helpers.ArgumentsHelper
	// so we verify that the command was called with some arguments
	if len(mockCommand.LastArguments) == 0 {
		t.Error("Expected some arguments to be passed to command")
	}
}

func TestRunFromRootController_ExecuteWithDifferentArgs(t *testing.T) {
	mockCommand := &MockRunFromRootCommand{}
	dependencies := []entities.Dependency{}

	controller := NewRunFromRootController(mockCommand, dependencies)
	cmd := &cobra.Command{}
	args := []string{"plan", "--dry-run"}

	// Execute the controller
	controller.Execute(cmd, args)

	// Verify that the command was called
	if mockCommand.ExecuteCallCount != 1 {
		t.Errorf("Expected Execute to be called once, got %d calls", mockCommand.ExecuteCallCount)
	}

	// Verify that some arguments were passed (after helper processing)
	if len(mockCommand.LastArguments) == 0 {
		t.Error("Expected some arguments to be passed to command")
	}
}

func TestRunFromRootController_ExecuteMultipleCalls(t *testing.T) {
	mockCommand := &MockRunFromRootCommand{}
	dependencies := []entities.Dependency{}

	controller := NewRunFromRootController(mockCommand, dependencies)
	cmd := &cobra.Command{}
	args := []string{"plan"}

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

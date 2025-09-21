package controllers_test

import (
	"testing"

	"github.com/rios0rios0/terra/internal/infrastructure/controllers"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockVersionCommand is a mock implementation of the Version interface
type MockVersionCommand struct {
	ExecuteCallCount int
}

func (m *MockVersionCommand) Execute() {
	m.ExecuteCallCount++
}

func TestNewVersionController_ShouldCreateInstance_WhenCommandProvided(t *testing.T) {
	// GIVEN: A mock version command
	mockCommand := &MockVersionCommand{}

	// WHEN: Creating a new version controller
	controller := controllers.NewVersionController(mockCommand)

	// THEN: Should create a valid controller instance
	require.NotNil(t, controller)
}

func TestVersionController_GetBind(t *testing.T) {
	mockCommand := &MockVersionCommand{}
	controller := NewVersionController(mockCommand)
	bind := controller.GetBind()

	expectedUse := "version"
	expectedShort := "Show Terra, Terraform, and Terragrunt versions"
	expectedLong := "Display the version information for Terra and its dependencies (Terraform and Terragrunt)."

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

func TestVersionController_Execute(t *testing.T) {
	mockCommand := &MockVersionCommand{}
	controller := NewVersionController(mockCommand)

	// Create a mock cobra command and args
	cmd := &cobra.Command{}
	args := []string{}

	// Execute the controller
	controller.Execute(cmd, args)

	// Verify that the command was called
	if mockCommand.ExecuteCallCount != 1 {
		t.Errorf("Expected Execute to be called once, got %d calls", mockCommand.ExecuteCallCount)
	}
}

func TestVersionController_ExecuteMultipleCalls(t *testing.T) {
	mockCommand := &MockVersionCommand{}
	controller := NewVersionController(mockCommand)
	cmd := &cobra.Command{}
	args := []string{}

	// Execute multiple times
	controller.Execute(cmd, args)
	controller.Execute(cmd, args)

	// Verify that the command was called the correct number of times
	if mockCommand.ExecuteCallCount != 2 {
		t.Errorf(
			"Expected Execute to be called 2 times, got %d calls",
			mockCommand.ExecuteCallCount,
		)
	}
}

package controllers_test

import (
	"testing"

	"github.com/spf13/cobra"
)

// MockVersionCommand is a mock implementation of the Version interface
type MockVersionCommand struct {
	ExecuteCallCount int
}

func (m *MockVersionCommand) Execute() {
	m.ExecuteCallCount++
}

func TestNewVersionController(t *testing.T) {
	mockCommand := &MockVersionCommand{}
	controller := NewVersionController(mockCommand)

	if controller == nil {
		t.Fatal("NewVersionController returned nil")
	}

	if controller.command != mockCommand {
		t.Error("Controller command was not set correctly")
	}
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

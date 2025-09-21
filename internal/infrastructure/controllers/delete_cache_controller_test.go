package controllers_test

import (
	"testing"

	"github.com/spf13/cobra"
)

// MockDeleteCacheCommand is a mock implementation of the DeleteCache interface
type MockDeleteCacheCommand struct {
	ExecuteCallCount int
	LastToBeDeleted  []string
}

func (m *MockDeleteCacheCommand) Execute(toBeDeleted []string) {
	m.ExecuteCallCount++
	m.LastToBeDeleted = toBeDeleted
}

func TestNewDeleteCacheController(t *testing.T) {
	mockCommand := &MockDeleteCacheCommand{}
	controller := NewDeleteCacheController(mockCommand)

	if controller == nil {
		t.Fatal("NewDeleteCacheController returned nil")
	}

	if controller.command != mockCommand {
		t.Error("Controller command was not set correctly")
	}
}

func TestDeleteCacheController_GetBind(t *testing.T) {
	mockCommand := &MockDeleteCacheCommand{}
	controller := NewDeleteCacheController(mockCommand)
	bind := controller.GetBind()

	expectedUse := "clear"
	expectedShort := "Clear all cache and modules directories"
	expectedLong := "Clear all temporary directories and cache folders created during the Terraform and Terragrunt execution."

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

func TestDeleteCacheController_Execute(t *testing.T) {
	mockCommand := &MockDeleteCacheCommand{}
	controller := NewDeleteCacheController(mockCommand)

	// Create a mock cobra command and args
	cmd := &cobra.Command{}
	args := []string{}

	// Execute the controller
	controller.Execute(cmd, args)

	// Verify that the command was called
	if mockCommand.ExecuteCallCount != 1 {
		t.Errorf("Expected Execute to be called once, got %d calls", mockCommand.ExecuteCallCount)
	}

	// Verify that the correct directories were passed
	expectedDirs := []string{".terraform", ".terragrunt-cache"}
	if len(mockCommand.LastToBeDeleted) != len(expectedDirs) {
		t.Errorf(
			"Expected %d directories to be deleted, got %d",
			len(expectedDirs),
			len(mockCommand.LastToBeDeleted),
		)
	}

	for i, expected := range expectedDirs {
		if i < len(mockCommand.LastToBeDeleted) && mockCommand.LastToBeDeleted[i] != expected {
			t.Errorf(
				"Expected directory %s at index %d, got %s",
				expected,
				i,
				mockCommand.LastToBeDeleted[i],
			)
		}
	}
}

func TestDeleteCacheController_ExecuteMultipleCalls(t *testing.T) {
	mockCommand := &MockDeleteCacheCommand{}
	controller := NewDeleteCacheController(mockCommand)
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

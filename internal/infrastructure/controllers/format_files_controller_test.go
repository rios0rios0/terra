package controllers

import (
	"testing"

	"github.com/rios0rios0/terra/internal/domain/entities"
	"github.com/spf13/cobra"
)

// MockFormatFilesCommand is a mock implementation of the FormatFiles interface
type MockFormatFilesCommand struct {
	ExecuteCallCount int
	LastDependencies []entities.Dependency
}

func (m *MockFormatFilesCommand) Execute(dependencies []entities.Dependency) {
	m.ExecuteCallCount++
	m.LastDependencies = dependencies
}

func TestNewFormatFilesController(t *testing.T) {
	mockCommand := &MockFormatFilesCommand{}
	dependencies := []entities.Dependency{
		{
			Name:              "Test Tool",
			CLI:               "test",
			FormattingCommand: []string{"format", "-recursive"},
		},
	}

	controller := NewFormatFilesController(mockCommand, dependencies)

	if controller == nil {
		t.Fatal("NewFormatFilesController returned nil")
	}

	if controller.command != mockCommand {
		t.Error("Controller command was not set correctly")
	}

	if len(controller.dependencies) != len(dependencies) {
		t.Errorf("Expected %d dependencies, got %d", len(dependencies), len(controller.dependencies))
	}

	if controller.dependencies[0].Name != dependencies[0].Name {
		t.Errorf("Expected dependency name %s, got %s", dependencies[0].Name, controller.dependencies[0].Name)
	}
}

func TestFormatFilesController_GetBind(t *testing.T) {
	mockCommand := &MockFormatFilesCommand{}
	dependencies := []entities.Dependency{}

	controller := NewFormatFilesController(mockCommand, dependencies)
	bind := controller.GetBind()

	expectedUse := "format"
	expectedShort := "Format all files in the current directory"
	expectedLong := "Format all the Terraform and Terragrunt files in the current directory."

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

func TestFormatFilesController_Execute(t *testing.T) {
	mockCommand := &MockFormatFilesCommand{}
	dependencies := []entities.Dependency{
		{
			Name: "Terraform",
			CLI:  "terraform",
		},
		{
			Name: "Terragrunt",
			CLI:  "terragrunt",
		},
	}

	controller := NewFormatFilesController(mockCommand, dependencies)

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

func TestFormatFilesController_ExecuteMultipleCalls(t *testing.T) {
	mockCommand := &MockFormatFilesCommand{}
	dependencies := []entities.Dependency{}

	controller := NewFormatFilesController(mockCommand, dependencies)
	cmd := &cobra.Command{}
	args := []string{}

	// Execute multiple times
	controller.Execute(cmd, args)
	controller.Execute(cmd, args)

	// Verify that the command was called the correct number of times
	if mockCommand.ExecuteCallCount != 2 {
		t.Errorf("Expected Execute to be called 2 times, got %d calls", mockCommand.ExecuteCallCount)
	}
}
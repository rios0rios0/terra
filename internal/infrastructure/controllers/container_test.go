package controllers_test

import (
	"testing"
)

func TestNewControllers(t *testing.T) {
	// Create mock controllers
	deleteCacheController := &DeleteCacheController{}
	formatFilesController := &FormatFilesController{}
	installDependenciesController := &InstallDependenciesController{}
	versionController := &VersionController{}

	// Test the NewControllers function
	controllers := NewControllers(
		deleteCacheController,
		formatFilesController,
		installDependenciesController,
		versionController,
	)

	if controllers == nil {
		t.Fatal("NewControllers returned nil")
	}

	// Verify the returned slice contains the correct number of controllers
	expectedCount := 4
	if len(*controllers) != expectedCount {
		t.Errorf("Expected %d controllers, got %d", expectedCount, len(*controllers))
	}

	// Verify the controllers are in the expected order
	controllerSlice := *controllers

	if controllerSlice[0] != deleteCacheController {
		t.Error("Expected first controller to be deleteCacheController")
	}

	if controllerSlice[1] != formatFilesController {
		t.Error("Expected second controller to be formatFilesController")
	}

	if controllerSlice[2] != installDependenciesController {
		t.Error("Expected third controller to be installDependenciesController")
	}

	if controllerSlice[3] != versionController {
		t.Error("Expected fourth controller to be versionController")
	}
}

func TestNewControllers_InterfaceConformance(t *testing.T) {
	// Create mock controllers
	deleteCacheController := &DeleteCacheController{}
	formatFilesController := &FormatFilesController{}
	installDependenciesController := &InstallDependenciesController{}
	versionController := &VersionController{}

	// Test the NewControllers function
	controllers := NewControllers(
		deleteCacheController,
		formatFilesController,
		installDependenciesController,
		versionController,
	)

	// Verify each controller implements the Controller interface
	for i, controller := range *controllers {
		// Check if the controller is not nil
		if controller == nil {
			t.Errorf("Controller at index %d is nil", i)
		}
	}
}

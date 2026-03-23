//go:build unit

package main

import (
	"testing"

	"github.com/rios0rios0/terra/internal"
	"github.com/rios0rios0/terra/internal/domain/entities"
	"github.com/rios0rios0/terra/internal/infrastructure/controllers"
	"github.com/rios0rios0/terra/test/domain/commanddoubles"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAppContext(t *testing.T) {
	t.Parallel()

	t.Run("should return AppContext when valid AppInternal provided", func(t *testing.T) {
		t.Parallel()
		// GIVEN
		controllersList := []entities.Controller{}
		appInternal := internal.NewAppInternal(&controllersList)

		// WHEN
		result := newAppContext(appInternal)

		// THEN
		require.NotNil(t, result)
		assert.Empty(t, result.GetControllers())
	})
}

func TestInjectAppContext(t *testing.T) {
	t.Run("should resolve full DIG container when all providers registered", func(t *testing.T) {
		// GIVEN: The full DIG container with all real providers

		// WHEN: Injecting the app context (exercises the entire registration chain)
		result := injectAppContext()

		// THEN: Should return a valid AppContext with controllers
		require.NotNil(t, result)
		assert.NotEmpty(t, result.GetControllers())
	})
}

func TestInjectRootController(t *testing.T) {
	t.Run("should resolve root controller from DIG container when all providers registered", func(t *testing.T) {
		// GIVEN: The full DIG container with all real providers

		// WHEN: Injecting the root controller
		result := injectRootController()

		// THEN: Should return a valid Controller
		require.NotNil(t, result)
		assert.Contains(t, result.GetBind().Use, "terra")
	})
}

func TestNewRootController(t *testing.T) {
	t.Parallel()

	t.Run("should return Controller when valid RunFromRootController provided", func(t *testing.T) {
		t.Parallel()
		// GIVEN
		stubCommand := &commanddoubles.StubRunFromRootCommand{}
		dependencies := []entities.Dependency{}
		ctrl := controllers.NewRunFromRootController(stubCommand, dependencies)

		// WHEN
		result := newRootController(ctrl)

		// THEN
		require.NotNil(t, result)
		assert.Equal(t, "terra [flags] [terragrunt command] [directory]", result.GetBind().Use)
	})
}

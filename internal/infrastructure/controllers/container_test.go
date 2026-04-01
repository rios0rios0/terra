//go:build unit

package controllers_test

import (
	"testing"

	"github.com/rios0rios0/terra/internal/domain/entities"
	"github.com/rios0rios0/terra/internal/infrastructure/controllers"
	"github.com/rios0rios0/terra/test/domain/commanddoubles"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/dig"
)

func TestRegisterProviders(t *testing.T) {
	t.Parallel()

	t.Run("should register all controller providers when container is valid", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A fresh DIG container
		container := dig.New()

		// WHEN: Registering all controller providers
		err := controllers.RegisterProviders(container)

		// THEN: Should succeed without error
		assert.NoError(t, err)
	})

	t.Run("should return error when called twice on same container", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A container that already has all providers registered
		container := dig.New()
		firstErr := controllers.RegisterProviders(container)
		require.NoError(t, firstErr)

		// WHEN: Attempting to register all providers again
		err := controllers.RegisterProviders(container)

		// THEN: Should return an error due to duplicate registration
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already provided")
	})
}

func TestNewControllers(t *testing.T) {
	t.Parallel()

	t.Run("should create controllers list when all controllers provided", func(t *testing.T) {
		t.Parallel()
		// given
		deps := []entities.Dependency{}
		deleteCache := controllers.NewDeleteCacheController(&commanddoubles.StubDeleteCacheCommand{})
		formatFiles := controllers.NewFormatFilesController(&commanddoubles.StubFormatFilesCommand{}, deps)
		installDeps := controllers.NewInstallDependenciesController(
			&commanddoubles.StubInstallDependenciesCommand{}, deps,
		)
		updateDeps := controllers.NewUpdateDependenciesController(
			&commanddoubles.StubInstallDependenciesCommand{}, deps,
		)
		selfUpdate := controllers.NewSelfUpdateController(&commanddoubles.StubSelfUpdateCommand{})
		version := controllers.NewVersionController(&commanddoubles.StubVersionCommand{})

		// when
		result := controllers.NewControllers(
			deleteCache, formatFiles, installDeps, updateDeps, selfUpdate, version,
		)

		// then
		require.NotNil(t, result)
		assert.Len(t, *result, 6)
	})
}

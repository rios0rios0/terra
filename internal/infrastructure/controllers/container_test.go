//go:build unit

package controllers_test

import (
	"testing"

	"github.com/rios0rios0/terra/internal/domain/entities"
	"github.com/rios0rios0/terra/internal/infrastructure/controllers"
	"github.com/rios0rios0/terra/test/domain/commanddoubles"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

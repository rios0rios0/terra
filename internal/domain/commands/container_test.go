//go:build unit

package commands_test

import (
	"testing"

	"github.com/rios0rios0/terra/internal/domain/commands"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/dig"
)

func TestRegisterProviders(t *testing.T) {
	t.Parallel()

	t.Run("should register all providers when container is valid", func(t *testing.T) {
		t.Parallel()
		// GIVEN
		container := dig.New()

		// WHEN
		err := commands.RegisterProviders(container)

		// THEN
		assert.NoError(t, err)
	})

	t.Run("should return error when duplicate provider conflicts with existing registration", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A container where the first provider type is already registered
		container := dig.New()
		// Pre-register a conflicting provider for *DeleteCacheCommand (first Provide call)
		preRegErr := container.Provide(commands.NewDeleteCacheCommand)
		require.NoError(t, preRegErr)

		// WHEN: Attempting to register all providers
		err := commands.RegisterProviders(container)

		// THEN: Should return an error due to conflicting provider
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already provided")
	})

	t.Run("should return error when called twice on same container", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A container that already has all providers registered
		container := dig.New()
		firstErr := commands.RegisterProviders(container)
		require.NoError(t, firstErr)

		// WHEN: Attempting to register all providers again
		err := commands.RegisterProviders(container)

		// THEN: Should return an error due to duplicate registration
		assert.Error(t, err)
	})

	t.Run("should return error when FormatFilesCommand conflicts with pre-registered provider", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A container where FormatFilesCommand is already registered
		container := dig.New()
		// Pre-register the first constructor to let it pass, then conflict on the second
		require.NoError(t, container.Provide(commands.NewDeleteCacheCommand))
		require.NoError(t, container.Provide(commands.NewFormatFilesCommand))

		// WHEN: Attempting to register all providers
		err := commands.RegisterProviders(container)

		// THEN: Should return an error when hitting the conflicting FormatFilesCommand
		assert.Error(t, err)
	})

	t.Run("should return error when ParallelStateCommand conflicts with pre-registered provider", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A container where multiple constructors are pre-registered
		container := dig.New()
		require.NoError(t, container.Provide(commands.NewDeleteCacheCommand))
		require.NoError(t, container.Provide(commands.NewFormatFilesCommand))
		require.NoError(t, container.Provide(commands.NewInstallDependenciesCommand))
		require.NoError(t, container.Provide(commands.NewParallelStateCommand))

		// WHEN: Attempting to register all providers
		err := commands.RegisterProviders(container)

		// THEN: Should return an error when hitting the conflicting constructor
		assert.Error(t, err)
	})

	t.Run("should return error when RunFromRootCommand conflicts with pre-registered provider", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A container where constructors up to RunFromRoot are pre-registered
		container := dig.New()
		require.NoError(t, container.Provide(commands.NewDeleteCacheCommand))
		require.NoError(t, container.Provide(commands.NewFormatFilesCommand))
		require.NoError(t, container.Provide(commands.NewInstallDependenciesCommand))
		require.NoError(t, container.Provide(commands.NewParallelStateCommand))
		require.NoError(t, container.Provide(commands.NewRunAdditionalBeforeCommand))
		require.NoError(t, container.Provide(commands.NewRunFromRootCommand))

		// WHEN: Attempting to register all providers
		err := commands.RegisterProviders(container)

		// THEN: Should return an error when hitting the conflicting constructor
		assert.Error(t, err)
	})

	t.Run("should return error when VersionCommand conflicts with pre-registered provider", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A container where all constructors are pre-registered (but not interface binds)
		container := dig.New()
		require.NoError(t, container.Provide(commands.NewDeleteCacheCommand))
		require.NoError(t, container.Provide(commands.NewFormatFilesCommand))
		require.NoError(t, container.Provide(commands.NewInstallDependenciesCommand))
		require.NoError(t, container.Provide(commands.NewParallelStateCommand))
		require.NoError(t, container.Provide(commands.NewRunAdditionalBeforeCommand))
		require.NoError(t, container.Provide(commands.NewRunFromRootCommand))
		require.NoError(t, container.Provide(commands.NewSelfUpdateCommand))
		require.NoError(t, container.Provide(commands.NewVersionCommand))

		// WHEN: Attempting to register all providers
		err := commands.RegisterProviders(container)

		// THEN: Should return an error when hitting the conflicting constructor
		assert.Error(t, err)
	})

	t.Run("should return error when DeleteCache interface binding conflicts", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A container with all struct constructors already registered + first interface
		container := dig.New()
		require.NoError(t, container.Provide(commands.NewDeleteCacheCommand))
		require.NoError(t, container.Provide(commands.NewFormatFilesCommand))
		require.NoError(t, container.Provide(commands.NewInstallDependenciesCommand))
		require.NoError(t, container.Provide(commands.NewParallelStateCommand))
		require.NoError(t, container.Provide(commands.NewRunAdditionalBeforeCommand))
		require.NoError(t, container.Provide(commands.NewRunFromRootCommand))
		require.NoError(t, container.Provide(commands.NewSelfUpdateCommand))
		require.NoError(t, container.Provide(commands.NewVersionCommand))
		// Pre-register the DeleteCache interface binding
		require.NoError(t, container.Provide(func(impl *commands.DeleteCacheCommand) commands.DeleteCache {
			return impl
		}))

		// WHEN: Attempting to register all providers
		err := commands.RegisterProviders(container)

		// THEN: Should return error on the first struct constructor duplicate
		assert.Error(t, err)
	})
}

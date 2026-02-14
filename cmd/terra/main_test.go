//go:build unit

package main

import (
	"testing"

	"github.com/rios0rios0/terra/internal/domain/entities"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// stubController is a minimal Controller implementation for testing.
type stubController struct {
	bind     entities.ControllerBind
	executed bool
}

func (s *stubController) GetBind() entities.ControllerBind {
	return s.bind
}

func (s *stubController) Execute(_ *cobra.Command, _ []string) {
	s.executed = true
}

// stubAppContext is a minimal AppContext implementation for testing.
type stubAppContext struct {
	controllers []entities.Controller
}

func (s *stubAppContext) GetControllers() []entities.Controller {
	return s.controllers
}

func TestBuildRootCommand(t *testing.T) {
	t.Parallel()

	t.Run("should create root command with flag parsing enabled", func(t *testing.T) {
		t.Parallel()
		// given
		controller := &stubController{
			bind: entities.ControllerBind{
				Use:   "terra",
				Short: "Terra CLI",
				Long:  "Terra CLI long description",
			},
		}

		// when
		cmd := buildRootCommand(controller, true)

		// then
		require.NotNil(t, cmd)
		assert.Equal(t, "terra", cmd.Use)
		assert.Equal(t, "Terra CLI", cmd.Short)
		assert.False(t, cmd.DisableFlagParsing)
	})

	t.Run("should create root command with flag parsing disabled", func(t *testing.T) {
		t.Parallel()
		// given
		controller := &stubController{
			bind: entities.ControllerBind{
				Use:   "terra",
				Short: "Terra CLI",
				Long:  "Terra CLI long description",
			},
		}

		// when
		cmd := buildRootCommand(controller, false)

		// then
		require.NotNil(t, cmd)
		assert.True(t, cmd.DisableFlagParsing)
	})
}

func TestAddSubcommands(t *testing.T) {
	t.Parallel()

	t.Run("should add subcommands to root command", func(t *testing.T) {
		t.Parallel()
		// given
		clearCtrl := &stubController{
			bind: entities.ControllerBind{Use: "clear", Short: "Clear cache"},
		}
		formatCtrl := &stubController{
			bind: entities.ControllerBind{Use: "format", Short: "Format files"},
		}
		appCtx := &stubAppContext{
			controllers: []entities.Controller{clearCtrl, formatCtrl},
		}
		//nolint:exhaustruct // minimal test setup
		rootCmd := &cobra.Command{Use: "terra"}

		// when
		addSubcommands(rootCmd, appCtx)

		// then
		assert.Len(t, rootCmd.Commands(), 2)
		assert.Equal(t, "clear", rootCmd.Commands()[0].Use)
		assert.Equal(t, "format", rootCmd.Commands()[1].Use)
	})

	t.Run("should add global flag to clear command", func(t *testing.T) {
		t.Parallel()
		// given
		clearCtrl := &stubController{
			bind: entities.ControllerBind{Use: "clear", Short: "Clear cache"},
		}
		appCtx := &stubAppContext{
			controllers: []entities.Controller{clearCtrl},
		}
		//nolint:exhaustruct // minimal test setup
		rootCmd := &cobra.Command{Use: "terra"}

		// when
		addSubcommands(rootCmd, appCtx)

		// then
		clearCmd := rootCmd.Commands()[0]
		globalFlag := clearCmd.Flags().Lookup("global")
		require.NotNil(t, globalFlag, "clear command should have --global flag")
		assert.Equal(t, "false", globalFlag.DefValue)
	})

	t.Run("should add dry-run and force flags to self-update command", func(t *testing.T) {
		t.Parallel()
		// given
		selfUpdateCtrl := &stubController{
			bind: entities.ControllerBind{Use: "self-update", Short: "Self update"},
		}
		appCtx := &stubAppContext{
			controllers: []entities.Controller{selfUpdateCtrl},
		}
		//nolint:exhaustruct // minimal test setup
		rootCmd := &cobra.Command{Use: "terra"}

		// when
		addSubcommands(rootCmd, appCtx)

		// then
		selfUpdateCmd := rootCmd.Commands()[0]
		dryRunFlag := selfUpdateCmd.Flags().Lookup("dry-run")
		forceFlag := selfUpdateCmd.Flags().Lookup("force")
		require.NotNil(t, dryRunFlag, "self-update command should have --dry-run flag")
		require.NotNil(t, forceFlag, "self-update command should have --force flag")
	})

	t.Run("should handle empty controllers", func(t *testing.T) {
		t.Parallel()
		// given
		appCtx := &stubAppContext{controllers: []entities.Controller{}}
		//nolint:exhaustruct // minimal test setup
		rootCmd := &cobra.Command{Use: "terra"}

		// when
		addSubcommands(rootCmd, appCtx)

		// then
		assert.Empty(t, rootCmd.Commands())
	})
}

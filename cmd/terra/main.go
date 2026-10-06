package main

import (
	"os"
	"slices"

	"github.com/joho/godotenv"
	"github.com/rios0rios0/cliforge/pkg/selfupdate"
	"github.com/rios0rios0/terra/internal/domain/commands"
	"github.com/rios0rios0/terra/internal/domain/entities"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// version is set at build time via ldflags by GoReleaser.
// During development, it defaults to "dev".
//

var version = "dev"

// runUpdateCheck queries the cliforge selfupdate command for a newer version,
// skipping local dev builds and the commands checksForUpdates leaves out.
func runUpdateCheck(command *cobra.Command) {
	if commands.TerraVersion == "dev" || !checksForUpdates(command) {
		return
	}
	selfupdate.NewCommand("rios0rios0", "terra", "terra", commands.TerraVersion).CheckForUpdates()
}

// checksForUpdates reports whether running command also checks for a newer
// release. The self-update and version subcommands skip it, to avoid redundant
// GitHub API calls and noisy warnings, and so does
// shell completion: `completion` runs from a shell's startup file every time a
// shell starts, and cobra's hidden `__complete` on every TAB press. Both exit at
// once, so a lookup started there would never be read and would only use up the
// day's update check. `completion bash` is named `bash`, so a command is judged
// by its ancestor directly under the root.
func checksForUpdates(command *cobra.Command) bool {
	for command.HasParent() && command.Parent().HasParent() {
		command = command.Parent()
	}
	return !slices.Contains([]string{
		"self-update", "version", "completion", cobra.ShellCompRequestCmd,
	}, command.Name())
}

// buildRootCommand creates and configures the root cobra command.
func buildRootCommand(rootController entities.Controller, enableFlagParsing bool) *cobra.Command {
	bind := rootController.GetBind()
	//nolint:exhaustruct // Minimal Command initialization with required fields only
	cmd := &cobra.Command{
		Use:   bind.Use,
		Short: bind.Short,
		Long:  bind.Long,
		PersistentPreRun: func(command *cobra.Command, _ []string) {
			runUpdateCheck(command)
		},
		Run: func(command *cobra.Command, arguments []string) {
			rootController.Execute(command, arguments)
		},
	}

	if !enableFlagParsing {
		cmd.Args = cobra.MinimumNArgs(1)
		cmd.DisableFlagParsing = true
	}

	return cmd
}

// addSubcommands adds all available subcommands to the provided root command.
func addSubcommands(rootCmd *cobra.Command, appContext entities.AppContext) {
	for _, controller := range appContext.GetControllers() {
		bind := controller.GetBind()
		//nolint:exhaustruct // Minimal Command initialization with required fields only
		subCmd := &cobra.Command{
			Use:   bind.Use,
			Short: bind.Short,
			Long:  bind.Long,
			Run: func(command *cobra.Command, arguments []string) {
				controller.Execute(command, arguments)
			},
		}

		// Add flags for self-update command
		if bind.Use == "self-update" {
			subCmd.Flags().Bool("dry-run", false, "Show what would be updated without performing it")
			subCmd.Flags().Bool("force", false, "Skip confirmation prompts")
		}

		// Add flags for clear command
		if bind.Use == "clear" {
			subCmd.Flags().Bool("global", false, "Also remove centralized module and provider cache directories")
		}

		rootCmd.AddCommand(subCmd)
	}
}

func main() {
	commands.TerraVersion = version //nolint:reassign // Bridge build-time ldflags to domain package

	//nolint:exhaustruct // Minimal TextFormatter initialization with required fields only
	logger.SetFormatter(&logger.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})
	if os.Getenv("DEBUG") == "true" {
		logger.SetLevel(logger.DebugLevel)
	}

	err := godotenv.Load()
	if err != nil {
		logger.Debugf("Error loading .env file: %s", err)
	}

	// Handle --version flag before cobra processing
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		// Inject the version controller and execute it directly
		appContext := injectAppContext()
		for _, controller := range appContext.GetControllers() {
			if controller.GetBind().Use == "version" {
				controller.Execute(nil, []string{})
				return
			}
		}
	}

	// Handle --help and -h flags before cobra processing
	if len(os.Args) > 1 && (os.Args[1] == "--help" || os.Args[1] == "-h") {
		// Create command structure for help display (without argument requirements)
		rootController := injectRootController()
		tempRoot := buildRootCommand(rootController, true) // enable flag parsing for help

		// Add subcommands for complete help
		appContext := injectAppContext()
		addSubcommands(tempRoot, appContext)

		// Set args and execute help
		tempRoot.SetArgs([]string{"--help"})
		err = tempRoot.Execute()
		if err != nil {
			logger.Fatalf("Error showing help: %s", err)
		}
		return
	}

	// "cobra" library needs to start with a cobraRoot command
	rootController := injectRootController()
	cobraRoot := buildRootCommand(
		rootController,
		false,
	) // disable flag parsing for normal execution

	// all other commands are added as subcommands
	appContext := injectAppContext()
	addSubcommands(cobraRoot, appContext)

	err = cobraRoot.Execute()
	if err != nil {
		logger.Fatalf("Error executing 'terra': %s", err)
	}
}

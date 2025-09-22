package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/rios0rios0/terra/internal/domain/entities"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// buildRootCommand creates and configures the root cobra command.
func buildRootCommand(rootController entities.Controller, enableFlagParsing bool) *cobra.Command {
	bind := rootController.GetBind()
	//nolint:exhaustruct // Minimal Command initialization with required fields only
	cmd := &cobra.Command{
		Use:   bind.Use,
		Short: bind.Short,
		Long:  bind.Long,
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
		rootCmd.AddCommand(subCmd)
	}
}

func main() {
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

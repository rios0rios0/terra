package main

import (
	"os"

	"github.com/joho/godotenv"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func main() {
	//nolint:exhaustruct
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
		// Create a temporary command structure to display proper help
		// without argument requirements for help display
		rootController := injectRootController()
		bind := rootController.GetBind()
		//nolint:exhaustruct
		tempRoot := &cobra.Command{
			Use:   bind.Use,
			Short: bind.Short,
			Long:  bind.Long,
			Run: func(command *cobra.Command, arguments []string) {
				rootController.Execute(command, arguments)
			},
		}

		// Add subcommands for complete help
		appContext := injectAppContext()
		for _, controller := range appContext.GetControllers() {
			bind = controller.GetBind()
			//nolint:exhaustruct
			tempRoot.AddCommand(&cobra.Command{
				Use:   bind.Use,
				Short: bind.Short,
				Long:  bind.Long,
				Run: func(command *cobra.Command, arguments []string) {
					controller.Execute(command, arguments)
				},
			})
		}

		// Set args and execute help
		tempRoot.SetArgs([]string{"--help"})
		if err := tempRoot.Execute(); err != nil {
			logger.Fatalf("Error showing help: %s", err)
		}
		return
	}

	// "cobra" library needs to start with a cobraRoot command
	rootController := injectRootController()
	bind := rootController.GetBind()
	//nolint:exhaustruct
	cobraRoot := &cobra.Command{
		Use:                bind.Use,
		Short:              bind.Short,
		Long:               bind.Long,
		Args:               cobra.MinimumNArgs(1),
		DisableFlagParsing: true,
		Run: func(command *cobra.Command, arguments []string) {
			rootController.Execute(command, arguments)
		},
	}

	// all other commands are added as subcommands
	appContext := injectAppContext()
	for _, controller := range appContext.GetControllers() {
		bind = controller.GetBind()
		//nolint:exhaustruct
		cobraRoot.AddCommand(&cobra.Command{
			Use:   bind.Use,
			Short: bind.Short,
			Long:  bind.Long,
			Run: func(command *cobra.Command, arguments []string) {
				controller.Execute(command, arguments)
			},
		})
	}

	if err := cobraRoot.Execute(); err != nil {
		logger.Fatalf("Error executing 'terra': %s", err)
	}
}

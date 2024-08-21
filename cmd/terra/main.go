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
		logger.Warnf("Error loading .env file: %s", err)
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

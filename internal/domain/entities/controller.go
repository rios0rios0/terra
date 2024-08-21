package entities

import "github.com/spf13/cobra"

type Controller interface {
	GetBind() ControllerBind
	Execute(command *cobra.Command, arguments []string)
}

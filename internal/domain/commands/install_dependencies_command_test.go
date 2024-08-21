//go:build unit

package commands_test

import (
	"testing"
)

func TestInstallDependenciesCommand_Execute(t *testing.T) {
	t.Run("should install dependencies if not available", func(t *testing.T) {
		//// given
		//dependencies := builders.NewDependencyBuilder().BuildMany()
		//command := commands.NewInstallDependenciesCommand()
		//
		//listeners := interfaces.InstallDependenciesListeners{OnSuccess: func() {
		//	// then
		//	assert.True(t, true, "the success listener should be called")
		//}}
		//
		//// when
		//command.Execute(dependencies, listeners)
	})

	t.Run("should throw an error when some unexpected condition happens", func(t *testing.T) {
		//// given
		//dependencies := builders.NewDependencyBuilder().BuildMany()
		//command := commands.NewInstallDependenciesCommand()
		//
		//listeners := interfaces.InstallDependenciesListeners{OnError: func(err error) {
		//	// then
		//	assert.Error(t, err, "the error listener should be called")
		//}}
		//
		//// when
		//command.Execute(dependencies, listeners)
	})
}

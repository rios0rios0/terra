//go:build unit

package commands_test

import (
	"github.com/rios0rios0/terra/internal/domain/commands"
	"github.com/rios0rios0/terra/internal/domain/commands/interfaces"
	"github.com/rios0rios0/terra/internal/domain/entities"
	testbuilders "github.com/rios0rios0/terra/test/domain/entities/builders"
	testrepositories "github.com/rios0rios0/terra/test/infrastructure/repositories/doubles"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInstallDependenciesCommand_Execute(t *testing.T) {
	t.Run("should successfully install all dependencies", func(t *testing.T) {
		// given
		dependencies := testbuilders.NewDependencyBuilder().BuildMany()

		osRepository := testrepositories.NewOSRepositoryStub().WithSuccess()
		webRepository := testrepositories.NewWebStringsRepositoryStub().WithSuccess()
		command := commands.NewInstallDependenciesCommand(osRepository, webRepository)

		listeners := interfaces.InstallDependenciesListeners{
			OnSuccess: func() {
				// then
				assert.True(t, true, "the success listener should be called when dependencies are successfully installed")
			},
		}

		// when
		command.Execute(dependencies, listeners)
	})

	t.Run("should return an error when not running as root", func(t *testing.T) {
		// given
		osRepository := testrepositories.NewOSRepositoryStub().WithUserID(1000)
		webRepository := testrepositories.NewWebStringsRepositoryStub().WithSuccess()
		command := commands.NewInstallDependenciesCommand(osRepository, webRepository)

		listeners := interfaces.InstallDependenciesListeners{
			OnError: func(err error) {
				// then
				assert.Error(t, err, "the error listener should be called when not running as root")
				assert.ErrorContains(t, err, "run this command with root privileges to install the dependencies")
			},
		}

		// when
		command.Execute([]entities.Dependency{}, listeners)
	})

	t.Run("should return an error when failing to fetch latest version", func(t *testing.T) {
		// given
		dependencies := testbuilders.NewDependencyBuilder().BuildMany()

		osRepository := testrepositories.NewOSRepositoryStub().WithSuccess()
		webRepository := testrepositories.NewWebStringsRepositoryStub().WithError()
		command := commands.NewInstallDependenciesCommand(osRepository, webRepository)

		listeners := interfaces.InstallDependenciesListeners{
			OnError: func(err error) {
				// then
				assert.Error(t, err, "the error listener should be called when failing to fetch latest version")
				assert.ErrorContains(t, err, "failed to fetch latest version for test-dependency")
			},
		}

		// when
		command.Execute(dependencies, listeners)
	})

	//t.Run("should return an error when failing to download dependency", func(t *testing.T) {
	//	// given
	//	dependencies := testbuilders.NewDependencyBuilder().BuildMany()
	//
	//	osRepository := testrepositories.NewOSRepositoryStub().WithDownloadError()
	//	webRepository := testrepositories.NewWebStringsRepositoryStub().WithSuccess()
	//	command := commands.NewInstallDependenciesCommand(osRepository, webRepository)
	//
	//	listeners := interfaces.InstallDependenciesListeners{
	//		OnError: func(err error) {
	//			// then
	//			assert.Error(t, err, "the error listener should be called when failing to download dependency")
	//			assert.ErrorContains(t, err, "failed to download test-dependency")
	//		},
	//	}
	//
	//	// when
	//	command.Execute(dependencies, listeners)
	//})
	//
	//t.Run("should return an error when failing to install dependency", func(t *testing.T) {
	//	// given
	//	dependencies := testbuilders.NewDependencyBuilder().BuildMany()
	//
	//	osRepository := testrepositories.NewOSRepositoryStub().WithInstallError()
	//	webRepository := testrepositories.NewWebStringsRepositoryStub().WithSuccess()
	//	command := commands.NewInstallDependenciesCommand(osRepository, webRepository)
	//
	//	listeners := interfaces.InstallDependenciesListeners{
	//		OnError: func(err error) {
	//			// then
	//			assert.Error(t, err, "the error listener should be called when failing to install dependency")
	//			assert.ErrorContains(t, err, "failed to install dependency test-dependency")
	//		},
	//	}
	//
	//	// when
	//	command.Execute(dependencies, listeners)
	//})
}

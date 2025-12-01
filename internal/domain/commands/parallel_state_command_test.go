//go:build unit

package commands_test

import (
	"os"
	"testing"

	"github.com/rios0rios0/terra/internal/domain/commands"
	"github.com/rios0rios0/terra/internal/domain/entities"
	"github.com/rios0rios0/terra/test/infrastructure/repositorydoubles"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper functions for file operations in tests
func mkdir(dir string) error {
	// nosemgrep: go.lang.correctness.permissions.file_permission.incorrect-default-permission
	return os.MkdirAll(dir, 0755) // Directory requires execute permission (0700) for traversal in tests
}

func writeFile(filename, content string) error {
	return os.WriteFile(filename, []byte(content), 0644)
}

func TestNewParallelStateCommand(t *testing.T) {
	t.Parallel()

	t.Run("should create instance when called", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A shell repository
		repository := &repositorydoubles.StubShellRepositoryForParallelState{}

		// WHEN: Creating a new parallel state command
		cmd := commands.NewParallelStateCommand(repository)

		// THEN: Should create a valid command instance
		require.NotNil(t, cmd)
	})
}

func TestParallelStateCommand_Execute(t *testing.T) {
	t.Run("should execute import command in parallel when --all flag present", func(t *testing.T) {
		// GIVEN: A parallel state command and import arguments with --all
		repository := &repositorydoubles.StubShellRepositoryForParallelState{}
		cmd := commands.NewParallelStateCommand(repository)
		arguments := []string{"import", "--all", "null_resource.test", "test-id"}
		dependencies := []entities.Dependency{}

		// Create test directories with terraform files
		tempDir := t.TempDir()
		testHelper := newTestDirectoryHelper(t)
		testHelper.createModuleDirectories(tempDir, []string{"module1", "module2"})

		// WHEN: Executing the command
		err := cmd.Execute(tempDir, arguments, dependencies)

		// THEN: Should execute successfully
		require.NoError(t, err)
		assert.GreaterOrEqual(t, repository.ExecuteCallCount, 2, "Should execute at least 2 commands")
	})

	t.Run("should execute state rm command in parallel when --all flag present", func(t *testing.T) {
		// GIVEN: A parallel state command and state rm arguments with --all
		repository := &repositorydoubles.StubShellRepositoryForParallelState{}
		cmd := commands.NewParallelStateCommand(repository)
		arguments := []string{"state", "rm", "--all", "null_resource.test"}
		dependencies := []entities.Dependency{}

		// Create test directories with terraform files
		tempDir := t.TempDir()
		testHelper := newTestDirectoryHelper(t)
		testHelper.createModuleDirectories(tempDir, []string{"module1", "module2", "module3"})

		// WHEN: Executing the command
		err := cmd.Execute(tempDir, arguments, dependencies)

		// THEN: Should execute successfully
		require.NoError(t, err)
		assert.GreaterOrEqual(t, repository.ExecuteCallCount, 3, "Should execute at least 3 commands")
	})

	t.Run("should return error when command is not state manipulation", func(t *testing.T) {
		// GIVEN: A parallel state command and non-state arguments without --parallel=N flag
		repository := &repositorydoubles.StubShellRepositoryForParallelState{}
		cmd := commands.NewParallelStateCommand(repository)
		targetPath := "/tmp/test-terraform"
		arguments := []string{"plan", "--all"}
		dependencies := []entities.Dependency{}

		// WHEN: Executing the command
		err := cmd.Execute(targetPath, arguments, dependencies)

		// THEN: Should return error
		require.Error(t, err)
		assert.Contains(t, err.Error(), "not a parallel command")
	})

	t.Run("should return error when no modules found", func(t *testing.T) {
		// GIVEN: A parallel state command and empty directory
		repository := &repositorydoubles.StubShellRepositoryForParallelState{}
		cmd := commands.NewParallelStateCommand(repository)
		arguments := []string{"import", "--all", "null_resource.test", "test-id"}
		dependencies := []entities.Dependency{}

		// Create empty test directory
		tempDir := t.TempDir()

		// WHEN: Executing the command
		err := cmd.Execute(tempDir, arguments, dependencies)

		// THEN: Should return error
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no terraform/terragrunt modules found")
	})

	t.Run("should handle state mv command correctly", func(t *testing.T) {
		// GIVEN: A parallel state command and state mv arguments
		repository := &repositorydoubles.StubShellRepositoryForParallelState{}
		cmd := commands.NewParallelStateCommand(repository)
		arguments := []string{"state", "mv", "--all", "old_resource", "new_resource"}
		dependencies := []entities.Dependency{}

		// Create test directories with terraform files
		tempDir := t.TempDir()
		testHelper := newTestDirectoryHelper(t)
		testHelper.createModuleDirectories(tempDir, []string{"module1"})

		// WHEN: Executing the command
		err := cmd.Execute(tempDir, arguments, dependencies)

		// THEN: Should execute successfully
		require.NoError(t, err)
		assert.GreaterOrEqual(t, repository.ExecuteCallCount, 1, "Should execute at least 1 command")

		// Verify --all flag was removed from arguments
		lastCall := repository.CallHistory[len(repository.CallHistory)-1]
		assert.NotContains(t, lastCall.Arguments, "--all", "Should remove --all flag from individual module execution")
	})
}

// testDirectoryHelper helps create test directories for parallel state tests
type testDirectoryHelper struct {
	t *testing.T
}

func newTestDirectoryHelper(t *testing.T) *testDirectoryHelper {
	t.Helper()
	return &testDirectoryHelper{t: t}
}

func (h *testDirectoryHelper) createModuleDirectories(baseDir string, moduleNames []string) {
	h.t.Helper()

	for _, moduleName := range moduleNames {
		testHelper := newModuleTestHelper(h.t, baseDir, moduleName)
		testHelper.createTerraformModule()
	}
}

// moduleTestHelper helps create individual terraform modules
type moduleTestHelper struct {
	t          *testing.T
	baseDir    string
	moduleName string
}

func newModuleTestHelper(t *testing.T, baseDir, moduleName string) *moduleTestHelper {
	t.Helper()
	return &moduleTestHelper{
		t:          t,
		baseDir:    baseDir,
		moduleName: moduleName,
	}
}

func (h *moduleTestHelper) createTerraformModule() {
	h.t.Helper()

	moduleDir := h.baseDir + "/" + h.moduleName
	err := mkdir(moduleDir)
	require.NoError(h.t, err, "Failed to create module directory")

	// Create a basic terraform file
	terraformContent := `resource "null_resource" "test" {
  lifecycle {
    create_before_destroy = true
  }
}`

	err = writeFile(moduleDir+"/main.tf", terraformContent)
	require.NoError(h.t, err, "Failed to create terraform file")
}

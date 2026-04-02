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
	t.Run("should execute import command in parallel when --parallel flag present", func(t *testing.T) {
		// GIVEN: A parallel state command and import arguments with --parallel=5
		repository := &repositorydoubles.StubShellRepositoryForParallelState{}
		cmd := commands.NewParallelStateCommand(repository)
		arguments := []string{"import", "--parallel=5", "null_resource.test", "test-id"}
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

	t.Run("should execute state rm command in parallel when --parallel flag present", func(t *testing.T) {
		// GIVEN: A parallel state command and state rm arguments with --parallel=5
		repository := &repositorydoubles.StubShellRepositoryForParallelState{}
		cmd := commands.NewParallelStateCommand(repository)
		arguments := []string{"state", "rm", "--parallel=5", "null_resource.test"}
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
		arguments := []string{"plan"}
		dependencies := []entities.Dependency{}

		// WHEN: Executing the command
		err := cmd.Execute(targetPath, arguments, dependencies)

		// THEN: Should return error
		require.Error(t, err)
		assert.Contains(t, err.Error(), "not a parallel command")
	})

	t.Run("should skip hidden directories when discovering modules", func(t *testing.T) {
		// GIVEN: A directory with hidden and non-hidden module directories
		repository := &repositorydoubles.StubShellRepositoryForParallelState{}
		cmd := commands.NewParallelStateCommand(repository)
		arguments := []string{"plan", "--parallel=2"}
		dependencies := []entities.Dependency{}

		tempDir := t.TempDir()
		testHelper := newTestDirectoryHelper(t)
		testHelper.createModuleDirectories(tempDir, []string{"visible_mod"})
		// Create a hidden directory with terraform files
		hiddenDir := tempDir + "/.hidden_mod"
		require.NoError(t, mkdir(hiddenDir))
		require.NoError(t, writeFile(hiddenDir+"/main.tf", "resource \"null_resource\" \"test\" {}"))

		// WHEN: Executing the command
		err := cmd.Execute(tempDir, arguments, dependencies)

		// THEN: Should only find the visible module, not the hidden one
		require.NoError(t, err)
		assert.Equal(t, 1, repository.ExecuteCallCount, "Should execute only 1 command (hidden dir skipped)")
	})

	t.Run("should not descend into .terragrunt-cache and discover cached dependencies", func(t *testing.T) {
		// GIVEN: A directory with a real module and a .terragrunt-cache containing cached dependency modules
		repository := &repositorydoubles.StubShellRepositoryForParallelState{}
		cmd := commands.NewParallelStateCommand(repository)
		arguments := []string{"plan", "--parallel=2"}
		dependencies := []entities.Dependency{}

		tempDir := t.TempDir()
		testHelper := newTestDirectoryHelper(t)
		testHelper.createModuleDirectories(tempDir, []string{"real_module"})

		// Simulate .terragrunt-cache structure with nested dependency modules
		cacheBase := tempDir + "/.terragrunt-cache/hashA/hashB/environments/01_shared"
		require.NoError(t, mkdir(cacheBase+"/01_tfstate"))
		require.NoError(t, writeFile(cacheBase+"/01_tfstate/main.tf", "resource \"null_resource\" \"cached\" {}"))
		require.NoError(t, mkdir(cacheBase+"/02_common"))
		require.NoError(t, writeFile(cacheBase+"/02_common/main.tf", "resource \"null_resource\" \"cached\" {}"))

		// WHEN: Executing the command
		err := cmd.Execute(tempDir, arguments, dependencies)

		// THEN: Should only find the real module, not the cached dependencies
		require.NoError(t, err)
		assert.Equal(t, 1, repository.ExecuteCallCount, "Should execute only 1 command (cached deps skipped)")
	})

	t.Run("should skip directories without terraform files", func(t *testing.T) {
		// GIVEN: A directory with one tf module and one non-tf directory
		repository := &repositorydoubles.StubShellRepositoryForParallelState{}
		cmd := commands.NewParallelStateCommand(repository)
		arguments := []string{"plan", "--parallel=2"}
		dependencies := []entities.Dependency{}

		tempDir := t.TempDir()
		testHelper := newTestDirectoryHelper(t)
		testHelper.createModuleDirectories(tempDir, []string{"real_mod"})
		// Create a directory without terraform files
		nonTfDir := tempDir + "/not_terraform"
		require.NoError(t, mkdir(nonTfDir))
		require.NoError(t, writeFile(nonTfDir+"/readme.md", "# Not a terraform module"))

		// WHEN: Executing the command
		err := cmd.Execute(tempDir, arguments, dependencies)

		// THEN: Should only find the real terraform module
		require.NoError(t, err)
		assert.Equal(t, 1, repository.ExecuteCallCount, "Should execute only 1 command (non-tf dir skipped)")
	})

	t.Run("should detect terragrunt.hcl files as valid modules", func(t *testing.T) {
		// GIVEN: A directory with a terragrunt.hcl module (no .tf files)
		repository := &repositorydoubles.StubShellRepositoryForParallelState{}
		cmd := commands.NewParallelStateCommand(repository)
		arguments := []string{"plan", "--parallel=2"}
		dependencies := []entities.Dependency{}

		tempDir := t.TempDir()
		tgDir := tempDir + "/tg_module"
		require.NoError(t, mkdir(tgDir))
		require.NoError(t, writeFile(tgDir+"/terragrunt.hcl", "terraform { source = \".\" }"))

		// WHEN: Executing the command
		err := cmd.Execute(tempDir, arguments, dependencies)

		// THEN: Should detect the terragrunt.hcl module
		require.NoError(t, err)
		assert.Equal(t, 1, repository.ExecuteCallCount, "Should execute 1 command for terragrunt.hcl module")
	})

	t.Run("should detect tfvars files as valid modules", func(t *testing.T) {
		// GIVEN: A directory with only .tfvars files
		repository := &repositorydoubles.StubShellRepositoryForParallelState{}
		cmd := commands.NewParallelStateCommand(repository)
		arguments := []string{"plan", "--parallel=2"}
		dependencies := []entities.Dependency{}

		tempDir := t.TempDir()
		tfvarsDir := tempDir + "/vars_module"
		require.NoError(t, mkdir(tfvarsDir))
		require.NoError(t, writeFile(tfvarsDir+"/terraform.tfvars", "region = \"us-east-1\""))

		// WHEN: Executing the command
		err := cmd.Execute(tempDir, arguments, dependencies)

		// THEN: Should detect the tfvars module
		require.NoError(t, err)
		assert.Equal(t, 1, repository.ExecuteCallCount, "Should execute 1 command for tfvars module")
	})

	t.Run("should return error when no modules found", func(t *testing.T) {
		// GIVEN: A parallel state command and empty directory
		repository := &repositorydoubles.StubShellRepositoryForParallelState{}
		cmd := commands.NewParallelStateCommand(repository)
		arguments := []string{"import", "--parallel=5", "null_resource.test", "test-id"}
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
		arguments := []string{"state", "mv", "--parallel=5", "old_resource", "new_resource"}
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

		// Verify --parallel=5 flag was removed from arguments
		lastCall := repository.CallHistory[len(repository.CallHistory)-1]
		assert.NotContains(t, lastCall.Arguments, "--parallel=5", "Should remove --parallel=5 flag from individual module execution")
	})

	t.Run("should execute with --parallel=2 for any command", func(t *testing.T) {
		// GIVEN: A parallel state command with --parallel=2 flag
		repository := &repositorydoubles.StubShellRepositoryForParallelState{}
		cmd := commands.NewParallelStateCommand(repository)
		arguments := []string{"plan", "--parallel=2"}
		dependencies := []entities.Dependency{}

		tempDir := t.TempDir()
		testHelper := newTestDirectoryHelper(t)
		testHelper.createModuleDirectories(tempDir, []string{"mod1", "mod2"})

		// WHEN: Executing the command
		err := cmd.Execute(tempDir, arguments, dependencies)

		// THEN: Should execute successfully
		require.NoError(t, err)
		assert.GreaterOrEqual(t, repository.ExecuteCallCount, 2)
	})

	t.Run("should execute with --only when only flag present", func(t *testing.T) {
		// GIVEN: A parallel state command with --only flag
		repository := &repositorydoubles.StubShellRepositoryForParallelState{}
		cmd := commands.NewParallelStateCommand(repository)
		arguments := []string{"import", "--parallel=5", "--only=mod1,mod2", "null_resource.test", "id"}
		dependencies := []entities.Dependency{}

		tempDir := t.TempDir()
		testHelper := newTestDirectoryHelper(t)
		testHelper.createModuleDirectories(tempDir, []string{"mod1", "mod2", "mod3"})

		// WHEN: Executing the command
		err := cmd.Execute(tempDir, arguments, dependencies)

		// THEN: Should execute only for included modules (mod1 and mod2, not mod3)
		require.NoError(t, err)
		assert.Equal(t, 2, repository.ExecuteCallCount, "Should execute exactly 2 commands for included modules")
	})

	t.Run("should execute with --skip when skip flag present", func(t *testing.T) {
		// GIVEN: A parallel state command with --skip flag
		repository := &repositorydoubles.StubShellRepositoryForParallelState{}
		cmd := commands.NewParallelStateCommand(repository)
		arguments := []string{"import", "--parallel=5", "--skip=mod3", "null_resource.test", "id"}
		dependencies := []entities.Dependency{}

		tempDir := t.TempDir()
		testHelper := newTestDirectoryHelper(t)
		testHelper.createModuleDirectories(tempDir, []string{"mod1", "mod2", "mod3"})

		// WHEN: Executing the command
		err := cmd.Execute(tempDir, arguments, dependencies)

		// THEN: Should execute only for non-excluded modules (mod1 and mod2)
		require.NoError(t, err)
		assert.Equal(t, 2, repository.ExecuteCallCount, "Should execute 2 commands excluding mod3")
	})

	t.Run("should return error when only matches no valid paths", func(t *testing.T) {
		// GIVEN: A parallel state command with --only pointing to nonexistent dirs
		repository := &repositorydoubles.StubShellRepositoryForParallelState{}
		cmd := commands.NewParallelStateCommand(repository)
		arguments := []string{"import", "--parallel=5", "--only=nonexistent", "null_resource.test", "id"}
		dependencies := []entities.Dependency{}

		tempDir := t.TempDir()

		// WHEN: Executing the command
		err := cmd.Execute(tempDir, arguments, dependencies)

		// THEN: Should return error about no valid module paths
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no valid module paths found")
	})

	t.Run("should execute with both --only and --skip when both flags present", func(t *testing.T) {
		// GIVEN: A parallel state command with both --only and --skip flags
		repository := &repositorydoubles.StubShellRepositoryForParallelState{}
		cmd := commands.NewParallelStateCommand(repository)
		arguments := []string{"import", "--parallel=5", "--only=mod1,mod2,mod3", "--skip=mod2", "null_resource.test", "id"}
		dependencies := []entities.Dependency{}

		tempDir := t.TempDir()
		testHelper := newTestDirectoryHelper(t)
		testHelper.createModuleDirectories(tempDir, []string{"mod1", "mod2", "mod3"})

		// WHEN: Executing the command
		err := cmd.Execute(tempDir, arguments, dependencies)

		// THEN: Should execute only for included modules minus excluded (mod1 and mod3)
		require.NoError(t, err)
		assert.Equal(t, 2, repository.ExecuteCallCount, "Should execute 2 commands (mod1 and mod3, excluding mod2)")
	})

	t.Run("should execute with --skip only discovering all modules first", func(t *testing.T) {
		// GIVEN: A parallel state command with only --skip flag (no --only)
		repository := &repositorydoubles.StubShellRepositoryForParallelState{}
		cmd := commands.NewParallelStateCommand(repository)
		arguments := []string{"plan", "--parallel=4", "--skip=mod2"}
		dependencies := []entities.Dependency{}

		tempDir := t.TempDir()
		testHelper := newTestDirectoryHelper(t)
		testHelper.createModuleDirectories(tempDir, []string{"mod1", "mod2", "mod3"})

		// WHEN: Executing the command
		err := cmd.Execute(tempDir, arguments, dependencies)

		// THEN: Should discover all modules then exclude mod2
		require.NoError(t, err)
		assert.Equal(t, 2, repository.ExecuteCallCount, "Should execute 2 commands (mod1 and mod3)")
	})

	t.Run("should return error when skip removes all discovered modules", func(t *testing.T) {
		// GIVEN: A parallel state command where --skip removes all modules
		repository := &repositorydoubles.StubShellRepositoryForParallelState{}
		cmd := commands.NewParallelStateCommand(repository)
		arguments := []string{"plan", "--parallel=2", "--skip=mod1,mod2"}
		dependencies := []entities.Dependency{}

		tempDir := t.TempDir()
		testHelper := newTestDirectoryHelper(t)
		testHelper.createModuleDirectories(tempDir, []string{"mod1", "mod2"})

		// WHEN: Executing the command
		err := cmd.Execute(tempDir, arguments, dependencies)

		// THEN: Should return error because no modules remain after exclusion
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no valid module paths found")
	})

	t.Run("should not pass --only or --skip flags to terragrunt", func(t *testing.T) {
		// GIVEN: A parallel state command with both --only and --skip flags
		repository := &repositorydoubles.StubShellRepositoryForParallelState{}
		cmd := commands.NewParallelStateCommand(repository)
		arguments := []string{"import", "--parallel=5", "--only=mod1", "--skip=mod2", "null_resource.test", "id"}
		dependencies := []entities.Dependency{}

		tempDir := t.TempDir()
		testHelper := newTestDirectoryHelper(t)
		testHelper.createModuleDirectories(tempDir, []string{"mod1"})

		// WHEN: Executing the command
		err := cmd.Execute(tempDir, arguments, dependencies)

		// THEN: Should not pass --only or --skip flags to terragrunt
		require.NoError(t, err)
		require.NotEmpty(t, repository.CallHistory)
		lastCall := repository.CallHistory[len(repository.CallHistory)-1]
		for _, arg := range lastCall.Arguments {
			assert.NotContains(t, arg, "--only=", "Should not pass --only flag to terragrunt")
			assert.NotContains(t, arg, "--skip=", "Should not pass --skip flag to terragrunt")
		}
	})

	t.Run("should strip --reply flag and inject --non-interactive and -auto-approve for apply", func(t *testing.T) {
		// GIVEN: A parallel state command with --reply=y flag on an apply command
		repository := &repositorydoubles.StubShellRepositoryForParallelState{}
		cmd := commands.NewParallelStateCommand(repository)
		arguments := []string{"apply", "--parallel=2", "--reply=y"}
		dependencies := []entities.Dependency{}

		tempDir := t.TempDir()
		testHelper := newTestDirectoryHelper(t)
		testHelper.createModuleDirectories(tempDir, []string{"mod1"})

		// WHEN: Executing the command
		err := cmd.Execute(tempDir, arguments, dependencies)

		// THEN: Should strip --reply and inject both --non-interactive and -auto-approve
		require.NoError(t, err)
		require.NotEmpty(t, repository.CallHistory)
		lastCall := repository.CallHistory[len(repository.CallHistory)-1]
		for _, arg := range lastCall.Arguments {
			assert.NotContains(t, arg, "--reply", "Should not pass --reply flag to terragrunt")
		}
		assert.Contains(t, lastCall.Arguments, "--non-interactive", "Should inject --non-interactive when --reply was present")
		assert.Contains(t, lastCall.Arguments, "-auto-approve", "Should inject -auto-approve for apply commands")
	})

	t.Run("should inject --non-interactive but not -auto-approve for plan with --reply", func(t *testing.T) {
		// GIVEN: A parallel state command with --reply=y flag on a plan command (non-interactive)
		repository := &repositorydoubles.StubShellRepositoryForParallelState{}
		cmd := commands.NewParallelStateCommand(repository)
		arguments := []string{"plan", "--parallel=2", "--reply=y"}
		dependencies := []entities.Dependency{}

		tempDir := t.TempDir()
		testHelper := newTestDirectoryHelper(t)
		testHelper.createModuleDirectories(tempDir, []string{"mod1"})

		// WHEN: Executing the command
		err := cmd.Execute(tempDir, arguments, dependencies)

		// THEN: Should inject --non-interactive but NOT -auto-approve (plan doesn't need it)
		require.NoError(t, err)
		require.NotEmpty(t, repository.CallHistory)
		lastCall := repository.CallHistory[len(repository.CallHistory)-1]
		assert.Contains(t, lastCall.Arguments, "--non-interactive", "Should inject --non-interactive")
		assert.NotContains(t, lastCall.Arguments, "-auto-approve", "Should not inject -auto-approve for plan")
	})

	t.Run("should not inject --non-interactive when reply not present", func(t *testing.T) {
		// GIVEN: A parallel state command without --reply flag
		repository := &repositorydoubles.StubShellRepositoryForParallelState{}
		cmd := commands.NewParallelStateCommand(repository)
		arguments := []string{"plan", "--parallel=2"}
		dependencies := []entities.Dependency{}

		tempDir := t.TempDir()
		testHelper := newTestDirectoryHelper(t)
		testHelper.createModuleDirectories(tempDir, []string{"mod1"})

		// WHEN: Executing the command
		err := cmd.Execute(tempDir, arguments, dependencies)

		// THEN: Should not inject --non-interactive
		require.NoError(t, err)
		require.NotEmpty(t, repository.CallHistory)
		lastCall := repository.CallHistory[len(repository.CallHistory)-1]
		assert.NotContains(t, lastCall.Arguments, "--non-interactive", "Should not inject --non-interactive without --reply")
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

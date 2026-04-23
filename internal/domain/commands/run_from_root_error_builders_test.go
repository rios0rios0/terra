//go:build unit

package commands_test

import (
	"testing"

	"github.com/rios0rios0/terra/internal/domain/commands"
	"github.com/stretchr/testify/assert"
)

func TestBuildSelectionFlagsError(t *testing.T) {
	t.Parallel()

	t.Run("should echo the user command verbatim in the error", func(t *testing.T) {
		t.Parallel()
		// GIVEN: The arguments a user would pass on the terra CLI
		arguments := []string{"apply", "--all", "--skip=excluded-mod"}
		targetPath := "environments/example-env/dev"

		// WHEN: Building the selection-flags error
		message := commands.BuildSelectionFlagsError(arguments, targetPath)

		// THEN: The error includes the exact command the user typed
		assert.Contains(
			t,
			message,
			"You used: terra apply --all --skip=excluded-mod environments/example-env/dev",
		)
	})

	t.Run("should suggest --parallel=5 with user's skip values", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A --skip invocation without --parallel
		arguments := []string{"plan", "--skip=mod1,mod2"}
		targetPath := "/infra"

		// WHEN: Building the error
		message := commands.BuildSelectionFlagsError(arguments, targetPath)

		// THEN: The --parallel suggestion preserves the skip values verbatim
		assert.Contains(t, message, "terra plan --parallel=5 --skip=mod1,mod2 /infra")
	})

	t.Run("should suggest --filter='!value' with negation for each skip value", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A --skip invocation
		arguments := []string{"plan", "--skip=mod1,mod2"}
		targetPath := "/infra"

		// WHEN: Building the error
		message := commands.BuildSelectionFlagsError(arguments, targetPath)

		// THEN: The --filter suggestion negates each skip value with '!'
		assert.Contains(t, message, "terra plan --all --filter='!mod1' --filter='!mod2' /infra")
	})

	t.Run("should suggest --filter='value' positively for each only value", func(t *testing.T) {
		t.Parallel()
		// GIVEN: An --only invocation
		arguments := []string{"plan", "--only=mod1,mod2"}
		targetPath := "/infra"

		// WHEN: Building the error
		message := commands.BuildSelectionFlagsError(arguments, targetPath)

		// THEN: The --filter suggestion uses positive values without '!'
		assert.Contains(t, message, "terra plan --all --filter='mod1' --filter='mod2' /infra")
	})

	t.Run("should include --yes in the --parallel suggestion for apply", func(t *testing.T) {
		t.Parallel()
		// GIVEN: apply is an interactive command and terra requires a confirmation flag for it
		arguments := []string{"apply", "--skip=mod1"}
		targetPath := "/infra"

		// WHEN: Building the error
		message := commands.BuildSelectionFlagsError(arguments, targetPath)

		// THEN: The --parallel suggestion includes --yes
		assert.Contains(t, message, "terra apply --parallel=5 --skip=mod1 --yes /infra")
	})

	t.Run("should include --yes in the --parallel suggestion for destroy", func(t *testing.T) {
		t.Parallel()
		// GIVEN: destroy is also interactive
		arguments := []string{"destroy", "--only=mod1"}
		targetPath := "/infra"

		// WHEN: Building the error
		message := commands.BuildSelectionFlagsError(arguments, targetPath)

		// THEN: The --parallel suggestion includes --yes
		assert.Contains(t, message, "terra destroy --parallel=5 --only=mod1 --yes /infra")
	})

	t.Run("should NOT include --yes in the --parallel suggestion for plan", func(t *testing.T) {
		t.Parallel()
		// GIVEN: plan is not an interactive command
		arguments := []string{"plan", "--skip=mod1"}
		targetPath := "/infra"

		// WHEN: Building the error
		message := commands.BuildSelectionFlagsError(arguments, targetPath)

		// THEN: The --parallel suggestion does NOT include --yes
		assert.Contains(t, message, "terra plan --parallel=5 --skip=mod1 /infra")
		assert.NotContains(t, message, "terra plan --parallel=5 --skip=mod1 --yes")
	})

	t.Run("should preserve the target path in both suggestion lines", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A nested target path
		arguments := []string{"plan", "--skip=mod1"}
		targetPath := "/home/vsts/work/1/s/environments/example-env/dev"

		// WHEN: Building the error
		message := commands.BuildSelectionFlagsError(arguments, targetPath)

		// THEN: Both suggestion lines end with the target path
		assert.Contains(
			t,
			message,
			"terra plan --parallel=5 --skip=mod1 /home/vsts/work/1/s/environments/example-env/dev",
		)
		assert.Contains(
			t,
			message,
			"terra plan --all --filter='!mod1' /home/vsts/work/1/s/environments/example-env/dev",
		)
	})

	t.Run("should handle --only and --skip together", func(t *testing.T) {
		t.Parallel()
		// GIVEN: Both --only and --skip are present
		arguments := []string{"plan", "--only=a,b", "--skip=c"}
		targetPath := "/infra"

		// WHEN: Building the error
		message := commands.BuildSelectionFlagsError(arguments, targetPath)

		// THEN: The --parallel suggestion includes both, and the --filter
		// suggestion emits positive filters for --only and negated filters for --skip
		assert.Contains(t, message, "terra plan --parallel=5 --only=a,b --skip=c /infra")
		assert.Contains(
			t,
			message,
			"terra plan --all --filter='a' --filter='b' --filter='!c' /infra",
		)
	})

	t.Run("should reference docs/parallel-execution.md for further reading", func(t *testing.T) {
		t.Parallel()
		// GIVEN: Any selection-flag error
		arguments := []string{"plan", "--skip=mod1"}
		targetPath := "/infra"

		// WHEN: Building the error
		message := commands.BuildSelectionFlagsError(arguments, targetPath)

		// THEN: The error points at the documentation
		assert.Contains(t, message, "docs/parallel-execution.md")
	})

	t.Run("should preserve 'state rm' in suggestions for state commands", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A state rm invocation with --skip (two-word subcommand)
		arguments := []string{"state", "rm", "--skip=mod1", "null_resource.foo"}
		targetPath := "/infra"

		// WHEN: Building the error
		message := commands.BuildSelectionFlagsError(arguments, targetPath)

		// THEN: Both suggestions preserve the full "state rm" prefix so the
		// user can copy-paste them without producing broken syntax
		assert.Contains(t, message, "terra state rm --parallel=5 --skip=mod1 /infra")
		assert.Contains(t, message, "terra state rm --all --filter='!mod1' /infra")
	})

	t.Run("should preserve 'state mv' in suggestions for state commands", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A state mv invocation with --skip
		arguments := []string{"state", "mv", "--only=mod1"}
		targetPath := "/infra"

		// WHEN: Building the error
		message := commands.BuildSelectionFlagsError(arguments, targetPath)

		// THEN: Both suggestions include "state mv"
		assert.Contains(t, message, "terra state mv --parallel=5 --only=mod1 /infra")
		assert.Contains(t, message, "terra state mv --all --filter='mod1' /infra")
	})

	t.Run("should redact -var secret values in the echoed command", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A command that passes a sensitive -var flag alongside --skip
		arguments := []string{"apply", "--all", "-var=db_password=s3cret", "--skip=mod1"}
		targetPath := "/infra"

		// WHEN: Building the error
		message := commands.BuildSelectionFlagsError(arguments, targetPath)

		// THEN: The echoed command does not contain the secret, but still keeps
		// the structure so the user can see what they typed
		assert.NotContains(t, message, "s3cret")
		assert.Contains(t, message, "-var=<redacted>")
		assert.Contains(
			t,
			message,
			"You used: terra apply --all -var=<redacted> --skip=mod1 /infra",
		)
	})

	t.Run("should redact space-separated -var secret values", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A command that uses the space-separated -var form
		arguments := []string{"apply", "--all", "-var", "db_password=s3cret", "--skip=mod1"}
		targetPath := "/infra"

		// WHEN: Building the error
		message := commands.BuildSelectionFlagsError(arguments, targetPath)

		// THEN: The token after -var is redacted, not the flag itself
		assert.NotContains(t, message, "s3cret")
		assert.Contains(t, message, "-var <redacted>")
	})

	t.Run("should redact -backend-config secret values", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A command that passes a sensitive -backend-config flag
		arguments := []string{
			"apply", "--all",
			"-backend-config=access_key=AKIA1234", "--skip=mod1",
		}
		targetPath := "/infra"

		// WHEN: Building the error
		message := commands.BuildSelectionFlagsError(arguments, targetPath)

		// THEN: The access key is redacted
		assert.NotContains(t, message, "AKIA1234")
		assert.Contains(t, message, "-backend-config=<redacted>")
	})

	t.Run("should NOT redact -var-file paths", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A command that passes -var-file (a filename, not a secret)
		arguments := []string{"apply", "--all", "-var-file=prod.tfvars", "--skip=mod1"}
		targetPath := "/infra"

		// WHEN: Building the error
		message := commands.BuildSelectionFlagsError(arguments, targetPath)

		// THEN: The -var-file argument is preserved verbatim
		assert.Contains(t, message, "-var-file=prod.tfvars")
		assert.NotContains(t, message, "-var-file=<redacted>")
	})
}

func TestBuildParallelAllConflictError(t *testing.T) {
	t.Parallel()

	t.Run("should echo the user command in the error", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A command that combines --parallel and --all
		arguments := []string{"plan", "--parallel=5", "--all"}
		targetPath := "/infra"

		// WHEN: Building the error
		message := commands.BuildParallelAllConflictError(arguments, targetPath)

		// THEN: The error includes the exact command the user typed
		assert.Contains(t, message, "You used: terra plan --parallel=5 --all /infra")
	})

	t.Run("should show both the --parallel=N and --all alternatives", func(t *testing.T) {
		t.Parallel()
		// GIVEN: A command that combines both strategies
		arguments := []string{"plan", "--parallel=5", "--all"}
		targetPath := "/infra"

		// WHEN: Building the error
		message := commands.BuildParallelAllConflictError(arguments, targetPath)

		// THEN: Both alternative forms are present, each ending with the target path
		assert.Contains(t, message, "terra plan --parallel=5 /infra")
		assert.Contains(t, message, "terra plan --all /infra")
	})
}

//go:build unit

package commands_test

import (
	"testing"

	"github.com/rios0rios0/terra/internal/domain/commands"
	"github.com/stretchr/testify/assert"
)

func TestIsStateManipulationCommand(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		arguments []string
		expected  bool
	}{
		{"should return true when import command", []string{"import", "resource.id"}, true},
		{"should return true when state rm command", []string{"state", "rm", "resource.id"}, true},
		{"should return true when state mv command", []string{"state", "mv", "a", "b"}, true},
		{"should return true when state pull command", []string{"state", "pull"}, true},
		{"should return true when state push command", []string{"state", "push"}, true},
		{"should return true when state show command", []string{"state", "show"}, true},
		{"should return false when plan command", []string{"plan"}, false},
		{"should return false when apply command", []string{"apply"}, false},
		{"should return false when empty arguments", []string{}, false},
		{"should return true when state command even with unknown subcommand", []string{"state", "list"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expected, commands.IsStateManipulationCommand(tt.arguments))
		})
	}
}

func TestHasAllFlag(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		arguments []string
		expected  bool
	}{
		{"should return true when --all present", []string{"import", "--all"}, true},
		{"should return false when --all absent", []string{"import"}, false},
		{"should return false when empty arguments", []string{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expected, commands.HasAllFlag(tt.arguments))
		})
	}
}

func TestHasParallelFlag(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		arguments []string
		expected  bool
	}{
		{"should return true when --parallel=5 present", []string{"plan", "--parallel=5"}, true},
		{"should return false when --parallel absent", []string{"plan"}, false},
		{"should return false when empty arguments", []string{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expected, commands.HasParallelFlag(tt.arguments))
		})
	}
}

func TestGetParallelValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		arguments     []string
		expectedValue int
		expectedFound bool
	}{
		{"should return value when --parallel=5", []string{"--parallel=5"}, 5, true},
		{"should return value when --parallel=10", []string{"plan", "--parallel=10"}, 10, true},
		{"should return false when invalid value", []string{"--parallel=abc"}, 0, false},
		{"should return false when zero value", []string{"--parallel=0"}, 0, false},
		{"should return false when negative value", []string{"--parallel=-1"}, 0, false},
		{"should return false when not present", []string{"plan"}, 0, false},
		{"should return false when empty arguments", []string{}, 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			value, found := commands.GetParallelValue(tt.arguments)
			assert.Equal(t, tt.expectedValue, value)
			assert.Equal(t, tt.expectedFound, found)
		})
	}
}

func TestRemoveParallelFlag(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		arguments []string
		expected  []string
	}{
		{"should remove --parallel=5", []string{"plan", "--parallel=5"}, []string{"plan"}},
		{"should return unchanged when no parallel flag", []string{"plan"}, []string{"plan"}},
		{"should return nil when only parallel flag", []string{"--parallel=3"}, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expected, commands.RemoveParallelFlag(tt.arguments))
		})
	}
}

func TestHasOnlyFlag(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		arguments []string
		expected  bool
	}{
		{"should return true when --only=mod1 present", []string{"plan", "--only=mod1"}, true},
		{"should return false when absent", []string{"plan"}, false},
		{"should return false when empty arguments", []string{}, false},
		{"should return false when --skip=mod1 present", []string{"--skip=mod1"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expected, commands.HasOnlyFlag(tt.arguments))
		})
	}
}

func TestGetOnlyValues(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		arguments     []string
		expectedValue []string
		expectedFound bool
	}{
		{"should return values when --only=a,b", []string{"--only=a,b"}, []string{"a", "b"}, true},
		{"should return single value", []string{"--only=mod1"}, []string{"mod1"}, true},
		{"should return false when empty value", []string{"--only="}, nil, false},
		{"should return false when not present", []string{"plan"}, nil, false},
		{"should trim whitespace", []string{"--only= a , b "}, []string{"a", "b"}, true},
		{"should return false when empty arguments", []string{}, nil, false},
		{"should handle hyphens and underscores", []string{"--only=my-module,other_module"}, []string{"my-module", "other_module"}, true},
		{"should skip empty and use later valid occurrence", []string{"--only=", "--only=a,b"}, []string{"a", "b"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			value, found := commands.GetOnlyValues(tt.arguments)
			assert.Equal(t, tt.expectedValue, value)
			assert.Equal(t, tt.expectedFound, found)
		})
	}
}

func TestRemoveOnlyFlag(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		arguments []string
		expected  []string
	}{
		{"should remove --only=a,b", []string{"plan", "--only=a,b"}, []string{"plan"}},
		{"should return unchanged when absent", []string{"plan"}, []string{"plan"}},
		{"should return nil when only only flag", []string{"--only=mod1"}, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expected, commands.RemoveOnlyFlag(tt.arguments))
		})
	}
}

func TestHasSkipFlag(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		arguments []string
		expected  bool
	}{
		{"should return true when --skip=mod1 present", []string{"plan", "--skip=mod1"}, true},
		{"should return false when absent", []string{"plan"}, false},
		{"should return false when empty arguments", []string{}, false},
		{"should return false when --only=mod1 present", []string{"--only=mod1"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expected, commands.HasSkipFlag(tt.arguments))
		})
	}
}

func TestGetSkipValues(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		arguments     []string
		expectedValue []string
		expectedFound bool
	}{
		{"should return values when --skip=a,b", []string{"--skip=a,b"}, []string{"a", "b"}, true},
		{"should return single value", []string{"--skip=mod1"}, []string{"mod1"}, true},
		{"should return false when empty value", []string{"--skip="}, nil, false},
		{"should return false when not present", []string{"plan"}, nil, false},
		{"should trim whitespace", []string{"--skip= a , b "}, []string{"a", "b"}, true},
		{"should return false when empty arguments", []string{}, nil, false},
		{"should handle hyphens and underscores", []string{"--skip=my-module,other_module"}, []string{"my-module", "other_module"}, true},
		{"should skip empty and use later valid occurrence", []string{"--skip=", "--skip=a,b"}, []string{"a", "b"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			value, found := commands.GetSkipValues(tt.arguments)
			assert.Equal(t, tt.expectedValue, value)
			assert.Equal(t, tt.expectedFound, found)
		})
	}
}

func TestRemoveSkipFlag(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		arguments []string
		expected  []string
	}{
		{"should remove --skip=a,b", []string{"plan", "--skip=a,b"}, []string{"plan"}},
		{"should return unchanged when absent", []string{"plan"}, []string{"plan"}},
		{"should return nil when only skip flag", []string{"--skip=mod1"}, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expected, commands.RemoveSkipFlag(tt.arguments))
		})
	}
}

func TestGetSelectionValues(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		arguments    []string
		expectedOnly []string
		expectedSkip []string
	}{
		{
			"should return both when both flags present",
			[]string{"plan", "--only=a,b", "--skip=c"},
			[]string{"a", "b"},
			[]string{"c"},
		},
		{
			"should return only only when only only present",
			[]string{"plan", "--only=a,b"},
			[]string{"a", "b"},
			nil,
		},
		{
			"should return only skip when only skip present",
			[]string{"plan", "--skip=x,y"},
			nil,
			[]string{"x", "y"},
		},
		{
			"should return empty when neither flag present",
			[]string{"plan"},
			nil,
			nil,
		},
		{
			"should handle flags mixed with other arguments",
			[]string{"plan", "--parallel=4", "--only=mod1", "--skip=mod2", "/path"},
			[]string{"mod1"},
			[]string{"mod2"},
		},
		{
			"should handle multiple values in both flags",
			[]string{"--only=a,b,c", "--skip=x,y"},
			[]string{"a", "b", "c"},
			[]string{"x", "y"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := commands.GetSelectionValues(tt.arguments)
			assert.Equal(t, tt.expectedOnly, result.Only)
			assert.Equal(t, tt.expectedSkip, result.Skip)
		})
	}
}

func TestRemoveSelectionFlags(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		arguments []string
		expected  []string
	}{
		{
			"should remove both only and skip flags",
			[]string{"plan", "--only=a,b", "--skip=c"},
			[]string{"plan"},
		},
		{
			"should remove only only when skip absent",
			[]string{"plan", "--only=a"},
			[]string{"plan"},
		},
		{
			"should remove only skip when only absent",
			[]string{"plan", "--skip=c"},
			[]string{"plan"},
		},
		{
			"should return unchanged when neither present",
			[]string{"plan", "--parallel=4"},
			[]string{"plan", "--parallel=4"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expected, commands.RemoveSelectionFlags(tt.arguments))
		})
	}
}

func TestIsInteractiveCommand(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		arguments []string
		expected  bool
	}{
		{"should return true when apply command", []string{"apply"}, true},
		{"should return true when destroy command", []string{"destroy"}, true},
		{"should return false when plan command", []string{"plan"}, false},
		{"should return false when init command", []string{"init"}, false},
		{"should return false when import command", []string{"import", "addr", "id"}, false},
		{"should return false when state rm command", []string{"state", "rm", "addr"}, false},
		{"should return true when apply has leading flags", []string{"--reply=y", "apply"}, true},
		{"should return true when destroy has leading flags", []string{"--parallel=4", "--reply=y", "destroy"}, true},
		{"should return false when plan has leading flags", []string{"--parallel=4", "plan"}, false},
		{"should return false when empty arguments", []string{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expected, commands.IsInteractiveCommand(tt.arguments))
		})
	}
}

func TestHasReplyFlag(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		arguments []string
		expected  bool
	}{
		{"should return true when --reply present", []string{"apply", "--reply"}, true},
		{"should return true when -r present", []string{"apply", "-r"}, true},
		{"should return true when --reply=y present", []string{"apply", "--reply=y"}, true},
		{"should return true when -r=n present", []string{"apply", "-r=n"}, true},
		{"should return false when absent", []string{"apply"}, false},
		{"should return false when empty arguments", []string{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expected, commands.HasReplyFlag(tt.arguments))
		})
	}
}

func TestRemoveReplyFlag(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		arguments []string
		expected  []string
	}{
		{"should remove --reply", []string{"apply", "--reply"}, []string{"apply"}},
		{"should remove -r", []string{"apply", "-r"}, []string{"apply"}},
		{"should remove --reply=y", []string{"apply", "--reply=y"}, []string{"apply"}},
		{"should remove -r=n", []string{"apply", "-r=n"}, []string{"apply"}},
		{"should return unchanged when absent", []string{"plan"}, []string{"plan"}},
		{"should return nil when only reply flag", []string{"--reply=y"}, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expected, commands.RemoveReplyFlag(tt.arguments))
		})
	}
}

func TestHasExplicitReplyValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		arguments []string
		expected  bool
	}{
		{"should return false for boolean --reply", []string{"apply", "--reply"}, false},
		{"should return false for boolean -r", []string{"apply", "-r"}, false},
		{"should return true for --reply=y", []string{"apply", "--reply=y"}, true},
		{"should return true for -r=n", []string{"apply", "-r=n"}, true},
		{"should return true for --reply=custom", []string{"apply", "--reply=custom"}, true},
		{"should return false for empty --reply=", []string{"apply", "--reply="}, false},
		{"should return false for empty -r=", []string{"apply", "-r="}, false},
		{"should return false when absent", []string{"apply"}, false},
		{"should return false for empty arguments", []string{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expected, commands.HasExplicitReplyValue(tt.arguments))
		})
	}
}

func TestGetReplyValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		arguments []string
		expected  string
	}{
		{"should return 'n' for boolean --reply", []string{"--reply", "apply"}, "n"},
		{"should return 'n' for boolean -r", []string{"-r", "apply"}, "n"},
		{"should return 'y' for --reply=y", []string{"--reply=y", "apply"}, "y"},
		{"should return 'n' for -r=n", []string{"-r=n", "apply"}, "n"},
		{"should return custom value for --reply=custom", []string{"--reply=custom", "plan"}, "custom"},
		{"should return empty when absent", []string{"plan"}, ""},
		{"should return empty for empty arguments", []string{}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expected, commands.GetReplyValue(tt.arguments))
		})
	}
}

func TestHasYesFlag(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		arguments []string
		expected  bool
	}{
		{"should return true when --yes present", []string{"apply", "--yes"}, true},
		{"should return true when -y present", []string{"apply", "-y"}, true},
		{"should return false when only --no present", []string{"apply", "--no"}, false},
		{"should return false when absent", []string{"apply"}, false},
		{"should return false when empty arguments", []string{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expected, commands.HasYesFlag(tt.arguments))
		})
	}
}

func TestHasNoFlag(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		arguments []string
		expected  bool
	}{
		{"should return true when --no present", []string{"apply", "--no"}, true},
		{"should return true when -n present", []string{"apply", "-n"}, true},
		{"should return false when only --yes present", []string{"apply", "--yes"}, false},
		{"should return false when absent", []string{"apply"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expected, commands.HasNoFlag(tt.arguments))
		})
	}
}

func TestHasConfirmationFlag(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		arguments []string
		expected  bool
	}{
		{"should return true when --yes present", []string{"apply", "--yes"}, true},
		{"should return true when --no present", []string{"apply", "--no"}, true},
		{"should return true when deprecated --reply present", []string{"apply", "--reply"}, true},
		{"should return true when deprecated --reply=y present", []string{"apply", "--reply=y"}, true},
		{"should return false when absent", []string{"apply"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expected, commands.HasConfirmationFlag(tt.arguments))
		})
	}
}

func TestRemoveConfirmationFlags(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		arguments []string
		expected  []string
	}{
		{
			"should strip --yes",
			[]string{"apply", "--yes", "-target=mod"},
			[]string{"apply", "-target=mod"},
		},
		{
			"should strip -y",
			[]string{"apply", "-y"},
			[]string{"apply"},
		},
		{
			"should strip --no",
			[]string{"apply", "--no"},
			[]string{"apply"},
		},
		{
			"should strip -n",
			[]string{"apply", "-n"},
			[]string{"apply"},
		},
		{
			"should strip deprecated --reply and its value forms",
			[]string{"apply", "--reply", "-r", "--reply=y", "-r=n"},
			[]string{"apply"},
		},
		{
			"should leave arguments unchanged when no confirmation flag present",
			[]string{"apply", "-target=mod"},
			[]string{"apply", "-target=mod"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expected, commands.RemoveConfirmationFlags(tt.arguments))
		})
	}
}

func TestResolveConfirmation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		arguments  []string
		expectYes  bool
		expectNo   bool
	}{
		{"should resolve --yes to yes", []string{"apply", "--yes"}, true, false},
		{"should resolve --no to no", []string{"apply", "--no"}, false, true},
		{"should map --reply=y to yes", []string{"apply", "--reply=y"}, true, false},
		{"should map --reply=n to no", []string{"apply", "--reply=n"}, false, true},
		{"should map bare --reply to yes (default)", []string{"apply", "--reply"}, true, false},
		{"should map bare -r to yes (default)", []string{"apply", "-r"}, true, false},
		{"should map --reply=no to no", []string{"apply", "--reply=no"}, false, true},
		{"should return false,false when no confirmation flag", []string{"apply"}, false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			yes, no := commands.ResolveConfirmation(tt.arguments)
			assert.Equal(t, tt.expectYes, yes, "yes")
			assert.Equal(t, tt.expectNo, no, "no")
		})
	}
}

func TestBuildConfirmationInjection(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		arguments []string
		expected  []string
	}{
		{
			"should inject --non-interactive and -auto-approve for --yes on apply",
			[]string{"apply", "--yes"},
			[]string{"--non-interactive", "-auto-approve"},
		},
		{
			"should inject --non-interactive and -auto-approve for --yes on destroy",
			[]string{"destroy", "--yes"},
			[]string{"--non-interactive", "-auto-approve"},
		},
		{
			"should inject only --non-interactive for --yes on plan",
			[]string{"plan", "--yes"},
			[]string{"--non-interactive"},
		},
		{
			"should inject only --non-interactive for --no on apply",
			[]string{"apply", "--no"},
			[]string{"--non-interactive"},
		},
		{
			"should inject --non-interactive and -auto-approve for --reply=y on apply",
			[]string{"apply", "--reply=y"},
			[]string{"--non-interactive", "-auto-approve"},
		},
		{
			"should inject only --non-interactive for --reply=n on apply",
			[]string{"apply", "--reply=n"},
			[]string{"--non-interactive"},
		},
		{
			"should return nil when no confirmation flag is present",
			[]string{"apply"},
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expected, commands.BuildConfirmationInjection(tt.arguments))
		})
	}
}

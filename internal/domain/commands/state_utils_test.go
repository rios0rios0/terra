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

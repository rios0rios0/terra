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

func TestHasNoParallelBypassFlag(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		arguments []string
		expected  bool
	}{
		{"should return true when present", []string{"--no-parallel-bypass", "plan"}, true},
		{"should return false when absent", []string{"plan"}, false},
		{"should return false when empty", []string{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expected, commands.HasNoParallelBypassFlag(tt.arguments))
		})
	}
}

func TestRemoveNoParallelBypassFlag(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		arguments []string
		expected  []string
	}{
		{
			"should remove flag",
			[]string{"plan", "--no-parallel-bypass"},
			[]string{"plan"},
		},
		{
			"should return unchanged when absent",
			[]string{"plan"},
			[]string{"plan"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expected, commands.RemoveNoParallelBypassFlag(tt.arguments))
		})
	}
}

func TestHasIncludeFlag(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		arguments []string
		expected  bool
	}{
		{"should return true when --include=mod1 present", []string{"plan", "--include=mod1"}, true},
		{"should return false when absent", []string{"plan"}, false},
		{"should return false when empty arguments", []string{}, false},
		{"should return false when --exclude=mod1 present", []string{"--exclude=mod1"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expected, commands.HasIncludeFlag(tt.arguments))
		})
	}
}

func TestGetIncludeValues(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		arguments     []string
		expectedValue []string
		expectedFound bool
	}{
		{"should return values when --include=a,b", []string{"--include=a,b"}, []string{"a", "b"}, true},
		{"should return single value", []string{"--include=mod1"}, []string{"mod1"}, true},
		{"should return false when empty value", []string{"--include="}, nil, false},
		{"should return false when not present", []string{"plan"}, nil, false},
		{"should trim whitespace", []string{"--include= a , b "}, []string{"a", "b"}, true},
		{"should return false when empty arguments", []string{}, nil, false},
		{"should handle hyphens and underscores", []string{"--include=my-module,other_module"}, []string{"my-module", "other_module"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			value, found := commands.GetIncludeValues(tt.arguments)
			assert.Equal(t, tt.expectedValue, value)
			assert.Equal(t, tt.expectedFound, found)
		})
	}
}

func TestRemoveIncludeFlag(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		arguments []string
		expected  []string
	}{
		{"should remove --include=a,b", []string{"plan", "--include=a,b"}, []string{"plan"}},
		{"should return unchanged when absent", []string{"plan"}, []string{"plan"}},
		{"should return nil when only include flag", []string{"--include=mod1"}, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expected, commands.RemoveIncludeFlag(tt.arguments))
		})
	}
}

func TestHasExcludeFlag(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		arguments []string
		expected  bool
	}{
		{"should return true when --exclude=mod1 present", []string{"plan", "--exclude=mod1"}, true},
		{"should return false when absent", []string{"plan"}, false},
		{"should return false when empty arguments", []string{}, false},
		{"should return false when --include=mod1 present", []string{"--include=mod1"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expected, commands.HasExcludeFlag(tt.arguments))
		})
	}
}

func TestGetExcludeValues(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		arguments     []string
		expectedValue []string
		expectedFound bool
	}{
		{"should return values when --exclude=a,b", []string{"--exclude=a,b"}, []string{"a", "b"}, true},
		{"should return single value", []string{"--exclude=mod1"}, []string{"mod1"}, true},
		{"should return false when empty value", []string{"--exclude="}, nil, false},
		{"should return false when not present", []string{"plan"}, nil, false},
		{"should trim whitespace", []string{"--exclude= a , b "}, []string{"a", "b"}, true},
		{"should return false when empty arguments", []string{}, nil, false},
		{"should handle hyphens and underscores", []string{"--exclude=my-module,other_module"}, []string{"my-module", "other_module"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			value, found := commands.GetExcludeValues(tt.arguments)
			assert.Equal(t, tt.expectedValue, value)
			assert.Equal(t, tt.expectedFound, found)
		})
	}
}

func TestRemoveExcludeFlag(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		arguments []string
		expected  []string
	}{
		{"should remove --exclude=a,b", []string{"plan", "--exclude=a,b"}, []string{"plan"}},
		{"should return unchanged when absent", []string{"plan"}, []string{"plan"}},
		{"should return nil when only exclude flag", []string{"--exclude=mod1"}, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expected, commands.RemoveExcludeFlag(tt.arguments))
		})
	}
}

func TestGetFilterValues(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name               string
		arguments          []string
		expectedInclusions []string
		expectedExclusions []string
	}{
		{
			"should return both when both flags present",
			[]string{"plan", "--include=a,b", "--exclude=c"},
			[]string{"a", "b"},
			[]string{"c"},
		},
		{
			"should return only inclusions when only include present",
			[]string{"plan", "--include=a,b"},
			[]string{"a", "b"},
			nil,
		},
		{
			"should return only exclusions when only exclude present",
			[]string{"plan", "--exclude=x,y"},
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
			[]string{"plan", "--parallel=4", "--include=mod1", "--exclude=mod2", "/path"},
			[]string{"mod1"},
			[]string{"mod2"},
		},
		{
			"should handle multiple values in both flags",
			[]string{"--include=a,b,c", "--exclude=x,y"},
			[]string{"a", "b", "c"},
			[]string{"x", "y"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := commands.GetFilterValues(tt.arguments)
			assert.Equal(t, tt.expectedInclusions, result.Inclusions)
			assert.Equal(t, tt.expectedExclusions, result.Exclusions)
		})
	}
}

func TestRemoveFilterFlags(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		arguments []string
		expected  []string
	}{
		{
			"should remove both include and exclude flags",
			[]string{"plan", "--include=a,b", "--exclude=c"},
			[]string{"plan"},
		},
		{
			"should remove only include when exclude absent",
			[]string{"plan", "--include=a"},
			[]string{"plan"},
		},
		{
			"should remove only exclude when include absent",
			[]string{"plan", "--exclude=c"},
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
			assert.Equal(t, tt.expected, commands.RemoveFilterFlags(tt.arguments))
		})
	}
}

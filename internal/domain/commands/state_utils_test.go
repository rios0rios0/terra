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

func TestHasFilterFlag(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		arguments []string
		expected  bool
	}{
		{"should return true when present", []string{"--filter=mod1"}, true},
		{"should return false when absent", []string{"plan"}, false},
		{"should return false when empty", []string{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expected, commands.HasFilterFlag(tt.arguments))
		})
	}
}

func TestGetFilterValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		arguments     []string
		expectedValue []string
		expectedFound bool
	}{
		{
			"should return values when --filter=a,b",
			[]string{"--filter=a,b"},
			[]string{"a", "b"},
			true,
		},
		{
			"should return single value",
			[]string{"--filter=mod1"},
			[]string{"mod1"},
			true,
		},
		{
			"should return false when empty value",
			[]string{"--filter="},
			nil,
			false,
		},
		{
			"should return false when not present",
			[]string{"plan"},
			nil,
			false,
		},
		{
			"should trim whitespace",
			[]string{"--filter= a , b "},
			[]string{"a", "b"},
			true,
		},
		{
			"should handle exclusion values",
			[]string{"--filter=a,!b"},
			[]string{"a", "!b"},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			value, found := commands.GetFilterValue(tt.arguments)
			assert.Equal(t, tt.expectedValue, value)
			assert.Equal(t, tt.expectedFound, found)
		})
	}
}

func TestParseFilterValues(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name               string
		filterValues       []string
		expectedInclusions []string
		expectedExclusions []string
	}{
		{
			"should separate inclusions and exclusions",
			[]string{"mod1", "!mod2", "mod3"},
			[]string{"mod1", "mod3"},
			[]string{"mod2"},
		},
		{
			"should handle only inclusions",
			[]string{"a", "b"},
			[]string{"a", "b"},
			nil,
		},
		{
			"should handle only exclusions",
			[]string{"!x", "!y"},
			nil,
			[]string{"x", "y"},
		},
		{
			"should skip empty values",
			[]string{"", "!"},
			nil,
			nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := commands.ParseFilterValues(tt.filterValues)
			assert.Equal(t, tt.expectedInclusions, result.Inclusions)
			assert.Equal(t, tt.expectedExclusions, result.Exclusions)
		})
	}
}

func TestRemoveFilterFlag(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		arguments []string
		expected  []string
	}{
		{
			"should remove filter flag",
			[]string{"plan", "--filter=a,b"},
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
			assert.Equal(t, tt.expected, commands.RemoveFilterFlag(tt.arguments))
		})
	}
}

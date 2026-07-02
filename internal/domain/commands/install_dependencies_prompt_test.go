//go:build unit

package commands_test

import (
	"testing"

	"github.com/rios0rios0/terra/internal/domain/commands"
	"github.com/stretchr/testify/require"
)

func TestIsTruthyPublic(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name  string
		value string
		want  bool
	}{
		{"true", "true", true},
		{"uppercase-true", "TRUE", true},
		{"one", "1", true},
		{"on", "on", true},
		{"yes", "yes", true},
		{"y", "y", true},
		{"padded", "  yes  ", true},
		{"false", "false", false},
		{"zero", "0", false},
		{"no", "no", false},
		{"off", "off", false},
		{"empty", "", false},
		{"arbitrary", "banana", false},
	}

	for _, testCase := range cases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			// GIVEN: an environment-variable value

			// WHEN: evaluating its truthiness
			got := commands.IsTruthyPublic(testCase.value)

			// THEN: it matches the expected result
			require.Equal(t, testCase.want, got)
		})
	}
}

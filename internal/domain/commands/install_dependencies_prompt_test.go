//go:build unit

package commands

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsTruthy(t *testing.T) {
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
			got := isTruthy(testCase.value)

			// THEN: it matches the expected result
			require.Equal(t, testCase.want, got)
		})
	}
}

func TestDecideUpdate(t *testing.T) {
	t.Parallel()

	t.Run("should update when TERRA_ASSUME_YES is truthy", func(t *testing.T) {
		t.Parallel()
		// GIVEN: assume-yes is set, even on a non-interactive session

		// WHEN: resolving the update decision
		decision := decideUpdate("true", false)

		// THEN: the update proceeds without prompting
		require.Equal(t, decisionUpdate, decision)
	})

	t.Run("should skip when non-interactive without assume-yes", func(t *testing.T) {
		t.Parallel()
		// GIVEN: no assume-yes and a non-interactive session (CI)

		// WHEN: resolving the update decision
		decision := decideUpdate("", false)

		// THEN: the update is skipped instead of blocking on stdin
		require.Equal(t, decisionSkip, decision)
	})

	t.Run("should prompt when interactive without assume-yes", func(t *testing.T) {
		t.Parallel()
		// GIVEN: an interactive session and no assume-yes

		// WHEN: resolving the update decision
		decision := decideUpdate("", true)

		// THEN: the caller prompts interactively
		require.Equal(t, decisionPrompt, decision)
	})
}

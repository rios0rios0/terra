//go:build unit

package repositories_test

import (
	"bytes"
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/rios0rios0/terra/internal/infrastructure/repositories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLinePrefixWriter_Write(t *testing.T) {
	t.Parallel()

	t.Run("should prefix a complete line when written with a trailing newline", func(t *testing.T) {
		t.Parallel()
		// given
		var dest bytes.Buffer
		writer := repositories.NewLinePrefixWriter(&dest, "mod1", &sync.Mutex{})

		// when
		n, err := writer.Write([]byte("hello world\n"))

		// then
		require.NoError(t, err)
		assert.Equal(t, len("hello world\n"), n)
		assert.Equal(t, "[mod1] hello world\n", dest.String())
	})

	t.Run("should prefix each line when multiple lines arrive in one write", func(t *testing.T) {
		t.Parallel()
		// given
		var dest bytes.Buffer
		writer := repositories.NewLinePrefixWriter(&dest, "mod1", &sync.Mutex{})

		// when
		_, err := writer.Write([]byte("line one\nline two\nline three\n"))

		// then
		require.NoError(t, err)
		assert.Equal(t, "[mod1] line one\n[mod1] line two\n[mod1] line three\n", dest.String())
	})

	t.Run("should buffer a partial line until its newline arrives", func(t *testing.T) {
		t.Parallel()
		// given
		var dest bytes.Buffer
		writer := repositories.NewLinePrefixWriter(&dest, "mod1", &sync.Mutex{})

		// when: a fragment without a newline is written
		_, err := writer.Write([]byte("partial "))

		// then: nothing is emitted yet
		require.NoError(t, err)
		assert.Empty(t, dest.String())

		// when: the rest of the line arrives
		_, err = writer.Write([]byte("line\n"))

		// then: the full line is emitted once, with a single prefix
		require.NoError(t, err)
		assert.Equal(t, "[mod1] partial line\n", dest.String())
	})

	t.Run("should not colorize the prefix when the destination is not a terminal", func(t *testing.T) {
		t.Parallel()
		// given
		var dest bytes.Buffer
		writer := repositories.NewLinePrefixWriter(&dest, "mod1", &sync.Mutex{})

		// when
		_, err := writer.Write([]byte("text\n"))

		// then: no ANSI escape sequences surround the label
		require.NoError(t, err)
		assert.Equal(t, "[mod1] text\n", dest.String())
		assert.NotContains(t, dest.String(), "\x1b[")
	})
}

func TestLinePrefixWriter_Flush(t *testing.T) {
	t.Parallel()

	t.Run("should emit the buffered partial line with a newline when flushed", func(t *testing.T) {
		t.Parallel()
		// given
		var dest bytes.Buffer
		writer := repositories.NewLinePrefixWriter(&dest, "mod1", &sync.Mutex{})
		_, err := writer.Write([]byte("no trailing newline"))
		require.NoError(t, err)

		// when
		writer.Flush()

		// then
		assert.Equal(t, "[mod1] no trailing newline\n", dest.String())
	})

	t.Run("should emit nothing when flushed with an empty buffer", func(t *testing.T) {
		t.Parallel()
		// given
		var dest bytes.Buffer
		writer := repositories.NewLinePrefixWriter(&dest, "mod1", &sync.Mutex{})
		_, err := writer.Write([]byte("complete\n"))
		require.NoError(t, err)

		// when
		writer.Flush()

		// then: only the already-complete line is present, no extra output
		assert.Equal(t, "[mod1] complete\n", dest.String())
	})
}

func TestLinePrefixWriter_Concurrency(t *testing.T) {
	t.Parallel()

	t.Run("should never interleave lines from writers sharing a mutex", func(t *testing.T) {
		t.Parallel()
		// given: several writers sharing one destination and mutex, like parallel workers
		const writerCount = 4
		const linesPerWriter = 200
		var dest bytes.Buffer
		mu := &sync.Mutex{}

		// when: each writer streams its lines from its own goroutine
		var wg sync.WaitGroup
		for w := range writerCount {
			label := fmt.Sprintf("mod%d", w)
			writer := repositories.NewLinePrefixWriter(&dest, label, mu)
			wg.Add(1)
			go func() {
				defer wg.Done()
				for line := range linesPerWriter {
					_, _ = writer.Write([]byte(fmt.Sprintf("%s-line-%d\n", label, line)))
				}
			}()
		}
		wg.Wait()

		// then: every emitted line's prefix matches its own content (no mid-line mixing)
		lines := strings.Split(strings.TrimRight(dest.String(), "\n"), "\n")
		require.Len(t, lines, writerCount*linesPerWriter)
		for _, line := range lines {
			closeIdx := strings.IndexByte(line, ']')
			require.Positive(t, closeIdx, "malformed line (no closing bracket): %q", line)
			label := line[1:closeIdx]
			assert.True(t, strings.HasPrefix(line, "["+label+"] "+label+"-line-"),
				"line prefix and content disagree, indicating interleaving: %q", line)
		}
	})

	t.Run("should be safe under concurrent Write and Flush on the same writer", func(t *testing.T) {
		t.Parallel()
		// given: a single writer targeted by many goroutines at once (guarded by -race)
		const goroutines = 8
		const linesPerGoroutine = 100
		var dest bytes.Buffer
		writer := repositories.NewLinePrefixWriter(&dest, "mod", &sync.Mutex{})

		// when: goroutines Write complete lines concurrently and each Flushes after
		var wg sync.WaitGroup
		for g := range goroutines {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for line := range linesPerGoroutine {
					_, _ = writer.Write([]byte(fmt.Sprintf("g%d-line-%d\n", g, line)))
					writer.Flush()
				}
			}()
		}
		wg.Wait()
		writer.Flush()

		// then: the per-writer mutex keeps the buffer consistent — every emitted line is
		// fully prefixed and none are lost or corrupted
		lines := strings.Split(strings.TrimRight(dest.String(), "\n"), "\n")
		require.Len(t, lines, goroutines*linesPerGoroutine)
		for _, line := range lines {
			assert.True(t, strings.HasPrefix(line, "[mod] "),
				"every line must carry the prefix even under concurrent access: %q", line)
		}
	})
}

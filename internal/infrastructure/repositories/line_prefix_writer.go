package repositories

import (
	"bytes"
	"hash/fnv"
	"io"
	"os"
	"sync"

	"golang.org/x/term"
)

// LinePrefixWriter is a thread-safe [io.Writer] that prepends a fixed prefix to every line
// it forwards to dest. It buffers partial lines until their terminating newline arrives.
// The per-writer bufMu guards that buffer so concurrent Write/Flush calls on the same
// writer are safe, while the shared mu serializes writes to dest so output from
// concurrently executing modules never interleaves mid-line.
type LinePrefixWriter struct {
	dest   io.Writer
	prefix []byte
	mu     *sync.Mutex // shared across writers: serializes writes to dest
	bufMu  sync.Mutex  // per-writer: guards buf for concurrent Write/Flush
	buf    []byte
}

// NewLinePrefixWriter builds a writer that labels every line with "[label] ". When dest is
// an interactive terminal (and NO_COLOR is unset), the label is colorized with a stable
// per-label color so parallel module streams are easy to tell apart. The mu is shared
// across all writers targeting the same console so their lines stay separated.
func NewLinePrefixWriter(dest io.Writer, label string, mu *sync.Mutex) *LinePrefixWriter {
	rendered := "[" + label + "]"
	if shouldColorize(dest) {
		start, reset := colorForLabel(label)
		rendered = start + rendered + reset
	}

	return &LinePrefixWriter{
		dest:   dest,
		prefix: []byte(rendered + " "),
		mu:     mu,
	}
}

// Write buffers p and emits every complete (newline-terminated) line with the prefix,
// keeping any trailing partial line buffered for the next call.
func (w *LinePrefixWriter) Write(p []byte) (int, error) {
	w.bufMu.Lock()
	defer w.bufMu.Unlock()

	w.buf = append(w.buf, p...)

	for {
		index := bytes.IndexByte(w.buf, '\n')
		if index < 0 {
			break
		}

		w.emit(w.buf[:index+1])
		w.buf = w.buf[index+1:]
	}

	if len(w.buf) == 0 {
		w.buf = nil
	}

	return len(p), nil
}

// Flush emits any buffered bytes that were not terminated by a newline, appending one so
// the trailing line is still prefixed and readable. Call it once the process exits.
func (w *LinePrefixWriter) Flush() {
	w.bufMu.Lock()
	defer w.bufMu.Unlock()

	if len(w.buf) == 0 {
		return
	}

	w.emit(append(w.buf, '\n'))
	w.buf = nil
}

// emit writes the prefix and the line to dest as a single locked write so lines from
// other writers sharing mu cannot split this one.
func (w *LinePrefixWriter) emit(line []byte) {
	out := make([]byte, 0, len(w.prefix)+len(line))
	out = append(out, w.prefix...)
	out = append(out, line...)

	w.mu.Lock()
	defer w.mu.Unlock()
	_, _ = w.dest.Write(out)
}

// shouldColorize reports whether ANSI colors should be emitted to w: only when w is an
// interactive terminal and the NO_COLOR convention (https://no-color.org) is not set.
func shouldColorize(w io.Writer) bool {
	if _, disabled := os.LookupEnv("NO_COLOR"); disabled {
		return false
	}

	file, ok := w.(*os.File)
	if !ok {
		return false
	}

	return term.IsTerminal(int(file.Fd()))
}

// colorForLabel returns a stable ANSI color and its reset sequence for label, hashing the
// label so the same module keeps the same color for the whole run. The palette is a fixed
// array so len(palette) is a compile-time constant.
func colorForLabel(label string) (string, string) {
	palette := [...]string{
		"\x1b[36m", // cyan
		"\x1b[32m", // green
		"\x1b[33m", // yellow
		"\x1b[35m", // magenta
		"\x1b[34m", // blue
		"\x1b[91m", // bright red
		"\x1b[96m", // bright cyan
		"\x1b[95m", // bright magenta
	}
	const reset = "\x1b[0m"

	hasher := fnv.New32a()
	_, _ = hasher.Write([]byte(label))

	return palette[hasher.Sum32()%uint32(len(palette))], reset
}

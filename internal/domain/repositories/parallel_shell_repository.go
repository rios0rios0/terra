package repositories

// ParallelShellRepository executes a shell command while streaming its output through a
// per-invocation line prefix, so concurrent module executions remain attributable in the
// combined console output. It is a separate, focused port from ShellRepository because
// only the parallel worker pool needs the prefixing behavior.
type ParallelShellRepository interface {
	ExecuteCommandWithPrefix(command string, arguments []string, directory, prefix string) error
}

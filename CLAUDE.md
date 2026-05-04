# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What is Terra?

Terra is a Go CLI tool that wraps Terraform and Terragrunt, providing simplified path-based infrastructure management inspired by Kubernetes. It adds parallel execution, centralized caching, auto-reply prompts, and self-update capabilities.

## Build & Development Commands

Requires Go 1.26+ (see `go.mod`). Make targets include `$(SCRIPTS_DIR)/makefiles/{common,golang}.mk` from the [pipelines](https://github.com/rios0rios0/pipelines) project — defaults to `$HOME/Development/github.com/rios0rios0/pipelines`, override with `SCRIPTS_DIR=...`.

```bash
make build              # Build binary to bin/terra (stripped, ~15-20s)
make install            # Build and copy to ~/.local/bin/terra
make debug              # Build with debug symbols (-N -l)
make run                # go run ./cmd/terra

# Quality gates (use pipelines project, auto-cloned via HTTPS)
make lint               # golangci-lint (~2-5 min, never cancel)
make test               # Tests with coverage (~1-2 min)
make sast               # Security scanning (CodeQL, Semgrep, Trivy, Gitleaks)
make all                # lint + sast + test

# Quick local checks
go fmt ./...
go vet ./...

# Run tests by category
go test -tags unit ./...         # Unit tests only
go test -tags integration ./...  # Integration tests only

# Run a single test / subtest (pass the build tag the file declares)
go test -tags unit ./internal/domain/commands -run TestRunFromRootCommand_Execute
go test -tags unit ./internal/domain/commands -run 'TestRunFromRootCommand_Execute/should_error_when.*'
go test -tags unit -v -run TestParallelStateCommand_Execute ./internal/domain/commands
```

Never call tool binaries directly — always use Makefile targets which load correct configs from the pipelines project.

## Architecture

Clean Architecture (Hexagonal/Ports & Adapters) with DDD and Uber DIG for dependency injection.

```
cmd/terra/
  main.go              # Cobra CLI setup, command routing
  dig.go               # DI container creation, injection functions
internal/
  app.go               # AppInternal orchestrator
  container.go         # Top-level DI provider registration
  domain/              # Contracts layer
    commands/           # Business logic (RunFromRoot, ParallelState, DeleteCache, etc.)
    entities/           # Framework-agnostic types (Controller, Settings, OS, CLI, Platform)
    repositories/       # Interface contracts
  infrastructure/      # Implementations layer
    controllers/        # Cobra command adapters
      helpers/          # Argument parsing
    repositories/       # Repository implementations
test/                  # Test helpers only (never in production folders)
  domain/
    commanddoubles/     # Command stubs
    entitybuilders/     # Fluent entity builders
    entitydoubles/      # Entity stubs (CLI, OS)
  infrastructure/
    repositorybuilders/ # HTTP test server builders
    repositorydoubles/  # Repository stubs
    repositoryhelpers/  # OS/network test helpers
```

**Dependency flow:** `infrastructure → domain` (never the reverse). Each layer has a `container.go` for DIG provider registration.

## Key Mechanisms

- **Centralized caching:** Sets `TG_DOWNLOAD_DIR` and `TG_PROVIDER_CACHE_DIR` automatically
- **CAS:** Enables `TG_EXPERIMENT=cas` by default (disable with `TERRA_NO_CAS=true`)
- **Provider caching:** Uses `TG_PROVIDER_CACHE` (Provider Cache Server) for concurrent-safe provider deduplication with file locking; disable with `TERRA_NO_PROVIDER_CACHE=true`. Terra also sets `TG_NO_AUTO_PROVIDER_CACHE_DIR=true` alongside it — without this, Terragrunt 0.99+'s CAS-auto-enabled `auto-provider-cache-dir` experiment silently overrides `TG_PROVIDER_CACHE_DIR` and writes providers per-module inside `TG_DOWNLOAD_DIR`, which bypasses the shared cache terra advertises. Keep these two env vars paired in `configureCacheEnvironment`
- **Partial Parse Config Cache:** Enables `TG_USE_PARTIAL_PARSE_CONFIG_CACHE=true` by default (disable with `TERRA_NO_PARTIAL_PARSE_CACHE=true`)
- **Auto-upgrade:** `UpgradeAwareShellRepository` detects backend/provider failures and retries with `init --upgrade`. When the original command included queue-scoping flags (`--all`, `--filter`, `--queue-include-dir`, `--queue-exclude-dir`, `--queue-include-units-reading`, `--queue-strict-include`, `--queue-include-external`, `--queue-exclude-external`), they are forwarded to the retry so init walks the same queue — without this the retry would be `terragrunt init --upgrade` in a parent directory that has no `terragrunt.hcl` and fail with "You attempted to run terragrunt in a folder that does not contain a terragrunt.hcl file"
- **Parallel execution (terra-managed):** `--parallel=N` runs across modules via `ParallelStateCommand`; use `--only=mod1,mod2` to select modules or `--skip=mod3` to exclude them. These are terra-managed flags and only work with `--parallel=N`.
- **Parallel execution (terragrunt-managed):** `--all`, `--parallelism=N`, and `--filter=query` are forwarded to terragrunt as-is. Filter modules on this path with terragrunt's `--filter='!mod'` (preferred) or `--queue-exclude-dir=mod`. `--parallel` and `--all` cannot be combined.
- **Selection-flag errors are educational:** Using `--only`/`--skip` without `--parallel` is fatal; the error echoes the user's command and prints both valid forms (`--parallel=5 --skip=mod` AND `--all --filter='!mod'`). Using terragrunt's `--filter`/`--queue-exclude-dir`/`--queue-include-dir` with `--parallel` is non-fatal; it logs a warning because terra's worker pool silently ignores those flags. When editing validation, update `BuildSelectionFlagsError` / `BuildParallelAllConflictError` in `internal/domain/commands/run_from_root_error_builders.go`, not the call sites in `run_from_root_command.go`.
- **Confirmation flags:** `--yes` / `-y` injects Terragrunt's `--non-interactive` plus Terraform's `-auto-approve`; `--no` / `-n` injects only `--non-interactive` (Terraform's apply prompt aborts). Required for `--parallel` with `apply`/`destroy`. The legacy `--reply` / `-r` flags still work but are deprecated and emit a migration warning; they map to `--yes` or `--no` internally

Deep-dive docs for the parallel subsystem live in `docs/parallel-execution.md` and `docs/parallel-git-clone-race.md` — read these before touching worker-pool, filter, or git-clone-race logic.

## Testing Conventions

- **Framework:** `stretchr/testify` with `github.com/rios0rios0/testkit` for shared builders
- **Build tags required on every test file:** `//go:build unit` or `//go:build integration`
- **Test helpers use:** `//go:build integration || unit || test`
- **BDD structure:** `// GIVEN` / `// WHEN` / `// THEN` comment blocks
- **Naming:** `TestStructName_MethodName` with subtests `"should [behavior] when [condition]"`
- **Parallel:** Use `t.Parallel()` unless test uses `t.Setenv()` or `t.Chdir()`
- **Builders:** Fluent API pattern in `/test/domain/entitybuilders/` (e.g., `NewSettingsBuilder().WithTerraModuleCacheDir(dir).BuildSettings()`)
- **One utility per file** in `/test` — each builder, stub, helper in its own file
- **Prefer stubs over mocks** — use mocks only when behavior verification is needed

## Adding New Commands

1. Create command in `internal/domain/commands/`
2. Create controller in `internal/infrastructure/controllers/`
3. Register DIG providers in respective `container.go` files
4. Add tests with build tags, following existing patterns

## File Standards

- **LF line endings** enforced via `.editorconfig`
- **Update `CHANGELOG.md`** under `[Unreleased]` for features/fixes (not doc-only changes)

# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

When a new release is proposed:

1. Create a new branch `bump/x.x.x` (this isn't a long-lived branch!!!);
2. The Unreleased section on `CHANGELOG.md` gets a version number and date;
3. Open a Pull Request with the bump version changes targeting the `main` branch;
4. When the Pull Request is merged, a new Git tag must be created using [GitHub environment](https://github.com/rios0rios0/terra/tags).

Releases to productive environments should run from a tagged version.
Exceptions are acceptable depending on the circumstances (critical bug fixes that can be cherry-picked, etc.).

## [Unreleased]

## [1.15.3] - 2026-04-30

### Changed

- changed the Go module dependencies to their latest versions

## [1.15.2] - 2026-04-29

### Changed

- changed the Go module dependencies to their latest versions

## [1.15.1] - 2026-04-28

### Changed

- refreshed `CLAUDE.md` and `.github/copilot-instructions.md` to document the `--yes`/`--no` confirmation flags and the `--reply` deprecation from v1.15.0, and removed references to non-existent test directories

## [1.15.0] - 2026-04-24

### Added

- added `--yes` / `-y` and `--no` / `-n` confirmation flags that translate to native Terraform and Terragrunt flags. `--yes` injects Terragrunt's `--non-interactive` plus Terraform's `-auto-approve` for `apply` / `destroy`; `--no` injects only `--non-interactive`, so Terraform's apply prompt aborts instead of proceeding. This aligns terra with the `apt` / `npm` / `az` convention and replaces the previous PTY-based approach with reliable flag injection.

### Changed

- changed the Go module dependencies to their latest versions
- changed the validation error for `--parallel` with `apply` / `destroy` from "`--reply` is required" to "`--yes` is required", matching the new flag names in error messages and copy-pasteable suggestions.

### Deprecated

- deprecated the `--reply` / `-r` flags. They still work and are translated to the new `--yes` / `--no` flag injection (`--reply=y` and bare `--reply` map to `--yes`; `--reply=n` maps to `--no`), but now emit a one-time migration warning. `--reply` will be removed in a future release.

### Fixed

- fixed `terra apply --reply=y` silently waiting forever on Terraform's "Do you want to perform these actions? Enter a value:" prompt. The previous PTY-based auto-responder only matched `[y/n]` and "external dependency" prompts, so Terraform's apply confirmation (which requires the literal word `yes`) was never answered. Users had to fall back to `-auto-approve`. The new flag-injection path invokes `-auto-approve` natively, so `terra apply --yes` (and the deprecated `--reply=y`) now work reliably.

## [1.14.3] - 2026-04-20

### Fixed

- fixed a flaky `TestUpgradeAwareShellRepository_ExecuteCommandWithUpgrade` subtest that intermittently failed in CI with `fork/exec ...: text file busy`. The test writes an executable shell script and immediately runs it; under `t.Parallel()` another goroutine could inherit to write fd during a fork, so `execve` saw the inode as still open-for-write and returned `ETXTBSY` (see `golang/go#22315`). Script creation now pipes through a subprocess, so no write fd for the script ever lives in the test process.
- fixed the centralized provider cache being silently bypassed: terra now sets `TG_NO_AUTO_PROVIDER_CACHE_DIR=true` whenever the Provider Cache Server is enabled, so Terragrunt `0.99+`'s CAS-auto-enabled `auto-provider-cache-dir` experiment stops overriding `TG_PROVIDER_CACHE_DIR`. Before this fix, providers were duplicated into `TG_DOWNLOAD_DIR/<hash>/.../.terraform/providers/` per go-getter source, `~/.cache/terra/providers/` stayed empty, and every new stack or concurrent terminal paid a full provider download. After the fix, providers download once and are shared across every stack, repo, and terminal until the version changes. Opt out with `TERRA_NO_PROVIDER_CACHE=true`.

## [1.14.2] - 2026-04-17

### Changed

- changed the Go module dependencies to their latest versions

## [1.14.1] - 2026-04-16

### Changed

- changed the Go module dependencies to their latest versions

## [1.14.0] - 2026-04-15

### Added

- added a non-fatal warning when Terragrunt-only flags (`--filter`, `--queue-exclude-dir`, `--queue-include-dir`) are combined with terra's `--parallel=N`, since they are silently ignored by terra's worker pool; the warning nudges users toward `--only`/`--skip` or toward switching to `--all`

### Changed

- changed the `--only`/`--skip` validation error to echo the user's command and show both valid forms (`--parallel=N --skip=mod1` and `--all --filter='!mod1'`), teaching the `--filter` alternative for the `--all` path instead of leaving users to discover terragrunt's native flags on their own
- changed the `--parallel` + `--all` conflict error to echo the user's command and offer both alternative forms as copy-pasteable examples
- changed the `terra --help` text to include a "Parallel execution strategies" block that summarizes when to use `--parallel=N` versus `--all`, making the split discoverable without reading the docs
- changed the Go version to `1.26.2` and updated all module dependencies

## [1.13.0] - 2026-04-14

### Added

- added automatic version check on CLI startup using `CheckForUpdates()`

### Changed

- changed the Go module dependencies to their latest versions

## [1.12.0] - 2026-04-03

### Added

- added `--reply` requirement when using `--parallel` with `apply` or `destroy` to prevent workers from hanging on interactive prompts; for terra-managed parallel, just `--reply` (no value) is sufficient since terra always injects `--non-interactive` when `--reply` is present and adds `-auto-approve` automatically for interactive commands like `apply` and `destroy`
- added documentation for the Git `refs/files-backend.c` race condition that occurs during parallel execution with shared dependencies, including root cause analysis and workarounds
- added validation requiring `--reply=<value>` (with explicit value) when used with `--all`, since the PTY auto-answering needs to know whether to respond "y" or "n"
- added warning when `--reply=<value>` is used with `--parallel`, informing the user the value is ignored and only meaningful with `--all` (Terragrunt-managed parallelism)

### Changed

- changed `--all` flag to always forward to Terragrunt (no longer intercepted by terra for state commands); use `--parallel=5` instead for terra-managed parallel state operations
- changed `--include` flag to `--only` and `--exclude` flag to `--skip` for terra's parallel module selection, eliminating name collisions with terragrunt's own `--include`/`--exclude` flags
- changed `cliforge` import paths to reflect upstream package restructuring
- changed self-update command to delegate to `cliforge/selfupdate` shared library, removing ~300 lines of duplicated GitHub API, archive extraction, and binary replacement logic

### Fixed

- fixed parallel `apply`/`destroy` with `--reply` not injecting `-auto-approve`, causing terraform to prompt for confirmation and hang workers
- fixed parallel module discovery descending into `.terragrunt-cache` and other hidden directories, which caused hundreds of cached dependency modules to be processed as actual targets

### Removed

- removed `--auto-answer` / `-a` flags; replaced with `--reply` / `-r` to avoid collision with Terragrunt's `-a` shorthand for `--all`
- removed `--no-parallel-bypass` flag (`--all` now always forwards to Terragrunt; use Terragrunt's `--parallelism=N` directly for Terragrunt-managed parallelism)
- removed legacy `--all` support for state commands (`import`, `state rm`, etc.); use `--parallel=N` instead

## [1.11.0] - 2026-04-01

### Added

- added validation for `--include`/`--exclude` flag combinations with `--parallel`, `--no-parallel-bypass`, and conflict detection

### Changed

- changed `--filter` flag to separate `--include` and `--exclude` flags for parallel execution, eliminating Bash shell escaping issues with the `!` exclusion prefix

## [1.10.1] - 2026-03-31

### Changed

- changed the Go module dependencies to their latest versions

## [1.10.0] - 2026-03-30

### Added

- added `TERRA_NO_PROVIDER_CACHE` environment variable to disable the Terragrunt Provider Cache Server (opt-out toggle)

### Changed

- changed provider caching strategy: replaced `TF_PLUGIN_CACHE_DIR` with Terragrunt Provider Cache Server (`TG_PROVIDER_CACHE=1` + `TG_PROVIDER_CACHE_DIR`) to fix "text file busy" errors during parallel execution (`--parallel=N`)

### Fixed

- fixed upgrade-aware retry triggering `init --upgrade` after user-canceled apply/plan/destroy operations

## [1.9.0] - 2026-03-24

### Added

- added `TERRA_NO_WORKSPACE` environment variable to disable automatic workspace selection from `TERRA_WORKSPACE`
- added unit tests for DIG container registration, self-update command, run-from-root command, upgrade-aware shell repository, version command, and OS operations

### Changed

- changed `clear` command to also remove `terragrunt-cache` (without leading dot) and `.terraform.lock.hcl` lock files

### Fixed

- fixed `clear` command not resetting found paths between iterations, causing already-deleted entries to be re-processed
- fixed `RunAdditionalBeforeCommand` tests using hard-coded `/test/path` instead of `t.TempDir()`, making them environment-dependent
- fixed auto-init running `terragrunt init` on every command even when `.terraform` directory already exists
- fixed auto-upgrade detection logging which pattern triggered the retry, aiding future debugging
- fixed overly broad auto-upgrade detection that triggered unnecessary `init --upgrade` on runtime provider errors (e.g., TLS failures)
- fixed proactive init not detecting `.terragrunt-cache` and legacy `terragrunt-cache` directories, causing unnecessary `terragrunt init` on every command when using Terragrunt
- fixed proactive init running unnecessarily when centralized caching (`TG_DOWNLOAD_DIR`) is active and cache already has content

## [1.8.0] - 2026-03-17

### Added

- added execution timing to command logs showing how long each Terragrunt invocation took

### Changed

- changed provider caching strategy: removed `TG_PROVIDER_CACHE` (Provider Cache Server) in favor of `TF_PLUGIN_CACHE_DIR` only, which benchmarks showed is faster (8.9s vs 10.6s warm) with identical disk savings via symlinks
- changed version management to use build-time `ldflags` injection instead of hardcoded constant

### Fixed

- fixed `terra self-update` failing due to incorrect asset name matching (expected `terra_os_arch` but releases use `terra-version-os-arch.tar.gz`) and missing archive extraction
- fixed Terragrunt deprecation warning by replacing `TERRAGRUNT_USE_PARTIAL_PARSE_CONFIG_CACHE` with `TG_USE_PARTIAL_PARSE_CONFIG_CACHE`

### Removed

- removed `TERRA_NO_PROVIDER_CACHE` environment variable (Provider Cache Server replaced by `TF_PLUGIN_CACHE_DIR`)

## [1.7.1] - 2026-03-14

### Changed

- changed the Go module dependencies to their latest versions

## [1.7.0] - 2026-03-12

### Added

- added Terragrunt Partial Parse Config Cache enabled by default (`TG_USE_PARTIAL_PARSE_CONFIG_CACHE=true`) for faster HCL config parsing across modules sharing the same root include; disabled with `TERRA_NO_PARTIAL_PARSE_CACHE=true`
- added Terragrunt Provider Cache Server enabled by default (`TG_PROVIDER_CACHE=1`) for localhost proxy-based provider deduplication via symlinks; disabled with `TERRA_NO_PROVIDER_CACHE=true`

### Changed

- changed the Go version to `1.26.1` and updated all module dependencies
- replaced raw struct literals in tests with `testkit` builders for consistent test data construction

### Fixed

- fixed opt-out toggles (`TERRA_NO_CAS`, `TERRA_NO_PROVIDER_CACHE`, `TERRA_NO_PARTIAL_PARSE_CACHE`) to explicitly unset the corresponding environment variables, ensuring deterministic behavior when the parent environment has pre-existing values

### Removed

- removed cross-platform file locking mechanism (`gofrs/flock`) that prevented running multiple terra instances simultaneously from the same repository; CAS and centralized caching make it unnecessary

## [1.6.1] - 2026-02-14

### Fixed

- fixed cross-compilation failure for `darwin` (macOS) targets by replacing platform-specific `os_linux.go` with `os_unix.go` using `//go:build !windows` constraint, and renamed `OSLinux` to `OSUnix`

## [1.6.0] - 2026-02-14

### Added

- added auto-initialization with upgrade detection: when Terragrunt commands fail due to uninitialized modules, backend changes, or provider version conflicts, terra automatically runs `init --upgrade` and retries the original command

### Changed

- migrated test builders to use the `testkit` library (`github.com/rios0rios0/testkit`) for standardized builder patterns
- renamed `cmd/terra/wire.go` to `cmd/terra/dig.go` to reflect the actual DI tool in use (Uber DIG, not Google Wire)

## [1.5.0] - 2026-02-13

### Added

- added `--auto-answer` flag (`-a`) to automatically handle Terragrunt prompts, with configurable responses (`--auto-answer=y` or `-a=n`; defaults to "n" for backward compatibility)
- added `--global` flag to the `clear` command to also remove centralized cache directories
- added `update` command as an alias for `install`
- added `version` command to display Terra, Terraform, and Terragrunt versions
- added centralized Terragrunt module and provider caching (`TG_DOWNLOAD_DIR`, `TF_PLUGIN_CACHE_DIR`) configured automatically before every invocation, with override via `TERRA_MODULE_CACHE_DIR` and `TERRA_PROVIDER_CACHE_DIR`
- added cross-platform file locking via `gofrs/flock` to prevent race conditions when multiple terra processes run concurrently from the same repository
- added dependency injection (first Wire, then DIG) and inverted all dependencies
- added installation shell script (`install.sh`) for automated terra installation from GitHub releases with platform detection, `--version`, `--force`, `--dry-run`, and `--install-dir` options
- added parallel execution for any command via `--parallel=N`, including state manipulation commands (`import`, `state rm`, `state mv`, `state pull`, `state push`, `state show`) with `--all`, automatic module discovery, configurable concurrency (default: 5 jobs), and error aggregation
- added self-update feature to update the CLI without any additional step
- added Terragrunt CAS (Content Addressable Store) enabled by default for Git clone deduplication via hard links; disabled with `TERRA_NO_CAS=true`
- added validation on the `settings` entity

### Changed

- changed the documentation with pipelines and minor change to template files
- corrected controllers responsibilities mapping the external to internal entities
- corrected the structure to follow best practices using DDD
- decoupled responsibilities from just one command to other layers
- moved all business logic to the domain structure
- replaced deprecated `run-all` command syntax with `--all` flag to align with Terragrunt's new syntax
- replaced Wire with DIG for dependency injection to support Go `1.25.1` and active maintenance
- updated `.editorconfig` to enforce LF line endings across all file types instead of just Go files
- updated Copilot instructions and contributing guide to enforce `LF` (Unix) line endings for all new and edited files
- updated Copilot instructions to use `rios0rios0/pipelines` project for linting and CI tools instead of direct tool installation
- updated documentation to require `CHANGELOG.md` updates for new features and bug fixes (not required for documentation-only changes)
- upgraded the project to Go `1.26` and all the dependencies

### Fixed

- fixed optional environment variables validation for `TERRA_CLOUD` to allow empty values while still enforcing valid values when provided
- fixed permission denied errors when normal users try to download dependencies via `terra install` or `terra update` by using unique temporary file creation instead of predictable file names
- fixed slice bounds error in `ArgumentsHelper` when no arguments are provided
- fixed the issue where version checks for Terraform and Terragrunt were triggered on every command execution, causing unnecessary network calls and slowdowns
- fixed version checks to only occur when explicitly running "install" or "update" commands

## [1.4.0] - 2024-08-07

### Added

- added a new environment variable to handle Azure subscriptions

### Fixed

- fixed the required workspace flag to be optional

## [1.3.0] - 2024-07-08

### Added

- added the `godotenv` to handle the environment variables

### Changed

- changed the main command to accept input from the user and wait for the `stdin` to be closed
- changed to forward unknown flags to Terraform and Terragrunt - [#3](https://github.com/rios0rios0/terra/issues/3)

## [1.2.0] - 2023-11-10

### Added

- added the `clear` command to remove the cache and temporary files
- added the `fmt` command to format all Terraform and Terragrunt files

## [1.0.0] - 2023-11-10

### Added

- created the first version working properly and installing all dependencies

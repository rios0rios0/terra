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

### Fixed

- fixed cross-compilation failure for darwin (macOS) targets by replacing platform-specific `os_linux.go` with `os_unix.go` using `//go:build !windows` constraint, and renamed `OSLinux` to `OSUnix`

## [1.6.0] - 2026-02-14

### Added

- added auto-initialization with upgrade detection: when Terragrunt commands fail due to uninitialized modules, backend changes, or provider version conflicts, terra automatically runs `init --upgrade` and retries the original command

### Changed

- migrated test builders to use the `testkit` library (`github.com/rios0rios0/testkit`) for standardized builder patterns
- renamed `cmd/terra/wire.go` to `cmd/terra/dig.go` to reflect the actual DI tool in use (Uber DIG, not Google Wire)

## [1.5.0] - 2026-02-13

### Added

- added Terragrunt CAS (Content Addressable Store) enabled by default for Git clone deduplication via hard links; disabled with `TERRA_NO_CAS=true`
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
- added validation on the `settings` entity

### Changed

- changed the documentation with pipelines and minor change to template files
- corrected controllers responsibilities mapping the external to internal entities
- corrected the structure to follow best practices using DDD
- decoupled responsibilities from just one command to other layers
- moved all business logic to the domain structure
- replaced Wire with DIG for dependency injection to support Go `1.25.1` and active maintenance
- replaced deprecated `run-all` command syntax with `--all` flag to align with Terragrunt's new syntax
- updated Copilot instructions and contributing guide to enforce `LF` (Unix) line endings for all new and edited files
- updated Copilot instructions to use `rios0rios0/pipelines` project for linting and CI tools instead of direct tool installation
- updated `.editorconfig` to enforce LF line endings across all file types instead of just Go files
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

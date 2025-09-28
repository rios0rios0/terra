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

### Added
- Enhanced `--auto-answer` feature to support configurable responses (`--auto-answer=y` or `--auto-answer=n`)
- Auto-initialization mode with upgrade detection for terraform/terragrunt commands
  - Automatically detects when `terraform init --upgrade` or `terragrunt init --upgrade` is needed
  - Handles common scenarios: module not initialized, backend changes, provider version conflicts
  - Supports both terraform and terragrunt error patterns
  - Automatically retries original command after successful upgrade initialization
- Installation shell script (`install.sh`) for automated terra installation from GitHub releases
  - Platform detection (Linux, macOS, Windows) with architecture support (amd64, arm64, 386, arm)
  - One-liner installation command: `curl -fsSL https://raw.githubusercontent.com/rios0rios0/terra/main/install.sh | sh`
  - Support for installation options: `--version`, `--force`, `--dry-run`, `--install-dir`
  - Environment variable support: `TERRA_INSTALL_DIR`, `TERRA_VERSION`, `TERRA_FORCE`, `TERRA_DRY_RUN`
  - Comprehensive error handling and user-friendly output with colors
  - Installation to `~/.local/bin` by default (follows existing terra pattern)
  - Verification of installation and PATH guidance
- Support for short form syntax (`-a=y` or `-a=n`) for auto-answer configuration
- Backward compatibility maintained - boolean `--auto-answer` and `-a` flags default to "n"
- Comprehensive unit and integration tests for auto-answer functionality

### Changed

- added parallel execution support for state manipulation commands (`import`, `state rm`, `state mv`, `state pull`, `state push`, `state show`) when using the `--all` flag
- added automatic discovery of Terraform/Terragrunt modules in subdirectories for parallel execution
- added configurable concurrency control (default: 5 parallel jobs) for state operations
- added error aggregation and reporting for parallel operations with detailed logging
- added "update" command as an alias for "install" command
- added `--auto-answer` flag to automatically handle Terragrunt prompts
- added `version` command to display Terra, Terraform, and Terragrunt versions
- added a self-update feature to update the CLI without any additional step
- added an interactive shell repository for auto-answering functionality
- added dependency injection with Wire and inverted all dependencies
- created validation on the `settings` entity

### Changed

- **BREAKING CHANGE**: Replaced deprecated `run-all` command syntax with `--all` flag to align with Terragrunt's new syntax
- changed the documentation with pipelines and minor change to template files
- corrected controllers responsibilities mapping the external to internal entities
- corrected dependency injection architecture
- corrected the structure to follow best practices using DDD
- decoupled responsibilities from just one command to other layers
- moved all business logic to the domain structure
- replaced Wire with DIG for dependency injection to support Go `1.25.1` and active maintenance
- updated Copilot instructions and contributing guide to enforce LF (Unix) line endings for all new and edited files
- updated Copilot instructions to use rios0rios0/pipelines project for linting and CI tools instead of direct tool installation
- updated `.editorconfig` to enforce LF line endings across all file types instead of just Go files
- updated documentation to require `CHANGELOG.md` updates for new features and bug fixes (not required for documentation-only changes)
- upgraded the project to Go `1.23` and all the dependencies

### Fixed

- fixed optional environment variables validation for TERRA_CLOUD to allow empty values while still enforcing valid values when provided
- fixed permission denied errors when normal users try to download dependencies via `terra install` or `terra update` by using unique temporary file creation instead of predictable file names
- fixed slice bounds error in ArgumentsHelper when no arguments are provided
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

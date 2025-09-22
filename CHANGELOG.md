# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

When a new release is proposed:

1. Create a new branch `bump/x.x.x` (this isn't a long-lived branch!!!);
2. The Unreleased section on `CHANGELOG.md` gets a version number and date;
3. Open a Pull Request with the bump version changes targeting the `main` branch;
4. When the Pull Request is merged, a new `git` tag must be created using [GitHub environment](https://github.com/rios0rios0/terra/tags).

Releases to productive environments should run from a tagged version.
Exceptions are acceptable depending on the circumstances (critical bug fixes that can be cherry-picked, etc.).

## [Unreleased]

### Added

- added dependency injection with Wire and inverted all dependencies
- created validation on the `settings` entity
- added `--auto-answer` flag to automatically handle Terragrunt prompts
- added interactive shell repository for auto-answering functionality
- added `version` command to display Terra, Terraform, and Terragrunt versions

### Changed

- Updated Copilot instructions and contributing guide to enforce LF (Unix) line endings for all new and edited files
- Updated documentation to require CHANGELOG.md updates for new features and bug fixes (not required for documentation-only changes)
- Updated .editorconfig to enforce LF line endings across all file types instead of just Go files
- Updated Copilot instructions to use rios0rios0/pipelines project for linting and CI tools instead of direct tool installation
- changed the documentation with pipelines and minor change to template files
- corrected controllers responsibilities mapping the external to internal entities
- corrected dependency injection architecture
- corrected the structure to follow best practices using DDD
- decoupled responsibilities from just one command to other layers
- moved all business logic to the domain structure
- upgraded the project to Go 1.23 and all the dependencies

### Fixed

- fixed slice bounds error in ArgumentsHelper when no arguments are provided
- fixed optional environment variables validation for TERRA_CLOUD to allow empty values while still enforcing valid values when provided

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

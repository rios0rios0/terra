# Terra - Terraform/Terragrunt CLI Wrapper

Terra is a Go CLI application that wraps Terraform and Terragrunt functionality, providing simplified path-based infrastructure management inspired by Kubernetes.

Always reference these instructions first and fallback to search or bash commands only when you encounter unexpected information that does not match the info here.

## Working Effectively

### Prerequisites and Environment Setup
- Ensure Go 1.26+ is installed and `go` is in PATH
- For linting and CI tools, use the pipelines project (https://github.com/rios0rios0/pipelines)
- NEVER CANCEL: Build takes 15-20 seconds. NEVER CANCEL. Set timeout to 60+ minutes for safety.

### Building and Installing
Bootstrap and build the repository:
```bash
# Build terra binary
make build
```

Install to `~/.local/bin/terra`:
```bash
make install
```

Run without installing:
```bash
# Build first, then run from bin/
make build
./bin/terra [command]
```

### Code Quality and Validation
- **NEVER CANCEL**: golangci-lint takes 2-5 minutes. Set timeout to 30+ minutes.
- Always run code formatting, linting and testing before committing using Makefile targets:
```bash
# Standard Go tools (fast)
go fmt ./...
go vet ./...

# Full pipeline using Makefile targets (automatically clones pipelines project via HTTPS)
make lint    # Runs golangci-lint via pipelines project
make test    # Runs tests via pipelines project
make sast    # Runs security scanning via pipelines project
make all     # Runs lint + sast + test
```

**File Standards:**
- **Use LF line endings** - All new and edited files must use LF (Unix) line endings, not CRLF (Windows)
- **Update CHANGELOG.md** - Add entries to the `[Unreleased]` section for new features and bug fixes (not required for documentation-only changes)
- The `.editorconfig` file enforces line ending standards for most editors

- **Tests exist** in this repository with good coverage across domain and infrastructure layers
- **All new tests must use testify framework** and follow BDD structure (Given/When/Then)
- **Bug fixes require unit tests** that reproduce the issue before fixing it
- **Never test private methods directly** - test through public interfaces with sufficient coverage
- The CI pipeline uses rios0rios0/pipelines project and runs golangci-lint, SAST security scanning (horusec/semgrep), and gitleaks
- Makefile automatically clones pipelines project using HTTPS (no SSH keys required)

### Environment Configuration
Terra requires environment variables for cloud provider configuration. Create a `.env` file:

```bash
# Required for AWS
TERRA_CLOUD=aws
TERRA_AWS_ROLE_ARN=arn:aws:iam::123456789012:role/terraform-role

# Required for Azure
TERRA_CLOUD=azure
TERRA_AZURE_SUBSCRIPTION_ID=12345678-1234-1234-1234-123456789012

# Optional
TERRA_WORKSPACE=dev

# Optional: Centralized cache directories (defaults shown below)
# TERRA_MODULE_CACHE_DIR=~/.cache/terra/modules
# TERRA_PROVIDER_CACHE_DIR=~/.cache/terra/providers

# Optional: Disable Terragrunt CAS (Content Addressable Store) experiment
# TERRA_NO_CAS=true

# Optional: Disable Terragrunt Partial Parse Config Cache
# TERRA_NO_PARTIAL_PARSE_CACHE=true

# Terraform variables (optional)
TF_VAR_environment=development
TF_VAR_region=us-west-2
```

**Note**: TERRA_CLOUD is optional. When set, it must be either "aws" or "azure". Cloud-specific features (account/subscription switching) require it to be set, but terra will run without it for commands like `clear`, `format`, `install`, `version`, and `self-update`.

## Pipelines Integration

This repository uses the [rios0rios0/pipelines](https://github.com/rios0rios0/pipelines) project for standardized CI/CD operations:

- **CI Pipeline**: Uses `rios0rios0/pipelines/.github/workflows/go-binary.yaml@main`
- **Local Development**: Use Makefile targets which automatically handle pipelines integration
- **Linting**: Run `make lint` instead of manual script execution
- **Testing**: Run `make test` to execute all tests with coverage reporting
- **Security Scanning**: Run `make sast` for security analysis
- **All Checks**: Run `make all` to execute lint + sast + test

### Automatic Pipelines Setup
The Makefile automatically handles pipelines project setup using HTTPS (no SSH configuration required):
```bash
# All these targets automatically clone/update pipelines project as needed
make lint    # Linting via pipelines/global/scripts/golangci-lint/run.sh
make test    # Testing via pipelines/global/scripts/golang/test/run.sh
make sast    # Security scanning via pipelines project
make all     # All quality checks
```

## Available Commands

### Standalone Commands (work without Terraform/Terragrunt installed)
```bash
# Clear local cache directories (.terraform, .terragrunt-cache)
terra clear

# Clear local AND centralized cache directories (~/.cache/terra/modules, ~/.cache/terra/providers)
terra clear --global

# Format code files (warns if terraform/terragrunt not in PATH)
terra format

# Install terraform and terragrunt dependencies
terra install

# Update/install terraform and terragrunt to latest versions (alias for install)
terra update

# Show terra, terraform, and terragrunt versions
terra version

# Also accessible via flags:
terra --version  # or terra -v

# Self-update terra to the latest version from GitHub releases
terra self-update
terra self-update --dry-run   # Show what would be updated without performing it
terra self-update --force     # Skip confirmation prompts
```

### Main Terra Commands (require terraform/terragrunt)
```bash
# Apply all subdirectories in path
terra apply --all /path/to/infrastructure

# Plan all subdirectories in path
terra plan --all /path/to/infrastructure

# Plan specific module
terra plan /path/to/infrastructure/module

# Apply specific module
terra apply /path/to/infrastructure/module

# Execute command in parallel across N modules (default: 5)
terra plan --parallel=5 /path/to/infrastructure

# Select specific modules in parallel execution (comma-separated subdirectory names)
terra apply --parallel=3 --only=module1,module2 /path/to/infrastructure

# Exclude specific modules from parallel execution
terra apply --parallel=3 --skip=excluded_module /path/to/infrastructure

# Terragrunt-managed run-all: filter with terragrunt's --filter (preferred, strictly more expressive)
terra apply --all --filter='!excluded_module' /path/to/infrastructure

# Terragrunt-managed run-all: legacy filter syntax
terra apply --all --queue-exclude-dir=excluded_module /path/to/infrastructure

# Auto-reply to interactive prompts (defaults to "n")
terra apply --reply=y /path/to/infrastructure/module
terra apply -r=y /path/to/infrastructure/module
```

## Validation and Testing

### Testing Guidelines (See CONTRIBUTING.md for detailed requirements)

**CRITICAL Testing Requirements:**
- **Bug fixes MUST include unit tests** that reproduce the issue before fixing it
- **All new tests MUST use testify framework** (`github.com/stretchr/testify`)
- **Follow test method organization pattern**: One function per public method with t.Run() for test cases
- **Never test private methods directly** - test through public interfaces with comprehensive coverage
- **Use testify assertions**: `assert.*` for non-critical, `require.*` for critical, `mock.*` for test doubles
- **ALL test files MUST include build tags** for proper categorization and test execution

**Build Tag Requirements:**
```go
// Unit tests (files in internal/ directories)
//go:build unit

// Integration tests (*_integration_test.go, BDD examples)
//go:build integration

// Test utilities (files in /test folder)
//go:build integration || unit || test
```

**Running tests by category:**
```bash
go test -tags unit ./...        # Unit tests only
go test -tags integration ./... # Integration tests only
go test ./...                   # All tests (default)
```

**Required Test Structure Pattern:**
Each test file must organize tests by grouping them around public methods:

```go
func TestStructName_MethodBeingTested(t *testing.T) {
    t.Parallel() // Use when no environment variables

    t.Run("should return error when invalid input provided", func(t *testing.T) {
        // GIVEN: Setup test data and mocks
        invalidInput := "bad-data"
        service := NewService()

        // WHEN: Execute the action being tested
        result, err := service.Process(invalidInput)

        // THEN: Assert expected outcomes
        assert.Error(t, err)
        assert.Nil(t, result)
        assert.Contains(t, err.Error(), "invalid input")
    })

    t.Run("should succeed when valid input provided", func(t *testing.T) {
        // Additional test case for same method
    })
}
```

**Test Naming Conventions:**
- **Test Methods**: `Test[StructName]_[MethodName]` (e.g., `TestFormatFilesCommand_Execute`)
- **Test Cases**: `"should [behavior] when [condition]"` (e.g., `"should return error when invalid input provided"`)

**Handling Conflicting Method Names:**
When multiple test files test the same method, use descriptive suffixes to avoid naming conflicts:
- `TestDependency_GetBinaryURL` (main functionality)
- `TestDependency_GetBinaryURL_AndroidPlatform` (Android-specific tests)
- `TestDependency_GetBinaryURL_BDDExamples` (BDD example tests)
- `TestInstallDependenciesCommand_Execute_Integration` (integration tests)

**Suffix Guidelines:** Use clear, concise suffixes that describe the test file's focus (_AndroidPlatform, _Integration, _BDDExamples, _EdgeCases, _Performance)

**Parallel Testing Rules:**
- Use `t.Parallel()` when tests don't use `t.Setenv()` or modify global state
- Avoid `t.Parallel()` when using environment variables or shared resources
- **NEVER use `t.Parallel()` with `t.Chdir()`** - This causes runtime panic: "testing: test using t.Setenv or t.Chdir can not use t.Parallel"

**CRITICAL Test Helper Rules:**
- **Test helpers MUST be placed in domain-specific subfolders within `/test` folder** - NEVER in production folders (internal/, cmd/, pkg/)
- **Organization by domain**:
  - `/test/domain/entitydoubles/` - Stubs implementing entity interfaces (CLI, OS, etc.)
  - `/test/domain/entitybuilders/` - Builders that create domain entities
  - `/test/domain/commanddoubles/` - Stubs implementing command interfaces
  - `/test/infrastructure/repositorydoubles/` - Stubs implementing repository interfaces (infrastructure layer)
  - `/test/infrastructure/repositorybuilders/` - Builders for infrastructure testing (HTTP servers, etc.)
  - `/test/infrastructure/repositoryhelpers/` - Helpers for testing repository/OS functionality
  - `/test/infrastructure/controllerdoubles/` - Stubs implementing controller interfaces
  - `/test/infrastructure/controllerhelpers/` - Helpers for controller testing
- **Test helpers in production folders affect coverage unnecessarily** and violate project standards
- **Use `t.Helper()` in all helper functions** for better test failure reporting
- **Name helpers with `Helper` prefix** - avoid `Test` prefix to prevent Go test runner conflicts
- **Follow "one per file" rule** - Each builder, mock, stub, in-memory implementation, dummy, or helper must be in its own file

**Test Utilities Organization Rules:**
- **One utility per file** - Never combine multiple builders, stubs, mocks, or helpers in a single file
- **Domain-specific organization** - All test utilities must be organized by their corresponding production packages
- **Package naming** - Use descriptive package names that reflect the organization:
  - `entitydoubles`, `entitybuilders`
  - `commanddoubles`
  - `repositorydoubles`, `repositorybuilders`, `repositoryhelpers`
  - `controllerdoubles`, `controllerhelpers`
- **Clear naming convention** - Use descriptive names that indicate the utility type and purpose:
  - Builders: `dependency_builder.go`, `test_server_builder.go`
  - Stubs: `stub_shell_repository.go`, `stub_install_dependencies.go`, `stub_cli.go`
  - Mocks: `mock_api_client.go` (when using behavioral verification with testify/mock)
  - In-memory implementations: `inmemory_cache.go`, `inmemory_storage.go`
  - Dummies: `dummy_config.go`, `dummy_logger.go`
  - Helpers: `os_helpers.go`, `network_helpers.go`

**Choosing Between Stubs and Mocks:**

**Use Stubs when:**
- Testing **state verification** (final output/result)
- Controlling dependency return values for different test scenarios
- Recording calls for later assertion (call count, parameters)
- Testing query operations (data retrieval without side effects)

**Use Mocks when:**
- Testing **behavior verification** (specific method calls with expected parameters)
- Ensuring interactions happen in correct order
- Testing command operations (operations with side effects)
- Wanting tests to fail immediately on unexpected interactions
- Using testify/mock package for behavioral expectations

**Example Test Helper Structure:**
```go
// File: /test/infrastructure/repositoryhelpers/os_helpers.go (Helpers in separate files)
package repositoryhelpers

import (
    "testing"
    "github.com/rios0rios0/terra/internal/domain/entities"
)

// HelperDownloadSuccess tests successful download scenarios
func HelperDownloadSuccess(t *testing.T, osImpl entities.OS, testPrefix string) {
    t.Helper() // Mark as test helper
    // ... implementation
}
```

**Example Builder Structure:**
```go
// File: /test/domain/entitybuilders/dependency_builder.go (Domain entity builders)
package entitybuilders

import "github.com/rios0rios0/terra/internal/domain/entities"

// DependencyBuilder helps create test dependencies with a fluent interface
type DependencyBuilder struct { /* ... */ }

func NewDependencyBuilder() *DependencyBuilder { /* ... */ }
func (b *DependencyBuilder) WithName(name string) *DependencyBuilder { /* ... */ }
func (b *DependencyBuilder) Build() entities.Dependency { /* ... */ }
```

**Example Infrastructure Builder Structure:**
```go
// File: /test/infrastructure/repositorybuilders/test_server_builder.go (Infrastructure builders)
package repositorybuilders

import "net/http/httptest"

// TestServerBuilder helps create mock HTTP servers for testing infrastructure
type TestServerBuilder struct { /* ... */ }

func NewTestServerBuilder() *TestServerBuilder { /* ... */ }
func (b *TestServerBuilder) WithVersionResponse(path, response string) *TestServerBuilder { /* ... */ }
func (b *TestServerBuilder) BuildServers() (*httptest.Server, *httptest.Server) { /* ... */ }
```

**Example Stub Structure:**
```go
// File: /test/infrastructure/repositorydoubles/stub_shell_repository.go (Stubs in separate files)
package repositorydoubles

// StubShellRepository for testing shell-related commands
type StubShellRepository struct { /* ... */ }

func (m *StubShellRepository) ExecuteCommand(/* ... */) error { /* ... */ }
```

### Manual Validation Requirements
Always test terra functionality after making changes:

1. **Build Validation**:
   ```bash
   export PATH=$PATH:~/go/bin
   make build
   ./bin/terra clear  # Should succeed
   ```

2. **Code Quality Validation**:
   ```bash
   make test    # Run all tests with coverage
   make lint    # Run linting checks
   make sast    # Run security scanning
   make all     # Run all quality checks
   ```

3. **Environment Validation**:
   ```bash
   # Create test .env file
   echo "TERRA_CLOUD=aws" > .env
   ./bin/terra clear  # Should work without warnings
   ```

4. **Format Command Validation**:
   ```bash
   ./bin/terra format  # Should run (may warn about missing terraform/terragrunt)
   ```

5. **Directory Handling Validation**:
   ```bash
   mkdir -p /tmp/test-terraform
   echo 'resource "null_resource" "test" {}' > /tmp/test-terraform/main.tf
   # This will fail without terraform installed, but tests argument parsing
   ./bin/terra plan /tmp/test-terraform
   ```

### Known Limitations and Issues
- **Network Restrictions**: `terra install` and `terra self-update` fail in environments with restricted internet access due to HashiCorp API calls and GitHub API calls respectively
- **Dependencies**: terraform and terragrunt must be manually installed if `terra install` fails
- **Validation Requirements**: Application requires TERRA_CLOUD to be set to "aws" or "azure" for cloud-specific features; commands that don't need cloud access (clear, format, install, version, self-update) work without it

## Project Structure

### Key Directories
```
cmd/terra/               # Main application entry point (main.go, dig.go)
internal/                # Application bootstrap (app.go, container.go)
internal/domain/         # Business logic: commands, entities, repositories (interfaces)
internal/infrastructure/ # Controllers and repository implementations
test/                    # Test helpers organized by domain/infrastructure layer
```

### Important Files
- `Makefile` - Build and install targets
- `CONTRIBUTING.md` - Comprehensive contributing guidelines including mandatory testing requirements
- `.golangci.yml` - Linting configuration
- `go.mod` - Go module dependencies (includes testify, cobra, and testkit for testing)
- `internal/domain/entities/settings.go` - Environment variable configuration
- `internal/domain/entities/app_context.go` - AppContext interface for DIG container
- `internal/domain/entities/controller.go` - Controller interface for all CLI controllers (uses cobra)
- `internal/domain/entities/controller_bind.go` - ControllerBind struct for cobra command bindings
- `internal/domain/entities/platform.go` - Cross-platform OS/arch detection utilities
- `internal/domain/commands/run_from_root_command.go` - Main command orchestration (caching, execution)
- `internal/domain/commands/run_additional_before_command.go` - Pre-execution setup (account switching, workspace selection)
- `internal/domain/commands/parallel_state_command.go` - Parallel execution of terragrunt commands across modules
- `internal/domain/commands/state_utils.go` - Flag utilities for state manipulation and parallel flags
- `internal/domain/commands/version_command.go` - Version display command
- `internal/domain/commands/self_update_command.go` - Self-update from GitHub releases
- `internal/infrastructure/controllers/helpers/arguments_helper.go` - Command argument parsing

### DIG Dependency Injection and Cobra CLI
- Uses Uber's DIG for runtime dependency injection
- Uses `github.com/spf13/cobra` for CLI command management
- Uses `github.com/rios0rios0/testkit` as a shared test utilities library
- Providers are registered in container.go files in each layer
- No code generation required

### Concurrency and Caching
- **Centralized module cache**: Terra sets `TG_DOWNLOAD_DIR` before invoking Terragrunt so all stacks share a single module download directory (default `~/.cache/terra/modules`). Override with `TERRA_MODULE_CACHE_DIR`.
- **Centralized provider cache**: Terra sets `TG_PROVIDER_CACHE_DIR` so provider plugins are downloaded once and reused (default `~/.cache/terra/providers`). Override with `TERRA_PROVIDER_CACHE_DIR`.
- **CAS (Content Addressable Store)**: Terra enables the Terragrunt CAS experiment by default (`TG_EXPERIMENT=cas`), which deduplicates Git clones via hard links. This reduces disk usage and speeds up subsequent clones. Disable with `TERRA_NO_CAS=true`.
- **Provider caching**: Terra uses the Terragrunt Provider Cache Server (`TG_PROVIDER_CACHE=1`) for concurrent-safe provider deduplication with file locking. This replaced `TF_PLUGIN_CACHE_DIR` which caused "text file busy" errors during parallel execution. Disable with `TERRA_NO_PROVIDER_CACHE=true`. Terra also sets `TG_NO_AUTO_PROVIDER_CACHE_DIR=true` whenever the Provider Cache Server is enabled; without this, Terragrunt 0.99+'s CAS experiment silently auto-enables `auto-provider-cache-dir`, which overrides `TG_PROVIDER_CACHE_DIR` and duplicates providers per go-getter source inside `TG_DOWNLOAD_DIR`. The two flags must stay paired in `configureCacheEnvironment`.
- **Partial Parse Config Cache**: Terra enables the Terragrunt Partial Parse Config Cache by default (`TG_USE_PARTIAL_PARSE_CONFIG_CACHE=true`), which caches parsed HCL configs across modules sharing the same root include. Disable with `TERRA_NO_PARTIAL_PARSE_CACHE=true`.
- **Auto-initialization with upgrade**: `UpgradeAwareShellRepository` wraps command execution. When a terragrunt command fails with output matching upgrade-needed patterns (backend changed, provider conflicts, uninitialized modules), it automatically runs `init --upgrade` and retries the original command. This is used in the normal (non-interactive) execution path of `RunFromRootCommand`.
- **Parallel execution**: Two independent strategies exist. **Terra-managed**: `--parallel=N` runs terragrunt across multiple modules using N goroutine workers; use `--only=mod1,mod2` to select modules or `--skip=mod3` to exclude. **Terragrunt-managed**: `--all`, `--parallelism=N`, and `--filter=query` (and the legacy `--queue-exclude-dir`/`--queue-include-dir`) are forwarded directly to terragrunt for its native run-all behavior. These two strategies cannot be combined (`--parallel` and `--all` together is an error). Terra's `--only`/`--skip` only work with `--parallel=N`; on the `--all` path you must use terragrunt's own filter flags (prefer `--filter='!mod'` which is strictly more expressive than `--queue-exclude-dir`).
- **Educational validation errors**: When a user passes `--only`/`--skip` without `--parallel`, terra fatalfs with a multi-line error that echoes the command they typed and prints both valid escape hatches (`--parallel=5 --skip=mod` AND `--all --filter='!mod'`) as copy-pasteable examples. Same treatment for the `--parallel` + `--all` conflict. When a user passes terragrunt-only queue/filter flags (`--filter`, `--queue-exclude-dir`, `--queue-include-dir`) alongside `--parallel=N`, terra logs a non-fatal warning because the flags are silently ignored by the worker pool. These error builders live in `internal/domain/commands/run_from_root_error_builders.go` as `BuildSelectionFlagsError` and `BuildParallelAllConflictError`; update them (and the unit tests in `run_from_root_error_builders_test.go`) when changing validation messages.
- **Reply mode**: Use `--reply=<value>` (or `-r=<value>`) to automatically answer interactive prompts from terragrunt. Uses `creack/pty` for PTY-based interaction.
- **Pre-execution steps**: `RunAdditionalBeforeCommand` runs before the main terragrunt command: switches cloud account (if `TERRA_CLOUD` is set), initializes terraform (if not already done), and selects the workspace (if `TERRA_WORKSPACE` is set).
- **Clearing caches**: `terra clear` removes local `.terraform` and `.terragrunt-cache` directories. Use `terra clear --global` to also remove the centralized cache directories.

## Common Tasks

### Adding New Commands
1. Create command in `internal/domain/commands/`
2. Create controller in `internal/infrastructure/controllers/`
3. Add providers to respective container.go files
4. Rebuild and test

### Debugging Build Issues
- Run `go mod tidy` after dependency changes
- Check container.go files have proper DIG provider registration

### Environment Variables Reference
```bash
# Cloud provider (optional - required for account/subscription switching features)
TERRA_CLOUD=aws|azure

# AWS specific (required if TERRA_CLOUD=aws and role switching needed)
TERRA_AWS_ROLE_ARN=arn:aws:iam::account:role/name

# Azure specific (required if TERRA_CLOUD=azure and subscription switching needed)
TERRA_AZURE_SUBSCRIPTION_ID=subscription-id

# Terraform workspace (optional)
TERRA_WORKSPACE=workspace-name

# Centralized module cache directory (optional, default: ~/.cache/terra/modules)
# Sets TG_DOWNLOAD_DIR so Terragrunt reuses downloaded modules across stacks
TERRA_MODULE_CACHE_DIR=/custom/path/to/modules

# Centralized provider cache directory (optional, default: ~/.cache/terra/providers)
# Sets TG_PROVIDER_CACHE_DIR so the Provider Cache Server reuses provider plugins across stacks
TERRA_PROVIDER_CACHE_DIR=/custom/path/to/providers

# Disable CAS experiment (optional, default: false = CAS enabled)
TERRA_NO_CAS=true

# Disable Provider Cache Server (optional, default: false = Provider Cache enabled)
TERRA_NO_PROVIDER_CACHE=true

# Disable Partial Parse Config Cache (optional, default: false = Partial Parse Cache enabled)
TERRA_NO_PARTIAL_PARSE_CACHE=true

# Terraform variables (optional, any TF_VAR_* variables)
TF_VAR_*=value
```

## CRITICAL Build and Timing Information
- **Build Time**: 15-20 seconds typical, NEVER CANCEL builds
- **Linting Time**: 2-5 minutes with `make lint`, NEVER CANCEL
- **Testing Time**: 1-2 minutes with `make test`, includes coverage reporting
- **Dependencies**: DIG-based dependency injection with cobra CLI, no code generation needed
- **Pipelines**: Use Makefile targets which automatically handle pipelines project via HTTPS
- **Install Failures**: `terra install` and `terra self-update` will fail in restricted network environments - this is expected behavior

Always validate changes by building, testing, and running the basic terra commands to ensure functionality is preserved.

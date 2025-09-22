# Terra - Terraform/Terragrunt CLI Wrapper

Terra is a Go CLI application that wraps Terraform and Terragrunt functionality, providing simplified path-based infrastructure management inspired by Kubernetes.

Always reference these instructions first and fallback to search or bash commands only when you encounter unexpected information that does not match the info here.

## Working Effectively

### Prerequisites and Environment Setup
- Ensure Go 1.23+ is installed and `go` is in PATH
- Add `~/go/bin` to PATH for wire tool: `export PATH=$PATH:~/go/bin`
- For linting and CI tools, use the pipelines project (https://github.com/rios0rios0/pipelines)
- NEVER CANCEL: Build takes 15-20 seconds. NEVER CANCEL. Set timeout to 60+ minutes for safety.

### Building and Installing
Bootstrap and build the repository:
```bash
# Build terra binary
export PATH=$PATH:~/go/bin
make build
```

Install system-wide (requires sudo):
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
make horusec # Runs security scanning via pipelines project
make all     # Runs lint + horusec + test
```

**File Standards:**
- **Use LF line endings** - All new and edited files must use LF (Unix) line endings, not CRLF (Windows)
- **Update CHANGELOG.md** - Add entries to the `[Unreleased]` section for new features and bug fixes (not required for documentation-only changes)
- The `.editorconfig` file enforces line ending standards for most editors

- **Tests exist** in this repository with good coverage across domain and infrastructure layers
- **All new tests must use testify framework** and follow BDD structure (Given/When/Then)
- **Bug fixes require unit tests** that reproduce the issue before fixing it
- **Never test private methods directly** - test through public interfaces with sufficient coverage
- The CI pipeline uses rios0rios0/pipelines project and runs golangci-lint, horusec security scanning, semgrep, and gitleaks
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

# Terraform variables (optional)
TF_VAR_environment=development
TF_VAR_region=us-west-2
```

**CRITICAL**: TERRA_CLOUD must be either "aws" or "azure" - the application will fail validation if empty or invalid.

## Pipelines Integration

This repository uses the [rios0rios0/pipelines](https://github.com/rios0rios0/pipelines) project for standardized CI/CD operations:

- **CI Pipeline**: Uses `rios0rios0/pipelines/.github/workflows/go-binary.yaml@main`
- **Local Development**: Use Makefile targets which automatically handle pipelines integration
- **Linting**: Run `make lint` instead of manual script execution  
- **Testing**: Run `make test` to execute all tests with coverage reporting
- **Security Scanning**: Run `make horusec` for security analysis
- **All Checks**: Run `make all` to execute lint + horusec + test

### Automatic Pipelines Setup
The Makefile automatically handles pipelines project setup using HTTPS (no SSH configuration required):
```bash
# All these targets automatically clone/update pipelines project as needed
make lint    # Linting via pipelines/global/scripts/golangci-lint/run.sh
make test    # Testing via pipelines/global/scripts/golang/test/run.sh  
make horusec # Security via pipelines/global/scripts/horusec/run.sh
make all     # All quality checks
```

## Available Commands

### Standalone Commands (work without Terraform/Terragrunt installed)
```bash
# Clear cache directories
terra clear

# Format code files (warns if terraform/terragrunt not in PATH)
terra format

# Install terraform and terragrunt dependencies
terra install
```

### Main Terra Commands (require terraform/terragrunt)
```bash
# Apply all subdirectories in path
terra run-all apply /path/to/infrastructure

# Plan all subdirectories in path  
terra run-all plan /path/to/infrastructure

# Plan specific module
terra plan /path/to/infrastructure/module

# Apply specific module
terra apply /path/to/infrastructure/module
```

## Validation and Testing

### Testing Guidelines (See CONTRIBUTING.md for detailed requirements)

**CRITICAL Testing Requirements:**
- **Bug fixes MUST include unit tests** that reproduce the issue before fixing it
- **All new tests MUST use testify framework** (`github.com/stretchr/testify`)
- **Follow test method organization pattern**: One function per public method with t.Run() for test cases
- **Never test private methods directly** - test through public interfaces with comprehensive coverage
- **Use testify assertions**: `assert.*` for non-critical, `require.*` for critical, `mock.*` for test doubles

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
  - `/test/domain/entity_doubles/` - Stubs implementing entity interfaces (CLI, OS, etc.)
  - `/test/domain/entity_builders/` - Builders that create domain entities
  - `/test/domain/command_doubles/` - Stubs implementing command interfaces
  - `/test/infrastructure/repository_doubles/` - Stubs implementing repository interfaces (infrastructure layer)
  - `/test/infrastructure/repository_builders/` - Builders for infrastructure testing (HTTP servers, etc.)
  - `/test/infrastructure/repository_helpers/` - Helpers for testing repository/OS functionality
  - `/test/infrastructure/controller_doubles/` - Stubs implementing controller interfaces
  - `/test/infrastructure/controller_helpers/` - Helpers for controller testing
- **Test helpers in production folders affect coverage unnecessarily** and violate project standards
- **Use `t.Helper()` in all helper functions** for better test failure reporting
- **Name helpers with `Helper` prefix** - avoid `Test` prefix to prevent Go test runner conflicts
- **Follow "one per file" rule** - Each builder, mock, stub, in-memory implementation, dummy, or helper must be in its own file

**Test Utilities Organization Rules:**
- **One utility per file** - Never combine multiple builders, stubs, mocks, or helpers in a single file
- **Domain-specific organization** - All test utilities must be organized by their corresponding production packages
- **Package naming** - Use descriptive package names that reflect the organization:
  - `entity_doubles`, `entity_builders`
  - `command_doubles`
  - `repository_doubles`, `repository_builders`, `repository_helpers`
  - `controller_doubles`, `controller_helpers`
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
// File: /test/infrastructure/repository_helpers/os_helpers.go (Helpers in separate files)
package repository_helpers

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
// File: /test/domain/entity_builders/dependency_builder.go (Domain entity builders)
package entity_builders

import "github.com/rios0rios0/terra/internal/domain/entities"

// DependencyBuilder helps create test dependencies with a fluent interface
type DependencyBuilder struct { /* ... */ }

func NewDependencyBuilder() *DependencyBuilder { /* ... */ }
func (b *DependencyBuilder) WithName(name string) *DependencyBuilder { /* ... */ }
func (b *DependencyBuilder) Build() entities.Dependency { /* ... */ }
```

**Example Infrastructure Builder Structure:**
```go
// File: /test/infrastructure/repository_builders/test_server_builder.go (Infrastructure builders)
package repository_builders

import "net/http/httptest"

// TestServerBuilder helps create mock HTTP servers for testing infrastructure
type TestServerBuilder struct { /* ... */ }

func NewTestServerBuilder() *TestServerBuilder { /* ... */ }
func (b *TestServerBuilder) WithVersionResponse(path, response string) *TestServerBuilder { /* ... */ }
func (b *TestServerBuilder) BuildServers() (*httptest.Server, *httptest.Server) { /* ... */ }
```

**Example Stub Structure:**
```go
// File: /test/infrastructure/repository_doubles/stub_shell_repository.go (Stubs in separate files)
package repository_doubles

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
   make horusec # Run security scanning
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
- **Network Restrictions**: `terra install` fails in environments with restricted internet access due to HashiCorp API calls
- **Dependencies**: terraform and terragrunt must be manually installed if `terra install` fails
- **Argument Parsing Bug**: Commands like `--help` cause runtime panic due to slice bounds error in argument parser
- **Validation Requirements**: Application requires TERRA_CLOUD to be set to "aws" or "azure" - cannot run without this

## Project Structure

### Key Directories
```
cmd/terra/           # Main application entry point
internal/domain/     # Business logic and entities  
internal/infrastructure/ # Controllers and repositories
```

### Important Files
- `Makefile` - Build and install targets
- `CONTRIBUTING.md` - Comprehensive contributing guidelines including mandatory testing requirements
- `.golangci.yml` - Linting configuration
- `go.mod` - Go module dependencies (includes testify for testing)
- `internal/domain/entities/settings.go` - Environment variable configuration
- `internal/infrastructure/helpers/arguments_helper.go` - Command argument parsing (has known bugs)

### Wire Dependency Injection
- Uses Google Wire for dependency injection
- `wire ./...` generates wire_gen.go files
- Requires wire tool in PATH: `export PATH=$PATH:~/go/bin`

## Common Tasks

### Adding New Commands
1. Create command in `internal/domain/commands/`
2. Create controller in `internal/infrastructure/controllers/`
3. Register in `internal/infrastructure/controllers/container.go`
4. Run `wire ./...` to regenerate dependencies
5. Rebuild and test

### Debugging Build Issues
- Ensure PATH includes `~/go/bin` for wire tool
- Run `go mod tidy` after dependency changes
- Check wire_gen.go files are properly generated

### Environment Variables Reference
```bash
# Cloud provider (required)
TERRA_CLOUD=aws|azure

# AWS specific (required if TERRA_CLOUD=aws and role switching needed)
TERRA_AWS_ROLE_ARN=arn:aws:iam::account:role/name

# Azure specific (required if TERRA_CLOUD=azure and subscription switching needed)  
TERRA_AZURE_SUBSCRIPTION_ID=subscription-id

# Terraform workspace (optional)
TERRA_WORKSPACE=workspace-name

# Terraform variables (optional, any TF_VAR_* variables)
TF_VAR_*=value
```

## CRITICAL Build and Timing Information
- **Build Time**: 15-20 seconds typical, NEVER CANCEL builds
- **Linting Time**: 2-5 minutes with `make lint`, NEVER CANCEL
- **Testing Time**: 1-2 minutes with `make test`, includes coverage reporting
- **Dependencies**: Requires wire tool in PATH for successful builds
- **Pipelines**: Use Makefile targets which automatically handle pipelines project via HTTPS
- **Install Failures**: `terra install` will fail in restricted network environments - this is expected behavior

Always validate changes by building, testing, and running the basic terra commands to ensure functionality is preserved.
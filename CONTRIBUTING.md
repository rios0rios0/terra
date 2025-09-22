# Contributing to Terra

Thank you for your interest in contributing to Terra! We welcome contributions of all kinds, including bug reports, feature requests, documentation improvements, and code contributions.

## Code of Conduct

Please be respectful and constructive when participating in discussions and contributing to the project.

## How to Contribute

### Reporting Issues

When reporting bugs or requesting features:

1. Check existing issues to avoid duplicates
2. Use clear, descriptive titles
3. Provide detailed information about the issue
4. Include steps to reproduce for bugs
5. Specify your environment (OS, Go version, etc.)

### Contributing Code

1. Fork the repository
2. Create a feature branch from `main`
3. Make your changes following our guidelines
4. Ensure all tests pass
5. Submit a pull request

## Testing Guidelines

Testing is a critical part of maintaining code quality in Terra. All contributors must follow these testing requirements:

### Bug Fixes Require Unit Tests

**When fixing a bug, you MUST:**

1. **Create unit tests that reproduce the bug** - Write failing tests that demonstrate the issue before fixing it
2. **Fix the bug** - Implement the minimal change needed to resolve the issue
3. **Verify tests pass** - Ensure your fix makes the new tests pass
4. **Cover edge cases** - Add additional tests for related scenarios that could cause similar issues

### Testing Framework Requirements

**All new tests MUST use the testify framework:**

```go
import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "github.com/stretchr/testify/mock"
)
```

Use:
- `assert.*` for non-critical assertions that allow tests to continue
- `require.*` for critical assertions that should stop test execution
- `mock.*` for creating test doubles and mocks

### Build Tag Requirements

**All test files MUST include appropriate build tags for categorization:**

#### Unit Tests
Unit tests (files in `internal/` directories) must use the `unit` build tag:

```go
//go:build unit

package package_name_test
```

#### Integration Tests  
Integration tests (files with `*_integration_test.go` pattern or BDD examples) must use the `integration` build tag:

```go
//go:build integration

package package_name_test
```

#### Test Utilities
All test helper files in the `/test` folder must use the inclusive build tag to ensure availability during both unit and integration testing:

```go
//go:build integration || unit || test

package package_name
```

#### Running Tests by Category
Use build tags to run specific test categories:

```bash
# Run only unit tests
go test -tags unit ./...

# Run only integration tests  
go test -tags integration ./...

# Run all tests (default behavior)
go test ./...
```

### Test Helper Organization

**Test helpers and utilities MUST be located in the `/test` folder at the root of the project:**

**CRITICAL RULE: Never create test helper files in production folders (internal/, cmd/, pkg/) as they unnecessarily affect code coverage calculations.**

#### Proper Test Helper Structure

```go
// ✅ DO: Place test helpers in /test folder
// File: /test/my_helpers.go
package test

import (
    "testing"
    "github.com/rios0rios0/terra/internal/domain/entities"
)

// HelperFunctionName provides testing utilities for X functionality
func HelperFunctionName(t *testing.T, param entities.SomeType) {
    t.Helper() // Mark as test helper
    // ... helper implementation
}
```

#### Using Test Helpers

```go
// File: internal/domain/entities/some_test.go
package entities_test

import (
    "testing"
    "github.com/rios0rios0/terra/internal/domain/entities"
    "github.com/rios0rios0/terra/test" // Import test helpers
)

func TestComponent_ShouldWork_WhenValidInput(t *testing.T) {
    // GIVEN: Setup using test helper
    component := entities.NewComponent()
    
    // WHEN: Using test helper for consistent testing
    test.HelperFunctionName(t, component)
    
    // THEN: Assertions...
}
```

#### Test Helper Guidelines

1. **All test helpers go in `/test` folder** - Never in production folders
2. **Use `t.Helper()` in helper functions** - This improves test failure reporting
3. **Name helpers with `Helper` prefix** - Avoid `Test` prefix to prevent conflicts with Go test runner
4. **Keep helpers focused** - Each helper should have a single, clear purpose
5. **Document helper purpose** - Include comments explaining what the helper does

### Test Builders, Stubs, Mocks, and Helpers Organization

**CRITICAL RULE: All test utilities (builders, stubs, mocks, in-memory implementations, dummies, and helpers) MUST be organized following the "one per file" rule in domain-specific subfolders within the `/test` folder.**

#### Organization Guidelines

1. **One utility per file** - Never combine multiple builders, stubs, mocks, or helpers in a single file
2. **Domain-specific organization** - All test utilities must be organized by their corresponding production packages:
   - **`test/domain/entity_doubles/`** - Stubs implementing domain entity interfaces (CLI, OS, etc.)
   - **`test/domain/entity_builders/`** - Builders that create domain entities for testing
   - **`test/domain/command_doubles/`** - Stubs implementing domain command interfaces
   - **`test/infrastructure/repository_doubles/`** - Stubs implementing repository interfaces (infrastructure layer)
   - **`test/infrastructure/repository_builders/`** - Builders for infrastructure testing (HTTP servers, etc.)
   - **`test/infrastructure/repository_helpers/`** - Helpers for testing repository/OS functionality
   - **`test/infrastructure/controller_doubles/`** - Stubs implementing infrastructure controller interfaces
   - **`test/infrastructure/controller_helpers/`** - Helpers for testing controller functionality

3. **Clear naming convention** - Use descriptive names that indicate the utility type and purpose:
   - Builders: `dependency_builder.go`, `test_server_builder.go`
   - Stubs: `stub_shell_repository.go`, `stub_install_dependencies.go`, `stub_cli.go`
   - Mocks: `mock_api_client.go` (when using behavioral verification with testify/mock)
   - In-memory implementations: `inmemory_cache.go`, `inmemory_storage.go`
   - Dummies: `dummy_config.go`, `dummy_logger.go`
   - Helpers: `os_helpers.go`, `network_helpers.go`

4. **Package naming** - Use descriptive package names that reflect the organization:
   - `entity_doubles`, `entity_builders`
   - `command_doubles`
   - `repository_doubles`, `repository_builders`, `repository_helpers`
   - `controller_doubles`, `controller_helpers`

#### Test Double Definitions

Following Martin Fowler's definitions:
- **Stubs**: Test doubles that return fixed responses and enable state verification (most test doubles in this project)
- **Mocks**: Test doubles with pre-configured behavioral expectations that fail tests if interactions don't match expectations

#### When to Choose Stubs vs Mocks

**Use Stubs when:**
- You want to verify the **state** of the system after an operation
- You need to control what a dependency returns to test different scenarios
- You want to record calls for later assertion (call count, last parameters, etc.)
- Testing the final output or result of the system under test
- The dependency is a **query** operation (returns data without side effects)

**Use Mocks when:**
- You want to verify **behavior** - that specific methods were called with expected parameters
- You need to ensure interactions happen in a specific order
- The dependency represents a **command** operation (causes side effects)
- Testing that the system properly communicates with external services
- You want the test to fail immediately if unexpected interactions occur

**Examples:**

```go
// ✅ Stub Example: State verification
func TestCalculateTotal_ShouldReturnCorrectSum_WhenValidItemsProvided(t *testing.T) {
    // GIVEN: Stub returns fixed tax rate
    taxStub := &test.StubTaxService{TaxRate: 0.1}
    calculator := NewCalculator(taxStub)
    
    // WHEN: Calculate total
    total := calculator.CalculateTotal(items)
    
    // THEN: Verify final state/result
    assert.Equal(t, 110.0, total)
    assert.Equal(t, 1, taxStub.CallCount) // State verification
}

// ✅ Mock Example: Behavior verification using testify/mock
func TestSendNotification_ShouldCallEmailService_WhenNotificationSent(t *testing.T) {
    // GIVEN: Mock with behavioral expectations
    emailMock := &mocks.EmailService{}
    emailMock.On("SendEmail", "user@example.com", "Hello").Return(nil)
    notifier := NewNotifier(emailMock)
    
    // WHEN: Send notification
    err := notifier.Send("user@example.com", "Hello")
    
    // THEN: Verify behavior occurred as expected
    require.NoError(t, err)
    emailMock.AssertExpectations(t) // Behavior verification
}
```

#### File Structure Examples

```go
// ✅ DO: Separate builder file
// File: /test/dependency_builder.go
package test

import "github.com/rios0rios0/terra/internal/domain/entities"

// DependencyBuilder helps create test dependencies with a fluent interface
type DependencyBuilder struct { /* ... */ }

func NewDependencyBuilder() *DependencyBuilder { /* ... */ }
func (b *DependencyBuilder) WithName(name string) *DependencyBuilder { /* ... */ }
func (b *DependencyBuilder) Build() entities.Dependency { /* ... */ }
```

```go
// ✅ DO: Separate stub file
// File: /test/stub_shell_repository.go
package test

// StubShellRepository for testing shell-related commands
type StubShellRepository struct { /* ... */ }

func (m *StubShellRepository) ExecuteCommand(/* ... */) error { /* ... */ }
```

```go
// ❌ DON'T: Multiple utilities in one file
// File: /test/test_utilities.go
package test

// Multiple builders, stubs, and helpers in same file - WRONG!
type DependencyBuilder struct { /* ... */ }
type StubShellRepository struct { /* ... */ }
type NetworkHelper struct { /* ... */ }
```

#### Benefits of This Organization

1. **Better maintainability** - Easy to locate and modify specific test utilities
2. **Improved readability** - Clear separation of concerns
3. **Enhanced discoverability** - Developers can quickly find the utility they need
4. **Reduced merge conflicts** - Changes to different utilities don't affect each other
5. **Consistent code organization** - Follows established patterns across the project

### BDD (Behavior Driven Design) Test Structure

**All tests MUST follow BDD structure with three distinct sections:**

#### Test Method Naming Convention

Test methods should follow BDD naming using clear, descriptive names:

```go
func TestCommandService_ShouldReturnError_WhenInvalidPathProvided(t *testing.T) {
func TestVersionCommand_ShouldDisplayTerraVersion_WhenCommandExecuted(t *testing.T) {
func TestSettingsValidator_ShouldRejectEmptyCloud_WhenValidationRuns(t *testing.T)
```

Pattern: `Test[ComponentName]_Should[ExpectedBehavior]_When[Condition]`

#### Test Structure: Given, When, Then

Each test MUST be organized into three clear sections:

```go
func TestCommandService_ShouldReturnError_WhenInvalidPathProvided(t *testing.T) {
    // GIVEN: Arrange test data and dependencies
    invalidPath := "/nonexistent/path/12345"
    mockOS := &MockOSInterface{}
    mockOS.On("FileExists", invalidPath).Return(false)
    service := NewCommandService(mockOS)
    
    // WHEN: Execute the action being tested
    result, err := service.ValidatePath(invalidPath)
    
    // THEN: Assert expected outcomes
    assert.Error(t, err)
    assert.Nil(t, result)
    assert.Contains(t, err.Error(), "path does not exist")
    mockOS.AssertExpectations(t)
}
```

**Section Guidelines:**
- **GIVEN**: Set up test data, mocks, and preconditions
- **WHEN**: Execute the specific action/method being tested  
- **THEN**: Verify all expected outcomes using testify assertions

#### Use Comments to Separate Sections

Always use comments to clearly separate the three sections:

```go
func TestExample(t *testing.T) {
    // GIVEN: description of setup
    // ... setup code ...
    
    // WHEN: description of action
    // ... action code ...
    
    // THEN: description of expectations
    // ... assertions ...
}
```

#### Avoid Loops in Test Cases

**NEVER use loops (`for range`) to create test cases inside a test method.** Instead:

1. **Create individual test methods** - Each test scenario should be a separate test function
2. **Use descriptive test names** - Each test should clearly indicate what it's testing
3. **Keep tests independent** - Each test should be able to run in isolation

```go
// ❌ DON'T: Use loops for test cases
func TestValidateInput(t *testing.T) {
    tests := []struct{
        name string
        input string
        expected bool
    }{
        {"valid input", "test", true},
        {"invalid input", "", false},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // test logic
        })
    }
}

// ✅ DO: Create separate test methods
func TestValidateInput_ShouldReturnTrue_WhenValidInputProvided(t *testing.T) {
    // GIVEN: Valid input
    // WHEN: Validation is called
    // THEN: Should return true
}

func TestValidateInput_ShouldReturnFalse_WhenEmptyInputProvided(t *testing.T) {
    // GIVEN: Empty input
    // WHEN: Validation is called  
    // THEN: Should return false
}
```

### Test Method Organization

**All test files MUST organize tests by grouping them around the public methods they test.** This provides better structure and makes it easier to understand test coverage.

#### Required Test Structure Pattern

Each test file should have **one test function per public method** being tested, with each function containing multiple test cases using `t.Run()`:

```go
// For a struct with public methods GetName() and Execute()
func TestMyStruct_GetName(t *testing.T) {
    t.Parallel() // Use when no environment variables
    
    t.Run("should return correct name when called", func(t *testing.T) {
        // GIVEN: A struct instance
        instance := NewMyStruct()
        
        // WHEN: Getting the name
        result := instance.GetName()
        
        // THEN: Should return expected name
        assert.Equal(t, "expected-name", result)
    })
    
    t.Run("should handle empty configuration when no config provided", func(t *testing.T) {
        // Another test case for the same method
    })
}

func TestMyStruct_Execute(t *testing.T) {
    t.Parallel() // Use when no environment variables
    
    t.Run("should complete successfully when valid input provided", func(t *testing.T) {
        // GIVEN: Valid input and struct instance
        instance := NewMyStruct()
        validInput := "test-input"
        
        // WHEN: Executing the method
        err := instance.Execute(validInput)
        
        // THEN: Should complete without error
        assert.NoError(t, err)
    })
    
    t.Run("should return error when invalid input provided", func(t *testing.T) {
        // Another test case for the same method
    })
}
```

#### Naming Convention for Test Methods

Use this pattern for test method names:
- **Pattern**: `Test[StructName]_[MethodName]`
- **Examples**: 
  - `TestFormatFilesCommand_Execute`
  - `TestVersionController_GetBind`
  - `TestDependency_GetBinaryURL`

#### Handling Conflicting Method Names

When multiple test files test the same method, use descriptive suffixes to maintain method-based organization while avoiding naming conflicts:

**Problem**: Multiple test files testing the same method create function name conflicts:
```go
// ❌ CONFLICT: This will cause compilation error
// File: dependency_test.go
func TestDependency_GetBinaryURL(t *testing.T) { ... }

// File: android_fix_test.go  
func TestDependency_GetBinaryURL(t *testing.T) { ... } // Duplicate name!
```

**Solution**: Use descriptive suffixes that identify the specific focus of each test file:
```go
// ✅ DO: Use descriptive suffixes to avoid conflicts
// File: dependency_test.go - General functionality tests
func TestDependency_GetBinaryURL(t *testing.T) { ... }

// File: android_fix_test.go - Android platform-specific tests
func TestDependency_GetBinaryURL_AndroidPlatform(t *testing.T) { ... }

// File: dependency_bdd_example_test.go - BDD example tests
func TestDependency_GetBinaryURL_BDDExamples(t *testing.T) { ... }
```

**Suffix Guidelines**:
- Use suffixes that clearly describe the test file's focus
- Keep suffixes concise but descriptive
- Examples of good suffixes:
  - `_AndroidPlatform` for Android-specific functionality
  - `_Integration` for integration tests
  - `_BDDExamples` for demonstration/example tests
  - `_EdgeCases` for edge case testing
  - `_Performance` for performance testing

#### Naming Convention for Test Cases

Use descriptive names that follow BDD pattern:
- **Pattern**: `"should [expected behavior] when [condition]"`
- **Examples**:
  - `"should return error when invalid path provided"`
  - `"should execute successfully when valid dependencies provided"`
  - `"should create instance when called with valid parameters"`

#### Parallel Testing Guidelines

1. **Use `t.Parallel()`** when tests don't modify environment variables or shared state
2. **Avoid `t.Parallel()`** when using `t.Setenv()` or modifying global state
3. **NEVER use `t.Parallel()` with `t.Chdir()`** - This causes a runtime panic as Go detects incompatible test modifications
4. **Each `t.Run()` test case** can run in parallel by default unless they conflict

**Example of incompatible usage:**
```go
// ❌ DON'T: This will cause runtime panic
func TestSomething(t *testing.T) {
    t.Parallel() // This line will cause panic
    
    t.Run("test case", func(t *testing.T) {
        tempDir := t.TempDir()
        t.Chdir(tempDir) // Incompatible with t.Parallel()
        // ... test logic
    })
}

// ✅ DO: Remove t.Parallel() when using t.Chdir()
func TestSomething(t *testing.T) {
    t.Run("test case", func(t *testing.T) {
        tempDir := t.TempDir()
        t.Chdir(tempDir) // Now safe without t.Parallel()
        // ... test logic
    })
}
```

#### Constructor Testing

For constructors (like `NewMyStruct()`), create a dedicated test function:

```go
func TestNewMyStruct(t *testing.T) {
    t.Parallel()
    
    t.Run("should create instance when valid parameters provided", func(t *testing.T) {
        // GIVEN: Valid constructor parameters
        param1 := "test"
        param2 := 42
        
        // WHEN: Creating instance
        instance := NewMyStruct(param1, param2)
        
        // THEN: Should return valid instance
        require.NotNil(t, instance)
        // Additional assertions about the instance state
    })
}
```

### Testing Private Methods

**NEVER test private methods directly.** Instead:

1. **Test through public interfaces** - Create comprehensive test scenarios for public methods
2. **Use sufficient test coverage** - Write enough test cases to exercise private methods indirectly  
3. **Test edge cases** - Include boundary conditions and error scenarios
4. **Verify behavior, not implementation** - Focus on what the public interface should do

Example of testing private methods indirectly:

```go
// ❌ DON'T: Test private method directly
func TestPrivateMethod(t *testing.T) {
    // This violates our guidelines
}

// ✅ DO: Test private method through public interface
func TestPublicMethod_ShouldHandleInvalidInput_WhenCalledWithBadData(t *testing.T) {
    // GIVEN: Invalid input that will exercise private validation method
    service := NewService()
    invalidInput := "invalid-data"
    
    // WHEN: Call public method that uses private method internally
    result, err := service.ProcessData(invalidInput)
    
    // THEN: Verify behavior demonstrates private method worked correctly
    assert.Error(t, err)
    assert.Nil(t, result)
    // This test covers the private validation method indirectly
}
```

### Test Coverage Requirements

- **Minimum coverage**: Aim for 80%+ code coverage
- **Critical paths**: 100% coverage for error handling and validation logic
- **Integration tests**: Include tests that verify component interactions
- **Edge cases**: Test boundary conditions, null inputs, and error scenarios

### Example: Complete BDD Test

```go
func TestInstallCommand_ShouldDownloadDependency_WhenValidURLProvided(t *testing.T) {
    // GIVEN: Valid dependency configuration and mock HTTP server
    dependency := entities.Dependency{
        Name:       "terraform",
        BinaryURL:  "https://example.com/terraform.zip",
        VersionURL: "https://example.com/version.json",
    }
    
    mockDownloader := &mocks.MockDownloader{}
    mockDownloader.On("Download", dependency.BinaryURL).Return([]byte("binary-content"), nil)
    
    mockOS := &mocks.MockOSInterface{}
    mockOS.On("WriteFile", mock.AnythingOfType("string"), mock.AnythingOfType("[]uint8")).Return(nil)
    mockOS.On("MakeExecutable", mock.AnythingOfType("string")).Return(nil)
    
    command := NewInstallCommand(mockDownloader, mockOS)
    
    // WHEN: Execute install command with valid dependency
    err := command.InstallDependency(dependency)
    
    // THEN: Verify successful installation
    assert.NoError(t, err)
    mockDownloader.AssertExpectations(t)
    mockOS.AssertExpectations(t)
    
    // Verify specific calls were made with expected parameters
    mockDownloader.AssertCalled(t, "Download", dependency.BinaryURL)
    mockOS.AssertCalled(t, "WriteFile", mock.AnythingOfType("string"), []byte("binary-content"))
    mockOS.AssertCalled(t, "MakeExecutable", mock.AnythingOfType("string"))
}
```

## Development Workflow

### Before Starting Development

1. Ensure Go 1.23+ is installed
2. Clone the repository and navigate to the project directory

```bash
make build
```

### Code Quality Checklist

Before submitting a pull request, ensure your code passes all quality checks:

```bash
# Format code
go fmt ./...

# Run linting (takes 2-5 minutes, do not cancel)
make lint

# Run tests with coverage
make test

# Run security scanning
make horusec

# Run all quality checks
make all
```

**File Standards:**
- Ensure all files use LF (Unix) line endings, not CRLF (Windows)
- The `.editorconfig` file enforces this standard for most editors
- Update CHANGELOG.md for new features and bug fixes (not required for documentation-only changes)

### Pull Request Guidelines

1. **Create focused PRs** - One feature or bug fix per pull request
2. **Write clear descriptions** - Explain what changes were made and why
3. **Include tests** - Follow all testing guidelines above
4. **Update documentation** - Update relevant docs if your changes require it
5. **Keep changes minimal** - Make the smallest possible changes to achieve your goal
6. **Use LF line endings** - All new and edited files must use LF (Unix) line endings, not CRLF (Windows)
7. **Update CHANGELOG.md** - Add entries to the `[Unreleased]` section for new features and bug fixes (not required for documentation-only changes)

### Environment Configuration

Terra requires specific environment variables. Create a `.env` file:

```bash
# Required: Cloud provider (must be "aws" or "azure")
TERRA_CLOUD=aws

# AWS specific (if using AWS)
TERRA_AWS_ROLE_ARN=arn:aws:iam::123456789012:role/terraform-role

# Azure specific (if using Azure)  
TERRA_AZURE_SUBSCRIPTION_ID=12345678-1234-1234-1234-123456789012

# Optional
TERRA_WORKSPACE=dev
TF_VAR_environment=development
TF_VAR_region=us-west-2
```

## Project Structure

```
cmd/terra/                    # Main application entry point
internal/domain/             # Business logic and entities
  ├── commands/              # Core business commands
  ├── entities/              # Domain entities and models
  └── repositories/          # Repository interfaces
internal/infrastructure/     # Controllers and repository implementations
  ├── controllers/           # HTTP/CLI controllers
  ├── repositories/          # Concrete repository implementations
  └── helpers/               # Utility functions
```

## Getting Help

- **Issues**: Use GitHub issues for bug reports and feature requests
- **Discussions**: Start a discussion for questions and ideas
- **Documentation**: Check README.md and this CONTRIBUTING.md for guidance

## License

By contributing to Terra, you agree that your contributions will be licensed under the MIT License.

---

Thank you for helping make Terra better!
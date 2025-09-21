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
2. Add `~/go/bin` to your PATH for the wire tool
3. Clone the repository and navigate to the project directory

```bash
export PATH=$PATH:~/go/bin
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

### Pull Request Guidelines

1. **Create focused PRs** - One feature or bug fix per pull request
2. **Write clear descriptions** - Explain what changes were made and why
3. **Include tests** - Follow all testing guidelines above
4. **Update documentation** - Update relevant docs if your changes require it
5. **Keep changes minimal** - Make the smallest possible changes to achieve your goal

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
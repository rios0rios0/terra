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
- Always run code formatting and linting before committing using pipelines project:
```bash
# Standard Go tools (fast)
go fmt ./...
go vet ./...

# Full linting using pipelines project (https://github.com/rios0rios0/pipelines)
# Clone pipelines project if not available locally
if [ ! -d "../pipelines" ]; then
  git clone https://github.com/rios0rios0/pipelines.git ../pipelines
fi
# Run linting using pipelines GoLang scripts
../pipelines/GoLang/golangci-lint.sh
```

- No unit tests exist in this repository (as of current state)
- The CI pipeline uses rios0rios0/pipelines project and runs golangci-lint, horusec security scanning, semgrep, and gitleaks

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
- **Local Development**: Use shell scripts from the pipelines GoLang folder instead of direct tool installation
- **Linting**: Run `../pipelines/GoLang/golangci-lint.sh` instead of installing golangci-lint directly
- **Security Scanning**: Pipeline handles horusec, semgrep, and gitleaks automatically

### Setting up Pipelines Locally
```bash
# Clone pipelines project alongside terra repository
cd ..
git clone https://github.com/rios0rios0/pipelines.git
cd terra

# Now you can use pipelines scripts
../pipelines/GoLang/golangci-lint.sh
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

### Manual Validation Requirements
Always test terra functionality after making changes:

1. **Build Validation**:
   ```bash
   export PATH=$PATH:~/go/bin
   make build
   ./bin/terra clear  # Should succeed
   ```

2. **Environment Validation**:
   ```bash
   # Create test .env file
   echo "TERRA_CLOUD=aws" > .env
   ./bin/terra clear  # Should work without warnings
   ```

3. **Format Command Validation**:
   ```bash
   ./bin/terra format  # Should run (may warn about missing terraform/terragrunt)
   ```

4. **Directory Handling Validation**:
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
- `.golangci.yml` - Linting configuration
- `go.mod` - Go module dependencies
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
- **Linting Time**: 2-5 minutes with pipelines golangci-lint script, NEVER CANCEL
- **No Tests**: Repository contains no unit tests
- **Dependencies**: Requires wire tool in PATH for successful builds
- **Pipelines**: Use rios0rios0/pipelines project scripts instead of direct tool installation
- **Install Failures**: `terra install` will fail in restricted network environments - this is expected behavior

Always validate changes by building and running the basic terra commands to ensure functionality is preserved.
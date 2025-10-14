# Terra CLI - Model Context Protocol (MCP) Documentation

**Version:** 1.0.0  
**Last Updated:** 2025-10-14

This documentation is optimized for AI agents and Model Context Protocol (MCP) systems to understand the Terra CLI capabilities, parameters, and usage patterns.

## Overview

Terra is a CLI wrapper for Terraform and Terragrunt that simplifies infrastructure management through path-based operations. It provides enhanced state management, parallel execution, automatic dependency installation, and cloud provider integration.

**Key Capabilities:**
- Wraps Terraform and Terragrunt commands with enhanced functionality
- Path-based infrastructure management (Kubernetes-style)
- Parallel execution for state manipulation commands
- Automatic dependency installation and updates
- Cloud provider account/subscription switching (AWS, Azure)
- Auto-answering for non-interactive execution
- Self-update mechanism

## Installation

Terra can be installed through multiple methods:

```bash
# Quick install (recommended)
curl -fsSL https://raw.githubusercontent.com/rios0rios0/terra/main/install.sh | sh

# With options
curl -fsSL https://raw.githubusercontent.com/rios0rios0/terra/main/install.sh | sh -s -- --version v1.0.0
curl -fsSL https://raw.githubusercontent.com/rios0rios0/terra/main/install.sh | sh -s -- --install-dir /usr/local/bin
curl -fsSL https://raw.githubusercontent.com/rios0rios0/terra/main/install.sh | sh -s -- --dry-run
curl -fsSL https://raw.githubusercontent.com/rios0rios0/terra/main/install.sh | sh -s -- --force

# Build from source
git clone https://github.com/rios0rios0/terra.git
cd terra
make install
```

## Command Reference

### Core Commands

#### Root Command (Terragrunt Wrapper)
**Usage:** `terra [flags] [terragrunt command] [directory]`

**Description:** Wraps any Terragrunt command with enhanced functionality including:
- Automatic format execution before running commands
- Cloud provider account/subscription switching
- Terraform workspace management
- Auto-answer support for non-interactive execution
- Parallel state manipulation

**Examples:**
```bash
# Apply all subdirectories
terra apply --all /path/to/infrastructure

# Plan specific module
terra plan /path/to/infrastructure/module

# Plan all subdirectories
terra plan --all /path/to/infrastructure

# With auto-answer (default "n")
terra --auto-answer apply --all /path
terra -a plan --all /path/to

# With explicit "y" responses
terra --auto-answer=y apply --all /path
terra -a=y plan --all /path/to

# With explicit "n" responses
terra --auto-answer=n apply --all /path
terra -a=n plan --all /path/to
```

**Flags:**
- `--auto-answer` or `-a` - Enable automatic responses to prompts (boolean, defaults to "n")
- `--auto-answer=y` or `-a=y` - Automatically answer "y" to prompts
- `--auto-answer=n` or `-a=n` - Automatically answer "n" to prompts

**Note:** The `--auto-answer` flag is filtered out before passing arguments to Terragrunt.

---

#### clear
**Usage:** `terra clear`

**Description:** Clear all temporary directories and cache folders created during Terraform and Terragrunt execution. This includes:
- `.terraform` directories
- `.terragrunt-cache` directories
- Module cache directories

**Examples:**
```bash
terra clear
```

**Flags:** None

---

#### format
**Usage:** `terra format`

**Description:** Format all Terraform and Terragrunt files in the current directory using the formatting commands configured for each dependency (terraform fmt, terragrunt hclfmt).

**Examples:**
```bash
terra format
```

**Flags:** None

**Notes:**
- Warns if terraform or terragrunt are not in PATH
- Automatically runs before other terra commands that execute Terragrunt

---

#### install
**Usage:** `terra install`

**Description:** Install all dependencies required to run Terra, or update them if newer versions are available. Dependencies are installed to `~/.local/bin` on Linux.

**Dependencies installed:**
- Terraform (latest version from HashiCorp)
- Terragrunt (latest version from Gruntwork)

**Examples:**
```bash
terra install
```

**Flags:** None

**Notes:**
- Checks for existing installations and prompts for updates if newer versions are available
- May fail in restricted network environments due to HashiCorp API calls
- Downloads appropriate binaries for the current platform (OS and architecture)

---

#### update
**Usage:** `terra update`

**Description:** Alias for the `install` command. Install or update Terraform and Terragrunt to the latest versions.

**Examples:**
```bash
terra update
```

**Flags:** None

---

#### self-update
**Usage:** `terra self-update [flags]`

**Description:** Download and install the latest version of terra from GitHub releases. Supports dry-run and force modes.

**Examples:**
```bash
# Interactive update (prompts for confirmation)
terra self-update

# Force update without prompts
terra self-update --force

# Dry run to see what would be updated
terra self-update --dry-run
```

**Flags:**
- `--dry-run` - Show what would be updated without performing it
- `--force` - Skip confirmation prompts

**Update Process:**
1. Checks current terra version
2. Fetches latest release from GitHub API
3. Compares versions
4. Downloads appropriate binary for current platform
5. Creates backup of current binary
6. Installs new binary
7. Removes backup on success

---

#### version
**Usage:** `terra version`

**Description:** Show Terra, Terraform, and Terragrunt versions. Displays:
- Terra version (current version installed)
- Terraform version (if installed, otherwise "not installed")
- Terragrunt version (if installed, otherwise "not installed")

**Examples:**
```bash
terra version
```

**Flags:** None

**Output Example:**
```
Terra version: 1.0.0
Terraform version: 1.6.0
Terragrunt version: 0.54.0
```

---

#### completion
**Usage:** `terra completion [shell]`

**Description:** Generate shell completion script for the specified shell (bash, zsh, fish, powershell).

**Examples:**
```bash
# Generate bash completion
terra completion bash

# Generate zsh completion
terra completion zsh
```

**Flags:** None

---

### Parallel State Management

Terra provides powerful parallel execution for state manipulation commands when using the `--all` flag. This feature automatically discovers Terraform/Terragrunt modules and executes state operations across them simultaneously.

**Supported State Commands:**
- `import` - Import existing infrastructure into Terraform state
- `state rm` - Remove resources from state
- `state mv` - Move resources in state
- `state pull` - Pull remote state
- `state push` - Push local state to remote
- `state show` - Show attributes of a resource in state

**How Parallel Execution Works:**
1. **Automatic Module Discovery**: Scans subdirectories for `.tf`, `.tfvars`, or `terragrunt.hcl` files
2. **Parallel Execution**: Runs up to 5 jobs concurrently (default)
3. **Error Aggregation**: Collects and reports errors from all parallel operations
4. **Progress Tracking**: Provides real-time logging of module processing
5. **Flag Filtering**: Removes `--all` flag for individual module execution

**Examples:**
```bash
# Import a resource across all modules in parallel
terra import --all null_resource.example resource-id /path/to/infrastructure

# Remove a resource from state across all modules
terra state rm --all null_resource.example /path/to/infrastructure

# Move a resource in state across all modules
terra state mv --all old_resource.name new_resource.name /path/to/infrastructure

# Pull state from remote across all modules
terra state pull --all /path/to/infrastructure
```

**Module Discovery:**
- Searches for directories containing `.tf`, `.tfvars`, or `terragrunt.hcl` files
- Skips hidden directories (starting with `.`)
- Does not traverse into subdirectories once a module is found

**Concurrency:** Default maximum of 5 concurrent jobs

**Error Handling:**
- Logs errors for individual modules
- Continues processing other modules on failure
- Reports summary at completion (successful/failed counts)
- Returns error if any module fails

---

## Environment Configuration

Terra uses environment variables for configuration. These can be set in a `.env` file in the project root.

### Environment Variables

#### TERRA_CLOUD
**Type:** String  
**Required:** No (but if specified, must be "aws" or "azure")  
**Validation:** Must be either "aws" or "azure" if provided  
**Default:** None

**Purpose:** Specifies the cloud provider for account/subscription switching.

**Example:**
```bash
TERRA_CLOUD=aws
# or
TERRA_CLOUD=azure
```

---

#### TERRA_AWS_ROLE_ARN
**Type:** String  
**Required:** No (required for AWS role switching if TERRA_CLOUD=aws)  
**Format:** ARN format (e.g., `arn:aws:iam::123456789012:role/terraform-role`)  
**Default:** None

**Purpose:** AWS IAM role ARN to assume before executing Terragrunt commands.

**Example:**
```bash
TERRA_AWS_ROLE_ARN=arn:aws:iam::123456789012:role/terraform-role
```

**Behavior:**
- When set with TERRA_CLOUD=aws, terra will execute `aws sts assume-role` before running Terragrunt
- Command executed: `aws sts assume-role --role-arn <ARN> --role-session-name session1`

---

#### TERRA_AZURE_SUBSCRIPTION_ID
**Type:** String  
**Required:** No (required for Azure subscription switching if TERRA_CLOUD=azure)  
**Format:** UUID (e.g., `12345678-1234-1234-1234-123456789012`)  
**Default:** None

**Purpose:** Azure subscription ID to switch to before executing Terragrunt commands.

**Example:**
```bash
TERRA_AZURE_SUBSCRIPTION_ID=12345678-1234-1234-1234-123456789012
```

**Behavior:**
- When set with TERRA_CLOUD=azure, terra will execute `az account set` before running Terragrunt
- Command executed: `az account set --subscription <SUBSCRIPTION_ID>`

---

#### TERRA_WORKSPACE
**Type:** String  
**Required:** No  
**Default:** None

**Purpose:** Terraform workspace to switch to before executing Terragrunt commands.

**Example:**
```bash
TERRA_WORKSPACE=dev
```

**Behavior:**
- When set, terra will execute `terraform workspace select <WORKSPACE>` before running Terragrunt

---

#### TF_VAR_* (Terraform Variables)
**Type:** String  
**Required:** No  
**Format:** Any environment variable starting with `TF_VAR_`  
**Default:** None

**Purpose:** Pass variables to Terraform/Terragrunt. Any environment variable starting with `TF_VAR_` will be automatically available as a Terraform variable.

**Examples:**
```bash
TF_VAR_environment=development
TF_VAR_region=us-west-2
TF_VAR_instance_type=t3.micro
```

**Behavior:**
- These variables are automatically passed to Terraform/Terragrunt
- Variable name is extracted after the `TF_VAR_` prefix
- Example: `TF_VAR_foo=bar` creates a Terraform variable named `foo` with value `bar`

**References:**
- [Terraform documentation](https://www.terraform.io/docs/language/values/variables.html#environment-variables)
- [Terragrunt documentation](https://terragrunt.gruntwork.io/docs/features/inputs/)

---

### Complete .env Example

```bash
# Cloud provider configuration (optional)
TERRA_CLOUD=aws

# AWS specific (required for role switching when using AWS)
TERRA_AWS_ROLE_ARN=arn:aws:iam::123456789012:role/terraform-role

# Azure specific (required for subscription switching when using Azure)
TERRA_AZURE_SUBSCRIPTION_ID=12345678-1234-1234-1234-123456789012

# Terraform workspace (optional)
TERRA_WORKSPACE=dev

# Terraform variables (optional)
TF_VAR_environment=development
TF_VAR_region=us-west-2
TF_VAR_instance_type=t3.micro
```

---

## Execution Flow

### Normal Terragrunt Command Execution

1. Load environment configuration from `.env` file
2. Validate environment variables (TERRA_CLOUD must be "aws" or "azure" if specified)
3. Execute format command (terraform fmt, terragrunt hclfmt)
4. Execute cloud provider account/subscription switching if configured
5. Execute Terraform workspace selection if configured
6. Execute Terragrunt command with provided arguments
7. Handle interactive mode with auto-answer if flag is present

### Parallel State Command Execution

1. Load environment configuration from `.env` file
2. Validate environment variables
3. Check if command is a state manipulation command with `--all` flag
4. Discover all modules in the target path (directories with .tf, .tfvars, or terragrunt.hcl)
5. Remove `--all` flag from arguments
6. Execute command in parallel across all modules (max 5 concurrent)
7. Collect and aggregate results
8. Report summary (successful/failed counts)

---

## Caveats and Known Issues

### Network Restrictions
- `terra install` requires internet access to HashiCorp and Gruntwork APIs
- May fail in restricted network environments
- Workaround: Manually install terraform and terragrunt

### Platform Support
- Automatically detects OS (Linux, Darwin/macOS, Windows) and architecture (amd64, arm64, 386)
- Downloads appropriate binaries for current platform
- Android platform support with special architecture handling

### Windows Path Limitations
- Windows has path size limitations (256 characters)
- WSL interoperability may encounter issues with long paths
- Error example: `/mnt/c/WINDOWS/system32/notepad.exe: Invalid argument`
- **Solution:** Move infrastructure projects closer to home directory to shorten paths

### Validation Requirements
- TERRA_CLOUD must be either "aws" or "azure" if specified
- Application will fail validation if TERRA_CLOUD is set to invalid value
- Empty TERRA_CLOUD is valid (no cloud provider switching)

### Parallel Execution Limitations
- Maximum of 5 concurrent jobs (hardcoded)
- No progress bar (only log messages)
- Cannot be cancelled mid-execution
- All errors are collected and reported at the end

### Auto-Answer Behavior
- Boolean `--auto-answer` or `-a` defaults to "n" for backward compatibility
- Explicit values can be set with `--auto-answer=y` or `--auto-answer=n`
- Flag is filtered out before passing to Terragrunt
- Only works with interactive Terragrunt prompts
- Confirmation prompts may still require manual intervention

---

## Command Patterns for AI Agents

### Pattern 1: Apply Infrastructure Changes
```bash
# Single module
terra apply /path/to/module

# All modules in path
terra apply --all /path/to/infrastructure

# With auto-answer for CI/CD
terra --auto-answer=y apply --all /path/to/infrastructure
```

### Pattern 2: Plan Infrastructure Changes
```bash
# Single module
terra plan /path/to/module

# All modules in path
terra plan --all /path/to/infrastructure
```

### Pattern 3: State Management
```bash
# Import resource to all modules
terra import --all aws_instance.example i-1234567890abcdef0 /path/to/infrastructure

# Remove resource from all modules
terra state rm --all aws_instance.deprecated /path/to/infrastructure

# Move resource in all modules
terra state mv --all aws_instance.old aws_instance.new /path/to/infrastructure
```

### Pattern 4: Dependency Management
```bash
# Check versions
terra version

# Install/update dependencies
terra install

# Update terra itself
terra self-update --force
```

### Pattern 5: Cleanup Operations
```bash
# Clear cache directories
terra clear

# Format code
terra format
```

---

## Technical Details

### Architecture
- **Language:** Go 1.23+
- **CLI Framework:** Cobra
- **Dependency Injection:** Uber's DIG
- **Structure:** Clean Architecture (Domain, Infrastructure layers)

### Build Information
- **Build Time:** 15-20 seconds typical
- **Linting Time:** 2-5 minutes
- **Binary Location:** `bin/terra` (after make build)
- **Install Location:** `~/go/bin/terra` (after make install) or `/usr/local/bin/terra` (system-wide)

### Dependencies
- **Runtime:** terraform, terragrunt (installed via `terra install`)
- **Cloud CLIs:** aws-cli (for AWS), azure-cli (for Azure) - must be installed separately

---

## Exit Codes

Terra uses standard exit codes:
- **0:** Success
- **1:** Failure (logged with logger.Fatalf)

Specific error conditions:
- Environment variable validation failure
- Command execution failure
- Parallel execution failure (one or more modules failed)
- Download/update failure
- File system operation failure

---

## Logging

Terra uses structured logging (logrus):
- **Info:** Normal operation messages
- **Warn:** Non-fatal issues (e.g., backup cleanup failure)
- **Error:** Individual operation errors (e.g., module failure in parallel execution)
- **Fatal:** Terminal errors that stop execution

---

## Best Practices

1. **Use .env files** for environment configuration instead of exporting variables
2. **Run terra install** after initial installation to ensure dependencies are available
3. **Use --auto-answer flag** in CI/CD pipelines to avoid hanging on prompts
4. **Test with terra plan** before using terra apply
5. **Clear cache regularly** with terra clear to avoid state inconsistencies
6. **Keep terra updated** with terra self-update
7. **Use parallel state commands** for bulk operations across multiple modules
8. **Validate paths** are short enough on Windows to avoid path limitations

---

## Resources

- **Repository:** https://github.com/rios0rios0/terra
- **Issues:** https://github.com/rios0rios0/terra/issues
- **License:** MIT
- **Contributing:** See CONTRIBUTING.md
- **Changelog:** See CHANGELOG.md

---

## Version History

- **1.0.0** (2025-10-14): Initial MCP documentation creation

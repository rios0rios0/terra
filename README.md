# Terra

## Introduction
Welcome to `terra` – a powerful wrapper for Terragrunt and Terraform that revolutionizes infrastructure as code management.
Inspired by the simplicity and efficiency of Kubernetes, `terra` allows you to apply Terraform code with the ease of specifying paths directly in your commands.
Our motivation: "Have you ever wondered about applying Terraform code like Kubernetes? Typing the path to be applied after the command you want to use?"

## Features
- Seamless integration with Terraform and Terragrunt.
- Intuitive command-line interface.
- Enhanced state management.
- Simplified module path specification.
- Cross-platform compatibility.
- Auto-answering for Terragrunt prompts to avoid manual intervention.
- Self-update capability to automatically update terra to the latest version.
- Version checking for Terra, Terraform, and Terragrunt dependencies.
- Automatic dependency installation and management.
- Support for AWS and Azure cloud provider switching.
- **Parallel execution for state manipulation commands** - Run `import`, `state rm`, `state mv`, and other state commands across multiple modules simultaneously using the `--all` flag.

## Installation

### Quick Install (Recommended)
Install `terra` with a single command:
```bash
curl -fsSL https://raw.githubusercontent.com/rios0rios0/terra/main/install.sh | sh
```

Or using wget:
```bash
wget -qO- https://raw.githubusercontent.com/rios0rios0/terra/main/install.sh | sh
```

### Installation Options
The installer supports several options:
```bash
# Install specific version
curl -fsSL https://raw.githubusercontent.com/rios0rios0/terra/main/install.sh | sh -s -- --version v1.0.0

# Install to custom directory
curl -fsSL https://raw.githubusercontent.com/rios0rios0/terra/main/install.sh | sh -s -- --install-dir /usr/local/bin

# Show what would be installed without doing it
curl -fsSL https://raw.githubusercontent.com/rios0rios0/terra/main/install.sh | sh -s -- --dry-run

# Force reinstallation
curl -fsSL https://raw.githubusercontent.com/rios0rios0/terra/main/install.sh | sh -s -- --force
```

### Alternative Installation Methods

#### Build from Source
```bash
git clone https://github.com/rios0rios0/terra.git
cd terra
make install
```

#### Download Pre-built Binaries
Download pre-built binaries from the [releases page](https://github.com/rios0rios0/terra/releases).

After installation, you can install Terraform and Terragrunt dependencies automatically:
```bash
terra install
```

To update terra itself to the latest version:
```bash
terra self-update
```

## Usage
Here's how to use `terra` with Terraform/Terragrunt:
```bash
# it's going to apply all subdirectories inside "path"
terra apply --all /path

# it's going to plan all subdirectories inside "to"
terra plan --all /path/to

# it's going to plan just the "module" subdirectory inside "to"
terra plan --all /path/to/module

# or using Terraform approach, plan just the "module" subdirectory inside "to"
terra plan /path/to/module

# with auto-answering to avoid manual prompts (defaults to "n" for backward compatibility)
terra --auto-answer apply --all /path
terra -a plan --all /path/to

# with explicit "y" responses to prompts
terra --auto-answer=y apply --all /path
terra -a=y plan --all /path/to

# with explicit "n" responses to prompts  
terra --auto-answer=n apply --all /path
terra -a=n plan --all /path/to
```

The commands available are:
```bash
clear       Clear all cache and modules directories
format      Format all files in the current directory
install     Install or update Terraform and Terragrunt to the latest versions
update      Install or update Terraform and Terragrunt to the latest versions (alias for install)
self-update Update terra to the latest version
version     Show Terra, Terraform, and Terragrunt versions
```

### Auto-Answer Feature

The `--auto-answer` (or `-a`) flag enables automatic responses to Terragrunt prompts, eliminating the need for manual intervention during long-running operations. This is particularly useful in CI/CD pipelines or when running multiple Terragrunt commands.

**What it does:**
- Automatically answers "n" to external dependency prompts (when using boolean flag or explicit --auto-answer=n)
- Automatically answers "y" to external dependency prompts (when using --auto-answer=y)
- Automatically answers general yes/no prompts with the specified value
- Switches to manual mode for confirmation prompts (like "Are you sure you want to run...")
- Filters out the auto-answer flag before passing arguments to Terragrunt

**Usage Options:**
- `--auto-answer` or `-a` - Boolean flag (defaults to "n" for backward compatibility)
- `--auto-answer=y` - Explicitly answer "y" to prompts
- `--auto-answer=n` - Explicitly answer "n" to prompts
- `-a=y` - Short form to answer "y"
- `-a=n` - Short form to answer "n"

**Example:**
```bash
# Without auto-answer - requires manual input for each prompt
terra apply --all /path

# With boolean auto-answer - automatically answers "n" (backward compatible)
terra --auto-answer apply --all /path

# With explicit "y" answer - automatically answers "y" to prompts
terra --auto-answer=y apply --all /path

# With explicit "n" answer - automatically answers "n" to prompts  
terra --auto-answer=n apply --all /path

# Short form syntax
terra -a=y apply --all /path
terra -a=n plan --all /path
```

### Parallel State Management

Terra provides powerful parallel execution capabilities for state manipulation commands when using the `--all` flag. This feature automatically discovers Terraform/Terragrunt modules in subdirectories and executes state operations across them simultaneously, significantly reducing execution time.

**Supported State Commands:**
- `import` - Import existing infrastructure into Terraform state
- `state rm` - Remove resources from state
- `state mv` - Move resources in state
- `state pull` - Pull remote state
- `state push` - Push local state to remote
- `state show` - Show attributes of a resource in state

**Examples:**
```bash
# Import a resource across all modules in parallel (default: 5 concurrent jobs)
terra import --all null_resource.example resource-id /path/to/infrastructure

# Remove a resource from state across all modules in parallel
terra state rm --all null_resource.example /path/to/infrastructure

# Move a resource in state across all modules in parallel
terra state mv --all old_resource.name new_resource.name /path/to/infrastructure

# Pull state from remote across all modules in parallel
terra state pull --all /path/to/infrastructure
```

**How it works:**
1. **Automatic Module Discovery**: Scans subdirectories for `.tf`, `.tfvars`, or `terragrunt.hcl` files
2. **Parallel Execution**: Runs up to 5 jobs concurrently by default (configurable internally)
3. **Error Aggregation**: Collects and reports errors from all parallel operations
4. **Progress Tracking**: Provides real-time logging of module processing status
5. **Flag Filtering**: Removes `--all` flag for individual module execution since Terragrunt doesn't support it for state commands

**Benefits:**
- **Performance**: Executes across multiple modules simultaneously instead of sequentially
- **Native Integration**: No need for external tools like GNU parallel
- **Error Handling**: Comprehensive error reporting and aggregation
- **Logging**: Detailed progress and completion status for each module

**Note:** Regular Terragrunt commands with `--all` (like `plan --all`, `apply --all`) continue to work normally through Terragrunt's native implementation. Parallel execution is only used for state manipulation commands that don't natively support `--all` in Terragrunt.

### Version Management

#### Checking Versions
Use the `version` command to check Terra, Terraform, and Terragrunt versions:
```bash
terra version
```
This displays:
- Terra version (current version installed)
- Terraform version (if installed, otherwise "not installed")
- Terragrunt version (if installed, otherwise "not installed")

#### Self-Update
Keep terra up to date with the `self-update` command:
```bash
# Interactive update (prompts for confirmation)
terra self-update

# Force update without prompts
terra self-update --force

# Dry run to see what would be updated
terra self-update --dry-run
```

#### Dependency Management
Install or update Terraform and Terragrunt dependencies:
```bash
# Install dependencies (prompts for updates if newer versions available)
terra install

# Alternative command (alias for install)
terra update
```

## Environment Configuration

Terra can be configured with environment variables for cloud provider integration. Create a `.env` file in your project root:

```bash
# Optional: Cloud provider (if specified, must be "aws" or "azure")
TERRA_CLOUD=aws

# AWS specific (required for role switching when using AWS)
TERRA_AWS_ROLE_ARN=arn:aws:iam::123456789012:role/terraform-role

# Azure specific (required for subscription switching when using Azure)
TERRA_AZURE_SUBSCRIPTION_ID=12345678-1234-1234-1234-123456789012

# Optional: Terraform workspace
TERRA_WORKSPACE=dev

# Optional: Terraform variables (any TF_VAR_* variables)
TF_VAR_environment=development
TF_VAR_region=us-west-2
```

**Note**: If `TERRA_CLOUD` is specified, it must be set to either "aws" or "azure". This enables cloud-specific features like role switching for AWS or subscription switching for Azure.

If you have some input variables, you can use environment variables (`.env`) with the prefix `TF_VAR_`:
```bash
# .env example for Terraform variables
TF_VAR_foo=bar
# command (that depends on the environment variable called "foo")
terra apply --all /path/to/module
```
More about it in:
- [Terraform documentation](https://www.terraform.io/docs/language/values/variables.html#environment-variables).
- [Terragrunt documentation](https://terragrunt.gruntwork.io/docs/features/inputs/).

## Known Issues
1. Notice that Windows has `path` size limitations (256 characters).
   If you are using WSL interoperability (calling `.exe` files inside WSL), you could have errors like this:
   ```bash
   /mnt/c/WINDOWS/system32/notepad.exe: Invalid argument
   ```
   That means, you exceeded the `path` size limitation on the current `path` you are running the command.
   To avoid this issue, move your infrastructure project to a shorter `path`. Closer to your "/home" directory, for example.

## Documentation

- **[README.md](README.md)** - Quick start guide and general usage
- **[CONTRIBUTING.md](CONTRIBUTING.md)** - Development guidelines and contribution process
- **[MCP.md](MCP.md)** - Comprehensive documentation for AI agents (Model Context Protocol)
- **[CHANGELOG.md](CHANGELOG.md)** - Version history and release notes

## Contributing
Contributions to `terra` are welcome! Whether it's bug reports, feature requests, or code contributions, please feel free to contribute.
Check out our [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines on how to contribute.

## License
`terra` is released under the [MIT License](LICENSE.md). See the LICENSE file for more details.

---

We hope `terra` makes your infrastructure management smoother and more intuitive. Happy Terraforming!

## Best Practices

1. https://www.terraform.io/docs/extend/best-practices/index.html
2. https://www.terraform-best-practices.com/naming
3. https://www.terraform.io/docs/state/workspaces.html
4. https://terragrunt.gruntwork.io/

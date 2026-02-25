<h1 align="center">Terra</h1>
<p align="center">
    <a href="https://github.com/rios0rios0/terra/releases/latest">
        <img src="https://img.shields.io/github/release/rios0rios0/terra.svg?style=for-the-badge&logo=github" alt="Latest Release"/></a>
    <a href="https://github.com/rios0rios0/terra/blob/main/LICENSE">
        <img src="https://img.shields.io/github/license/rios0rios0/terra.svg?style=for-the-badge&logo=github" alt="License"/></a>
    <a href="https://github.com/rios0rios0/terra/actions/workflows/default.yaml">
        <img src="https://img.shields.io/github/actions/workflow/status/rios0rios0/terra/default.yaml?branch=main&style=for-the-badge&logo=github" alt="Build Status"/></a>
    <a href="https://sonarcloud.io/summary/overall?id=rios0rios0_terra">
        <img src="https://img.shields.io/sonar/coverage/rios0rios0_terra?server=https%3A%2F%2Fsonarcloud.io&style=for-the-badge&logo=sonarqubecloud" alt="Coverage"/></a>
    <a href="https://sonarcloud.io/summary/overall?id=rios0rios0_terra">
        <img src="https://img.shields.io/sonar/quality_gate/rios0rios0_terra?server=https%3A%2F%2Fsonarcloud.io&style=for-the-badge&logo=sonarqubecloud" alt="Quality Gate"/></a>
    <a href="https://www.bestpractices.dev/projects/12031">
        <img src="https://img.shields.io/cii/level/12031?style=for-the-badge&logo=opensourceinitiative" alt="OpenSSF Best Practices"/></a>
</p>

A powerful wrapper for Terragrunt and Terraform that revolutionizes infrastructure as code management. Inspired by the simplicity and efficiency of Kubernetes, `terra` allows you to apply Terraform code with the ease of specifying paths directly in your commands.

## Features

- Seamless integration with Terraform and Terragrunt
- Intuitive command-line interface with path-based syntax
- Enhanced state management
- Simplified module path specification
- Cross-platform compatibility
- Auto-answering for Terragrunt prompts to avoid manual intervention
- Self-update capability to automatically update terra to the latest version
- Version checking for Terra, Terraform, and Terragrunt dependencies
- Automatic dependency installation and management
- Support for AWS and Azure cloud provider switching
- **Parallel execution for any command** - Run any Terragrunt command across multiple modules simultaneously using the `--parallel=N` flag, where N is the number of concurrent threads. State commands also support the legacy `--all` flag for backward compatibility.
- **Cross-platform file locking** - Prevents race conditions when multiple terra processes run concurrently from the same repository
- **Centralized module and provider caching** - Automatically configures `TG_DOWNLOAD_DIR` and `TF_PLUGIN_CACHE_DIR` so Terragrunt modules and providers are downloaded once and reused across all stacks. Override defaults with `TERRA_MODULE_CACHE_DIR` and `TERRA_PROVIDER_CACHE_DIR` environment variables.
- **CAS (Content Addressable Store)** - Enables Terragrunt's experimental CAS by default (`TG_EXPERIMENT=cas`), which deduplicates Git clones via hard links for faster subsequent clones and reduced disk usage. Disable with `TERRA_NO_CAS=true`.
- **Auto-initialization with upgrade detection** - Automatically detects when terraform/terragrunt needs `init --upgrade` (backend changes, provider conflicts, uninitialized modules) and runs it transparently before retrying the original command.

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

### Command Reference

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
- Automatically answers "n" to external dependency prompts (when using boolean flag or explicit `--auto-answer=n`)
- Automatically answers "y" to external dependency prompts (when using `--auto-answer=y`)
- Automatically answers general yes/no prompts with the specified value
- Switches to manual mode for confirmation prompts (like "Are you sure you want to run...")
- Filters out the auto-answer flag before passing arguments to Terragrunt

**Usage Options:**
- `--auto-answer` or `-a` - Boolean flag (defaults to "n" for backward compatibility)
- `--auto-answer=y` or `-a=y` - Explicitly answer "y" to prompts
- `--auto-answer=n` or `-a=n` - Explicitly answer "n" to prompts

**Example:**
```bash
# Without auto-answer - requires manual input for each prompt
terra apply --all /path

# With boolean auto-answer - automatically answers "n" (backward compatible)
terra --auto-answer apply --all /path

# With explicit "y" answer - automatically answers "y" to prompts
terra --auto-answer=y apply --all /path

# Short form syntax
terra -a=y apply --all /path
terra -a=n plan --all /path
```

### Parallel Execution

Terra provides powerful parallel execution capabilities for **any Terragrunt command** using the `--parallel=N` flag, where N is the number of concurrent threads.

**Basic usage:**
```bash
# Run init across all modules with 4 parallel threads
terra init --parallel=4 /path/to/infrastructure

# Run plan with filtering specific directories
terra plan --parallel=4 --filter=dev,staging,prod /path/to/infrastructure

# Exclude specific directories
terra apply --parallel=4 --filter=!test,!backup /path/to/infrastructure

# State commands (--parallel automatically implies --all behavior)
terra import --parallel=4 null_resource.example resource-id /path/to/infrastructure
terra state rm --parallel=2 null_resource.example /path/to/infrastructure

# Legacy --all flag (backward compatibility for state commands)
terra import --all null_resource.example resource-id /path/to/infrastructure

# Forward --parallel to Terragrunt instead of Terra handling it
terra init --parallel=4 --no-parallel-bypass /path/to/infrastructure
```

For comprehensive documentation on parallel execution, including filtering, thread optimization, command scenarios, and all supported commands, see [docs/parallel-execution.md](docs/parallel-execution.md).

### Version Management

#### Checking Versions
```bash
terra version
```
This displays Terra, Terraform, and Terragrunt versions.

#### Self-Update
```bash
# Interactive update (prompts for confirmation)
terra self-update

# Force update without prompts
terra self-update --force

# Dry run to see what would be updated
terra self-update --dry-run
```

#### Dependency Management
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

# Optional: Centralized cache directories (defaults shown below)
# TERRA_MODULE_CACHE_DIR=~/.cache/terra/modules
# TERRA_PROVIDER_CACHE_DIR=~/.cache/terra/providers

# Optional: Disable Terragrunt CAS (Content Addressable Store) experiment
# TERRA_NO_CAS=true
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
- [Terraform documentation](https://www.terraform.io/docs/language/values/variables.html#environment-variables)
- [Terragrunt documentation](https://terragrunt.gruntwork.io/docs/features/inputs/)

## Known Issues

1. Notice that Windows has `path` size limitations (256 characters).
   If you are using WSL interoperability (calling `.exe` files inside WSL), you could have errors like this:
   ```bash
   /mnt/c/WINDOWS/system32/notepad.exe: Invalid argument
   ```
   That means you exceeded the `path` size limitation on the current `path` you are running the command.
   To avoid this issue, move your infrastructure project to a shorter `path`, closer to your "/home" directory, for example.

## Contributing

Contributions are welcome. See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

`terra` is released under the [MIT License](LICENSE.md). See the LICENSE file for more details.

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
- Auto-replying to Terragrunt prompts via `--reply` / `-r` to avoid manual intervention
- Self-update capability to automatically update terra to the latest version
- Version checking for Terra, Terraform, and Terragrunt dependencies
- Automatic dependency installation and management
- Support for AWS and Azure cloud provider switching
- **Parallel execution for any command** - Run any Terragrunt command across multiple modules simultaneously using the `--parallel=N` flag, where N is the number of concurrent threads. Use `--only=mod1,mod2` to select specific modules or `--skip=mod3` to exclude modules.
- **Centralized module and provider caching** - Automatically configures `TG_DOWNLOAD_DIR` and `TG_PROVIDER_CACHE_DIR` so Terragrunt modules and providers are downloaded once and reused across all stacks. Enables the Terragrunt Provider Cache Server (`TG_PROVIDER_CACHE=1`) for concurrent-safe provider deduplication with file locking. Override defaults with `TERRA_MODULE_CACHE_DIR` and `TERRA_PROVIDER_CACHE_DIR` environment variables. Disable the Provider Cache Server with `TERRA_NO_PROVIDER_CACHE=true`.
- **CAS (Content Addressable Store)** - Enables Terragrunt's experimental CAS by default (`TG_EXPERIMENT=cas`), which deduplicates Git clones via hard links for faster subsequent clones and reduced disk usage. Disable with `TERRA_NO_CAS=true`.
- **Partial Parse Config Cache** - Enables Terragrunt's Partial Parse Config Cache by default (`TG_USE_PARTIAL_PARSE_CONFIG_CACHE=true`), which caches parsed HCL configs across modules sharing the same root include for faster config parsing. Disable with `TERRA_NO_PARTIAL_PARSE_CACHE=true`.
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

# with auto-replying "y" to avoid manual prompts
terra --reply=y apply --all /path
terra -r=y plan --all /path/to

# with auto-replying "n" to reject prompts
terra --reply=n apply --all /path
terra -r=n plan --all /path/to
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

### Reply Feature

The `--reply` (or `-r`) flag enables automatic responses to Terragrunt prompts, eliminating the need for manual intervention during long-running operations. This is particularly useful in CI/CD pipelines or when running multiple Terragrunt commands.

**What it does:**
- Automatically answers "n" to external dependency prompts (when using boolean flag or explicit `--reply=n`)
- Automatically answers "y" to external dependency prompts (when using `--reply=y`)
- Automatically answers general yes/no prompts with the specified value
- Switches to manual mode for confirmation prompts (like "Are you sure you want to run...")
- Filters out the reply flag before passing arguments to Terragrunt

**Usage with `--all` (terragrunt-managed parallelism):** requires an explicit value (`--reply=y` or `--reply=n`) because the PTY auto-answering needs to know how to respond.

```bash
# With explicit "y" answer - automatically answers "y" to prompts
terra --reply=y apply --all /path

# Short form syntax
terra -r=y apply --all /path
terra -r=n plan --all /path
```

**Usage with `--parallel` (terra-managed parallelism):** just `--reply` (no value) is sufficient. Terra always injects `--non-interactive` when `--reply` is present, and adds `-auto-approve` automatically for interactive commands like `apply` and `destroy`; the reply value is ignored in this mode.

```bash
# Just --reply is enough for terra-managed parallel
terra apply --parallel=4 --reply /path
terra destroy --parallel=4 -r /path
```

### Parallel Execution

Terra provides two independent parallel execution strategies:

**Terra-managed parallel** (`--parallel=N`) -- terra discovers modules and runs N goroutine workers:
```bash
# Run init across all modules with 4 parallel threads
terra init --parallel=4 /path/to/infrastructure

# Select specific directories with --only
terra plan --parallel=4 --only=dev,staging,prod /path/to/infrastructure

# Skip specific directories with --skip
terra apply --parallel=4 --skip=test,backup /path/to/infrastructure

# State commands
terra import --parallel=4 null_resource.example resource-id /path/to/infrastructure
terra state rm --parallel=2 null_resource.example /path/to/infrastructure
```

**Terragrunt-managed parallel** (`--all` + `--parallelism`) -- forwarded directly to terragrunt:
```bash
# Terragrunt's native run-all
terra apply --all /path/to/infrastructure

# Terragrunt's run-all with concurrency control
terra plan --all --parallelism=4 /path/to/infrastructure

# Terragrunt's run-all with filter
terra apply --all --filter="region-us-east" /path/to/infrastructure
```

> **Note:** `--parallel` and `--all` cannot be used together -- they represent competing execution strategies.

For comprehensive documentation, see [docs/parallel-execution.md](docs/parallel-execution.md).

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

# Optional: Disable Terragrunt Partial Parse Config Cache
# TERRA_NO_PARTIAL_PARSE_CACHE=true

# Optional: Disable automatic workspace selection from TERRA_WORKSPACE
# TERRA_NO_WORKSPACE=true
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

## Benchmarks

Provider caching strategy comparison using `terragrunt init` with `azurerm` provider `v4.42.0` on a single module:

### Speed

| Strategy                   | Cold cache | Warm cache (median) |
|----------------------------|------------|---------------------|
| No cache                   | 12.4s      | 12.4s               |
| `TF_PLUGIN_CACHE_DIR` only | 11.9s      | **8.9s**            |
| `TG_PROVIDER_CACHE` only   | 11.8s      | 10.6s               |
| Both combined              | 11.4s      | 9.5s                |

### Disk usage

| Strategy                   | Shared cache | Per module         | Total (N modules)  |
|----------------------------|--------------|--------------------|--------------------|
| No cache                   | 0            | 238 MB (full copy) | 238 MB x N         |
| `TF_PLUGIN_CACHE_DIR` only | 219 MB       | 19 MB (symlink)    | 219 MB + 19 MB x N |
| `TG_PROVIDER_CACHE` only   | 219 MB       | 19 MB (symlink)    | 219 MB + 19 MB x N |
| Both combined              | 219 MB       | 19 MB (symlink)    | 219 MB + 19 MB x N |

While `TF_PLUGIN_CACHE_DIR` provides slightly better single-module warm-cache performance (8.9s vs 10.6s), it causes "text file busy" (`ETXTBSY`) errors during parallel execution (`--parallel=N`) because Terraform creates symlinks without file locking. Terra uses `TG_PROVIDER_CACHE` (Provider Cache Server) by default because it serializes provider downloads with file locking, making it safe for concurrent access from parallel goroutines. Disable with `TERRA_NO_PROVIDER_CACHE=true`.

## Contributing

Contributions are welcome. See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

## License

`terra` is released under the [MIT License](LICENSE.md). See the LICENSE file for more details.

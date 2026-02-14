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
- **Parallel execution for any command** - Run any Terragrunt command across multiple modules simultaneously using the `--parallel=N` flag, where N is the number of concurrent threads. State commands also support the legacy `--all` flag for backward compatibility.
- **Cross-platform file locking** - Prevents race conditions when multiple terra processes run concurrently from the same repository.
- **Centralized module and provider caching** - Automatically configures `TG_DOWNLOAD_DIR` and `TF_PLUGIN_CACHE_DIR` so Terragrunt modules and providers are downloaded once and reused across all stacks. Override defaults with `TERRA_MODULE_CACHE_DIR` and `TERRA_PROVIDER_CACHE_DIR` environment variables.
- **CAS (Content Addressable Store)** - Enables Terragrunt's experimental CAS by default (`TG_EXPERIMENT=cas`), which deduplicates Git clones via hard links for faster subsequent clones and reduced disk usage. Disable with `TERRA_NO_CAS=true`.

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

### Parallel Execution

Terra provides powerful parallel execution capabilities for **any Terragrunt command** using the `--parallel=N` flag, where N is the number of concurrent threads. This feature automatically discovers Terraform/Terragrunt modules in subdirectories and executes operations across them simultaneously, significantly reducing execution time.

#### Basic Usage

**For any command:**
```bash
# Run init across all modules with 4 parallel threads
terra init --parallel=4 /path/to/infrastructure

# Run plan across all modules with 2 parallel threads
terra plan --parallel=2 /path/to/infrastructure

# Run apply across all modules with 8 parallel threads
terra apply --parallel=8 /path/to/infrastructure
```

**For state commands (--parallel automatically implies --all behavior):**
```bash
# Import a resource across all modules in parallel (4 threads)
# Note: --all flag is NOT needed when using --parallel
terra import --parallel=4 null_resource.example resource-id /path/to/infrastructure

# Remove a resource from state across all modules in parallel (2 threads)
terra state rm --parallel=2 null_resource.example /path/to/infrastructure

# Move a resource in state across all modules in parallel
terra state mv --parallel=4 old_resource.name new_resource.name /path/to/infrastructure

# Pull state from remote across all modules in parallel
terra state pull --parallel=4 /path/to/infrastructure
```

**Legacy --all flag (backward compatibility for state commands only):**
```bash
# State commands still support --all flag (defaults to 5 concurrent jobs)
terra import --all null_resource.example resource-id /path/to/infrastructure
terra state rm --all null_resource.example /path/to/infrastructure
```

#### Filtering Specific Directories

Use the `--filter` flag to specify which subdirectories should be processed in parallel. This is useful when you only want to run commands on specific modules rather than all discovered modules. You can also exclude specific directories by prefixing them with `!`.

**Basic usage (inclusions):**
```bash
# Run init on specific folders (test1, test2, test3) within the target path
terra init --parallel=4 --filter=test1,test2,test3 environments/xpto/prod
# This executes init on:
# - environments/xpto/prod/test1
# - environments/xpto/prod/test2
# - environments/xpto/prod/test3
```

**Exclusion usage:**
```bash
# Run apply on all folders except folder2
terra apply --parallel=4 --filter=!folder2 /path/to/infrastructure

# Run plan on all folders except folder1 and folder3
terra plan --parallel=4 --filter=!folder1,!folder3 /path/to/infrastructure

# Combine inclusions and exclusions: include folder1, but exclude folder2
terra init --parallel=4 --filter=folder1,!folder2 /path/to/infrastructure
```

**How it works:**
- **Inclusions**: Filter values (without `!`) are concatenated with the target path using `filepath.Join()`
- **Exclusions**: Filter values prefixed with `!` exclude matching directories from processing
- When only exclusions are provided, all subdirectories are discovered first, then exclusions are removed
- When both inclusions and exclusions are provided, inclusions are processed first, then exclusions are removed
- Only directories that exist and are valid will be processed
- Non-existent or invalid filter paths are logged as warnings and skipped
- If the number of threads (`--parallel=N`) exceeds the number of filter items, the thread count is automatically reduced to match

**Examples:**
```bash
# Apply changes to specific environments only
terra apply --parallel=3 --filter=dev,staging,prod /path/to/infrastructure

# Plan all modules except test environments
terra plan --parallel=2 --filter=!test,!testing /path/to/infrastructure

# Import resources in specific directories, excluding backup folders
terra import --parallel=4 --filter=region1,region2,!backup null_resource.example resource-id /path/to/infrastructure
```

**Thread count optimization:**
```bash
# If you specify --parallel=4 but only provide 3 filter items, 
# the thread count is automatically reduced to 3
terra init --parallel=4 --filter=test1,test2,test3 /path
# Logs: "Reducing thread count to 3 (number of modules)"
```

#### Forwarding --parallel to Terragrunt

If you want Terragrunt to handle the `--parallel` flag instead of Terra, use the `--no-parallel-bypass` flag:

```bash
# Forward --parallel=4 to terragrunt (Terra won't handle parallel execution)
terra init --parallel=4 --no-parallel-bypass /path/to/infrastructure

# This is useful when you want Terragrunt's native parallel execution behavior
terra plan --parallel=2 --no-parallel-bypass /path/to/infrastructure
```

#### Supported Commands

**All Terragrunt commands support `--parallel=N`:**
- `init` - Initialize Terraform working directory
- `plan` - Generate and show execution plan
- `apply` - Build or change infrastructure
- `destroy` - Destroy Terraform-managed infrastructure
- `validate` - Validate Terraform files
- `fmt` - Format Terraform files
- `import` - Import existing infrastructure into Terraform state
- `state rm` - Remove resources from state
- `state mv` - Move resources in state
- `state pull` - Pull remote state
- `state push` - Push local state to remote
- `state show` - Show attributes of a resource in state
- And any other Terragrunt command

#### How It Works

1. **Automatic Module Discovery**: Scans subdirectories for `.tf`, `.tfvars`, or `terragrunt.hcl` files (unless `--filter` is specified)
2. **Selective Filtering**: When `--filter` is used, only the specified subdirectories are processed (concatenated with the target path)
3. **Parallel Execution**: Runs N jobs concurrently (where N is specified in `--parallel=N`, default is 5 for `--all` flag)
4. **Thread Optimization**: Automatically reduces thread count if it exceeds the number of modules to process
5. **Error Aggregation**: Collects and reports errors from all parallel operations
6. **Progress Tracking**: Provides real-time logging of module processing status
7. **Flag Filtering**: Removes Terra-specific flags (`--parallel=N`, `--all`, `--no-parallel-bypass`, `--filter=`) before passing to Terragrunt

#### Command Scenarios

| Command | Behavior |
|---------|----------|
| `terra init --parallel=4` | Terra handles parallel execution with 4 threads across all modules |
| `terra init --parallel=4 --filter=test1,test2,test3 /path` | Terra handles parallel execution with 4 threads on specified folders only |
| `terra import --parallel=4` | Terra handles parallel execution (equivalent to `--all` for state commands) |
| `terra import --all` | Terra handles parallel execution with 5 threads (backward compatibility) |
| `terra init --parallel=4 --no-parallel-bypass` | `--parallel=4` is forwarded to Terragrunt, Terra doesn't handle parallel execution |
| `terra plan --all` | Terragrunt handles `--all` flag natively (not handled by Terra) |
| `terra apply --all` | Terragrunt handles `--all` flag natively (not handled by Terra) |

#### Benefits

- **Performance**: Executes across multiple modules simultaneously instead of sequentially
- **Flexibility**: Control the number of parallel threads with `--parallel=N`
- **Selective Execution**: Use `--filter` to target specific directories instead of all discovered modules
- **Thread Optimization**: Automatically adjusts thread count to match the number of modules
- **Native Integration**: No need for external tools like GNU parallel
- **Error Handling**: Comprehensive error reporting and aggregation
- **Logging**: Detailed progress and completion status for each module
- **Backward Compatibility**: State commands still support the `--all` flag

#### Notes

- When using `--parallel=N` (without `--no-parallel-bypass`), Terra automatically handles parallel execution for **all commands**, including state commands. For state commands, you don't need to provide `--all` when using `--parallel`.
- Regular Terragrunt commands with `--all` (like `plan --all`, `apply --all`) continue to work normally through Terragrunt's native implementation and are not handled by Terra's parallel execution.
- State commands that don't natively support `--all` in Terragrunt (like `import`, `state rm`, etc.) are handled by Terra's parallel execution when using either `--all` or `--parallel=N`.

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

# Parallel Execution

Terra provides powerful parallel execution capabilities for **any Terragrunt command** using the `--parallel=N` flag, where N is the number of concurrent threads. This feature automatically discovers Terraform/Terragrunt modules in subdirectories and executes operations across them simultaneously, significantly reducing execution time.

## Basic Usage

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

## Filtering Specific Directories

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

## Forwarding --parallel to Terragrunt

If you want Terragrunt to handle the `--parallel` flag instead of Terra, use the `--no-parallel-bypass` flag:

```bash
# Forward --parallel=4 to terragrunt (Terra won't handle parallel execution)
terra init --parallel=4 --no-parallel-bypass /path/to/infrastructure

# This is useful when you want Terragrunt's native parallel execution behavior
terra plan --parallel=2 --no-parallel-bypass /path/to/infrastructure
```

## Supported Commands

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

## How It Works

1. **Automatic Module Discovery**: Scans subdirectories for `.tf`, `.tfvars`, or `terragrunt.hcl` files (unless `--filter` is specified)
2. **Selective Filtering**: When `--filter` is used, only the specified subdirectories are processed (concatenated with the target path)
3. **Parallel Execution**: Runs N jobs concurrently (where N is specified in `--parallel=N`, default is 5 for `--all` flag)
4. **Thread Optimization**: Automatically reduces thread count if it exceeds the number of modules to process
5. **Error Aggregation**: Collects and reports errors from all parallel operations
6. **Progress Tracking**: Provides real-time logging of module processing status
7. **Flag Filtering**: Removes Terra-specific flags (`--parallel=N`, `--all`, `--no-parallel-bypass`, `--filter=`) before passing to Terragrunt

## Command Scenarios

| Command | Behavior |
|---------|----------|
| `terra init --parallel=4` | Terra handles parallel execution with 4 threads across all modules |
| `terra init --parallel=4 --filter=test1,test2,test3 /path` | Terra handles parallel execution with 4 threads on specified folders only |
| `terra import --parallel=4` | Terra handles parallel execution (equivalent to `--all` for state commands) |
| `terra import --all` | Terra handles parallel execution with 5 threads (backward compatibility) |
| `terra init --parallel=4 --no-parallel-bypass` | `--parallel=4` is forwarded to Terragrunt, Terra doesn't handle parallel execution |
| `terra plan --all` | Terragrunt handles `--all` flag natively (not handled by Terra) |
| `terra apply --all` | Terragrunt handles `--all` flag natively (not handled by Terra) |

## Benefits

- **Performance**: Executes across multiple modules simultaneously instead of sequentially
- **Flexibility**: Control the number of parallel threads with `--parallel=N`
- **Selective Execution**: Use `--filter` to target specific directories instead of all discovered modules
- **Thread Optimization**: Automatically adjusts thread count to match the number of modules
- **Native Integration**: No need for external tools like GNU parallel
- **Error Handling**: Comprehensive error reporting and aggregation
- **Logging**: Detailed progress and completion status for each module
- **Backward Compatibility**: State commands still support the `--all` flag

## Notes

- When using `--parallel=N` (without `--no-parallel-bypass`), Terra automatically handles parallel execution for **all commands**, including state commands. For state commands, you don't need to provide `--all` when using `--parallel`.
- Regular Terragrunt commands with `--all` (like `plan --all`, `apply --all`) continue to work normally through Terragrunt's native implementation and are not handled by Terra's parallel execution.
- State commands that don't natively support `--all` in Terragrunt (like `import`, `state rm`, etc.) are handled by Terra's parallel execution when using either `--all` or `--parallel=N`.

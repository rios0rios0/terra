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

**For state commands:**
```bash
# Import a resource across all modules in parallel (4 threads)
terra import --parallel=4 null_resource.example resource-id /path/to/infrastructure

# Remove a resource from state across all modules in parallel (2 threads)
terra state rm --parallel=2 null_resource.example /path/to/infrastructure

# Move a resource in state across all modules in parallel
terra state mv --parallel=4 old_resource.name new_resource.name /path/to/infrastructure

# Pull state from remote across all modules in parallel
terra state pull --parallel=4 /path/to/infrastructure
```

## Selecting Specific Directories

Use the `--only` and `--skip` flags to control which subdirectories should be processed in parallel. This is useful when you only want to run commands on specific modules rather than all discovered modules.

**Selecting specific directories with `--only`:**
```bash
# Run init on specific folders (test1, test2, test3) within the target path
terra init --parallel=4 --only=test1,test2,test3 environments/xpto/prod
# This executes init on:
# - environments/xpto/prod/test1
# - environments/xpto/prod/test2
# - environments/xpto/prod/test3
```

**Skipping specific directories with `--skip`:**
```bash
# Run apply on all folders except folder2
terra apply --parallel=4 --skip=folder2 /path/to/infrastructure

# Run plan on all folders except folder1 and folder3
terra plan --parallel=4 --skip=folder1,folder3 /path/to/infrastructure
```

**Combining `--only` and `--skip`:**
```bash
# Select specific folders but skip a subset
terra init --parallel=4 --only=folder1,folder2,folder3 --skip=folder2 /path/to/infrastructure
```

**How it works:**
- **`--only=`**: Values are concatenated with the target path using `filepath.Join()`. Only these directories are processed.
- **`--skip=`**: Matching directories are removed from processing.
- When only `--skip` is provided, all subdirectories are discovered first, then skipped modules are removed.
- When both `--only` and `--skip` are provided, `--only` is processed first, then `--skip` is applied.
- Only directories that exist and are valid will be processed.
- Non-existent or invalid paths are logged as warnings and skipped.
- If the number of threads (`--parallel=N`) exceeds the number of modules, the thread count is automatically reduced to match.
- The same module cannot appear in both `--only` and `--skip` (validation error).

**Examples:**
```bash
# Apply changes to specific environments only
terra apply --parallel=3 --only=dev,staging,prod /path/to/infrastructure

# Plan all modules except test environments
terra plan --parallel=2 --skip=test,testing /path/to/infrastructure

# Import resources in specific directories, skipping backup folders
terra import --parallel=4 --only=region1,region2 --skip=backup null_resource.example resource-id /path/to/infrastructure
```

**Thread count optimization:**
```bash
# If you specify --parallel=4 but only provide 3 items in --only,
# the thread count is automatically reduced to 3
terra init --parallel=4 --only=test1,test2,test3 /path
# Logs: "Reducing thread count to 3 (number of modules)"
```

## Terragrunt's `--all` and `--parallelism`

Terra's `--parallel=N` is separate from Terragrunt's native `--all` and `--parallelism` flags. They serve different purposes:

| Flag | Owner | Purpose |
|------|-------|---------|
| `--parallel=N` | Terra | Terra manages goroutine workers across module directories |
| `--all` | Terragrunt | Terragrunt's native run-all behavior (forwarded as-is) |
| `--parallelism=N` | Terragrunt | Terragrunt's concurrency for `--all` (forwarded as-is) |
| `--filter=query` | Terragrunt | Terragrunt's config filter language (forwarded as-is) |

```bash
# Terra-managed parallel execution (terra discovers modules, runs N workers)
terra plan --parallel=4 /path/to/infrastructure

# Terragrunt-managed run-all (forwarded directly to terragrunt)
terra apply --all /path/to/infrastructure

# Terragrunt-managed run-all with parallelism and filter
terra apply --all --parallelism=4 --filter="region-us-east" /path/to/infrastructure
```

**Important:** `--parallel` and `--all` cannot be used together -- they represent competing execution strategies.

## Interactive Commands Require `--reply`

When using `--parallel` with `apply` or `destroy`, you **must** provide `--reply` because parallel workers cannot share a single stdin for interactive prompts. Each worker gets its own PTY (pseudo-terminal) that automatically replies to terragrunt prompts with the specified value.

```bash
# ERROR: apply prompts for confirmation, but parallel workers can't share stdin
terra apply --parallel=4 /path/to/infrastructure

# CORRECT: each worker gets its own PTY that auto-replies "y" to prompts
terra apply --parallel=4 --reply=y /path/to/infrastructure

# Decline external dependency prompts across all modules
terra apply --parallel=4 --reply=n /path/to/infrastructure

# Short form
terra destroy --parallel=4 -r=y /path/to/infrastructure

# OK: plan never prompts, so --reply is not required
terra plan --parallel=4 /path/to/infrastructure
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

1. **Automatic Module Discovery**: Scans subdirectories for `.tf`, `.tfvars`, or `terragrunt.hcl` files (unless `--only` is specified)
2. **Selective Filtering**: When `--only` or `--skip` is used, only the matching subdirectories are processed
3. **Parallel Execution**: Runs N jobs concurrently (where N is specified in `--parallel=N`, default is 5)
4. **Thread Optimization**: Automatically reduces thread count if it exceeds the number of modules to process
5. **Error Aggregation**: Collects and reports errors from all parallel operations
6. **Progress Tracking**: Provides real-time logging of module processing status
7. **Flag Filtering**: Removes Terra-specific flags (`--parallel=N`, `--only=`, `--skip=`) before passing to Terragrunt

## Command Scenarios

| Command | Behavior |
|---------|----------|
| `terra init --parallel=4` | Terra handles parallel execution with 4 threads across all modules |
| `terra init --parallel=4 --only=test1,test2,test3 /path` | Terra handles parallel execution with 4 threads on specified folders only |
| `terra import --parallel=4` | Terra handles parallel execution with 4 threads |
| `terra plan --all` | Terragrunt handles `--all` flag natively (forwarded by Terra) |
| `terra apply --all --parallelism=4` | Both flags forwarded to Terragrunt |
| `terra apply --all --filter=mod1` | Both flags forwarded to Terragrunt |

## Benefits

- **Performance**: Executes across multiple modules simultaneously instead of sequentially
- **Flexibility**: Control the number of parallel threads with `--parallel=N`
- **Selective Execution**: Use `--only` and `--skip` to target specific directories instead of all discovered modules
- **Thread Optimization**: Automatically adjusts thread count to match the number of modules
- **Native Integration**: No need for external tools like GNU parallel
- **Error Handling**: Comprehensive error reporting and aggregation
- **Logging**: Detailed progress and completion status for each module
- **Clean Separation**: Terra's `--parallel` and Terragrunt's `--all`/`--parallelism`/`--filter` are independent and unambiguous

## Notes

- When using `--parallel=N`, Terra automatically handles parallel execution for **all commands**, including state commands.
- Regular Terragrunt commands with `--all` (like `plan --all`, `apply --all`) are forwarded to Terragrunt's native implementation and are not handled by Terra's parallel execution.
- `--parallel=N` and `--all` cannot be used together.

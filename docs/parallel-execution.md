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

Terra's `--parallel=N` is separate from Terragrunt's native `--all` and `--parallelism` flags. They serve different purposes and own different filter flags:

| Flag                   | Owner      | Purpose                                                              | Filter flags                    |
|------------------------|------------|----------------------------------------------------------------------|---------------------------------|
| `--parallel=N`         | Terra      | Terra manages goroutine workers across module directories           | `--only=mod1,mod2`, `--skip=mod1,mod2` |
| `--all`                | Terragrunt | Terragrunt's native run-all behavior (forwarded as-is)              | `--filter`, `--queue-exclude-dir`, `--queue-include-dir` |
| `--parallelism=N`      | Terragrunt | Terragrunt's concurrency for `--all` (forwarded as-is)              | n/a                             |
| `--filter=query`       | Terragrunt | Terragrunt's expressive module filter (forwarded as-is)             | n/a                             |

```bash
# Terra-managed parallel execution (terra discovers modules, runs N workers)
terra plan --parallel=4 /path/to/infrastructure
terra plan --parallel=4 --skip=mod1,mod2 /path/to/infrastructure

# Terragrunt-managed run-all (forwarded directly to terragrunt)
terra apply --all /path/to/infrastructure
terra apply --all --filter='!mod1' /path/to/infrastructure

# Terragrunt-managed run-all with parallelism and filter
terra apply --all --parallelism=4 --filter='region-us-east' /path/to/infrastructure
```

**Important:** `--parallel` and `--all` cannot be used together -- they represent competing execution strategies.

## Choosing a strategy

Use this checklist to pick the right strategy before writing the command:

1. **State operation across multiple modules from a root directory?** → must use `--parallel=N`. Terragrunt's `--all` does not support state commands, and terra rejects that combination. Single-module state commands (e.g., `terra state rm <addr> /path/to/one/module`) still work without `--parallel` — terra just forwards them directly to terragrunt.
2. **Need terragrunt DAG ordering / `dependencies` block awareness?** → must use `--all`. Terra's worker pool runs modules in parallel without resolving dependencies between them.
3. **Flat stack, want simple basename filtering?** → either works. `--parallel=N` is slightly faster and its `--only`/`--skip` syntax is shorter for flat stacks.
4. **Need glob, graph, or git-diff filtering?** → must use `--all` with terragrunt's `--filter`. Terra's `--only`/`--skip` only match literal basenames.

## Filter equivalence table

`--only`/`--skip` (terra-managed) and `--filter`/`--queue-exclude-dir` (terragrunt-managed) are **not interchangeable at runtime**, but they have equivalents for common cases. Use the column that matches the strategy you picked above:

| Intent                  | With `--parallel=N`  | With `--all`                          |
|-------------------------|----------------------|---------------------------------------|
| Skip one module         | `--skip=mod1`        | `--filter='!mod1'`                    |
| Skip multiple           | `--skip=mod1,mod2`   | `--filter='!mod1' --filter='!mod2'`   |
| Only specific modules   | `--only=mod1,mod2`   | `--filter='mod1' --filter='mod2'`     |
| Path glob               | *(not supported)*    | `--filter='./prod/**'`                |
| Graph expression        | *(not supported)*    | `--filter='service...'`               |
| Git-diff expression     | *(not supported)*    | `--filter='[main...HEAD]'`            |

## Known differences between the two strategies

Terra deliberately does **not** translate `--skip`/`--only` into `--queue-exclude-dir`/`--filter` on the `--all` path, and does not translate terragrunt's filter flags into terra's worker-pool selection on the `--parallel` path. There are three reasons:

1. **Matching semantics differ.** Terra's `--skip=lab3` matches modules by basename anywhere in the subtree. Terragrunt's `--queue-exclude-dir=lab3` matches paths relative to the working directory and follows terragrunt-specific glob rules. Any automatic translation would either lie about semantics or require terra to walk the module tree itself, which defeats the purpose of `--all`.
2. **Upstream parsing quirks.** [gruntwork-io/terragrunt#5124](https://github.com/gruntwork-io/terragrunt/issues/5124) documents that `--queue-exclude-dir` still parses excluded directories during dependency discovery. A module skipped by terra's native `--skip` is never touched; the same name forwarded to terragrunt's queue flag still goes through DAG parsing. Surfacing that divergence as a silent translation would cause subtle incidents in CI.
3. **`--filter` is strictly more expressive.** Translating `--skip=a,b` to `--queue-exclude-dir` would downgrade capability for `--all` users who have access to graph and git-diff expressions.

Instead, terra provides **discoverability**: when you use `--skip` with `--all`, the validation error echoes your command and shows the exact `--filter` form you should type. When you use terragrunt's `--filter`/`--queue-exclude-dir` with `--parallel=N`, terra logs a warning pointing you at `--only`/`--skip`.

## Interactive Commands Require `--yes` (or `--no`)

When using `--parallel` with `apply` or `destroy`, you **must** provide a confirmation flag because parallel workers cannot share a single stdin for interactive prompts. Terra translates these flags into native Terraform and Terragrunt flags before forwarding the command, so there is no PTY pattern matching involved:

- `--yes` / `-y` injects Terragrunt's `--non-interactive` plus Terraform's `-auto-approve` for `apply` / `destroy`.
- `--no` / `-n` injects only Terragrunt's `--non-interactive`, which causes Terraform's apply prompt to abort instead of proceeding -- matching a "no" answer.

```bash
# ERROR: apply prompts for confirmation, but parallel workers can't share stdin
terra apply --parallel=4 /path/to/infrastructure

# CORRECT: --yes maps to --non-interactive -auto-approve
terra apply --parallel=4 --yes /path/to/infrastructure

# Short form
terra destroy --parallel=4 -y /path/to/infrastructure

# Non-interactive, but abort instead of auto-approving
terra apply --parallel=4 --no /path/to/infrastructure

# OK: plan never prompts, so no confirmation flag is required
terra plan --parallel=4 /path/to/infrastructure
```

> **Note:** The legacy `--reply` / `-r` flags still work and emit a deprecation warning. `--reply=y` and bare `--reply` map to `--yes`; `--reply=n` maps to `--no`. Migrate to `--yes` / `--no` at your earliest convenience; `--reply` will be removed in a future release.

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
- **Clean Separation**: Terra's `--parallel` and Terragrunt's `--all`/`--parallelism`/`--filter` are independent and unambiguous; mixing them produces educational validation errors that show the correct form for your command rather than silently doing the wrong thing

## Notes

- When using `--parallel=N`, Terra automatically handles parallel execution for **all commands**, including state commands.
- Regular Terragrunt commands with `--all` (like `plan --all`, `apply --all`) are forwarded to Terragrunt's native implementation and are not handled by Terra's parallel execution.
- `--parallel=N` and `--all` cannot be used together.
- When modules share Terragrunt dependencies, concurrent `terraform init` may trigger a Git ref backend race condition. See [parallel-git-clone-race.md](parallel-git-clone-race.md) for details and workarounds.

# Parallel Git Clone Race Condition

When using `--parallel` with modules that share common Terragrunt dependencies, concurrent `terraform init` processes may attempt to `git clone` the same module source into the same `.terraform/modules/` directory simultaneously. This triggers a fatal Git assertion:

```
BUG: refs/files-backend.c:3179: initial ref transaction called with existing refs
```

## Root Cause

The crash originates in Git's **files ref backend**. During `git clone`, Git creates a fresh repository and populates its references via an "initial ref transaction" (flagged with `REF_TRANSACTION_FLAG_INITIAL`). This transaction is designed to operate on a **completely empty reference store** -- it writes directly to `packed-refs` and deliberately skips locking on loose refs for performance.

When two parallel workers both trigger `terraform init` for the same shared dependency (e.g., both `1021-lab1` and `1041-lab2` depend on `01_shared/01_tfstate`), Terragrunt resolves the dependency to the **same filesystem path**. Two concurrent `git clone` operations then write to the same `.git` directory:

1. Process A creates the `.git` directory and starts its initial ref transaction
2. Process B also writes refs into the same `.git` directory
3. When Process A tries to commit, it finds refs that Process B already wrote
4. The BUG assertion fires because the initial transaction expected an empty ref store

This is not a bug in Terra or Terragrunt -- it is a fundamental limitation of Git's files ref backend when two processes clone into the same destination concurrently.

## Why No Git Flag Fixes This

- **No `GIT_*` environment variable** makes the files backend's initial ref transaction concurrent-safe. The design explicitly assumes single-process access during clone.
- **`GIT_CLONE_PROTECTION_ACTIVE`** is unrelated -- it is a security measure for `core.hooksPath` introduced in Git 2.39.4+.
- **`--ref-format=reftable`** (Git 2.45+) would avoid the assertion entirely, since the reftable backend has no separate "initial transaction" path and was designed with concurrency in mind. However, Terraform's module installer (go-getter) calls `git` internally and does not pass this flag, so it cannot be used today.
- **go-getter** (HashiCorp's library used by both Terraform and Terragrunt for module downloads) has **no locking mechanism** for concurrent operations to the same destination directory.

## What Terra Already Does

Terra mitigates related parallel issues through several mechanisms:

| Mechanism | What It Solves |
|-----------|---------------|
| **Provider Cache Server** (`TG_PROVIDER_CACHE=1`) | Prevents "text file busy" errors from concurrent provider downloads via a localhost proxy with file locking |
| **CAS** (`TG_EXPERIMENT=cas`) | Deduplicates Git clones via hard links with `.lock` files, reducing contention on module downloads |
| **Centralized module cache** (`TG_DOWNLOAD_DIR`) | Shares downloaded modules across stacks to avoid redundant downloads |
| **Unset `TF_PLUGIN_CACHE_DIR`** | Prevents Terraform's native plugin cache from causing symlink races |

However, none of these fully prevent the race condition when two parallel workers trigger `terraform init` for the same shared dependency at the same instant. The CAS feature helps in many cases by deduplicating clones, but the race window still exists when both workers request the same module before CAS can acquire its lock.

## Workarounds

### 1. Pre-warm caches with sequential init (recommended)

Run `terra init` for each module before running the parallel apply. This ensures all Git clones and provider downloads complete without contention:

```bash
# Initialize each module sequentially first
terra init environments/06_opensearch/dev/1021-lab1
terra init environments/06_opensearch/dev/1041-lab2

# Then apply in parallel (no git clones needed, everything is cached)
terra apply --parallel=2 --reply environments/06_opensearch/dev
```

### 2. Use `--parallel=1` for modules with shared dependencies

When modules share the same Terragrunt dependencies, serializing execution avoids the race entirely:

```bash
terra apply --parallel=1 --reply environments/06_opensearch/dev
```

### 3. Use `--all` for terragrunt-managed parallelism

Terragrunt's native `run-all` manages its own concurrency and serializes init steps internally when needed:

```bash
terra apply --all --reply=y environments/06_opensearch/dev
```

## Community References

- [Terragrunt issue #2542: Init concurrency/parallelism issues](https://github.com/gruntwork-io/terragrunt/issues/2542)
- [Terragrunt issue #3093: run-all init fails with parallelism](https://github.com/gruntwork-io/terragrunt/issues/3093)
- [Terragrunt issue #4535: Race condition with nested stacks and go-getter](https://github.com/gruntwork-io/terragrunt/issues/4535)
- [go-getter source: no locking in get_git.go](https://github.com/hashicorp/go-getter/blob/master/get_git.go)
- [Git source: refs/files-backend.c initial ref transaction](https://github.com/git/git/blob/master/refs/files-backend.c)
- [Git reftable format (Git 2.45+)](https://github.blog/open-source/git/highlights-from-git-2-45/)

## Future

The Git project plans to make the **reftable backend** the default in Git 3.0. Reftable was designed for concurrency with atomic writes and append-only semantics, which would eliminate this class of race condition. Once go-getter and the broader ecosystem adopt reftable, this issue will resolve itself. Until then, the workarounds above are the recommended approach.

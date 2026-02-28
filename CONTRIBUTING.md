# Contributing

Contributions are welcome. By participating, you agree to maintain a respectful and constructive environment.

For coding standards, testing patterns, architecture guidelines, commit conventions, and all
development practices, refer to the **[Development Guide](https://github.com/rios0rios0/guide/wiki)**.

## Prerequisites

- [Go](https://go.dev/dl/) 1.26+
- [Make](https://www.gnu.org/software/make/)
- [Terraform](https://developer.hashicorp.com/terraform/install) (runtime dependency)
- [Terragrunt](https://terragrunt.gruntwork.io/docs/getting-started/install/) (runtime dependency)
- [golangci-lint](https://golangci-lint.run/) (for linting)
- [Pipelines repo](https://github.com/rios0rios0/pipelines) cloned at `~/Development/github.com/rios0rios0/pipelines` (for shared Makefile targets)

## Development Workflow

1. Fork and clone the repository
2. Create a branch: `git checkout -b feat/my-change`
3. Set up the shared pipelines (one-time):
   ```bash
   make setup
   ```
4. Install dependencies:
   ```bash
   go mod download
   ```
5. Build the binary:
   ```bash
   make build
   ```
   This compiles the CLI to `bin/terra` from `./cmd/terra`.
6. Run the application (without building):
   ```bash
   make run
   ```
7. Build a debug binary (with symbols for debuggers):
   ```bash
   make debug
   ```
8. Install locally:
   ```bash
   make install
   ```
   This builds and copies `bin/terra` to `~/.local/bin/terra`.
9. Run the linter:
   ```bash
   make lint
   ```
10. Run tests:
    ```bash
    make test
    ```
11. Run static analysis (SAST):
    ```bash
    make sast
    ```
12. Commit following the [commit conventions](https://github.com/rios0rios0/guide/wiki/Life-Cycle/Git-Flow)
13. Open a pull request against `main`

## Local Environment

Copy `.env.example` to `.env` and fill in the required values:

```bash
cp .env.example .env
```

| Variable | Description | Required |
|----------|-------------|----------|
| `TERRA_CLOUD` | Cloud provider (`aws` or `azure`) | Yes |
| `TERRA_AWS_ROLE_ARN` | AWS IAM role ARN for Terraform | If AWS |
| `TERRA_AZURE_SUBSCRIPTION_ID` | Azure subscription ID | If Azure |
| `TERRA_WORKSPACE` | Terraform workspace name | No |

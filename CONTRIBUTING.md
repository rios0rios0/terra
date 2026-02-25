# Contributing

Contributions are welcome. By participating, you agree to maintain a respectful and constructive environment.

For coding standards, testing patterns, architecture guidelines, commit conventions, and all
development practices, refer to the **[Development Guide](https://github.com/rios0rios0/guide/wiki)**.

## Prerequisites

- [Go](https://go.dev/dl/) 1.26+
- [Make](https://www.gnu.org/software/make/)

## Development Workflow

1. Fork and clone the repository
2. Create a branch: `git checkout -b feat/my-change`
3. Install dependencies:
   ```bash
   go mod download
   ```
4. Build the binary:
   ```bash
   make build
   ```
5. Make your changes
6. Validate:
   ```bash
   make lint
   make test
   make sast
   ```
7. Update `CHANGELOG.md` under `[Unreleased]`
8. Commit following the [commit conventions](https://github.com/rios0rios0/guide/wiki/Life-Cycle/Git-Flow)
9. Open a pull request against `main`

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

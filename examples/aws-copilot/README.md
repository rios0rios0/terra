# AWS Copilot Examples for Terra

This directory contains examples for configuring Terra in AWS Copilot environments.

## Network Configuration

### terra-network.yml
AWS CloudFormation template to be placed in your Copilot service's `addons/` directory. This template creates:

- Security group allowing outbound HTTPS (port 443) for dependency downloads
- Optional VPC endpoint for GitHub access in private subnets
- Exports for use in other CloudFormation templates

### Usage

1. Copy `terra-network.yml` to your service's addons directory:
   ```bash
   cp examples/aws-copilot/terra-network.yml copilot/[your-service]/addons/
   ```

2. Deploy your Copilot service:
   ```bash
   copilot svc deploy --name [your-service] --env [environment]
   ```

3. The security group will be automatically applied to your service, allowing terra dependency downloads.

### Required Domains

The security group allows outbound HTTPS to all destinations. Terra specifically requires access to:

- `releases.hashicorp.com` - Terraform downloads
- `checkpoint-api.hashicorp.com` - Terraform version checks
- `github.com` - Terragrunt downloads  
- `api.github.com` - Terragrunt version checks

### Alternative: Environment Variables

Instead of modifying security groups, you can configure terra to use alternative URLs:

```bash
# In your Copilot service configuration or .env file
TERRAFORM_VERSION_URL=https://your-internal-mirror.com/terraform/version
TERRAFORM_BINARY_URL=https://your-internal-mirror.com/terraform/%[1]s/terraform_%[1]s_linux_amd64.zip
TERRAGRUNT_VERSION_URL=https://your-internal-mirror.com/terragrunt/version  
TERRAGRUNT_BINARY_URL=https://your-internal-mirror.com/terragrunt/v%s/terragrunt_linux_amd64
```
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

## Installation
To install `terra`, ensure you have Terraform and Terragrunt installed on your system
(you don't need, we install for you with `terra install`!!), then run the following command:
```bash
# Replace with the actual installation command for `terra`
git clone https://github.com/rios0rios0/terra.git
cd terra
make install
```

## Network Requirements

When using `terra install` to download dependencies, the following network access is required:

### Required Firewall Rules
- **Outbound HTTPS (port 443)** access to:
  - `releases.hashicorp.com` - Terraform binary downloads
  - `checkpoint-api.hashicorp.com` - Terraform version checks  
  - `github.com` - Terragrunt binary downloads
  - `api.github.com` - Terragrunt version checks

### Environment-Specific Configuration

#### Restricted Networks/Copilot Environments
If you're running terra in environments with restricted egress (like AWS Copilot services), configure your security groups or firewall rules to allow the above endpoints.

#### Proxy Configuration  
Set proxy environment variables if your network requires proxy access:
```bash
export TERRA_HTTPS_PROXY=https://your-proxy:8080
export TERRA_HTTP_PROXY=http://your-proxy:8080
```

#### URL Overrides
Override download URLs to use internal mirrors or alternative sources:
```bash
# Override Terraform URLs
export TERRAFORM_VERSION_URL=https://your-internal-mirror.com/terraform/version
export TERRAFORM_BINARY_URL=https://your-internal-mirror.com/terraform/%[1]s/terraform_%[1]s_linux_amd64.zip

# Override Terragrunt URLs  
export TERRAGRUNT_VERSION_URL=https://your-internal-mirror.com/terragrunt/version
export TERRAGRUNT_BINARY_URL=https://your-internal-mirror.com/terragrunt/v%s/terragrunt_linux_amd64
```

#### AWS Copilot Integration
For AWS Copilot services, add the following to your `addons/` directory to automatically configure the required firewall rules:

**addons/terra-network.yml**:
```yaml
Parameters:
  App:
    Type: String
  Env:
    Type: String

Resources:
  TerraEgressSecurityGroup:
    Type: AWS::EC2::SecurityGroup
    Properties:
      GroupDescription: Allow terra dependency downloads
      VpcId:
        Fn::ImportValue: !Sub ${App}-${Env}-VpcId
      SecurityGroupEgress:
        - IpProtocol: tcp
          FromPort: 443
          ToPort: 443
          CidrIp: 0.0.0.0/0
          Description: HTTPS for terra dependency downloads

Outputs:
  TerraEgressSecurityGroupId:
    Description: Security group for terra network access
    Value: !Ref TerraEgressSecurityGroup
```

Then reference this security group in your service's `copilot/[service]/addons/` or main service configuration.

## Usage
Here's how to use `terra` with Terraform/Terragrunt:
```bash
# it's going to apply all subdirectories inside "path"
terra run-all apply /path

# it's going to plan all subdirectories inside "to"
terra run-all plan /path/to

# it's going to plan just the "module" subdirectory inside "to"
terra run-all plan /path/to/module

# or using Terraform approach, plan just the "module" subdirectory inside "to"
terra plan /path/to/module
```

The commands available are:
```bash
clear       Clear all cache and modules directories
fmt         Format all files in the current directory
install     Install Terraform and Terragrunt (they are pre-requisites)
```

If you have some input variables, you can use environment variables (`.env`) with the prefix `TF_VAR_`:
```bash
# .env
TF_VAR_foo=bar


# command (that depends on the environment variable called "foo")
terra run-all apply /path /path/to/module
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

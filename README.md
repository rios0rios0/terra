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

## Usage
Here's how to use `terra` (Terraform basics):
```bash
# Applying a specific Terraform module
terra apply /path/to/module

# Destroying a module
terra destroy /path/to/module

# Plan a module
terra plan /path/to/module
```

Here's how to use `terra` with Terragrunt:
```bash
# it's going to apply all subdirectories inside "path"
terra run-all apply /path

# it's going to plan all subdirectories inside "to"
terra run-all plan /path/to

# it's going to plan just the "module" subdirectory inside "to"
terra run-all plan /path/to/module
```

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

## TODO
- Ensure the newest dependencies when Terraform and Terragrunt are installed with the previous versions

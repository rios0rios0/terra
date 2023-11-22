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
- Forward unknown flags to Terraform and Terragrunt
- It's not deleting all cache directories, some of them are still there
- Add the feature to read the env variable and set terraform workspace

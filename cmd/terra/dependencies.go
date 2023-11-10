package main

import (
	"fmt"
	logger "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"strings"
)

const (
	terraformURL         = "https://releases.hashicorp.com/terraform/%[1]s/terraform_%[1]s_linux_amd64.zip"
	terraformVersionURL  = "https://checkpoint-api.hashicorp.com/v1/check/terraform"
	terragruntURL        = "https://github.com/gruntwork-io/terragrunt/releases/download/v%s/terragrunt_linux_amd64"
	terragruntVersionURL = "https://api.github.com/repos/gruntwork-io/terragrunt/releases/latest"
)

// Check if a command exists
func isCommandAvailable(name string) bool {
	cmd := exec.Command(name, "-v")
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

// Install a command using curl
func installCommand(url, name string) {
	tempFilePath := "/tmp/" + name
	logger.Infof("Downloading %s from %s...", name, url)
	curlCmd := exec.Command("curl", "-Ls", "-o", tempFilePath, url)
	curlCmd.Stderr = os.Stderr
	curlCmd.Stdout = os.Stdout

	if err := curlCmd.Run(); err != nil {
		logger.Fatalf("Failed to download %s: %s", name, err)
	}

	// Check if the file is a zip
	fileTypeCmd := exec.Command("file", tempFilePath)
	fileTypeOutput, err := fileTypeCmd.Output()
	if err != nil {
		logger.Fatalf("Failed to determine file type of %s: %s", name, err)
	}

	if strings.Contains(string(fileTypeOutput), "Zip archive data") {
		logger.Infof("%s is a zip file, extracting...", name)
		unzipCmd := exec.Command("unzip", "-o", tempFilePath, "-d", "/usr/local/bin")
		unzipCmd.Stderr = os.Stderr
		unzipCmd.Stdout = os.Stdout

		if err := unzipCmd.Run(); err != nil {
			logger.Fatalf("Failed to extract %s: %s", name, err)
		}

		// Remove the zip file
		rmCmd := exec.Command("rm", tempFilePath)
		if err := rmCmd.Run(); err != nil {
			logger.Fatalf("Failed to remove %s: %s", name, err)
		}
	} else {
		// Move file to /usr/local/bin if not a zip
		mvCmd := exec.Command("mv", tempFilePath, "/usr/local/bin/"+name)
		if err := mvCmd.Run(); err != nil {
			logger.Fatalf("Failed to move %s to /usr/local/bin: %s", name, err)
		}
	}

	// Make the binary executable
	err = os.Chmod("/usr/local/bin/"+name, 0755)
	if err != nil {
		logger.Fatalf("Failed to make %s executable: %s", name, err)
	}
}

func ensureRootPrivileges() {
	// Check if the program has root privileges
	if os.Geteuid() != 0 {
		logger.Fatalf("Run this command with root privileges to install the dependencies")
		return
	}
}

// Check if Terraform and Terragrunt are installed, install if not install them
func ensureToolsInstalled() {
	terraformVersion := fetchLatestVersion(terraformVersionURL, `"current_version":"([^"]+)"`)
	terragruntVersion := fetchLatestVersion(terragruntVersionURL, `"tag_name":"v([^"]+)"`)

	// TODO: this could be a for using mapper
	if !isCommandAvailable("terraform") {
		ensureRootPrivileges()
		logger.Warn("Terraform is not installed, installing now...")
		installCommand(fmt.Sprintf(terraformURL, terraformVersion), "terraform")
	}

	if !isCommandAvailable("terragrunt") {
		ensureRootPrivileges()
		logger.Warn("Terragrunt is not installed, installing now...")
		installCommand(fmt.Sprintf(terragruntURL, terragruntVersion), "terragrunt")
	}
}

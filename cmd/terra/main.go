package main

import (
	logger "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"strings"
)

// Run commands in the specified directory
func runInDir(command string, args []string, dir string) error {
	logger.Infof("Running [%s %s] in %s", command, strings.Join(args, " "), dir)
	cmd := exec.Command(command, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func main() {
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(formatCmd)
	rootCmd.AddCommand(clearCmd)
	if err := rootCmd.Execute(); err != nil {
		logger.Fatalf("Error executing terra: %s", err)
	}
}

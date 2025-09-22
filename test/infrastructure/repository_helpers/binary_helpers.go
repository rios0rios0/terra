package repository_helpers

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// HelperCreateMockTerraformBinary creates a temporary binary that mimics terraform --version output.
func HelperCreateMockTerraformBinary(t *testing.T, version string) string {
	t.Helper()
	
	tempDir := t.TempDir()
	
	binaryPath := filepath.Join(tempDir, "terraform")
	
	// Create a shell script that outputs proper terraform version format.
	scriptContent := fmt.Sprintf(`#!/bin/bash
if [ "$1" = "--version" ]; then
    echo "Terraform v%s"
elif [ "$1" = "-v" ]; then
    echo "Terraform v%s"
else
    echo "mock binary"
fi
`, version, version)
	
	err := os.WriteFile(binaryPath, []byte(scriptContent), 0755)
	require.NoError(t, err)
	
	return tempDir
}

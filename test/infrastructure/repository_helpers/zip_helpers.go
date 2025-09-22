package repository_helpers

import (
	"archive/zip"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

// HelperCreateZipWithBinary creates a zip file containing a binary file
func HelperCreateZipWithBinary(t *testing.T, binaryName string) []byte {
	t.Helper()
	
	// Create temporary file to write zip data
	tmpfile, err := os.CreateTemp("", "test-*.zip")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name())
	defer tmpfile.Close()
	
	// Create zip writer
	zipWriter := zip.NewWriter(tmpfile)
	defer zipWriter.Close()
	
	// Add binary file to zip
	binaryFile, err := zipWriter.Create(binaryName)
	require.NoError(t, err)
	
	// Write some binary content
	binaryContent := []byte("#!/bin/bash\necho 'mock binary'\n")
	_, err = binaryFile.Write(binaryContent)
	require.NoError(t, err)
	
	zipWriter.Close()
	
	// Read zip data
	zipData, err := os.ReadFile(tmpfile.Name())
	require.NoError(t, err)
	
	return zipData
}

// HelperCreateNestedZipWithBinary creates a zip with binary in nested directory
func HelperCreateNestedZipWithBinary(t *testing.T, nestedPath, binaryName string) []byte {
	t.Helper()
	
	tmpfile, err := os.CreateTemp("", "test-nested-*.zip")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name())
	defer tmpfile.Close()
	
	zipWriter := zip.NewWriter(tmpfile)
	defer zipWriter.Close()
	
	// Create nested directory structure and add binary
	binaryFile, err := zipWriter.Create(nestedPath)
	require.NoError(t, err)
	
	binaryContent := []byte("#!/bin/bash\necho 'nested binary'\n")
	_, err = binaryFile.Write(binaryContent)
	require.NoError(t, err)
	
	zipWriter.Close()
	
	zipData, err := os.ReadFile(tmpfile.Name())
	require.NoError(t, err)
	
	return zipData
}

// HelperCreateMixedContentZip creates a zip with various file types
func HelperCreateMixedContentZip(t *testing.T, binaryName string) []byte {
	t.Helper()
	
	tmpfile, err := os.CreateTemp("", "test-mixed-*.zip")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name())
	defer tmpfile.Close()
	
	zipWriter := zip.NewWriter(tmpfile)
	defer zipWriter.Close()
	
	// Add README.md file
	readmeFile, err := zipWriter.Create("README.md")
	require.NoError(t, err)
	_, err = readmeFile.Write([]byte("# Test Tool\nThis is a test tool."))
	require.NoError(t, err)
	
	// Add config.json file
	configFile, err := zipWriter.Create("config.json")
	require.NoError(t, err)
	_, err = configFile.Write([]byte(`{"version": "1.0.0"}`))
	require.NoError(t, err)
	
	// Add LICENSE.txt file
	licenseFile, err := zipWriter.Create("LICENSE.txt")
	require.NoError(t, err)
	_, err = licenseFile.Write([]byte("MIT License..."))
	require.NoError(t, err)
	
	// Add changelog.yml file
	changelogFile, err := zipWriter.Create("CHANGELOG.yml")
	require.NoError(t, err)
	_, err = changelogFile.Write([]byte("v1.0.0:\n  - Initial release"))
	require.NoError(t, err)
	
	// Add the actual binary (no extension)
	binaryFile, err := zipWriter.Create(binaryName)
	require.NoError(t, err)
	binaryContent := []byte("#!/bin/bash\necho 'actual binary'\n")
	_, err = binaryFile.Write(binaryContent)
	require.NoError(t, err)
	
	zipWriter.Close()
	
	zipData, err := os.ReadFile(tmpfile.Name())
	require.NoError(t, err)
	
	return zipData
}

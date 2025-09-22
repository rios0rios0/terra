//go:build integration || unit || test

package repositorydoubles //nolint:staticcheck // Test package naming follows established project structure

import (
	"github.com/rios0rios0/terra/internal/infrastructure/repositories"
)

// NewTestInteractiveShellRepository creates an InteractiveShellRepository for testing
func NewTestInteractiveShellRepository() *repositories.InteractiveShellRepository {
	return repositories.NewInteractiveShellRepository()
}
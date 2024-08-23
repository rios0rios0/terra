package repositories

import (
	"github.com/google/wire"
	"github.com/rios0rios0/terra/internal/domain/repositories"
)

//nolint:gochecknoglobals
var Container = wire.NewSet(
	NewDefaultOSRepository,
	wire.Bind(new(repositories.OSRepository), new(*DefaultOSRepository)),
	NewHttpWebStringsRepository,
	wire.Bind(new(repositories.WebStringsRepository), new(*HttpWebStringsRepository)),
)

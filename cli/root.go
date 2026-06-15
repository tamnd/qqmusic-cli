// Package cli assembles the qqmusic command tree from the qqmusic domain on
// top of the any-cli/kit framework.
package cli

import (
	"github.com/tamnd/any-cli/kit"
	"github.com/tamnd/qqmusic-cli/qqmusic"
)

// Build metadata, set via -ldflags at release time.
var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

// NewApp assembles the kit application from the qqmusic domain.
func NewApp() *kit.App {
	id := qqmusic.BaseIdentity()
	id.Version = Version

	app := kit.New(id, kit.WithDefaults(qqmusic.Defaults))
	qqmusic.Register(app)

	return app
}

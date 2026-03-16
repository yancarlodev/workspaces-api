package assets

import "embed"

//go:embed scripts/**/*.sql
var Fsys embed.FS

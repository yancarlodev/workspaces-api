package assets

import "embed"

//go:embed scripts/migrations/*.sql
var Migrations embed.FS

// Package migrations contains the reviewed, transactional database migrations.
package migrations

import "embed"

// Files is the immutable migration input shipped with the database tool.
//
//go:embed orders/*.sql inventory/*.sql payments/*.sql
var Files embed.FS

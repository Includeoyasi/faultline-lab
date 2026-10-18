package database

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io/fs"
	"regexp"
	"strconv"
)

type migration struct {
	version             int
	name, checksum, sql string
}

var migrationName = regexp.MustCompile(`^([0-9]{4})_[a-z0-9_]+\.sql$`)

func loadMigrations(files fs.FS, service string) ([]migration, error) {
	entries, err := fs.ReadDir(files, service)
	if err != nil {
		return nil, errors.New("cannot read migrations")
	}
	var result []migration
	previous := 0
	for _, entry := range entries {
		match := migrationName.FindStringSubmatch(entry.Name())
		if entry.IsDir() || match == nil {
			return nil, errors.New("invalid migration filename")
		}
		version, _ := strconv.Atoi(match[1])
		if version <= previous {
			return nil, errors.New("migration versions must be positive and strictly increasing")
		}
		data, err := fs.ReadFile(files, service+"/"+entry.Name())
		if err != nil || len(data) == 0 {
			return nil, errors.New("cannot read nonempty migration")
		}
		hash := sha256.Sum256(data)
		result = append(result, migration{version, entry.Name(), hex.EncodeToString(hash[:]), string(data)})
		previous = version
	}
	if len(result) == 0 {
		return nil, errors.New("service has no migrations")
	}
	return result, nil
}

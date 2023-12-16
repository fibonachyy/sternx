package repository

import (
	"context"
	"os"
	"path/filepath"
)

func (p postgres) Migrate(path string) error {

	rootDir, err := os.Getwd()
	if err != nil {
		return err
	}
	migrationDir := filepath.Join(rootDir, path)

	files, err := os.ReadDir(migrationDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".sql" {
			migrationSQL, err := os.ReadFile(filepath.Join(migrationDir, file.Name()))
			if err != nil {
				return err
			}

			_, err = p.conn.Exec(context.Background(), string(migrationSQL))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

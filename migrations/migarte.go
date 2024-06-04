// migrations/migrate.go

package migrations

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func RunMigrations(db *sql.DB, migrationsPath string) error {
	files, err := os.ReadDir(migrationsPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".sql") {
			filePath := filepath.Join(migrationsPath, file.Name())
			content, err := os.ReadFile(filePath)
			if err != nil {
				return err
			}

			_, err = db.Exec(string(content))
			if err != nil {
				return err
			}

			log.Printf("Миграции выполнены по файлу: %s", file.Name())
		}
	}

	return nil
}

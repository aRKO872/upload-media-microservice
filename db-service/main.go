package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	// db "github.com/upload-media-db/database"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq" // PostgreSQL driver
	db "github.com/upload-media-db/database"
)

func applyMigrationsHandler(c echo.Context) error {
	DBInstance := db.GetDBInstance()

	files, err := os.ReadDir("migrations")
	if err != nil {
		log.Println("Error reading migration files:", err)
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".sql") {
			migrationPath := filepath.Join("migrations", file.Name())
			migrationContent, err := os.ReadFile(migrationPath)
			if err != nil {
				log.Println("Error reading migration content:", err)
				return c.String(http.StatusInternalServerError, "Internal Server Error")
			}

			// Execute the migration SQL
			_, err = DBInstance.Exec(string(migrationContent))
			if err != nil {
				log.Printf("Error applying migration %s: %v\n", file.Name(), err)
				return c.String(http.StatusInternalServerError, "Internal Server Error")
			}

			log.Printf("Applied migration: %s\n", file.Name())
		}
	}

	return c.String(http.StatusOK, "Migrations applied successfully!")
}

func main () {
	e := echo.New()

	db.CreateDatabase()
	e.GET("/apply-migrations", applyMigrationsHandler)

	e.Start(":8086")
}

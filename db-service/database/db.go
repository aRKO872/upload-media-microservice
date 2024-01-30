// create_database.go

package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // Import PostgreSQL driver
	"github.com/upload-media-db/config"
)

var (
	DBInstance     *sql.DB
	DBConnectError error
)

func CreateDatabase() {
	// Connect to the existing database
	user := config.GetConfig().DATABASE_USER
	password := config.GetConfig().DATABASE_PASSWORD
	dbname := config.GetConfig().DATABASE_NAME
	
	connStr := fmt.Sprintf("host=postgres port=5432 user=%s password=%s dbname=%s sslmode=disable", user, password, dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("Unable to open")
		panic(err)
	}

	_, err = db.Exec(`
	DO $$ 
	BEGIN
			IF NOT EXISTS (SELECT FROM pg_database WHERE datname = 'media-db') THEN
					CREATE DATABASE "media-db";
			END IF;
	END $$;`)
	
	if err != nil {
		fmt.Println("Unable to create")
		panic(err)
	}
	
	DBInstance = db
	fmt.Println("Database 'media-db' created successfully!")
}

func GetDBInstance() *sql.DB {
	return DBInstance
}


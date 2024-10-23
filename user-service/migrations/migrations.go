package migrations

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
)

func readSQLFile(filepath string) (string, error) {
	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
func SetUpMigrating(action string, db *sql.DB) {
	var sqlFile string
	var err error

	if action == "up" {
		sqlFile, err = readSQLFile("migrations/user_service_up.sql")
		if err != nil {
			log.Fatalf("Failed to read up migration file: %v", err)
		}
	} else if action == "down" {
		sqlFile, err = readSQLFile("migrations/user_service_down.sql")
		if err != nil {
			log.Fatalf("Failed to read down migration file: %v", err)
		}
	} else {
		log.Println("no action")
		return
	}

	_, err = db.Exec(sqlFile)
	if err != nil {
		log.Fatalf("Failed to execute migration: %v", err)
	}

	fmt.Printf("Migration '%s' executed successfully\n", action)
}

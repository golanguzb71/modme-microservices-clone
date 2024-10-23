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
	//var sqlFile string
	//var err error
	//
	//if action == "up" {
	//	sqlFile, err = readSQLFile("migrations/user_service_up.sql")
	//	if err != nil {
	//		log.Fatalf("Failed to read up migration file: %v", err)
	//	}
	//} else if action == "down" {
	//	sqlFile, err = readSQLFile("migrations/user_service_down.sql")
	//	if err != nil {
	//		log.Fatalf("Failed to read down migration file: %v", err)
	//	}
	//} else {
	//	log.Println("no action")
	//	return
	//}

	_, err := db.Exec(`
CREATE TABLE IF NOT EXISTS users
(
    id           uuid primary key,
    full_name    varchar                                                          NOT NULL,
    phone_number varchar UNIQUE                                                   NOT NULL,
    password     varchar                                                          NOT NULL,
    role         varchar check ( role in ('CEO', 'TEACHER', 'ADMIN', 'EMPLOYEE')) NOT NULL,
    birth_date   date                                                             NOT NULL,
    gender       boolean                                                          NOT NULL DEFAULT TRUE,
    is_deleted   boolean                                                          NOT NULL DEFAULT FALSE,
    created_at   timestamp                                                                 DEFAULT NOW()
);


CREATE TABLE IF NOT EXISTS users_history
(
    id            uuid primary key,
    user_id       uuid references users (id),
    updated_field varchar   NOT NULL,
    old_value     varchar   NOT NULL,
    current_value varchar   NOT NULL,
    created_at    timestamp NOT NULL DEFAULT NOW()
);
`)
	if err != nil {
		log.Fatalf("Failed to execute migration: %v", err)
	}

	fmt.Printf("Migration '%s' executed successfully\n", action)
}

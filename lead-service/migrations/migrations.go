package migrations

import (
	"database/sql"
	"io/ioutil"
)

func readSQLFile(filepath string) (string, error) {
	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func SetUpMigrating(action string, db *sql.DB) {

}

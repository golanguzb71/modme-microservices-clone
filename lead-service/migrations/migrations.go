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
	// var sqlFile string
	// var err error

	// if action == "up" {
	// 	sqlFile, err = readSQLFile("migrations/lid_service_up.sql")
	// 	if err != nil {
	// 		log.Fatalf("Failed to read up migration file: %v", err)
	// 	}
	// } else if action == "down" {
	// 	sqlFile, err = readSQLFile("migrations/lid_service_down.sql")
	// 	if err != nil {
	// 		log.Fatalf("Failed to read down migration file: %v", err)
	// 	}
	// } else {
	// 	log.Println("no action")
	// 	return
	// }

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS lead_section
(
    id         serial PRIMARY KEY,
    title      varchar NOT NULL UNIQUE,
    created_at timestamp DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS expect_section
(
    id         serial UNIQUE,
    title      varchar NOT NULL UNIQUE,
    created_at timestamp DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS set_section
(
    id         serial PRIMARY KEY,
    title      varchar                                               NOT NULL,
    course_id  int                                                   NOT NULL,
    teacher_id uuid                                                  NOT NULL,
    date_type  varchar check (date_type in ('JUFT', 'TOQ', 'OTHER')) NOT NULL,
    days       TEXT[]                                                NOT NULL,
    start_time varchar                                               NOT NULL,
    created_at timestamp DEFAULT NOW(),
    CONSTRAINT valid_days CHECK (array_length(days, 1) > 0 AND days <@
                                                               ARRAY ['DUSHANBA', 'SESHANBA', 'CHORSHANBA', 'PAYSHANBA', 'JUMA', 'SHANBA', 'YAKSHANBA'])
);

CREATE TABLE IF NOT EXISTS lead_user
(
    id           serial PRIMARY KEY,
    phone_number varchar NOT NULL,
    full_name     varchar NOT NULL,
    lead_id      int references lead_section (id),
    expect_id    int references expect_section (id),
    set_id       int references set_section (id),
    comment      varchar,
    created_at   timestamp DEFAULT now()
);		 
    `)
	if err != nil {
		log.Fatalf("Failed to execute migration: %v", err)
	}

	fmt.Printf("Migration '%s' executed successfully\n", action)
}

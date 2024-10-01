package utils

import (
	"database/sql"
	"fmt"
	"log"
)

func MigrateUp(db *sql.DB) {
	//filePath := "./migrations/user_service_up.sql"
	//sqlFile, err := os.Open(filePath)
	//if err != nil {
	//	log.Fatalf("Error opening SQL migration file: %s", err)
	//}
	//defer sqlFile.Close()
	//sqlContent, err := ioutil.ReadAll(sqlFile)
	//if err != nil {
	//	log.Fatalf("error while migrating %v", err.Error())
	//}
	sqlContent := `CREATE TABLE IF NOT EXISTS users
(
    id           serial primary key,
    name         varchar NOT NULL,
    surname      varchar NOT NULL DEFAULT '********',
    phone_number varchar NOT NULL UNIQUE,
    chat_id      float NOT NULL UNIQUE,
    role         varchar check ( role in ('ADMIN', 'USER') ),
    created_at timestamp DEFAULT now()
);

INSERT INTO users(name, phone_number, chat_id, role) values ('Abdulaziz' , '+998950960153' , 7777 , 'ADMIN');
`
	_, err := db.Exec(sqlContent)
	if err != nil {
		log.Fatalf("Error executing SQL migration: %s", err)
	}

	fmt.Println("Database migration ran successfully!")
}

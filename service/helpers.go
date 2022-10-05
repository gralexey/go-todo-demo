package service

import (
	"database/sql"
	"log"
)

func PrepareDb(db *sql.DB) {
	sqlStmt := `
	create table if not exists users (
		id integer not null primary key,
		user_token text, 
		is_admin integer
	);
	`
	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	sqlStmt = `
	create table if not exists todos (
		id integer not null primary key, 
		content text,
		user_id integer,
		foreign key(user_id) references users(id)
	);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
}

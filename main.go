package main

import (
	"database/sql"
	"log"
	"net/http"
	"testcode/test3/notifier"
	"testcode/test3/service"
)

func main() {
	log.Println("Starting...")

	db, err := sql.Open("sqlite3", "./todo.db")
	if err != nil {
		log.Fatalln("can't open db", err)
		panic(err)
	}
	defer db.Close()

	service.PrepareDb(db)

	s := service.Service{Db: db, Notifier: notifier.NotifierService{}}

	http.Handle("/todos/", &s)
	http.Handle("/users/", &s)

	if err = http.ListenAndServe(":8081", nil); err != nil {
		log.Fatalln("error ListenAndServe", err)
		panic(err)
	}
}

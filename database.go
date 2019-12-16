package main

import (
	"database/sql"
	"fmt" 
  _ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "m123"
	dbname   = "store"
  )


type Email struct{
	id int
	send_to string
	from_user string
	subject string
	content string
	created_at string
}


func getCon() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
    "password=%s dbname=%s sslmode=disable",
    host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
	panic(err)
	}
	err = db.Ping()
	if err != nil {
	panic(err)
	}
	return db
}


func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
package main

import (
	"database/sql"
	"strings"
	"strconv"
  _ "github.com/lib/pq"
	"time"
)

func insertEmail(db *sql.DB, e Email) {
	var lastID int
	
	sql := `SELECT id FROM emails order by id desc limit 1`
	rows, err := db.Query(sql)
	if rows.Next(){
		err = rows.Scan(&lastID)
		checkErr(err)
		t := time.Now()
		s := strings.Split(t.Format(time.RFC3339),"T")
		e.created_at = s[0]
		sql = `INSERT INTO emails(send_to,subject,content,id,from_user,created_at) VALUES ('`+e.send_to+`', '{`+e.subject+`}', '{`+e.content+`}', `+strconv.Itoa(lastID+1)+`, '`+e.from_user+`', '`+e.created_at+`' )`
		_, err = db.Exec(sql)
		checkErr(err)
	}
}
func getEmails(db *sql.DB, email string) []Email {
  var emails []Email
  sql := `SELECT * FROM emails WHERE from_user='`+email+`'`
  rows, err := db.Query(sql)
  if err != nil {
    panic(err)
  }
  defer rows.Close()
  for rows.Next() {
    var id int
	var sub string
	var cont string
	var to string
	var from string 
	var date string
	err = rows.Scan(&to, &sub, &id, &from, &date, &cont)
	s := strings.Split(date,"T")
	date = s[0]
    if err != nil {
      panic(err)
    }
	// fmt.Println(id, from, sub, to, date, user_id)
	emails = append(emails,Email{
			id: id,
			from_user: from,
			send_to: to,
			subject: sub,
			content: cont,
			created_at: date,
		})
	
  }
  // get any error encountered during iteration
  err = rows.Err()
  if err != nil {
    panic(err)
  }
  return emails
}


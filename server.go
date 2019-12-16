package main

import (
	"log"
	"fmt"
	"net/smtp"
	"net/http"
	"github.com/bytbox/go-pop3"
	"strings"
	"time"
	"database/sql"
	"sync"
	"strconv"
	"bufio"
	"os"
)

var db *sql.DB
const address = "pop.gmail.com:995"
func main() {
	db = getCon()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		fmt.Println("Server is running...")
		wg.Done()
		http.HandleFunc("/", request)
		err := http.ListenAndServe(":8080", nil) 
		if err != nil {
			log.Fatal("FatalErr: ", err)
		}
		
	}()
	wg.Wait()
	
	var op int 
	op = 0
	for op != 4 {
		fmt.Println("-----------MENU-----------")
		fmt.Println("1 - Enviar email")
		fmt.Println("2 - Ver emails recebidos")
		fmt.Println("3 - Ver emails enviados")
		fmt.Println("4 - Sair")
		fmt.Scanf("%d", &op)
		
		handle(op)
	}
	
	defer db.Close()
}


func handle(option int){
	switch option {
	case 01: 
			sendEmail()
			break;
	case 02: 
			e := getReceivedEmails()
			fmt.Println("Emails Recebidos:")
			for _, m := range e {
				fmt.Println("-------------------------------------------")
				fmt.Println(m)
			}
			
			break;
	case 03:
			e := getSendedEmails()
			fmt.Println("Emails Enviados:")
			for _, m := range e {
				fmt.Println("-------------------------------------------")
				fmt.Println("To:"+m.send_to)
				fmt.Println("Subject:"+m.subject+" | Date:"+m.created_at)
				fmt.Println("Body:"+m.content)
			}
			break;
	case 04:
			fmt.Println("Saindo...")
			break;
	default:
			errorMsg()
			break;
		
	}
}


func getClient(adrs string, user string, pass string) (*pop3.Client){
	c, err := pop3.DialTLS(adrs)
	checkErr(err)
	
	err = c.Auth(user,pass)
	checkErr(err)

	return c
}


func sendEmail(){
	var useremail string
	var userpass string
	var receivers []string
	var subject string
	var body string
	var aux string
	var num int
	bio := bufio.NewReader(os.Stdin)
	
	fmt.Println("Email do Usuário:")
	fmt.Scanf("%s", &useremail)
	fmt.Println("Senha:")
	fmt.Scanf("%s", &userpass)
	fmt.Println("Quantos destinatários:")
	fmt.Scanf("%d", &num)
	for i := 1; i <= num; i++ {
		fmt.Printf("Destinatário %d:", i)
		fmt.Scanf("%s", &aux)
		receivers = append(receivers, aux)
	}
	fmt.Println("Assunto:")
	line, _, _ := bio.ReadLine()
	subject = string(line)
	fmt.Println("Corpo do email:")
	line, _, _ = bio.ReadLine()
	body = string(line)

	send(useremail, userpass, receivers, subject, body)
}

func getReceivedEmails() []string{
	var useremail string
	var userpass string


	fmt.Println("Email para busca:")
	fmt.Scanf("%s", &useremail)
	fmt.Println("Senha:")
	fmt.Scanf("%s", &userpass)
	client := getClient(address, useremail, userpass)

	msg,_, err := client.ListAll()
	checkErr(err)
	var msgs []string
	for _,m := range msg {
		s, err := client.Retr(m)
		checkErr(err)
		msgs = append(msgs, s)
	}
	return msgs
}

func getSendedEmails() []Email{
	var useremail string

	fmt.Println("Email para busca:")
	fmt.Scanf("%s", &useremail)

	emails := getEmails(db, useremail)
	return emails
}

func send(from string, pass string, to []string, sub string, body string) {
	
	date := strings.Split(time.Now().String()," ")[0]
	toHeader := strings.Join(to, ",")

	e := Email{
		id: 1,
		from_user: from,
		send_to: toHeader,
		subject: sub,
		content: body,
		created_at: date,
	}

	insertEmail(db, e)

	msg := "From: " + from + "\n" +
		"To: " + toHeader + "\n" +
		"Subject:" + sub + "\n\n" +
		body

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, to, []byte(msg))

	checkErr(err)
	fmt.Println("OK - 200")
}


func errorMsg(){
	fmt.Println("Error occur!!!")
	panic("Err")
}

func request(w http.ResponseWriter, r *http.Request){
	if (r.Method == "POST" || r.Method == "post" || r.Method == "Post"){
		w.Write([]byte("POST REQUEST MADE\n"))
		w.Write([]byte("Sending email...\n"))

		user := r.FormValue("user")
		pass := r.FormValue("pass")
		to := r.FormValue("to")
		sub := r.FormValue("sub")
		cont := r.FormValue("cont")
		var valid bool
		valid = true
		if (user == ""){
			valid = false
			w.Write([]byte("Argument 'user' is missing! - Bad Request(400)\n"))
		}
		if (pass == ""){
			valid = false
			w.Write([]byte("Argument 'pass' is missing! - Bad Request(400)\n"))
		}
		if (to == ""){
			valid = false
			w.Write([]byte("Argument 'to' is missing! - Bad Request(400)\n"))
		}
		if (sub == ""){
			valid = false
			w.Write([]byte("Argument 'sub' is missing! - Bad Request(400)\n"))
		}
		if (cont == ""){
			valid = false
			w.Write([]byte("Argument 'cont' is missing! - Bad Request(400)\n"))
		}
		//TUDO Ok
		if	(valid){
			send(user, pass, []string{to}, sub, cont)
			w.Write([]byte("OK - 200\n"))
		}

	}else
	if (r.Method == "GET" || r.Method == "get" || r.Method == "Get"){
		w.Write([]byte("GET REQUEST MADE\n"))
		var valid bool 
		valid = true
		user := ""
		pass := ""
		option := ""
		s := strings.Split(r.URL.String(), "/")
		if (len(s) < 4){
			valid = false
			w.Write([]byte("Wrong number of arguments. Excepted 3, given "+strconv.Itoa(len(s)-1)+" - Bad Request(400)\n"))
		}else{
			user = s[1]
			pass = s[2]
			option = s[3]
		}
		if (user == ""){
			valid = false
			w.Write([]byte("Argument 'user' is missing! - Bad Request(400)\n"))
		} 
		if (pass == ""){
			valid = false
			w.Write([]byte("Argument 'pass' is missing! - Bad Request(400)\n"))
		} 
		if (option == ""){
			valid = false
			w.Write([]byte("Argument 'type' is missing! - Bad Request(400)\n"))
		} 
		if (valid) {
			
			client := getClient(address, user, pass)
			if option == "SENT"{
				
				emails := getEmails(db, user)
				
				for _,e := range emails {
					
					w.Write([]byte("TO:"+e.send_to+" "+
						"Subject:"+e.subject+" "+
						"Date:"+e.created_at+" "+
						"Body"+e.content+" "+
						"\n"))
				}
				w.Write([]byte("OK - 200\n"))
			}else if option == "RECEIVED"{
				
				msg,_, err := client.ListAll()
				checkErr(err)
				for _,m := range msg {
					s, err := client.Retr(m)
					checkErr(err)
					w.Write([]byte("--------------------------------------------\n"))	
					w.Write([]byte(s+"\n"))									
				}
				w.Write([]byte("OK - 200\n"))
			}else{
				w.Write([]byte("Argument 'type' is invalid! - Internal Server Error(500)\n"))
			}
		
		}
		
	}
}
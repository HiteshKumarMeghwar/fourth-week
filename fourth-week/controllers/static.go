package controllers

import (
	"fmt"
	"fourth-week/cmd/database"
	"html/template"
	"net/http"

	"github.com/gorilla/sessions"
)

// var db = database.Connect()
var store = sessions.NewCookieStore([]byte("super-secret"))

func StaticHandler(tpl Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// fmt.Println(tpl.HTMLtpl.Name())
		// fmt.Println(r.URL.Path)
		// tpl.Execute(w, nil)

		/* Requiring Database */
		db := database.Connect()
		defer db.Close()

		if r.URL.Path == "/" {
			session, _ := store.Get(r, "session")
			id, ok := session.Values["userId"]
			fmt.Println("ok: ", ok)
			if !ok {
				http.Redirect(w, r, "/login", http.StatusFound) // http.StatusFound is 302
				return
			}
			fmt.Println(id)
			var fulname string
			var username string

			err := db.QueryRow("SELECT name, username FROM users WHERE id = $1", id).Scan(&fulname, &username)
			if err != nil {
				panic(err)
			}
			// var data map[string] string
			/* fmt.Println(fulname)
			fmt.Println(username) */
			tpl.Execute(w, nil)
			return
		} else if r.URL.Path == "/login" {
			session, _ := store.Get(r, "session")
			_, ok := session.Values["userId"]
			fmt.Println("ok: ", ok)
			if ok {
				http.Redirect(w, r, "/", http.StatusFound) // http.StatusFound is 302
				return
			}
		}
		tpl.Execute(w, nil)
	}
}

func FAQ(tpl Template) http.HandlerFunc {
	questions := []struct {
		Question string
		Answer   template.HTML
	}{
		{
			Question: "Is there a free version?",
			Answer:   "yes! We offer a free trail for 30 days.",
		}, {
			Question: "What are your support hours?",
			Answer:   "We have support staff answering emails .....",
		}, {
			Question: "How do I contact support?",
			Answer:   `Email us - <a href="/email.com" > Email Send </a>`,
		},
	}
	return func(w http.ResponseWriter, r *http.Request) {
		tpl.Execute(w, questions)
	}
}

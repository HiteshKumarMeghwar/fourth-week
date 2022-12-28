package controllers

import (
	"fmt"
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
		if r.URL.Path == "/" {
			session, _ := store.Get(r, "session")
			_, ok := session.Values["userId"]
			fmt.Println("ok: ", ok)
			if !ok {
				http.Redirect(w, r, "/login", http.StatusFound) // http.StatusFound is 302
				return
			}
			tpl.Execute(w, nil)
		} else if r.URL.Path == "/login" {
			session, _ := store.Get(r, "session")
			_, ok := session.Values["userId"]
			fmt.Println("ok: ", ok)
			if ok {
				http.Redirect(w, r, "/dashboard", http.StatusFound) // http.StatusFound is 302
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

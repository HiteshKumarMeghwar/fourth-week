package controllers

import (
	"assignment/views"
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
)

// var db = database.Connect()
var store = sessions.NewCookieStore([]byte("super-secret"))

func StaticHandler(tpl views.Template) http.HandlerFunc {
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

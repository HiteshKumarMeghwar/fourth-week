package controllers

import (
	"fmt"
	"fourth-week/bcryptPassword"
	"fourth-week/cmd/database"
	"html/template"
	"net/http"

	"github.com/gorilla/sessions"
)

type Product struct {
	ID       int
	Name     string
	Username string
	Password string
}

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

			var data struct {
				Fulname  string
				Username string
			}
			data.Fulname = fulname
			data.Username = username

			cookie := http.Cookie{
				Name:     "username",
				Value:    username,
				Path:     "/",
				HttpOnly: true,
			}
			http.SetCookie(w, &cookie)
			// SetCookie(w, CookieSession, session.Token)
			tpl.Execute(w, data)
			return
		} else if r.URL.Path == "/login" {
			session, _ := store.Get(r, "session")
			_, ok := session.Values["userId"]
			fmt.Println("ok: ", ok)
			if ok {
				username, _ := r.Cookie("username")
				fmt.Println("Printing Username: ", username.Value)
				http.Redirect(w, r, "/", http.StatusFound) // http.StatusFound is 302
				return
			}
		} else if r.URL.Path == "/all_users" {
			session, _ := store.Get(r, "session")
			id, ok := session.Values["userId"]
			fmt.Println("ok: ", ok)
			if !ok {
				http.Redirect(w, r, "/login", http.StatusFound) // http.StatusFound is 302
				return
			}
			fmt.Println(id)

			rows, err := db.Query("SELECT * FROM users")
			if err != nil {
				panic(err)
			}
			defer rows.Close()
			var products []Product
			for rows.Next() {
				var p Product
				err = rows.Scan(&p.ID, &p.Name, &p.Username, &p.Password)
				if err != nil {
					panic(err)
				}
				products = append(products, p)
			}
			// fmt.Println(products)
			tpl.Execute(w, products)
			return
		} else if r.URL.Path == "/edit" {
			session, _ := store.Get(r, "session")
			_, ok := session.Values["userId"]
			fmt.Println("ok: ", ok)
			if !ok {
				http.Redirect(w, r, "/login", http.StatusFound) // http.StatusFound is 302
				return
			}
			r.ParseForm()
			id := r.FormValue("id")
			fmt.Println(id)

			var data struct {
				ID       int
				Name     string
				Username string
				Password string
			}

			var index int
			var name string
			var username string
			var password string

			err := db.QueryRow("SELECT * FROM users WHERE id = $1", id).Scan(&index, &name, &username, &password)
			if err != nil {
				panic(err)
			}

			data.ID = index
			data.Name = name
			data.Username = username
			data.Password = password

			tpl.Execute(w, data)
			return
		} else if r.URL.Path == "/update" {
			session, _ := store.Get(r, "session")
			_, ok := session.Values["userId"]
			fmt.Println("ok: ", ok)
			if !ok {
				http.Redirect(w, r, "/login", http.StatusFound) // http.StatusFound is 302
				return
			}
			r.ParseForm()
			id := r.FormValue("id")
			name := r.FormValue("name")
			username := r.FormValue("username")
			password := r.FormValue("password")
			hash, _ := bcryptPassword.HashPassword(password)

			value, err := db.Exec(`UPDATE users SET name = $1, username = $2, password = $3 WHERE id = $4`, name, username, hash, id)
			if err != nil {
				panic(err)
			}

			if value != nil {
				http.Redirect(w, r, "/all_users", http.StatusFound) // http.StatusFound is 302
				return
			}
		} else if r.URL.Path == "/delete" {
			session, _ := store.Get(r, "session")
			_, ok := session.Values["userId"]
			fmt.Println("ok: ", ok)
			if !ok {
				http.Redirect(w, r, "/login", http.StatusFound) // http.StatusFound is 302
				return
			}

			r.ParseForm()
			id := r.FormValue("id")

			value, err := db.Exec(`DELETE FROM users WHERE id=$1;`, id)
			if err != nil {
				panic(err)
			}

			if value != nil {
				http.Redirect(w, r, "/all_users", http.StatusFound) // http.StatusFound is 302
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

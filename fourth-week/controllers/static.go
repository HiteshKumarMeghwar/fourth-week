package controllers

import (
	"fmt"
	"fourth-week/bcryptPassword"
	"fourth-week/cmd/database"
	"html/template"
	"net/http"
	"strings"

	"github.com/gorilla/sessions"
)

type Product struct {
	ID         int
	Name       string
	Username   string
	Password   string
	RoleID     int
	Session_id interface{}
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
			var role_id int

			err := db.QueryRow("SELECT name, username, role_id FROM users WHERE id = $1", id).Scan(&fulname, &username, &role_id)
			if err != nil {
				panic(err)
			}

			var data struct {
				Fulname  string
				Username string
				RoleID   int
			}
			data.Fulname = fulname
			data.Username = username
			data.RoleID = role_id

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
			// fmt.Println(id)

			rows, err := db.Query("SELECT * FROM users")
			if err != nil {
				panic(err)
			}
			defer rows.Close()

			var products []Product
			// products = append(products, Product{Session_id: id})
			for rows.Next() {
				var p Product
				err = rows.Scan(&p.ID, &p.Name, &p.Username, &p.Password, &p.RoleID)
				if err != nil {
					panic(err)
				}
				p.Session_id = id
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
				RoleID   int
			}

			var index int
			var name string
			var username string
			var password string
			var role_id int

			err := db.QueryRow("SELECT * FROM users WHERE id = $1", id).Scan(&index, &name, &username, &password, &role_id)
			if err != nil {
				panic(err)
			}

			data.ID = index
			data.Name = name
			data.Username = username
			data.Password = password
			data.RoleID = role_id

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
		} else if r.URL.Path == "/loginAuth" {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			r.ParseForm()
			username := r.FormValue("username")
			password := r.FormValue("password")
			fmt.Println("Username: ", username, "Password: ", password)

			/* tpl, err := template.ParseFiles(filepath.Join("templates", "login.gohtml"))

			if err != nil {
				panic(err)
			} */

			var ValidateError struct {
				Username string
				Password string
			}

			if len(strings.TrimSpace(username)) == 0 || len(strings.TrimSpace(password)) == 0 {
				if len(strings.TrimSpace(username)) == 0 {
					ValidateError.Username = "Username is Mandate"
				}
				if len(strings.TrimSpace(password)) == 0 {
					ValidateError.Password = "Password is Mandate"
				}

				tpl.Execute(w, ValidateError)
				return
			}

			var userId int
			var pass string
			// stmt := "SELECT id, password FROM users WHERE username = ?"
			// row := db.QueryRow(stmt, username)
			// err = row.Scan(&userId, &pass)
			err := db.QueryRow("SELECT id, password FROM users WHERE username = $1", username).Scan(&userId, &pass)
			fmt.Println("hash from db: ", pass)
			if err != nil {
				fmt.Println("error selecting Hash in db by Username")
				http.Redirect(w, r, "/login", http.StatusFound) // http.StatusFound is 302
				return
			}

			match := bcryptPassword.CheckPasswordHash(password, pass)
			fmt.Println(match)
			if match {
				session, _ := store.Get(r, "session")
				session.Values["userId"] = userId
				session.Save(r, w)
				http.Redirect(w, r, "/", http.StatusFound) // http.StatusFound is 302
				return
			}

			fmt.Println("incorrect password")
			http.Redirect(w, r, "/login", http.StatusFound) // http.StatusFound is 302
		} else if r.URL.Path == "/register-process" {
			err := r.ParseForm()
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			fullname := r.FormValue("name")
			username := r.FormValue("username")
			password := r.FormValue("password")

			var ValidateError struct {
				Submitted string
				Fullname  string
				Username  string
				Password  string
			}

			if len(strings.TrimSpace(fullname)) == 0 || len(strings.TrimSpace(username)) == 0 || len(strings.TrimSpace(password)) == 0 {
				if len(strings.TrimSpace(fullname)) == 0 {
					ValidateError.Fullname = "Fullname is Mandate"
				}
				if len(strings.TrimSpace(username)) == 0 {
					ValidateError.Username = "Username is Mandate"
				}
				if len(strings.TrimSpace(password)) == 0 {
					ValidateError.Password = "Password is Mandate"
				}

				tpl.Execute(w, ValidateError)
				return
			}

			fmt.Println(password)
			hash, _ := bcryptPassword.HashPassword(password)
			fmt.Println(hash)

			value, err := db.Exec(`INSERT INTO users(name, username, password, role_id) VALUES ($1, $2, $3, $4); `, fullname, username, hash, 3)

			if err != nil {
				panic(err)
			}

			if value != nil {
				http.Redirect(w, r, "/login", http.StatusFound) // http.StatusFound is 302
				return
			}
		} else if r.URL.Path == "/logout" {
			fmt.Println("logouting .........! ")
			session, _ := store.Get(r, "session")
			delete(session.Values, "userId")
			session.Save(r, w)
			http.Redirect(w, r, "/login", http.StatusFound) // http.StatusFound is 302
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

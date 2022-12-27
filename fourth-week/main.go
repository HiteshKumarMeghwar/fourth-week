package main

import (
	"assignment/cmd/database"
	"assignment/controllers"
	"assignment/views"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	/* Requiring Database */
	db := database.Connect()
	defer db.Close()

	/* Initialization of Session */
	var store = sessions.NewCookieStore([]byte("super-secret"))

	/* _, err := db.Exec(`

	CREATE TABLE IF NOT EXISTS users(
		id SERIAL PRIMARY KEY,
		name TEXT,
		username TEXT UNIQUE NOT NULL,
		password TEXT
	)`)

	if err != nil {
		panic(err)
	} */

	route := chi.NewRouter()

	route.Get("/", controllers.StaticHandler(views.Must(views.Parse(filepath.Join("templates", "dashboard.gohtml")))))

	route.Get("/login", controllers.StaticHandler(views.Must(views.Parse(filepath.Join("templates", "login.gohtml")))))

	route.Post("/loginAuth", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		r.ParseForm()
		username := r.FormValue("username")
		password := r.FormValue("password")
		fmt.Println("Username: ", username, "Password: ", password)

		tpl, err := template.ParseFiles(filepath.Join("templates", "login.gohtml"))
		if err != nil {
			panic(err)
		}

		var userId int
		var pass string
		// stmt := "SELECT id, password FROM users WHERE username = ?"
		// row := db.QueryRow(stmt, username)
		// err = row.Scan(&userId, &pass)
		err = db.QueryRow("SELECT id, password FROM users WHERE username = ?", username).Scan(&userId, &pass)
		fmt.Println("hash from db: ", pass)
		if err != nil {
			fmt.Println("error selecting Hash in db by Username")
			tpl.Execute(w, nil)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(pass), []byte(password))
		if err == nil {
			tpl, err := template.ParseFiles(filepath.Join("templates", "dashboard.gohtml"))
			if err != nil {
				panic(err)
			}
			session, _ := store.Get(r, "session")
			session.Values["userId"] = userId
			session.Save(r, w)
			tpl.Execute(w, nil)
			return
		}

		fmt.Println("incorrect password")
		tpl.Execute(w, nil)
	})

	route.Get("/logout", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("logouting .........! ")
		session, _ := store.Get(r, "session")
		delete(session.Values, "userId")
		session.Save(r, w)
		tpl, err := template.ParseFiles(filepath.Join("templates", "login.gohtml"))
		if err != nil {
			panic(err)
		}
		tpl.Execute(w, nil)
	})

	route.Get("/register", controllers.StaticHandler(views.Must(views.Parse(filepath.Join("templates", "register.gohtml")))))

	route.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page Not Found...!", http.StatusNotFound)
	})
	fmt.Println("Stating the server on :8080 port ... !")
	http.ListenAndServe(":8080", route)
}

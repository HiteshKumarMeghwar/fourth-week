package main

import (
	"fmt"
	"fourth-week/bcryptPassword"
	"fourth-week/cmd/database"
	"fourth-week/controllers"
	"fourth-week/migrations"
	"fourth-week/views"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/gorilla/sessions"
)

func main() {
	/* Requiring Database */
	db := database.Connect()
	defer db.Close()
	err := database.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	}
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
	// route.Get("/login", controllers.StaticHandler(views.Must(views.ParseFS(templates.FS, "templates/login.gohtml"))))

	route.Post("/loginAuth", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		r.ParseForm()
		username := r.FormValue("username")
		password := r.FormValue("password")
		fmt.Println("Username: ", username, "Password: ", password)

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
	})

	route.Get("/logout", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("logouting .........! ")
		session, _ := store.Get(r, "session")
		delete(session.Values, "userId")
		session.Save(r, w)
		http.Redirect(w, r, "/login", http.StatusFound) // http.StatusFound is 302
		return
	})

	usersC := controllers.Users{}
	usersC.Templates.New = views.Must(views.Parse(filepath.Join("templates", "register.gohtml")))
	route.Get("/register", usersC.New)
	route.Post("/register-process", usersC.Create)
	route.Get("/faq", controllers.FAQ(views.Must(views.Parse(filepath.Join("templates", "faq.gohtml")))))
	// route.Get("/home", controllers.StaticHandler(views.Must(views.Parse(filepath.Join("templates", "home.gohtml")))))

	route.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page Not Found...!", http.StatusNotFound)
	})
	fmt.Println("Stating the server on :8080 port ... !")
	http.ListenAndServe(":8080", route)
}

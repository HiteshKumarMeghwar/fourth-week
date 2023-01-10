package main

import (
	"fmt"
	"fourth-week/cmd/database"
	"fourth-week/controllers"
	"fourth-week/views"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
)

func main() {

	database.LoadEnvVariables()
	/* Requiring Database */
	db := database.Connect()
	defer db.Close()

	route := chi.NewRouter()

	route.Get("/", controllers.StaticHandler(views.Must(views.Parse(filepath.Join("templates", "dashboard.gohtml")))))

	route.Get("/login", controllers.StaticHandler(views.Must(views.Parse(filepath.Join("templates", "login.gohtml")))))
	route.Get("/register", controllers.StaticHandler(views.Must(views.Parse(filepath.Join("templates", "register.gohtml")))))
	route.Get("/all_users", controllers.StaticHandler(views.Must(views.Parse(filepath.Join("templates", "users.gohtml")))))
	route.Get("/edit", controllers.StaticHandler(views.Must(views.Parse(filepath.Join("templates", "edit.gohtml")))))
	route.Post("/update", controllers.StaticHandler(views.Must(views.Parse(filepath.Join("templates", "users.gohtml")))))
	route.Get("/delete", controllers.StaticHandler(views.Must(views.Parse(filepath.Join("templates", "users.gohtml")))))
	route.Post("/loginAuth", controllers.StaticHandler(views.Must(views.Parse(filepath.Join("templates", "login.gohtml")))))
	route.Post("/register-process", controllers.StaticHandler(views.Must(views.Parse(filepath.Join("templates", "register.gohtml")))))
	route.Get("/logout", controllers.StaticHandler(views.Must(views.Parse(filepath.Join("templates", "dashboard.gohtml")))))
	route.Get("/faq", controllers.FAQ(views.Must(views.Parse(filepath.Join("templates", "faq.gohtml")))))
	route.Get("/posts", controllers.HandlePosts(views.Must(views.Parse(filepath.Join("templates", "posts.gohtml")))))
	route.Get("/create_post", controllers.HandlePosts(views.Must(views.Parse(filepath.Join("templates", "create_post.gohtml")))))
	route.Post("/create_post_process", controllers.HandlePosts(views.Must(views.Parse(filepath.Join("templates", "create_post.gohtml")))))
	route.Get("/edit_post", controllers.HandlePosts(views.Must(views.Parse(filepath.Join("templates", "edit_post.gohtml")))))
	route.Get("/update_post", controllers.HandlePosts(views.Must(views.Parse(filepath.Join("templates", "posts.gohtml")))))
	route.Get("/delete_post", controllers.HandlePosts(views.Must(views.Parse(filepath.Join("templates", "posts.gohtml")))))
	route.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page Not Found...!", http.StatusNotFound)
	})

	fmt.Println("Stating the server on :8080 port ... !")
	http.ListenAndServe(":8080", route)

	/* ................................   Non workable code   .............................. */

	// usersC := controllers.Users{}
	// usersC.Templates.New = views.Must(views.Parse(filepath.Join("templates", "register.gohtml")))
	// route.Get("/register", usersC.New)
	// route.Post("/register-process", usersC.Create)

	/* err := database.MigrateFS(db, migrations.FS, ".")
	if err != nil {
		panic(err)
	} */
	/* Initialization of Session */
	// var store = sessions.NewCookieStore([]byte("super-secret"))

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
}

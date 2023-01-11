package controllers

import (
	"fmt"
	"fourth-week/cmd/database"
	"net/http"
	"strings"

	"github.com/gorilla/sessions"
)

type Post struct {
	ID          int
	Title       string
	Summary     string
	Description string
	Session_id  interface{}
	RoleID      int
	UserID      int
}

func HandlePosts(tpl Template) http.HandlerFunc {
	var store = sessions.NewCookieStore([]byte("super-secret"))
	return func(w http.ResponseWriter, r *http.Request) {
		// database connection on ........
		db := database.Connect()
		defer db.Close()

		session, _ := store.Get(r, "session")
		id, ok := session.Values["userId"]
		fmt.Println("ok: ", ok)
		if !ok {
			http.Redirect(w, r, "/login", http.StatusFound) // http.StatusFound is 302
			return
		}
		fmt.Println(id)

		if r.URL.Path == "/posts" {
			session, _ := store.Get(r, "session")
			id, ok := session.Values["userId"]
			fmt.Println("ok: ", ok)
			if !ok {
				http.Redirect(w, r, "/login", http.StatusFound) // http.StatusFound is 302
				return
			}

			var role_id int
			err := db.QueryRow("SELECT role_id FROM users WHERE id = $1", id).Scan(&role_id)
			if err != nil {
				panic(err)
			}

			rows, err := db.Query("SELECT * FROM posts")
			if err != nil {
				panic(err)
			}
			defer rows.Close()

			var posts []Post
			for rows.Next() {
				var p Post
				err = rows.Scan(&p.ID, &p.Title, &p.Summary, &p.Description, &p.UserID)
				if err != nil {
					panic(err)
				}
				p.RoleID = role_id
				p.Session_id = id
				posts = append(posts, p)
			}
			// fmt.Println(posts)
			tpl.Execute(w, posts)
			return

		} else if r.URL.Path == "/create_post" {
			tpl.Execute(w, nil)
			return
		} else if r.URL.Path == "/create_post_process" {
			err := r.ParseForm()
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			title := r.FormValue("post_title")
			summary := r.FormValue("post_summary")
			description := r.FormValue("post_description")

			var ValidateError struct {
				Submitted   string
				Title       string
				Summary     string
				Description string
			}

			if len(strings.TrimSpace(title)) == 0 || len(strings.TrimSpace(summary)) == 0 || len(strings.TrimSpace(description)) == 0 {
				if len(strings.TrimSpace(title)) == 0 {
					ValidateError.Title = "Title is Mandate"
				}
				if len(strings.TrimSpace(summary)) == 0 {
					ValidateError.Summary = "Summary is Mandate"
				}
				if len(strings.TrimSpace(description)) == 0 {
					ValidateError.Description = "Description is Mandate"
				}

				tpl.Execute(w, ValidateError)
				return
			}

			value, err := db.Exec(`INSERT INTO posts(title_post, summary_post, descritption_post, user_id) VALUES ($1, $2, $3, $4); `, title, summary, description, id)

			if err != nil {
				panic(err)
			}

			if value != nil {
				http.Redirect(w, r, "/posts", http.StatusFound) // http.StatusFound is 302
				return
			}
		} else if r.URL.Path == "/edit_post" {
			r.ParseForm()
			id := r.FormValue("id")

			var data struct {
				ID          int
				Title       string
				Summary     string
				Description string
				UserID      int
			}

			var index int
			var title string
			var summary string
			var description string
			var user_id int

			err := db.QueryRow("SELECT * FROM posts WHERE id = $1", id).Scan(&index, &title, &summary, &description, &user_id)
			if err != nil {
				panic(err)
			}

			data.ID = index
			data.Title = title
			data.Summary = summary
			data.Description = description
			data.UserID = user_id

			tpl.Execute(w, data)
			return
		} else if r.URL.Path == "/update_post" {
			r.ParseForm()
			id := r.FormValue("id")
			title := r.FormValue("post_title")
			summary := r.FormValue("post_summary")
			description := r.FormValue("post_description")

			value, err := db.Exec(`UPDATE posts SET title_post = $1, summary_post = $2, descritption_post = $3 WHERE id = $4`, title, summary, description, id)
			if err != nil {
				panic(err)
			}

			if value != nil {
				http.Redirect(w, r, "/posts", http.StatusFound) // http.StatusFound is 302
				return
			}
		} else if r.URL.Path == "/delete_post" {
			r.ParseForm()
			id := r.FormValue("id")

			value, err := db.Exec(`DELETE FROM posts WHERE id=$1;`, id)
			if err != nil {
				panic(err)
			}

			if value != nil {
				http.Redirect(w, r, "/posts", http.StatusFound) // http.StatusFound is 302
				return
			}
		}

		tpl.Execute(w, nil)
	}
}

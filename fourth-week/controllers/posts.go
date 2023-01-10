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
}

func HandlePosts(tpl Template) http.HandlerFunc {
	var store = sessions.NewCookieStore([]byte("super-secret"))
	return func(w http.ResponseWriter, r *http.Request) {
		// database connection on ........
		db := database.Connect()
		defer db.Close()

		if r.URL.Path == "/posts" {
			session, _ := store.Get(r, "session")
			id, ok := session.Values["userId"]
			fmt.Println("ok: ", ok)
			if !ok {
				http.Redirect(w, r, "/login", http.StatusFound) // http.StatusFound is 302
				return
			}
			fmt.Println(id)

			rows, err := db.Query("SELECT * FROM posts")
			if err != nil {
				panic(err)
			}
			defer rows.Close()

			var posts []Post
			for rows.Next() {
				var p Post
				err = rows.Scan(&p.ID, &p.Title, &p.Summary, &p.Description)
				if err != nil {
					panic(err)
				}
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

			value, err := db.Exec(`INSERT INTO posts(title_post, summary_post, descritption_post) VALUES ($1, $2, $3); `, title, summary, description)

			if err != nil {
				panic(err)
			}

			if value != nil {
				http.Redirect(w, r, "/posts", http.StatusFound) // http.StatusFound is 302
				return
			}
		} else if r.URL.Path == "/edit_post" {
			tpl.Execute(w, nil)
			return
		} else if r.URL.Path == "/update_post" {
			tpl.Execute(w, nil)
			return
		} else if r.URL.Path == "/delete_post" {
			tpl.Execute(w, nil)
			return
		}

		tpl.Execute(w, nil)
	}
}

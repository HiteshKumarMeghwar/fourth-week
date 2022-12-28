package controllers

import (
	"fmt"
	"fourth-week/bcryptPassword"
	"fourth-week/cmd/database"
	"net/http"
)

type Users struct {
	Templates UsersTemplates
}

type UsersTemplates struct {
	New Template
}

func (u Users) New(w http.ResponseWriter, r *http.Request) {
	// we need to render view
	u.Templates.New.Execute(w, nil)
}

func (u Users) Create(w http.ResponseWriter, r *http.Request) {
	/* Requiring Database */
	db := database.Connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fullname := r.FormValue("name")
	username := r.FormValue("username")
	password := r.FormValue("password")
	fmt.Println(password)
	hash, _ := bcryptPassword.HashPassword(password)
	fmt.Println(hash)

	/* match := bcryptPassword.CheckPasswordHash(password, hash)
	fmt.Println(match) */

	value, err := db.Exec(`INSERT INTO users(name, username, password) VALUES ($1, $2, $3); `, fullname, username, hash)

	if err != nil {
		panic(err)
	}

	if value != nil {
		http.Redirect(w, r, "/login", http.StatusFound) // http.StatusFound is 302
		return
	}
}

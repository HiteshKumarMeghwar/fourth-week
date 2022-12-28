package controllers

import (
	"fmt"
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
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprint(w, "Username: ", r.FormValue("username"))

	fmt.Fprint(w, "Password: ", r.FormValue("password"))
}

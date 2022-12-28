package controllers

import (
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

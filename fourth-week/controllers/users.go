package controllers

import (
	"fourth-week/views"
	"net/http"
)

type Users struct {
	Templates UsersTemplates
}

type UsersTemplates struct {
	New views.Template
}

func (u Users) New(w http.ResponseWriter, r *http.Request) {
	// we need to render view
	u.Templates.New.Execute(w, nil)
}

package main

import (
	"html/template"
	"io/fs"
	"net/http"
	"path/filepath"
	"time"

	"github.com/justinas/nosurf"
	"snippetbox.fepg.org/internal/models"
	"snippetbox.fepg.org/ui"
)

type templateData struct {
	CurrentYear     int
	Snippet         *models.Snippet
	Snippets        []*models.Snippet
	Form            any
	Flash           string
	IsAuthenticated bool
	CSRFToken       string
	UserData        *models.UserData
}

var templateFunctions = template.FuncMap{
	"readableDate": readableDate,
}

func (app *application) newTemplateData(r *http.Request) *templateData {
	td := &templateData{
		CurrentYear:     time.Now().Year(),
		IsAuthenticated: false,
		// NOTE(Farid): Se necesita el CSRFToken para recibir la respuesta de los form
		CSRFToken: nosurf.Token(r),
	}

	if r != nil {
		td.Flash = app.sessionManager.PopString(r.Context(), "flash")
		td.IsAuthenticated = app.isAuthenticated(r)
	}

	return td
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	//pages, err := filepath.Glob("./ui/html/pages/*.tmpl")
	pages, err := fs.Glob(ui.Files, "html/pages/*.tmpl")
	if err != nil {
		app.errorLog.Fatalf("Couldn't glob those bitches\n")
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		patterns := []string{
			"html/base.tmpl",
			"html/partials/*.tmpl",
			page,
		}

		//app.infoLog.Printf("parsing '%s'\n", name)
		//ts, err := template.New(name).Funcs(templateFunctions).ParseFiles("./ui/html/base.tmpl")
		ts, err := template.New(name).Funcs(templateFunctions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}
		/*
			if err != nil {
				app.errorLog.Printf("'%s' template generation error\n", page)
				return nil, err
			}

			ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl")
			if err != nil {
				app.errorLog.Fatalf("'%s' partials parsing error\n", page)
				return nil, err
			}

			ts, err = ts.ParseFiles(page)
			if err != nil {
				app.errorLog.Fatalf("'%s' parsing error\n", page)
				return nil, err
			}
		*/

		cache[name] = ts
	}

	return cache, nil
}

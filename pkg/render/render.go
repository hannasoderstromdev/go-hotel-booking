package render

import (
	"bytes"
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/hannasoderstromdev/go-web-2/models"
	"github.com/hannasoderstromdev/go-web-2/pkg/config"
)

var app *config.AppConfig

// NewTemplates set new templates in app config
func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

// RenderTemplate parses a template file and renders it as HTML
func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
	var tc map[string]*template.Template

	// Check app config for cache on/off
	if app.UseCache {
		// get template from app.Config cache
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	// get template from cache
	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("could not get template from template cache")
	}

	buf := new(bytes.Buffer)
	
	td = AddDefaultData(td)

	err := t.Execute(buf, td)
	if err != nil {
		log.Println(err)
	}
	
	// render template
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println(err)
	}
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	//myCache := make(map[string]*template.Template)

	myCache := map[string]*template.Template{}

	// get all files *.page.tmpl from ./templates folder
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {
		fileName := filepath.Base(page)
		ts, err := template.New(fileName).ParseFiles(page)
		if err !=  nil {
			return myCache, err
		}

		// Get all layout files
		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err !=  nil {
			return myCache, err
		}

		// If there's a match, parse the layout file
		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err !=  nil {
				return myCache, err
			}
		}

		myCache[fileName] = ts
	}

	return myCache, nil
}
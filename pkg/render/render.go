package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/AMFSouza/bookings/pkg/config"
	"github.com/AMFSouza/bookings/pkg/models"
)


var functions = template.FuncMap{

}

var app *config.AppConfig

// NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
	var tc map[string]*template.Template
	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	// get the template cache from the app config
	//tc := app.TemplateCache

	// create a template cache
	//tc, err := CreateTemplateCache()
	// get request template from cache
	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("Could not get from template cache")
	}

	buf := new(bytes.Buffer)
	td = AddDefaultData(td)


	_ = t.Execute(buf, td)

	// render the template
	_, err := buf.WriteTo(w)
	if err != nil {
		log.Println("Error writing template to browser", err)
	}
}

// We gonna do create all templates that exists in the folder once
func CreateTemplateCache() (map[string]*template.Template, error) {
	//myCache := make(map[string]*template.Template) this time we gonna do that another way
	myCache := map[string]*template.Template{} // the result is the same above but we have a context reason to make easy to code, this time

	// get all of the files named *.page.tmpl from ./templates
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}

	// range through all files ending with *.page.tmpl
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return myCache, err
		}
		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}
		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}
		myCache[name] = ts
	}
	return myCache, nil
}


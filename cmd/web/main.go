package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/AMFSouza/bookings/pkg/config"
	"github.com/AMFSouza/bookings/pkg/handlers"
	"github.com/AMFSouza/bookings/pkg/render"
	"github.com/alexedwards/scs/v2"
)

const PORT_NUMBER = ":8080"
var app config.AppConfig
var session *scs.SessionManager


// main is the main application function
func main() {

	// change this to true when in production
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache")
	}
	app.TemplateCache = tc
	app.UseCache = false
	repo := handlers.NewRepo(&app);
	handlers.NewHandler(repo)
	render.NewTemplates(&app)

	fmt.Println(fmt.Sprintf("Starting application on port %s", PORT_NUMBER))
	//_ = http.ListenAndServe(PORT_NUMBER, nil)
	srv := &http.Server {
		Addr: PORT_NUMBER,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}
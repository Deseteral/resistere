package main

import (
	"fmt"
	"github.com/deseteral/resistere/internal/configuration"
	"log"
	"net/http"

	"github.com/a-h/templ"

	"github.com/deseteral/resistere/internal/view"
)

func startApplication() error {
	config, err := configuration.ReadConfig()
	if err != nil {
		return err
	}

	err = startWebServer(config)
	if err != nil {
		return err
	}

	return nil
}

func startWebServer(config *configuration.Config) error {
	http.Handle("/", templ.Handler(view.Index()))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	log.Printf("Starting web server on port %v.\n", config.Web.Port)

	var addr = fmt.Sprintf(":%v", config.Web.Port)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	err := startApplication()
	if err != nil {
		log.Println("Could not start application.")
		log.Fatal(err)
	}
}

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/a-h/templ"

	"github.com/deseteral/resistere/internal/view"
)

func run() error {
	fmt.Println("Starting a server")

	http.Handle("/", templ.Handler(view.Index()))

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	http.ListenAndServe(":80", nil)

	return nil
}

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

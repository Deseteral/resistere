package webapp

import (
	"fmt"
	"github.com/a-h/templ"
	"github.com/deseteral/resistere/internal/configuration"
	"github.com/deseteral/resistere/internal/webapp/view"
	"log"
	"net/http"
)

func StartWebServerBlocking(config *configuration.Config) error {
	addr := fmt.Sprintf(":%v", config.Web.Port)

	http.Handle("/", templ.Handler(view.Index()))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	log.Printf("Web server starting on port %v.\n", config.Web.Port)

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Println("Error starting web server:")
		return err
	}

	return nil
}

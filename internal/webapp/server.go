package webapp

import (
	"embed"
	"fmt"
	"github.com/a-h/templ"
	"github.com/deseteral/resistere/internal/configuration"
	"github.com/deseteral/resistere/internal/webapp/view"
	"io/fs"
	"log"
	"net/http"
)

//go:embed static
var staticFiles embed.FS

func StartWebServerBlocking(config *configuration.Config) error {
	addr := fmt.Sprintf(":%v", config.Web.Port)

	var staticFs = fs.FS(staticFiles)
	htmlContent, err := fs.Sub(staticFs, "static")
	if err != nil {
		log.Printf("Error starting web server: %v.\n", err)
		return err
	}

	http.Handle("/", templ.Handler(view.Index()))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(htmlContent))))

	log.Printf("Web server starting on port %v.\n", config.Web.Port)

	err = http.ListenAndServe(addr, nil)
	if err != nil {
		log.Printf("Error starting web server: %v.\n", err)
		return err
	}

	return nil
}

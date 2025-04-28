package webapp

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"

	"github.com/a-h/templ"
	"github.com/deseteral/resistere/internal/configuration"
	"github.com/deseteral/resistere/internal/controller"
	"github.com/deseteral/resistere/internal/webapp/view"
)

//go:embed static
var staticFiles embed.FS

func StartWebServerBlocking(config *configuration.Config, controller *controller.Controller) error {
	addr := fmt.Sprintf(":%v", config.Web.Port)

	var staticFs = fs.FS(staticFiles)
	htmlContent, err := fs.Sub(staticFs, "static")
	if err != nil {
		log.Printf("Error starting web server: %v.\n", err)
		return err
	}

	router := http.NewServeMux()

	router.Handle("GET /", templ.Handler(view.Index()))
	router.HandleFunc("POST /controller/mode/toggle", postToggleControllerMode(controller))
	router.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.FS(htmlContent))))

	log.Printf("Web server starting on port %v.\n", config.Web.Port)

	err = http.ListenAndServe(addr, router)
	if err != nil {
		log.Printf("Error starting web server: %v.\n", err)
		return err
	}

	return nil
}

func postToggleControllerMode(controller *controller.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		controller.ToggleMode()
		w.WriteHeader(http.StatusOK)
	}
}

package webapp

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/a-h/templ"
	"github.com/deseteral/resistere/internal/configuration"
	"github.com/deseteral/resistere/internal/controller"
	"github.com/deseteral/resistere/internal/metrics"
	"github.com/deseteral/resistere/internal/webapp/view"
)

//go:embed static
var staticFiles embed.FS

func StartWebServerBlocking(config *configuration.Config, c *controller.Controller, m *metrics.Registry) error {
	addr := fmt.Sprintf(":%v", config.Web.Port)

	var staticFs = fs.FS(staticFiles)
	htmlContent, err := fs.Sub(staticFs, "static")
	if err != nil {
		log.Printf("Error starting web server: %v.\n", err)
		return err
	}

	router := http.NewServeMux()

	router.Handle("GET /", templ.Handler(view.Index(c, m)))
	router.HandleFunc("POST /controller/mode", postChangeControllerMode(c))
	router.Handle("GET /metrics/prometheus", getPrometheus(m))
	router.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.FS(htmlContent))))

	log.Printf("Web server starting on port %v.\n", config.Web.Port)

	err = http.ListenAndServe(addr, router)
	if err != nil {
		log.Printf("Error starting web server: %v.\n", err)
		return err
	}

	return nil
}

func postChangeControllerMode(c *controller.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		value := r.URL.Query().Get("value")

		if value == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		castedValue, err := strconv.Atoi(value)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		nextMode, err := controller.ParseIntToMode(castedValue)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		c.ChangeMode(nextMode)

		view.ControllerModeSection(c).Render(r.Context(), w)
	}
}

func getPrometheus(m *metrics.Registry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var sb strings.Builder

		sb.WriteString(fmt.Sprintf("inverter_power_production_watts %f\n", m.LatestFrame.PowerProductionWatts))
		sb.WriteString(fmt.Sprintf("inverter_power_consumption_watts %f\n", m.LatestFrame.PowerConsumptionWatts))

		w.Write([]byte(sb.String()))
	}
}

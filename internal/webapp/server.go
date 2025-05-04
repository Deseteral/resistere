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
	"github.com/invopop/ctxi18n"
)

//go:embed static
var staticFiles embed.FS

//go:embed locale
var locale embed.FS

func StartWebServerBlocking(config *configuration.Config, c *controller.Controller, m *metrics.Registry) error {
	addr := fmt.Sprintf(":%v", config.Web.Port)

	var staticFs = fs.FS(staticFiles)
	htmlContent, err := fs.Sub(staticFs, "static")
	if err != nil {
		log.Printf("Error starting web server: %v.\n", err)
		return err
	}

	err = ctxi18n.Load(locale)
	if err != nil {
		log.Printf("Error starting web server: %v.\n", err)
		return err
	}

	router := http.NewServeMux()

	router.Handle("GET /", i18nMiddleware(templ.Handler(view.Index(c, m))))
	router.Handle("GET /view/stats", i18nMiddleware(templ.Handler(view.StatsSection(m))))
	router.Handle("POST /controller/mode", i18nMiddleware(postChangeControllerMode(c)))
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

func i18nMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lang := "en"

		acceptLang := r.Header.Get("Accept-Language")
		if len(acceptLang) > 0 {
			lang = acceptLang
		}

		ctx, err := ctxi18n.WithLocale(r.Context(), lang)
		if err != nil {
			log.Printf("Error setting locale: %v", err)
			http.Error(w, "Error setting locale", http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
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

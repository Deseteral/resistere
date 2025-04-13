package evse

import (
	"encoding/json"
	"fmt"
	"github.com/deseteral/resistere/internal/configuration"
	"io"
	"log"
	"net/http"
	"time"
)

type TeslaWallConnector struct {
	baseUrl string
	client  *http.Client
}

func (w *TeslaWallConnector) IsVehicleConnected() (isVehicleConnected bool, error error) {
	url := fmt.Sprintf("%s/api/1/vitals", w.baseUrl)

	resp, err := w.client.Get(url)
	if err != nil {
		log.Printf("Error fetching Tesla Wall Connector vitals: %v", err)
		return false, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error parsing body of Tesla Wall Connector vitals response: %v", err)
		return false, err
	}

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Printf("Error parsing JSON body of Tesla Wall Connector vitals response: %v", err)
		return false, err
	}

	connected := data["vehicle_connected"].(bool)
	return connected, nil
}

func NewTeslaWallConnector(config *configuration.TeslaWallConnector) *TeslaWallConnector {
	return &TeslaWallConnector{
		baseUrl: fmt.Sprintf("http://%s", config.Ip),
		client: &http.Client{
			Timeout: time.Second * 3,
		},
	}
}

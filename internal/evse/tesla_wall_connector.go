package evse

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/deseteral/resistere/internal/configuration"
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

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("Tesla Wall Connector returned HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error parsing body of Tesla Wall Connector vitals response: %v", err)
		return false, err
	}

	var data struct {
		VehicleConnected *bool `json:"vehicle_connected"`
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Printf("Error parsing JSON body of Tesla Wall Connector vitals response: %v", err)
		return false, err
	}

	if data.VehicleConnected == nil {
		return false, fmt.Errorf("Tesla Wall Connector vitals missing vehicle_connected")
	}
	return *data.VehicleConnected, nil
}

func NewTeslaWallConnector(config *configuration.TeslaWallConnector) *TeslaWallConnector {
	return &TeslaWallConnector{
		baseUrl: fmt.Sprintf("http://%s", config.Ip),
		client: &http.Client{
			Timeout: time.Second * 3,
		},
	}
}

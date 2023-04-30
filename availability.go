package campscraper

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"
)

type CampsiteAvailabilityProvider interface {
	Fetch( /** some fields common to both APIs **/ ) (map[string]*CampsiteAvailabilities, error)
}

type CampsiteAvailabilities struct {
	CampsiteID     string
	CampsiteName   string
	CampsiteType   string
	MaxNumPeople   int
	MinNumPeople   int
	Availabilities map[time.Time]string
}

func FetchData(req *http.Request, target CampsiteAvailabilityProvider) (CampsiteAvailabilityProvider, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Failed to fetch data, status code: " + resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &target)
	if err != nil {
		return nil, err
	}

	return target, nil
}

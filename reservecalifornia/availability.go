package reservecalifornia

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/brensch/campscraper"
)

type FacilityAvailabilityResponse struct {
	Facility Facility `json:"Facility"`
}

type Facility struct {
	Units map[string]Unit `json:"Units"`
}

type Unit struct {
	UnitId int              `json:"unitId"`
	Name   string           `json:"name"`
	Slices map[string]Slice `json:"slices"`
}

type Slice struct {
	Date   string `json:"date"`
	IsFree bool   `json:"isFree"`
}

func (res FacilityAvailabilityResponse) ToCampsiteAvailabilities() map[string]*campscraper.CampsiteAvailabilities {
	campsites := make(map[string]*campscraper.CampsiteAvailabilities)
	for _, data := range res.Facility.Units {
		availabilities := make(map[time.Time]string)

		for _, slice := range data.Slices {
			date, _ := time.Parse("2006-01-02", slice.Date)
			status := "Reserved"
			if slice.IsFree {
				status = "Available"
			}
			availabilities[date] = status
		}

		campsites[fmt.Sprintf("%d", data.UnitId)] = &campscraper.CampsiteAvailabilities{
			CampsiteID:     fmt.Sprintf("%d", data.UnitId),
			CampsiteName:   data.Name,
			Availabilities: availabilities,
		}
	}
	return campsites
}

const (
	availabilityEndpoint = "https://calirdr.usedirect.com/RDR/rdr/search/grid"
)

func (rc *ReserveCalifornia) Fetch(startDate, endDate time.Time) (map[string]*campscraper.CampsiteAvailabilities, error) {
	requestBody := map[string]interface{}{
		"IsADA":             false,
		"MinVehicleLength":  0,
		"UnitCategoryId":    1,
		"StartDate":         startDate.Format("2006-01-02"),
		"WebOnly":           true,
		"UnitTypesGroupIds": []int{},
		"SleepingUnitId":    83,
		"EndDate":           endDate.Format("2006-01-02"),
		"UnitSort":          "orderby",
		"InSeasonOnly":      true,
		"FacilityId":        "683",
		"RestrictADA":       false,
	}

	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://calirdr.usedirect.com/RDR/rdr/search/grid", bytes.NewReader(requestBodyBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", "https://reservecalifornia.com")
	req.Header.Set("Referer", "https://reservecalifornia.com/")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36")

	resp, err := httputil.FetchData[ReserveCaliforniaResponse](req)
	if err != nil {
		return nil, err
	}

	// Convert the response into the map[string]*CampsiteAvailabilities format
	availabilities := make(map[string]*campscraper.CampsiteAvailabilities)
	// ... (convert the response data into the desired format)

	return availabilities, nil
}

package recreationgov

import (
	"time"

	"github.com/brensch/campscraper"
)

type CampsiteAvailabilityResponse struct {
	Campsites map[string]CampsiteAvailabilities `json:"campsites"`
}

type CampsiteAvailabilities struct {
	Availabilities map[string]string `json:"availabilities"`
	CampsiteID     string            `json:"campsite_id"`
	CampsiteType   string            `json:"campsite_type"`
	MaxNumPeople   int               `json:"max_num_people"`
	MinNumPeople   int               `json:"min_num_people"`
}

func (res CampsiteAvailabilityResponse) ToCampsiteAvailabilities() map[string]*campscraper.CampsiteAvailabilities {
	campsites := make(map[string]*campscraper.CampsiteAvailabilities)
	for _, data := range res.Campsites {
		availabilities := make(map[time.Time]string)
		for dateStr, status := range data.Availabilities {
			date, _ := time.Parse(time.RFC3339, dateStr)
			availabilities[date] = status
		}

		campsites[data.CampsiteID] = &campscraper.CampsiteAvailabilities{
			CampsiteID:     data.CampsiteID,
			CampsiteType:   data.CampsiteType,
			MaxNumPeople:   data.MaxNumPeople,
			MinNumPeople:   data.MinNumPeople,
			Availabilities: availabilities,
		}
	}
	return campsites
}

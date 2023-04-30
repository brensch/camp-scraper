package recreationgov_test

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/brensch/campscraper"
	"github.com/brensch/campscraper/recreationgov"
)

func TestToCampsiteAvailabilities(t *testing.T) {
	sampleJSON, err := os.ReadFile("./availability_test.json")
	if err != nil {
		t.Fatal(err)
	}

	var res recreationgov.CampsiteAvailabilityResponse
	err = json.Unmarshal(sampleJSON, &res)
	if err != nil {
		t.Fatalf("Failed to unmarshal sample JSON: %v", err)
	}

	// Ensure it implements the CampsiteAvailabilityProvider interface
	var _ campscraper.CampsiteAvailabilityProvider = res

	availabilities := res.ToCampsiteAvailabilities()

	// Check if the length of the availabilities map is as expected
	expectedLen := 1 // Replace with the expected length based on the sample JSON
	if len(availabilities) != expectedLen {
		t.Errorf("Expected %d campsites, got %d", expectedLen, len(availabilities))
	}

	// Check if specific values are as expected
	expectedCampsiteID := "906"
	expectedCampsiteType := "STANDARD NONELECTRIC"
	expectedMaxNumPeople := 6
	expectedMinNumPeople := 0
	expectedAvailability := "Reserved"
	expectedDate, _ := time.Parse(time.RFC3339, "2022-10-02T00:00:00Z")

	if campsite, ok := availabilities[expectedCampsiteID]; ok {
		if campsite.CampsiteID != expectedCampsiteID {
			t.Errorf("Expected CampsiteID %s, got %s", expectedCampsiteID, campsite.CampsiteID)
		}
		if campsite.CampsiteType != expectedCampsiteType {
			t.Errorf("Expected CampsiteType %s, got %s", expectedCampsiteType, campsite.CampsiteType)
		}
		if campsite.MaxNumPeople != expectedMaxNumPeople {
			t.Errorf("Expected MaxNumPeople %d, got %d", expectedMaxNumPeople, campsite.MaxNumPeople)
		}
		if campsite.MinNumPeople != expectedMinNumPeople {
			t.Errorf("Expected MinNumPeople %d, got %d", expectedMinNumPeople, campsite.MinNumPeople)
		}
		if status, ok := campsite.Availabilities[expectedDate]; ok {
			if status != expectedAvailability {
				t.Errorf("Expected availability %s, got %s", expectedAvailability, status)
			}
		} else {
			t.Errorf("Expected date %s not found in availabilities", expectedDate)
		}
	} else {
		t.Errorf("Expected campsite with ID %s not found", expectedCampsiteID)
	}
}

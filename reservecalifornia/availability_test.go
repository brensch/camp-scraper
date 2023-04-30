package reservecalifornia_test

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/brensch/campscraper/reservecalifornia"
)

func TestToCampsiteAvailabilities(t *testing.T) {
	sampleJSON, err := os.ReadFile("./availability_test.json")
	if err != nil {
		t.Fatal(err)
	}

	var response reservecalifornia.FacilityAvailabilityResponse
	err = json.Unmarshal(sampleJSON, &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal sample JSON: %v", err)
	}

	campsites := response.ToCampsiteAvailabilities()
	if len(campsites) != 1 {
		t.Fatalf("Expected 1 campsite, got %d", len(campsites))
	}

	campsite, ok := campsites["43841"]
	if !ok {
		t.Fatal("Campsite ID 43841 not found")
	}

	expectedCampsiteName := "Campsite #N002"
	if campsite.CampsiteName != expectedCampsiteName {
		t.Errorf("Expected campsite name %q, got %q", expectedCampsiteName, campsite.CampsiteName)
	}

	expectedAvailabilities := map[string]string{
		"2023-10-11": "Available",
		"2023-10-12": "Reserved",
		"2023-10-13": "Reserved",
		"2023-10-14": "Reserved",
		"2023-10-15": "Available",
		"2023-10-16": "Available",
		"2023-10-17": "Available",
		"2023-10-18": "Available",
	}

	for dateStr, expectedStatus := range expectedAvailabilities {
		date, _ := time.Parse("2006-01-02", dateStr)
		status, ok := campsite.Availabilities[date]
		if !ok {
			t.Errorf("Availability for %s not found", dateStr)
			continue
		}

		if status != expectedStatus {
			t.Errorf("Expected availability status %q for %s, got %q", expectedStatus, dateStr, status)
		}
	}
}

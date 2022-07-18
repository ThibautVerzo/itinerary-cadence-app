package managers

import (
	"training/hello-cadence/app/models"
)

func ComputeItinerary(itinerary *models.Itinerary) {
	var road []models.Point
	road = append(road, models.Point{
		Lat:  (itinerary.Departure.Lat + itinerary.Arrival.Lat) / 2,
		Long: (itinerary.Departure.Long + itinerary.Arrival.Long) / 2,
	})
	itinerary.Road = road
}

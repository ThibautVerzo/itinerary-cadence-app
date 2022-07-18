package models

import "fmt"

type Itinerary struct {
	Departure, Arrival *Point
	Road               []Point
}

type Point struct {
	Lat, Long float32
}

func (p Point) String() string {
	return fmt.Sprintf("Lat: %f, Long: %f", p.Lat, p.Long)
}

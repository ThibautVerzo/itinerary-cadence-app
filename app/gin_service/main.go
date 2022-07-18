package main

import (
	"github.com/gin-gonic/gin"

	"training/hello-cadence/app/adapters/cadenceAdapter"
	"training/hello-cadence/app/config"
	"training/hello-cadence/app/gin_service/controllers"
)

func main() {
	router := gin.Default()

	var appConfig config.AppConfig
	appConfig.Setup()
	var cadenceClient cadenceAdapter.CadenceAdapter
	cadenceClient.Setup(&appConfig.Cadence)

	itineraryController := controllers.ItineraryController{
		CadenceAdapter: &cadenceClient,
		Logger:         appConfig.Logger,
	}

	router.POST("/itineraries.set-departure", itineraryController.SendItineraryDeparture)
	router.POST("/itineraries.set-arrival", itineraryController.SendItineraryArrival)
	router.GET("/itineraries.status", itineraryController.GetItineraryStatus)
	router.GET("/itineraries", itineraryController.GetItinerary)

	router.Run("localhost:3030")
}

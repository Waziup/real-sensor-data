package api

import (
	"log"
	"net/http"
	"os"

	routing "github.com/julienschmidt/httprouter"
)

/*-------------------------*/

func setupRouter() *routing.Router {

	var router = routing.New()

	router.GET("/", HomeLink)
	router.GET("/ui/*file_path", UI)

	router.POST("/auth", PostAuth)
	router.POST("/auth/logout", PostLogout)

	// router.GET("/docs/", APIDocs)
	// router.GET("/docs/:file_path", APIDocs)

	router.GET("/dataCollection/status", GetDataCollectionStatus)
	router.GET("/dataCollection/statistics", GetDataCollectionStatistics)

	router.GET("/sensors", GetSensors)
	router.GET("/sensors/:sensor_id", GetSensor)
	router.GET("/sensors/:sensor_id/values", GetSensorValues)

	router.GET("/sensors/:sensor_id/pushSettings", GetSensorPushSettings)
	router.POST("/sensors/:sensor_id/pushSettings", PostSensorPushSettings)
	router.DELETE("/sensors/:sensor_id/pushSettings/:id", DeleteSensorPushSettings)
	router.GET("/myPushSettings/sensors", GetMyPushSensors)

	router.GET("/search/sensors/:query", GetSearchSensors)

	router.GET("/channels", GetChannels)
	router.GET("/channels/:channel_id", GetChannel)
	router.GET("/channels/:channel_id/sensors", GetChannelSensors)
	router.GET("/channels/:channel_id/sensors/:sensor_id/values", GetSensorValues)

	router.GET("/user", GetUser)
	router.GET("/userDevices", GetUserDevicesAndSensors)

	return router
}

/*-------------------------*/

// ListenAndServeHTTP serves the APIs and the ui
func ListenAndServeHTTP() {

	log.Printf("Initializing...")

	router := setupRouter()

	addr := os.Getenv("SERVING_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	log.Printf("[Info  ] Serving on %s", addr)

	log.Fatal(http.ListenAndServe(addr, router))
}

/*-------------------------*/

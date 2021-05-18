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

	// router.GET("/docs/", APIDocs)
	// router.GET("/docs/:file_path", APIDocs)

	router.GET("/dataCollection/status", GetDataCollectionStatus)

	router.GET("/sensors", GetSensors)
	router.GET("/sensors/search/:query", GetSearchSensors)

	router.GET("/channels", GetChannels)
	router.GET("/channels/:channel", GetChannel)
	router.GET("/channels/:channel/sensors", GetChannelSensors)
	router.GET("/channels/:channel/sensors/:name/values", GetSensorValues)

	// router.GET("/docker/:cId", DockerStatusById)
	// router.POST("/docker/:cId/:action", DockerAction)
	// router.PUT("/docker/:cId/:action", DockerAction)
	// router.GET("/docker/:cId/logs", DockerLogs)
	// router.GET("/docker/:cId/logs/:tail", DockerLogs)

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

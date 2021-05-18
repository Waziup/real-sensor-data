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

// HomeLink implements GET / Just a test msg to see if it works
func HomeLink(resp http.ResponseWriter, req *http.Request, params routing.Params) {

	homeHTML := `
	<style>
		.box{
			-moz-border-radius: 6px;
			-webkit-border-radius: 6px;
			background-color: #fbf8ff;
			background-image: url(../Images/icons/Pencil-48.png);
			background-position: 9px 0px;
			background-repeat: no-repeat;
			border: solid 1px #3498db;
			border-radius: 6px;
			line-height: 18px;
			overflow: hidden;
			padding: 15px 60px;
			width: 300px;
			margin: auto;
				margin-top: auto;
			box-shadow: rgba(0, 0, 0, 0.25) 0px 0.0625em 0.0625em, rgba(0, 0, 0, 0.25) 0px 0.125em 0.5em, rgba(255, 255, 255, 0.1) 0px 0px 0px 1px inset;
			margin-top: 200px;
			text-align:center;
		}
	</style>
	<div class="box">
		Salam Goloooo, It works!
		<p>
			Navigate to the <a href="/ui/" >Web UI</a>
		</p>
	</div>
	<script>
		setTimeout( () => {window.location.href="/ui/"}, 1000)
	</script>
	`
	resp.Header().Set("Content-Type", "text/html; charset=utf-8")
	resp.Write([]byte(homeHTML))
}

/*-------------------------*/

func UI(resp http.ResponseWriter, req *http.Request, params routing.Params) {

	rootPath := os.Getenv("EXEC_PATH")
	if rootPath == "" {
		rootPath = "./"
	}

	http.FileServer(http.Dir(rootPath)).ServeHTTP(resp, req)
}

/*-------------------------*/

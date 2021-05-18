package main

import (
	"real-sensor-data/api"
	"real-sensor-data/database"
	"real-sensor-data/datacollection"
	"real-sensor-data/global"
	// "github.com/influxdata/influxdb-client-go/v3"
)

func main() {

	global.DB = database.New(database.Postgres)
	defer global.DB.Close()

	datacollection.Init()

	api.ListenAndServeHTTP()
}

/*--------------------------------*/

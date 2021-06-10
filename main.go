package main

import (
	"fmt"
	"sensor-data-simulator/api"
	"sensor-data-simulator/database"
	"sensor-data-simulator/datacollection"
	"sensor-data-simulator/datapush"
	"sensor-data-simulator/dbinit"
	"sensor-data-simulator/global"
)

func main() {

	psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		global.ENV.POSTGRES_HOST,
		global.ENV.POSTGRES_PORT,
		global.ENV.POSTGRES_USER,
		global.ENV.POSTGRES_PASSWORD,
		global.ENV.POSTGRES_DB,
	)

	global.DB = database.New(database.Postgres, psqlconn)
	defer global.DB.Close()

	/*--------*/

	dbinit.DatabaseInit()

	/*--------*/

	datapush.Init()

	datacollection.Init()

	/*--------*/

	api.ListenAndServeHTTP()
}

/*--------------------------------*/

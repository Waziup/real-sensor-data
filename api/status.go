package api

import (
	"log"
	"net/http"
	"sensor-data-simulator/database"
	"sensor-data-simulator/global"
	"sensor-data-simulator/tools"

	routing "github.com/julienschmidt/httprouter"
)

/*-------------*/

/*
* This function implements GET /dataCollection/status
 */

func GetDataCollectionStatus(resp http.ResponseWriter, req *http.Request, params routing.Params) {

	tools.SendJSON(resp, global.DataCollectorProgress)
}

/*-------------*/
/*
* This function implements GET /dataCollection/statistics
 */

func GetDataCollectionStatistics(resp http.ResponseWriter, req *http.Request, params routing.Params) {

	totalChannels := int64(0)
	{
		SQL := `SELECT COUNT(*) AS "total" FROM "channels"`
		rows, err := global.DB.Query(SQL, database.QueryParams{})
		if err != nil {
			log.Printf("Error in db query: %v", err)
			http.Error(resp, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		totalChannels = rows[0]["total"].(int64)
	}

	/*-------*/

	totalSensors := int64(0)
	{
		SQL := `SELECT COUNT(*) AS "total" FROM "sensors"`
		rows, err := global.DB.Query(SQL, database.QueryParams{})
		if err != nil {
			log.Printf("Error in db query: %v", err)
			http.Error(resp, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		totalSensors = rows[0]["total"].(int64)
	}

	/*-------*/

	totalSensorValues := int64(0)
	{
		SQL := `SELECT COUNT(*) AS "total" FROM "sensor_values"`
		rows, err := global.DB.Query(SQL, database.QueryParams{})
		if err != nil {
			log.Printf("Error in db query: %v", err)
			http.Error(resp, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		totalSensorValues = rows[0]["total"].(int64)
	}

	/*-------*/

	tools.SendJSON(resp, map[string]interface{}{
		"totalChannels":     totalChannels,
		"totalSensors":      totalSensors,
		"totalSensorValues": totalSensorValues,
	})
}

/*-------------*/

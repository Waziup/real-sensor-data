package api

import (
	"log"
	"net/http"
	"real-sensor-data/database"
	"real-sensor-data/global"
	"real-sensor-data/tools"
	"strconv"

	routing "github.com/julienschmidt/httprouter"
)

/*-------------*/

/*
* This function implements GET /sensors
 */
func GetSensors(resp http.ResponseWriter, req *http.Request, params routing.Params) {

	limit, offset, page := tools.GetLimitOffset(req)

	SQL := `SELECT DISTINCT "name", "channel_id" FROM "sensor_values" LIMIT $1 OFFSET $2`

	rows, err := global.DB.Query(SQL, database.QueryParams{limit, offset})
	if err != nil {
		log.Printf("Error in db query: %v", err)
		http.Error(resp, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tools.SendJSON(resp, map[string]interface{}{"page": page, "rows": rows})
}

/*-------------*/
/*
* This function implements GET /sensors/search
 */
func GetSearchSensors(resp http.ResponseWriter, req *http.Request, params routing.Params) {

	query := params.ByName("query")

	limit, offset, page := tools.GetLimitOffset(req)

	SQL := `SELECT DISTINCT "name", "channel_id" FROM "sensor_values" WHERE "name" ILIKE $1 LIMIT $2 OFFSET $3`

	rows, err := global.DB.Query(SQL, database.QueryParams{"%" + query + "%", limit, offset})
	if err != nil {
		log.Printf("Error in db query: %v", err)
		http.Error(resp, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tools.SendJSON(resp, map[string]interface{}{"query": query, "page": page, "rows": rows})
}

/*-------------*/
/*
* This function implements GET /sensors/:channel/sensors/:name/values
 */
func GetSensorValues(resp http.ResponseWriter, req *http.Request, params routing.Params) {

	name := params.ByName("name")
	channel := params.ByName("channel")

	channel_id, err := strconv.Atoi(channel)
	if err != nil {
		channel_id = 0
	}

	limit, offset, page := tools.GetLimitOffset(req)

	SQL := `SELECT * 
			FROM "sensor_values" 
			WHERE 
				"name" = $1 AND
				"channel_id" = $2
			LIMIT $3 OFFSET $4`

	rows, err := global.DB.Query(SQL, database.QueryParams{name, channel_id, limit, offset})
	if err != nil {
		log.Printf("Error in db query: %v", err)
		http.Error(resp, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tools.SendJSON(resp, map[string]interface{}{"page": page, "rows": rows})
}

/*-------------*/

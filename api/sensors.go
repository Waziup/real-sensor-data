package api

import (
	"log"
	"math"
	"net/http"
	"sensor-data-simulator/database"
	"sensor-data-simulator/global"
	"sensor-data-simulator/tools"
	"strconv"

	routing "github.com/julienschmidt/httprouter"
)

/*
* This function implements GET /sensors
 */
func GetSensors(resp http.ResponseWriter, req *http.Request, params routing.Params) {

	limit, offset, page := tools.GetLimitOffset(req)

	/*------*/

	totalRows := int64(0)
	{
		SQL := `SELECT COUNT(*) AS "total" FROM "sensors"`
		rows, err := global.DB.Query(SQL, database.QueryParams{})
		if err != nil {
			log.Printf("Error in db query: %v", err)
			http.Error(resp, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		totalRows = rows[0]["total"].(int64)
	}

	totalPages := int64(math.Ceil(float64(totalRows) / float64(global.RowsPerPage)))
	pagination := map[string]interface{}{
		"current_page":  page,
		"total_pages":   totalPages,
		"total_entries": totalRows,
	}

	/*------*/

	SQL := `SELECT 
					s.*,
					c."name"	AS "channel_name"
			FROM 
				"sensors"	AS s,
				"channels"	AS c
			WHERE
				c."id" = s."channel_id"
			LIMIT $1 OFFSET $2`

	rows, err := global.DB.Query(SQL, database.QueryParams{limit, offset})
	if err != nil {
		log.Printf("Error in db query: %v", err)
		http.Error(resp, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tools.SendJSON(resp, map[string]interface{}{"pagination": pagination, "rows": rows})
}

/*-------------*/

/*
* This function implements GET /sensors/:sensor_id
 */
func GetSensor(resp http.ResponseWriter, req *http.Request, params routing.Params) {

	sensorIdStr := params.ByName("sensor_id")

	sensor_id, err := strconv.Atoi(sensorIdStr)
	if err != nil {
		sensor_id = 0
	}

	/*------*/

	SQL := `SELECT * FROM "sensors" WHERE "id" = $1`

	rows, err := global.DB.Query(SQL, database.QueryParams{sensor_id})
	if err != nil {
		log.Printf("Error in db query: %v", err)
		http.Error(resp, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if rows == nil || len(rows) == 0 {
		http.Error(resp, "Sensor not found!", http.StatusNotFound)
		return
	}

	tools.SendJSON(resp, rows[0])
}

/*-------------*/
/*
* This function implements GET /search/sensors/:query
 */
func GetSearchSensors(resp http.ResponseWriter, req *http.Request, params routing.Params) {

	query := params.ByName("query")
	limit, offset, page := tools.GetLimitOffset(req)

	/*------*/

	totalRows := int64(0)
	{
		SQL := `SELECT COUNT(*) AS "total" FROM "sensors" WHERE "name" ILIKE $1`
		rows, err := global.DB.Query(SQL, database.QueryParams{"%" + query + "%"})
		if err != nil {
			log.Printf("Error in db query: %v", err)
			http.Error(resp, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		totalRows = rows[0]["total"].(int64)
	}

	totalPages := int64(math.Ceil(float64(totalRows) / float64(global.RowsPerPage)))
	pagination := map[string]interface{}{
		"current_page":  page,
		"total_pages":   totalPages,
		"total_entries": totalRows,
	}

	/*------*/

	SQL := `SELECT 
				s.*,
				c."name"	AS "channel_name"
			FROM 
				"sensors"	AS s,
				"channels"		AS c
			WHERE
				c."id" = s."channel_id" AND
				s."name" ILIKE $1
			LIMIT $2 OFFSET $3`

	rows, err := global.DB.Query(SQL, database.QueryParams{"%" + query + "%", limit, offset})
	if err != nil {
		log.Printf("Error in db query: %v", err)
		http.Error(resp, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tools.SendJSON(resp, map[string]interface{}{"query": query, "pagination": pagination, "rows": rows})
}

/*-------------*/
/*
* This function implements GET /sensors/:sensor_id/values
* and GET /channels/:channel_id/sensors/:sensor_id/values (keep it for legacy)
 */
func GetSensorValues(resp http.ResponseWriter, req *http.Request, params routing.Params) {

	sensorIdStr := params.ByName("sensor_id")

	// Sensor ID is unique enough now, so let's ignore the channel
	// channelIdStr := params.ByName("channel_id")

	// channel_id, err := strconv.Atoi(channelIdStr)
	// if err != nil {
	// 	channel_id = 0
	// }

	sensor_id, err := strconv.Atoi(sensorIdStr)
	if err != nil {
		sensor_id = 0
	}

	limit, offset, page := tools.GetLimitOffset(req)

	/*------*/

	totalRows := int64(0)
	{
		SQL := `SELECT 
					COUNT(*) AS "total" 
				FROM 
					"sensors"			AS	s,
					"sensor_values"		AS	v
				WHERE 
					s."id" = $1 AND
					s."id" = v."sensor_id" AND
					v."value" != ''`
		rows, err := global.DB.Query(SQL, database.QueryParams{sensor_id})
		if err != nil {
			log.Printf("Error in db query: %v", err)
			http.Error(resp, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		totalRows = rows[0]["total"].(int64)
	}

	totalPages := int64(math.Ceil(float64(totalRows) / float64(global.RowsPerPage)))
	pagination := map[string]interface{}{
		"current_page":  page,
		"total_pages":   totalPages,
		"total_entries": totalRows,
	}

	/*------*/
	SQL := `SELECT s."name", v.* 
			FROM 
				"sensors"			AS	s,
				"sensor_values"		AS	v
			WHERE 
				s."id" = $1 AND
				s."id" = v."sensor_id" AND
				v."value" != ''
			ORDER BY "entry_id" DESC
			LIMIT $2 OFFSET $3`

	rows, err := global.DB.Query(SQL, database.QueryParams{sensor_id, limit, offset})
	if err != nil {
		log.Printf("Error in db query: %v", err)
		http.Error(resp, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tools.SendJSON(resp, map[string]interface{}{"pagination": pagination, "rows": rows})
}

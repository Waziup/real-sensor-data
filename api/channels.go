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

/*-------------*/

/*
* This function implements GET /channels
 */
func GetChannels(resp http.ResponseWriter, req *http.Request, params routing.Params) {

	limit, offset, page := tools.GetLimitOffset(req)

	/*------*/

	totalRows := int64(0)
	{
		SQL := `SELECT COUNT(*) AS "total" FROM "channels"`
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

	SQL := `SELECT * FROM "channels" LIMIT $1 OFFSET $2`

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
* This function implements GET /channels/:channel_id
 */
func GetChannel(resp http.ResponseWriter, req *http.Request, params routing.Params) {

	channelIdStr := params.ByName("channel_id")

	channel_id, err := strconv.Atoi(channelIdStr)
	if err != nil {
		channel_id = 0
	}

	SQL := `SELECT * FROM "channels" WHERE "id" = $1`
	channelRows, err := global.DB.Query(SQL, database.QueryParams{channel_id})
	if err != nil {
		log.Printf("Error in db query: %v", err)
		http.Error(resp, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if channelRows == nil || len(channelRows) == 0 {
		http.Error(resp, "Channel not found!", http.StatusNotFound)
		return
	}

	tools.SendJSON(resp, channelRows[0])
}

/*-------------*/

/*
* This function implements GET /channels/:channel_id/sensors
 */
func GetChannelSensors(resp http.ResponseWriter, req *http.Request, params routing.Params) {

	channelIdStr := params.ByName("channel_id")

	channel_id, err := strconv.Atoi(channelIdStr)
	if err != nil {
		channel_id = 0
	}

	limit, offset, page := tools.GetLimitOffset(req)

	SQL := `SELECT * FROM "channels" WHERE "id" = $1`
	channelRows, err := global.DB.Query(SQL, database.QueryParams{channel_id})
	if err != nil {
		log.Printf("Error in db query: %v", err)
		http.Error(resp, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if channelRows == nil || len(channelRows) == 0 {
		http.Error(resp, "Channel not found!", http.StatusNotFound)
		return
	}

	/*------*/

	totalRows := int64(0)
	{
		SQL := `SELECT COUNT(*) AS "total" FROM "sensors" WHERE "channel_id" = $1`
		rows, err := global.DB.Query(SQL, database.QueryParams{channel_id})
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

	SQL = `SELECT "id", "name"
			FROM "sensors" 
			WHERE
				"channel_id" = $1
			LIMIT $2 OFFSET $3`

	rows, err := global.DB.Query(SQL, database.QueryParams{channel_id, limit, offset})
	if err != nil {
		log.Printf("Error in db query: %v", err)
		http.Error(resp, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tools.SendJSON(resp, map[string]interface{}{
		"channel":    channelRows[0],
		"pagination": pagination,
		"rows":       rows,
	})
}

/*-------------*/

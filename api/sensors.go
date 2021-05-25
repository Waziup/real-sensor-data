package api

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"sensor-data-simulator/database"
	"sensor-data-simulator/global"
	"sensor-data-simulator/tools"
	"strconv"
	"time"

	routing "github.com/julienschmidt/httprouter"
)

/*-------------*/
type SensorPushSettings struct {
	ID int64 `json:"id"`
	// UserId            int64     `json:"user_id"` //should not be exposed
	// SensorId          int64     `json:"sensor_id"`
	TargetDeviceId    string    `json:"target_device_id"`
	TargetSensorId    string    `json:"target_sensor_id"`
	Active            bool      `json:"active"`
	LastPushedEntryId int64     `json:"last_pushed_entry_id"`
	PushInterval      int       `json:"push_interval"`
	LastPushTime      time.Time `json:"last_push_time"`
}

/*-------------*/

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

/*-------------*/
/*
* This function implements POST /sensors/:sensor_id/pushSettings
* It Adds or Modify the setting
 */
func PostSensorPushSettings(resp http.ResponseWriter, req *http.Request, params routing.Params) {

	userId, err := getAuthorizedUserID(resp, req)
	if err != nil {
		http.Error(resp, "Unauthorized", http.StatusUnauthorized)
		return
	}

	/*------------*/

	sensorIdStr := params.ByName("sensor_id")

	sensorId, err := strconv.Atoi(sensorIdStr)
	if err != nil {
		sensorId = 0
	}

	/*------------*/

	body, err := tools.ReadAll(req.Body)
	if err != nil {
		log.Printf("[ERR  ] PostSensorPushSettings: %s", err.Error())
		http.Error(resp, "bad request", http.StatusBadRequest)
		return
	}

	var inputRecord SensorPushSettings

	err = json.Unmarshal(body, &inputRecord)
	if err != nil {
		log.Printf("[ERR  ] PostSensorPushSettings: %s", err.Error())
		http.Error(resp, "bad request", http.StatusBadRequest)
		return
	}

	/*------------*/

	row := database.RowType{
		"user_id":          userId,
		"sensor_id":        sensorId,
		"target_device_id": inputRecord.TargetDeviceId,
		"target_sensor_id": inputRecord.TargetSensorId,
		"active":           inputRecord.Active,
		"push_interval":    inputRecord.PushInterval,
	}

	if inputRecord.ID == 0 { // New record

		_, err := global.DB.Insert("push_settings", row)
		if err != nil {
			log.Printf("\nError in `push_settings` insertion: %v \nRow: \n%v", err, row)
			http.Error(resp, "something went wrong", http.StatusInternalServerError)
			return
		}

	} else {

		condRows := database.RowType{
			"id":      inputRecord.ID,
			"user_id": userId,
		}
		_, err := global.DB.Update("push_settings", row, condRows)
		if err != nil {
			log.Printf("\nError in `push_settings` insertion: %v \nRow: \n%v", err, row)
			http.Error(resp, "something went wrong", http.StatusInternalServerError)
			return
		}
	}
	resp.Write([]byte("OK"))

}

/*-------------*/
/*
* This function implements POST /sensors/:sensor_id/pushSettings/:id
* It Adds or Modify the setting
 */
func DeleteSensorPushSettings(resp http.ResponseWriter, req *http.Request, params routing.Params) {

	userId, err := getAuthorizedUserID(resp, req)
	if err != nil {
		http.Error(resp, "Unauthorized", http.StatusUnauthorized)
		return
	}

	/*------------*/

	sensorIdStr := params.ByName("sensor_id")
	sensorId, err := strconv.Atoi(sensorIdStr)
	if err != nil {
		sensorId = 0
	}

	/*------------*/

	recordIdStr := params.ByName("id")
	recordId, err := strconv.Atoi(recordIdStr)
	if err != nil {
		recordId = 0
	}

	/*------------*/

	condRows := database.RowType{
		"id":        recordId,
		"user_id":   userId,
		"sensor_id": sensorId,
	}
	_, err = global.DB.Delete("push_settings", condRows)
	if err != nil {
		log.Printf("\nError in `push_settings` Deletion: %v \ncondRows: \n%v", err, condRows)
		http.Error(resp, "something went wrong", http.StatusInternalServerError)
		return
	}

	resp.Write([]byte("OK"))
}

/*-------------*/ /*
* This function implements GET /sensors/:sensor_id/pushSettings
 */
func GetSensorPushSettings(resp http.ResponseWriter, req *http.Request, params routing.Params) {

	userId, err := getAuthorizedUserID(resp, req)
	if err != nil {
		http.Error(resp, "Unauthorized", http.StatusUnauthorized)
		return
	}

	/*------------*/

	sensorIdStr := params.ByName("sensor_id")

	sensorId, err := strconv.Atoi(sensorIdStr)
	if err != nil {
		sensorId = 0
	}

	/*------*/

	limit, offset, page := tools.GetLimitOffset(req)

	/*------*/

	totalRows := int64(0)
	{
		SQL := `SELECT COUNT(*) AS "total" 
				FROM	"push_settings"
				WHERE
					"sensor_id" = $1 AND 
					"user_id" = $2`
		rows, err := global.DB.Query(SQL, database.QueryParams{sensorId, userId})
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
					"id",
					"target_device_id",
					"target_sensor_id",
					"active",
					"push_interval",
					"last_push_time"
					
			FROM	"push_settings"
			WHERE
				"sensor_id" = $1 AND 
				"user_id" = $2
			LIMIT $3 OFFSET $4`

	rows, err := global.DB.Query(SQL, database.QueryParams{sensorId, userId, limit, offset})
	if err != nil {
		log.Printf("Error in db query: %v", err)
		http.Error(resp, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tools.SendJSON(resp, map[string]interface{}{"pagination": pagination, "rows": rows})
}

/*-------------*/

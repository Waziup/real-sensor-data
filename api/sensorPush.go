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
	UseOriginalTime   bool      `json:"use_original_time"`
	PushedCount       bool      `json:"pushed_count"`
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
		"user_id":           userId,
		"sensor_id":         sensorId,
		"target_device_id":  inputRecord.TargetDeviceId,
		"target_sensor_id":  inputRecord.TargetSensorId,
		"active":            inputRecord.Active,
		"push_interval":     inputRecord.PushInterval,
		"use_original_time": inputRecord.UseOriginalTime,
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
					"last_push_time",
					"use_original_time",
					"pushed_count"
					
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
/*
* This function implements GET /myPushSettings/sensors
* This function retrieves the sensors that the logged-in user has configured push-settings for
 */
func GetMyPushSensors(resp http.ResponseWriter, req *http.Request, params routing.Params) {

	userId, err := getAuthorizedUserID(resp, req)
	if err != nil {
		http.Error(resp, "Unauthorized", http.StatusUnauthorized)
		return
	}

	/*------------*/

	limit, offset, page := tools.GetLimitOffset(req)

	/*------*/

	totalRows := int64(0)
	{
		SQL := `SELECT COUNT(*) AS "total" 
				FROM ( SELECT s.* FROM 
							"push_settings" AS p, 
							"sensors" AS s 
						WHERE 
							s."id" = p."sensor_id"	AND
							p."user_id" = $1
						GROUP BY s."id") AS "tmp"`
		rows, err := global.DB.Query(SQL, database.QueryParams{userId})
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

	SQL := `SELECT s.* 
			FROM 
				"push_settings" AS p, 
				"sensors" AS s 
			WHERE 
				s."id" = p."sensor_id"	AND
				p."user_id" = $1
			GROUP BY s."id"
			LIMIT $2 OFFSET $3`

	rows, err := global.DB.Query(SQL, database.QueryParams{userId, limit, offset})
	if err != nil {
		log.Printf("Error in db query: %v", err)
		http.Error(resp, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	tools.SendJSON(resp, map[string]interface{}{"pagination": pagination, "rows": rows})
}

/*-------------*/

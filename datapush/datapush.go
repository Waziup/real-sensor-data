package datapush

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"sensor-data-simulator/api"
	"sensor-data-simulator/database"
	"sensor-data-simulator/global"
	"strconv"
	"time"
)

func Init() {

	// These intervals are the only ones that we support for data push and take care of
	pushIntervalsInMinutes := []int{
		1,
		3,
		5,
		10,
		30,
		60,
		2 * 60,
		3 * 60,
		5 * 60,
		24 * 60,
		2 * 24 * 60,
		3 * 24 * 60,
	}

	for _, intervalInMinutes := range pushIntervalsInMinutes {
		go handlePushInterval(intervalInMinutes)
	}

}

/*--------------*/

func handlePushInterval(intervalInMinutes int) {

	for {
		time.Sleep(time.Duration(intervalInMinutes) * time.Minute)

		/*-------*/

		SQL := `SELECT p.*, u."token" 
				FROM 
					"push_settings" AS p, 
					"users" AS u 
				WHERE 
					p."active" = true		AND
					u."id" = p."user_id"	AND 
					"push_interval" = $1`

		pushRows, err := global.DB.Query(SQL, database.QueryParams{intervalInMinutes})
		if err != nil {
			log.Printf("\nError in `push_settings` Load: %v \nSQL: \n%v\nintervalInMinutes: %v", err, SQL, intervalInMinutes)
			return
		}

		/*-------*/

		for _, pushRow := range pushRows {

			/*---------*/

			if pushRow["last_pushed_entry_id"] == nil {
				pushRow["last_pushed_entry_id"] = int64(0)
			}

			sourceSensorRow, err := GetTheNextValueToPush(pushRow["sensor_id"].(int64), pushRow["last_pushed_entry_id"].(int64))
			if err != nil {
				continue
			}

			if sourceSensorRow == nil {
				// log.Printf("No new values to push. device: %v, sensor: %v", pushRow["target_device_id"], pushRow["target_sensor_id"])
				continue
			}

			/*---------*/

			pushTime := time.Now()
			value := sourceSensorRow["value"].(string)

			sensorTimestamp := pushTime
			if pushRow["use_original_time"] != nil && pushRow["use_original_time"].(bool) {
				sensorTimestamp = sourceSensorRow["created_at"].(time.Time)
			}

			statusCode, err := PushDataToWaziup(pushRow["token"].(string), pushRow["target_device_id"].(string), pushRow["target_sensor_id"].(string), value, sensorTimestamp)
			if err != nil {
				if statusCode == 403 { // The token is not valid anymore, let's refresh it
					// log.Printf("[PUSH ] The token is expired for `%v` Renewing token", pushRow["user_id"])
					newToken, err := RefreshWaziupToken(pushRow["user_id"].(int64))
					if err != nil {
						log.Printf("[PUSH ] Error in token acquisition: %v", err)
						continue
					}

					// log.Printf("[PUSH ] New token acquired, pushing again...")

					// Let's repeat the push process one more time
					_, err = PushDataToWaziup(newToken, pushRow["target_device_id"].(string), pushRow["target_sensor_id"].(string), value, sensorTimestamp)
					if err != nil {
						log.Printf("[PUSH ] error in data push: `%v` \nUserId: %v", err, pushRow["user_id"])
						continue
					}

				} else {
					continue
				}
			}

			/*---------*/

			// Update the push row:
			UpdatePushSettingLastEntry(pushRow["id"].(int64), sourceSensorRow["entry_id"].(int64), pushTime)

		}
	}
}

/*--------------*/

func UpdatePushSettingLastEntry(id int64, lastPushedEntryId int64, lastPushTime time.Time) error {

	SQL := `UPDATE "push_settings" 
			SET 
				"last_pushed_entry_id" = $1,
				"last_push_time" = $2,
				"pushed_count" = "pushed_count" + 1
			WHERE 
				"id" = $3`

	params := database.QueryParams{lastPushedEntryId, lastPushTime, id}
	_, err := global.DB.Exec(SQL, params)
	if err != nil {
		log.Printf("\nError in updating `push_settings`: %v \nSQL: %v\nParams: %v", err, SQL, params)
	}

	return err
}

/*--------------*/

func GetTheNextValueToPush(sensorId int64, lastPushedEntryId int64) (database.RowType, error) {

	SQL := `SELECT * 
			FROM "sensor_values" 
			WHERE 
				"sensor_id" = $1 AND 
				"entry_id" > $2 
			ORDER BY "entry_id" ASC 
			LIMIT 1`
	params := database.QueryParams{sensorId, lastPushedEntryId}
	rows, err := global.DB.Query(SQL, params)
	if err != nil {
		log.Printf("\nError in Query: %v \nSQL: \n%v\nParams: %v", err, SQL, params)
		return nil, err
	}

	if rows == nil || len(rows) == 0 {
		return nil, nil
	}

	return rows[0], nil
}

/*--------------*/

func PushDataToWaziup(token string, deviceId string, sensorId string, value string, timestamp time.Time) (int, error) {

	apiPath := fmt.Sprintf(global.ENV.WAZIUP_API_PATH+`devices/%s/sensors/%s/value`, deviceId, sensorId)

	// Attempting to post number value if possible

	var postBody []byte
	floatValue, err := strconv.ParseFloat(value, 64)
	if err == nil {
		postBody = []byte(fmt.Sprintf(`{"value":%v, "timestamp": "%v"}`, floatValue, timestamp.Format(time.RFC3339)))
	} else {
		postBody = []byte(fmt.Sprintf(`{"value":"%s", "timestamp": "%v"}`, value, timestamp.Format(time.RFC3339)))
	}

	/*--------*/

	req, err := http.NewRequest("POST", apiPath, bytes.NewBuffer(postBody))
	if err != nil {
		log.Printf("[PUSH ] could not make the request: %v", err)
		return 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("[PUSH ] did not receive a response from Waziup Server: %v", err)
		return 0, err
	}

	if resp.StatusCode != 204 {
		err := fmt.Errorf("waziup api error (%v): %v \n\tAPI path: %v", resp.StatusCode, resp.Status, apiPath)
		// log.Printf("[PUSH ] Waziup API Error: %v \nAPI path: %v", err, apiPath)
		return resp.StatusCode, err
	}

	return resp.StatusCode, nil
}

/*--------------*/

func RefreshWaziupToken(userId int64) (string, error) {

	user, err := api.GetUserById(userId)
	if err != nil {
		return "", err
	}

	newToken, err := api.CheckUserCredentials(user.Username, user.Password)
	if err != nil {
		return "", err
	}

	user.Token = newToken
	user.TokenHash = ""

	err = api.SaveUserInfo(user)
	if err != nil {
		return "", err
	}

	return user.Token, nil
}

/*--------------*/

package datacollection

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"runtime"
	"sensor-data-simulator/database"
	"sensor-data-simulator/global"
	"strconv"
	"strings"
	"time"
)

/*--------------------------------*/

// This has to be initiated manually (e.g. in main)
func Init() {

	go func() {
		for {

			/*---------*/

			ExtractChannelsData()

			time.Sleep(1 * time.Second)
			ExtractSensorsData()

			/*---------*/

			global.DataCollectorProgress.LastExtractionTime = time.Now()

			/*---------*/

			dataExtractionInterval := 60

			if val := global.ENV.DATA_EXTRACTION_INTERVAL; val != "" {
				dataExtractionInterval, _ = strconv.Atoi(val)
				if dataExtractionInterval <= 0 {
					dataExtractionInterval = 60
				}
			}

			time.Sleep(3 * time.Second)

			fmt.Printf("\nThe next run is in %d minutes", dataExtractionInterval)
			time.Sleep(time.Duration(dataExtractionInterval) * time.Minute)

			/*---------*/
		}
	}()
}

/*--------------------------------*/

// this global var is used in go routines to inform the main func (ExtractChannelsData)
// that the routines hit the last page of the channel data extraction
// and it is time to break the main loop
var channelsHitTheLastPage bool

func ExtractChannelsData() {

	fmt.Print("\n\t\t* * * Extracting new channels * * *\n\n")

	global.DataCollectorProgress.NewExtractedChannels = 0
	global.DataCollectorProgress.ChannelsRunning = true
	defer func() { global.DataCollectorProgress.ChannelsRunning = false }()

	channelsHitTheLastPage = false
	for page := 1; ; page++ {

		if channelsHitTheLastPage {
			break
		}

		// processChannelDataExtraction(page) // Execution Time: 50.27s
		go processChannelDataExtraction(page) // Execution Time: 04.21s

		// Let's wait for the routins to finish
		for runtime.NumGoroutine() > global.MaxNumGoRoutines {
			time.Sleep(200 * time.Millisecond)
		}

	}

	fmt.Printf("\n\nAll Done [ New channels: %d ] :)\n\n---------------------------------------------------------\n", global.DataCollectorProgress.NewExtractedChannels)
}

/*--------------------------------*/

func ExtractSensorsData() {

	global.DataCollectorProgress.SensorsRunning = true
	global.DataCollectorProgress.SensorsProgress = 0
	global.DataCollectorProgress.NewExtractedSensorValues = 0
	defer func() { global.DataCollectorProgress.SensorsRunning = false }()

	fmt.Print("\n\t\t* * * Extracting new sensor data * * *\n\n")

	channels, err := global.DB.Load("channels", nil)
	if err != nil {
		log.Fatal(err)
	}

	totalChannels := float64(len(channels))
	for chIndex, channel := range channels {

		progress := int(math.Round(100 * (float64(chIndex) + 1) / totalChannels))

		// processChannelSensors(channel) // Execution Time: 9m23.08s
		go processChannelSensors(channel) // Execution Time: 57.08s

		// Let's wait for the routins to finish
		for runtime.NumGoroutine() > global.MaxNumGoRoutines {
			global.DataCollectorProgress.SensorsProgress = progress
			time.Sleep(1 * time.Second)
		}
	}

	global.DataCollectorProgress.SensorsProgress = 100

	fmt.Printf("\n\nAll Done [ New sensor values: %d ] :)\n\n---------------------------------------------------------\n", global.DataCollectorProgress.NewExtractedSensorValues)
}

/*--------------------------------*/

func processChannelSensors(channel database.RowType) {

	apiURL := "https://thingspeak.com/channels/%v/feed.json"

	fmt.Printf("\r\t[ %3v %% ]\tProcessing Channel: %-20v", global.DataCollectorProgress.SensorsProgress, channel["id"])

	req, err := http.NewRequest("GET", fmt.Sprintf(apiURL, channel["id"]), nil)
	if err != nil {
		log.Printf("Error in sensor data extraction: %v", err)
		return
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("\nError in sensor data extraction: %v", err)
		return
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("\nCould not read the content: %v", err)
	}

	var sensorFeedJSON struct {
		Channel struct {
			Id          json.Number `json:"id"`
			LastEntryId json.Number `json:"last_entry_id"`
			Field1      string      `json:"field1"`
			Field2      string      `json:"field2"`
			Field3      string      `json:"field3"`
			Field4      string      `json:"field4"`
			Field5      string      `json:"field5"`
			Field6      string      `json:"field6"`
			Field7      string      `json:"field7"`
			Field8      string      `json:"field8"`
		} `json:"channel"`
		Feeds []struct {
			EntryId   json.Number `json:"entry_id"`
			CreatedAt time.Time   `json:"created_at"`
			Field1    string      `json:"field1"`
			Field2    string      `json:"field2"`
			Field3    string      `json:"field3"`
			Field4    string      `json:"field4"`
			Field5    string      `json:"field5"`
			Field6    string      `json:"field6"`
			Field7    string      `json:"field7"`
			Field8    string      `json:"field8"`
		} `json:"feeds"`
	}

	if err := json.Unmarshal(content, &sensorFeedJSON); err != nil {
		// log.Printf("\nChannel: %v, Err: %v", channel["id"], err)
		return
	}

	if sensorFeedJSON.Channel.LastEntryId.String() == channel["last_entry_id"] {
		// fmt.Printf("Already updated")
		return
	}

	ch := &sensorFeedJSON.Channel
	fieldNames := []string{ch.Field1, ch.Field2, ch.Field3, ch.Field4, ch.Field5, ch.Field6, ch.Field7, ch.Field8}

	dataPointsCounts := int64(0)
	extractedSensorsCount := int64(0)
	for _, rec := range sensorFeedJSON.Feeds {

		// ideally a combination of fields (e.g. name, channel_id, ...) is required, but since the entry_id is unique in thingspeak, it is sufficient
		rows, _ := global.DB.Load("sensor_values", database.RowType{"entry_id": rec.EntryId})
		if rows != nil && len(rows) > 0 {
			// log.Printf("Already exist rec.EntryId: %v\n", rec.EntryId)
			continue
		}

		fieldValues := []string{rec.Field1, rec.Field2, rec.Field3, rec.Field4, rec.Field5, rec.Field6, rec.Field7, rec.Field8}
		for i := 0; i != 8; i++ {

			if fieldNames[i] == "" {
				continue
			}

			// Process the sensor details:

			sensorId := int64(0)
			sensorRow := database.RowType{
				"channel_id": sensorFeedJSON.Channel.Id,
				"name":       fieldNames[i],
			}

			rows, _ := global.DB.Load("sensors", sensorRow)
			if rows != nil && len(rows) > 0 {
				sensorId = rows[0]["id"].(int64)
			} else {

				insRes, err := global.DB.Insert("sensors", sensorRow)
				if err != nil {
					if !strings.Contains(err.Error(), "duplicate key") {
						log.Printf("\nError in sensor insertion: %v \nRow: \n%v", err, sensorRow)
					}
				}
				sensorId = insRes.LastInsertId
				extractedSensorsCount += insRes.RowsAffected
			}

			if sensorId == 0 {
				log.Printf("\nError: sensor id (LastInsertId) is zero! \nRow: \n%v", sensorRow)
				continue
			}

			/*------------*/

			row := database.RowType{
				"entry_id":   rec.EntryId,
				"created_at": rec.CreatedAt,
				"value":      fieldValues[i],
				"sensor_id":  sensorId,
			}
			insRes, err := global.DB.Insert("sensor_values", row)
			if err != nil {
				// Let's ignore duplicate key as some user's have used same sensor name twice or more
				if !strings.Contains(err.Error(), "duplicate key") {
					log.Printf("\nError in sensor_value insertion: %v", err)
				}
			}

			dataPointsCounts += insRes.RowsAffected
		}
	}

	if dataPointsCounts > 0 {
		_, err := global.DB.Update("channels", database.RowType{"last_entry_id": sensorFeedJSON.Channel.LastEntryId}, database.RowType{"id": sensorFeedJSON.Channel.Id})
		if err != nil {
			log.Printf("\nError in data update: %v", err)
		}

		global.DataCollectorProgress.NewExtractedSensorValues += dataPointsCounts
	}

	if extractedSensorsCount > 0 {
		global.DataCollectorProgress.NewExtractedSensors += extractedSensorsCount
	}

}

/*--------------------------------*/

func processChannelDataExtraction(page int) {

	apiURL := "https://api.thingspeak.com/channels/public.json?page="

	req, err := http.NewRequest("GET", fmt.Sprintf("%s%d", apiURL, page), nil)
	if err != nil {
		log.Printf("\nError in channel extraction: %v", err)
		return
	}

	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("\nError in channel extraction: %v", err)
		return
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("\nCould not read the content: %v", err)
		return
	}

	var channelsJSON struct {
		Channels []struct {
			Id json.Number `json:"id"`
			// LastEntryId json.Number `json:"last_entry_id"`
			Name        string    `json:"name"`
			Description string    `json:"description"`
			Latitude    string    `json:"latitude"`
			Longitude   string    `json:"longitude"`
			CreatedAt   time.Time `json:"created_at"`
			URL         string    `json:"url"`
		} `json:"channels"`
	}
	if err := json.Unmarshal(content, &channelsJSON); err != nil {

		log.Printf("Page: %v, Err: %v", page, err)
		return
	}

	if len(channelsJSON.Channels) == 0 {
		// fmt.Printf("\nAll Done [Page: %d ] \n\n", page)
		channelsHitTheLastPage = true
		return
	}

	for _, rec := range channelsJSON.Channels {

		fields := database.RowType{
			"id":            rec.Id.String(),
			"last_entry_id": "0", // We keep this Zero for the first time, later it will be updated through sensor data extraction
			"name":          rec.Name,
			"description":   rec.Description,
			"latitude":      rec.Latitude,
			"longitude":     rec.Longitude,
			"url":           rec.URL,
			"created_at":    rec.CreatedAt,
		}

		rows, _ := global.DB.Load("channels", database.RowType{"id": fields["id"]})
		if rows != nil && len(rows) > 0 {
			// log.Printf("Already exist\n")
			continue
		}

		insRes, err := global.DB.Insert("channels", fields)
		if err != nil {
			log.Printf("\nError in data insertion: %v", err)
		}
		global.DataCollectorProgress.NewExtractedChannels += insRes.RowsAffected

	}
	fmt.Printf("\rPage %-5d done", page)

}

/*--------------------------------*/

package datacollection

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"real-sensor-data/database"
	"real-sensor-data/global"
	"strings"
	"time"
)

/*--------------------------------*/

// This has to be initiated manually (e.g. in main)
func Init() {

	go func() {
		for {

			ExtractChannelsData()

			time.Sleep(30 * time.Second)
			ExtractSensorsData()

			time.Sleep(60 * time.Minute)
		}
	}()
}

/*--------------------------------*/

func ExtractChannelsData() {

	fmt.Print("\n\t\t* * * Extracting new channels * * *\n\n")

	apiURL := "https://api.thingspeak.com/channels/public.json?page="
	global.DataCollectorProgress.ChannelsRunning = true
	defer func() { global.DataCollectorProgress.ChannelsRunning = false }()

	newChannelsCount := 0
	for page := 1; ; page++ {

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
			panic(err)
		}

		if len(channelsJSON.Channels) == 0 {
			fmt.Printf("\nDone [new channels: %d ] \n\n", newChannelsCount)
			break
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

			newChannelsCount += int(insRes.RowsAffected)
		}
		fmt.Printf("\rPage %-5d done", page)
	}

	fmt.Printf("\n\nAll Done :)\n\n---------------------------------------------------------\n")
}

/*--------------------------------*/

func ExtractSensorsData() {

	apiURL := "https://thingspeak.com/channels/%v/feed.json"

	global.DataCollectorProgress.SensorsRunning = true
	defer func() { global.DataCollectorProgress.SensorsRunning = false }()

	fmt.Print("\n\t\t* * * Extracting new sensor data * * *\n\n")

	channels, err := global.DB.Load("channels", nil)
	if err != nil {
		log.Fatal(err)
	}

	totalChannels := float64(len(channels))
	totalDataPointsCounts := 0
	global.DataCollectorProgress.SensorsProgress = 0
	for chIndex, channel := range channels {

		progress := int(math.Round(100 * (float64(chIndex) + 1) / totalChannels))
		global.DataCollectorProgress.SensorsProgress = progress
		fmt.Printf("\r\t[ %3v %% ]\tProcessing ( id: %-20v) ...", progress, channel["id"])
		req, err := http.NewRequest("GET", fmt.Sprintf(apiURL, channel["id"]), nil)
		if err != nil {
			log.Printf("Error in sensor data extraction: %v", err)
			return
		}

		req.Header.Set("Content-Type", "application/json; charset=UTF-8")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("Error in sensor data extraction: %v", err)
			return
		}
		defer resp.Body.Close()

		content, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Could not read the content: %v", err)
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
			panic(err)
		}

		if sensorFeedJSON.Channel.LastEntryId.String() == channel["last_entry_id"] {
			fmt.Printf("Already updated")
			continue
		}

		ch := &sensorFeedJSON.Channel
		fieldNames := []string{ch.Field1, ch.Field2, ch.Field3, ch.Field4, ch.Field5, ch.Field6, ch.Field7, ch.Field8}

		dataPointsCounts := 0
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

				row := database.RowType{
					"entry_id":   rec.EntryId,
					"channel_id": sensorFeedJSON.Channel.Id,
					"created_at": rec.CreatedAt,
					"name":       fieldNames[i],
					"value":      fieldValues[i],
				}
				insRes, err := global.DB.Insert("sensor_values", row)
				if err != nil {
					// Let's ignore duplicate key as some user's have used same sensor name twice or more
					if !strings.Contains(err.Error(), "duplicate key") {
						log.Printf("\nError in data insertion: %v", err)
					}
				}

				dataPointsCounts += int(insRes.RowsAffected)
			}
		}

		if dataPointsCounts > 0 {
			_, err := global.DB.Update("channels", database.RowType{"last_entry_id": sensorFeedJSON.Channel.LastEntryId}, database.RowType{"id": sensorFeedJSON.Channel.Id})
			if err != nil {
				log.Printf("\nError in data update: %v", err)
			}
		}

		fmt.Printf("\tDone [new values: %d ]", dataPointsCounts)
		totalDataPointsCounts += dataPointsCounts

		// fmt.Printf("Page %d done\n", page)
	}

	global.DataCollectorProgress.SensorsProgress = 100

	fmt.Printf("\n\nAll Done [ total new data points: %-6d ] :)\n\n---------------------------------------------------------\n", totalDataPointsCounts)
}

/*--------------------------------*/

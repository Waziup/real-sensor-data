package global

import (
	"real-sensor-data/database"
	"time"
)

/*-------------*/

var DB *database.Database // initiated in the main package

const RowsPerPage = 200 // This is the number of rows that APIs show per page

const MaxNumGoRoutines = 64 // Max number of concurent threads (mostly for data collection)

/*-------------*/

var DataCollectorProgress struct {
	ChannelsRunning bool
	SensorsRunning  bool
	SensorsProgress int

	NewExtractedChannels     int64
	NewExtractedSensorValues int64
	LastExtractionTime       time.Time
}

/*-------------*/

func init() {

	DataCollectorProgress.ChannelsRunning = false
	DataCollectorProgress.SensorsRunning = false
	DataCollectorProgress.SensorsProgress = 0
	DataCollectorProgress.NewExtractedChannels = 0
	DataCollectorProgress.NewExtractedSensorValues = 0
}

/*-------------*/

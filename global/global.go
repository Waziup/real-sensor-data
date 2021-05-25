package global

import (
	"os"
	"sensor-data-simulator/database"
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
	NewExtractedSensors      int64
	NewExtractedSensorValues int64
	LastExtractionTime       time.Time
}

/*-------------*/

var ENV struct {
	SERVING_ADDR             string
	DATA_EXTRACTION_INTERVAL string
	POSTGRES_DB              string
	POSTGRES_USER            string
	POSTGRES_PASSWORD        string
	POSTGRES_PORT            string
	POSTGRES_HOST            string
	WAZIUP_API_PATH          string
}

/*-------------*/

func init() {

	/*----------*/

	DataCollectorProgress.ChannelsRunning = false
	DataCollectorProgress.SensorsRunning = false
	DataCollectorProgress.SensorsProgress = 0

	DataCollectorProgress.NewExtractedChannels = 0
	DataCollectorProgress.NewExtractedSensors = 0
	DataCollectorProgress.NewExtractedSensorValues = 0

	/*----------*/

	ENV.SERVING_ADDR = os.Getenv("SERVING_ADDR")
	ENV.DATA_EXTRACTION_INTERVAL = os.Getenv("DATA_EXTRACTION_INTERVAL")
	ENV.POSTGRES_DB = os.Getenv("POSTGRES_DB")
	ENV.POSTGRES_USER = os.Getenv("POSTGRES_USER")
	ENV.POSTGRES_PASSWORD = os.Getenv("POSTGRES_PASSWORD")
	ENV.POSTGRES_PORT = os.Getenv("POSTGRES_PORT")
	ENV.POSTGRES_HOST = os.Getenv("POSTGRES_HOST")
	ENV.WAZIUP_API_PATH = os.Getenv("WAZIUP_API_PATH")

	/*----------*/
}

/*-------------*/

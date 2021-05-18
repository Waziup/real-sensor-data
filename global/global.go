package global

import "real-sensor-data/database"

/*-------------*/

var DB *database.Database // initiated in the main package

/*-------------*/

var DataCollectorProgress struct {
	ChannelsRunning bool
	SensorsRunning  bool
	SensorsProgress int
}

/*-------------*/

func init() {

	DataCollectorProgress.ChannelsRunning = false
	DataCollectorProgress.SensorsRunning = false
	DataCollectorProgress.SensorsProgress = 0
}

/*-------------*/

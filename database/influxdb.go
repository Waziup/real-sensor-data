package database

import (
	"context"
	"fmt"
	"os"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	// influxdb2 "github.com/influxdata/influxdb/client/v2"
)

/*-----------------*/

func (db *Database) InfluxClose() {
	db.InfluxClient.Close()
}

/*-----------------*/

func NewInfluxDB() influxdb2.Client {
	client := influxdb2.NewClient(os.Getenv("INFLUXDB_ADDR"), os.Getenv("INFLUXDB_TOKEN"))

	return client
}

/*-----------------*/

// var writeAPI influxdb2.Client

func (db *Database) InfluxInsert(measurement string, tags map[string]string, fields map[string]interface{}) (ExecResult, error) {

	// db.InfluxClient = influxdb2.NewClient(os.Getenv("INFLUXDB_ADDR"), os.Getenv("INFLUXDB_TOKEN"))
	//
	writeAPI := db.InfluxClient.WriteAPI(os.Getenv("INFLUXDB_ORG"), os.Getenv("INFLUXDB_BUCKET"))

	var recordTime time.Time
	if fields["time"] != nil {
		recordTime = fields["time"].(time.Time)
		delete(fields, "time")

	} else {

		recordTime = time.Now()
	}

	if len(fields) == 0 {
		fields["_empty"] = 0
	}

	// fmt.Printf("\nTable: %v\n\n", measurement)
	// fmt.Printf("\nTAGS: %v\n\n", tags)
	// fmt.Printf("FIELDS: %v\n\n", fields)
	// fmt.Printf("recordTime: %v\n", recordTime)

	newPoint := influxdb2.NewPoint(measurement,
		tags,
		fields,
		recordTime)

	writeAPI.WritePoint(newPoint)
	writeAPI.Flush()

	return ExecResult{}, nil
}

/*-----------------*/

func (db *Database) InfluxLoad(measurement string, tags RowType) (QueryResult, error) {
	//TODO: Not finished yet!
	query := fmt.Sprintf("from(bucket:\"%v\") |> range(start:-1000y) |> filter(fn: (r) => r._measurement == \"%v\")", os.Getenv("INFLUXDB_BUCKET"), measurement)
	return db.InfluxQuery(query)
}

/*-----------------*/

func (db *Database) InfluxQuery(query string) (QueryResult, error) {

	var output QueryResult

	queryAPI := db.InfluxClient.QueryAPI(os.Getenv("INFLUXDB_ORG"))

	result, err := queryAPI.Query(context.Background(), query)

	if err != nil {
		return nil, err
	}

	for result.Next() {
		// Notice when group key has changed
		// if result.TableChanged() {
		// 	fmt.Printf("table: %s\n", result.TableMetadata().String())
		// }

		output = append(output, result.Record().Values())
	}

	if result.Err() != nil {
		return nil, result.Err()
		// fmt.Printf("query parsing error: %v\n", result.Err().Error())
	}

	return output, nil
}

/*-----------------*/

package database

import (
	"database/sql"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type DBType int

const (
	InfluxDB DBType = iota
	Postgres
	// MySQL
	// MongoDB
)

type Database struct {
	Type         DBType // influx, mysql, sqlite,...
	InfluxClient influxdb2.Client
	SQLConn      *sql.DB
	// MySQLConn ...
}

type ExecResult struct {
	RowsAffected int64
	LastInsertId int64
}

type RowType map[string]interface{}
type QueryResult []RowType
type QueryParams []interface{}

/*-----------------------*/

func New(DatabaseType DBType, params ...string) *Database {
	var newDB Database

	newDB.Type = DatabaseType

	switch DatabaseType {
	case InfluxDB:
		newDB.InfluxClient = NewInfluxDB()
	case Postgres:
		if len(params) == 0 {
			return nil
		}
		newDB.SQLConn = NewPostgresDB(params[0])
		newDB.PostgresInit()
	}

	return &newDB
}

/*-----------------------*/

func (db *Database) Close() {
	switch db.Type {
	case InfluxDB:
		db.InfluxClose()
	case Postgres:
		db.PostgresClose()
	}
}

/*-----------------------*/

func (db *Database) Insert(table string, fields RowType, tags ...map[string]string) (ExecResult, error) {

	tagsForInflux := make(map[string]string)
	if len(tags) > 0 {
		tagsForInflux = tags[0]
	}

	switch db.Type {
	case InfluxDB:
		return db.InfluxInsert(table, tagsForInflux, fields)
	case Postgres:
		return db.PostgresInsert(table, fields)
	}

	return ExecResult{}, nil //TODO: provide a useful error here
}

/*-----------------------*/

func (db *Database) Update(table string, fields RowType, conditions RowType) (ExecResult, error) {

	switch db.Type {
	case InfluxDB:
		return ExecResult{}, nil // Not implemented
	case Postgres:
		return db.PostgresUpdate(table, fields, conditions)
	}

	return ExecResult{}, nil //TODO: provide a useful error here
}

/*-----------------------*/

func (db *Database) Delete(table string, conditions RowType) (ExecResult, error) {

	switch db.Type {
	case InfluxDB:
		return ExecResult{}, nil // Not implemented
	case Postgres:
		return db.PostgresDelete(table, conditions)
	}

	return ExecResult{}, nil //TODO: provide a useful error here
}

/*-----------------------*/

func (db *Database) Load(table string, searchOnFields RowType) (QueryResult, error) {

	switch db.Type {
	case InfluxDB:
		return db.InfluxLoad(table, searchOnFields)
	case Postgres:
		return db.PostgresLoad(table, searchOnFields)
	}

	return QueryResult{}, nil //TODO: provide a useful error here

}

/*-----------------------*/

func (db *Database) Query(query string, params QueryParams) (QueryResult, error) {

	switch db.Type {
	case InfluxDB:
		return db.InfluxQuery(query)
	case Postgres:
		return db.PostgresQuery(query, params)
	}

	return QueryResult{}, nil //TODO: provide a useful error here

}

/*-----------------------*/

func (db *Database) Exec(query string, params QueryParams) (ExecResult, error) {

	switch db.Type {
	case InfluxDB:
		return ExecResult{}, nil // Not implemented
	case Postgres:
		return db.PostgresExec(query, params)
	}

	return ExecResult{}, nil //TODO: provide a useful error here

}

/*-----------------------*/

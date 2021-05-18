package database

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

/*-----------------*/

func NewPostgresDB() *sql.DB {

	psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)

	// open database
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		panic(err)
	}

	return db
}

/*-----------------*/

func (db *Database) PostgresClose() {
	db.SQLConn.Close()
}

/*-----------------*/

func (db *Database) PostgresInsert(table string, fields RowType) (ExecResult, error) {

	SQL := fmt.Sprintf(`INSERT INTO "%s" (`, table)

	var params QueryParams
	values := ""
	paramCounter := 1
	for fieldName, value := range fields {
		SQL += fmt.Sprintf(`"%s",`, fieldName)
		values += fmt.Sprintf(`$%d,`, paramCounter)
		paramCounter++
		params = append(params, value)
	}

	SQL = strings.Trim(SQL, ",")
	values = strings.Trim(values, ",")
	SQL += ") VALUES ( " + values + ")"

	return db.PostgresExec(SQL, params)
}

/*-----------------*/

func (db *Database) PostgresUpdate(table string, fields RowType, conditions RowType) (ExecResult, error) {

	SQL := fmt.Sprintf(`UPDATE "%s" SET `, table)

	var params QueryParams
	paramCounter := 1

	for fieldName, value := range fields {
		SQL += fmt.Sprintf(`"%s" = $%d,`, fieldName, paramCounter)
		paramCounter++
		params = append(params, value)
	}

	SQL = strings.Trim(SQL, ",")
	SQL += " WHERE 1 = 1 "

	for fieldName, value := range conditions {
		SQL += fmt.Sprintf(` AND "%s" = $%d `, fieldName, paramCounter)
		paramCounter++
		params = append(params, value)
	}

	return db.PostgresExec(SQL, params)

}

/*-----------------*/

func (db *Database) PostgresExec(query string, params QueryParams) (ExecResult, error) {

	res, err := db.SQLConn.Exec(query, params...)
	if err != nil {
		return ExecResult{}, err
	}

	var output ExecResult

	output.RowsAffected, _ = res.RowsAffected()
	output.LastInsertId, _ = res.LastInsertId()

	return output, nil
}

/*-----------------*/

func (db *Database) PostgresLoad(table string, searchOnFields RowType) (QueryResult, error) {

	SQL := fmt.Sprintf(`SELECT * FROM "%s" WHERE 1 = 1 `, table)

	var params QueryParams
	paramCounter := 1
	for fieldName, value := range searchOnFields {
		SQL += fmt.Sprintf(` AND "%s" = $%d `, fieldName, paramCounter)
		paramCounter++
		params = append(params, value)
	}

	// query := fmt.Sprintf("from(bucket:\"%v\") |> range(start:-1000y) |> filter(fn: (r) => r._measurement == \"%v\")", os.Getenv("PostgresDB_BUCKET"), measurement)
	// return db.PostgresQuery(query)
	return db.PostgresQuery(SQL, params)
}

/*-----------------*/

func (db *Database) PostgresQuery(query string, params QueryParams) (QueryResult, error) {

	var output QueryResult

	rows, err := db.SQLConn.Query(query, params...)
	if err != nil {
		return output, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return output, err
	}

	colCounts := len(columns)
	values := make([]interface{}, colCounts)
	scanArgs := make([]interface{}, colCounts)

	for i := range values {
		scanArgs[i] = &values[i]
	}

	rowCount := 0
	for rows.Next() {

		err = rows.Scan(scanArgs...)
		if err != nil {
			return output, err
		}

		output = append(output, make(RowType, colCounts))
		for i, v := range values {
			output[rowCount][columns[i]] = v
		}
		rowCount++
	}

	return output, nil
}

/*-----------------*/

func (db *Database) PostgresInit() error {

	// fmt.Print("Postgres Init")
	/*
		`-- Table: public.channels

		-- DROP TABLE public.channels;

		CREATE TABLE public.channels
		(
			created_at timestamp without time zone NOT NULL,
			description character varying(400) COLLATE pg_catalog."default" NOT NULL,
			id bigint NOT NULL,
			latitude double precision NOT NULL,
			longitude double precision NOT NULL,
			name character varying(255) COLLATE pg_catalog."default" NOT NULL,
			url character varying(255) COLLATE pg_catalog."default" NOT NULL,
			last_entry_id bigint NOT NULL,
			CONSTRAINT channels_pkey PRIMARY KEY (id)
		)

		TABLESPACE pg_default;

		ALTER TABLE public.channels
			OWNER to root;


		-- Table: public.sensor_values

		-- DROP TABLE public.sensor_values;

		CREATE TABLE public.sensor_values
		(
			entry_id bigint NOT NULL,
			name character varying(255) COLLATE pg_catalog."default" NOT NULL,
			value character varying(100) COLLATE pg_catalog."default",
			created_at timestamp without time zone NOT NULL,
			channel_id bigint NOT NULL,
			CONSTRAINT sensor_values_pkey PRIMARY KEY (entry_id, name)
		)

		TABLESPACE pg_default;

		ALTER TABLE public.sensor_values
			OWNER to root;
		-- Index: channel_Id

		-- DROP INDEX public."channel_Id";

		CREATE INDEX "channel_Id"
			ON public.sensor_values USING btree
			(channel_id ASC NULLS LAST)
			TABLESPACE pg_default;
		-- Index: entry_id

		-- DROP INDEX public.entry_id;

		CREATE INDEX entry_id
			ON public.sensor_values USING btree
			(entry_id ASC NULLS LAST)
			TABLESPACE pg_default;		`
	*/

	return nil

}

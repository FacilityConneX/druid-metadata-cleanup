package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var db *sql.DB
var err error

func initDB(dbURL string) error {
	db, err = sql.Open("pgx", dbURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return err
	}

	pingErr := db.Ping()
	if pingErr != nil {
		fmt.Fprintf(os.Stderr, "Ping failed: %v\n", pingErr)
		return err
	}

	fmt.Println("Connected to Postgres!")
	return nil
}

type TaskMetadata struct {
	id string
	createdDate string
}

func queryMetadata(endTime string) ([]TaskMetadata, error) {
	query := fmt.Sprintf(
		"SELECT id, created_date FROM public.druid_tasks WHERE created_date < '%v' AND active = false;",
		endTime,
	)
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := []TaskMetadata{}
	for rows.Next() {
		var r TaskMetadata
		err := rows.Scan(&r.id, &r.createdDate)
		if err != nil {
			return nil, err
		}

		results = append(results, r)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func deleteMetadata(endTime string) error {
	query := fmt.Sprintf(
		"DELETE FROM public.druid_tasks WHERE created_date < '%v' AND active = false;", 
		endTime,
	)
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}
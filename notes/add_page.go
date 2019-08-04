package main

import (
	"database/sql"
	"net/http"

	requests "./request"
)

var addPageStatement = `
	insert into page (
		name,
		createdAt,
		updatedAt,
	deletedAt	
	) values (
		$1,
		date('now'),
		date('now'),
		null
	)`

func addPage(db *sql.DB) requests.RouteHandler {
	return func(req requests.Request, w http.ResponseWriter) (interface{}, error) {
		res, err := db.Exec(addPageStatement, req.Body["name"])

		if err != nil {
			return nil, err
		}

		id, _ := res.LastInsertId()
		rows, err := db.Query(`select rowid as id, name, createdAt from page where id is $1`, id)

		if err != nil {
			return nil, err
		}

		rows.Next()
		defer rows.Close()
		return mapPage(rows), nil
	}
}

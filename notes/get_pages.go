package main

import (
	"database/sql"
	"net/http"

	requests "./request"
)

func getPages(db *sql.DB) requests.RouteHandler {
	return func(req requests.Request, w http.ResponseWriter) (interface{}, error) {
		var rows *sql.Rows
		var err error
		if req.Parameters["parent_id"] != "" {
			rows, err = db.Query(`select rowid as id, name, createdAt from page where parentID is $1`, req.Parameters["parent_id"])
		} else {
			rows, err = db.Query(`select rowid as id, name, createdAt from page where parentID is null`)
		}

		if err != nil {
			return nil, err
		}

		pages := mapPages(rows)
		return pages, nil
	}
}

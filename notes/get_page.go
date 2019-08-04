package main

import (
	"database/sql"
	"net/http"

	requests "./request"
)

func getPage(db *sql.DB) requests.RouteHandler {
	return func(req requests.Request, w http.ResponseWriter) (interface{}, error) {
		page, err := fetchPage(db, req.Parameters["id"])
		if err != nil {
			return nil, err
		}
		return page, nil
	}
}

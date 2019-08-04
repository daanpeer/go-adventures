package main

import (
	"database/sql"
	"net/http"

	requests "./request"
)

func deletePage(db *sql.DB) requests.RouteHandler {
	return func(req requests.Request, w http.ResponseWriter) (interface{}, error) {
		page, err := fetchPage(db, req.Parameters["id"])
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, &requests.NotFoundError{}
			}
			return nil, err
		}

		_, error := db.Exec(`
			delete
			from page
			where _ROWID_ = $1
		`, &page.ID)

		if error != nil {
			return nil, error
		}

		return page, nil
	}
}

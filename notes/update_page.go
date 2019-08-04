package main

import (
	"database/sql"
	"net/http"

	requests "./request"
)

func updatePage(db *sql.DB) requests.RouteHandler {
	return func(req requests.Request, w http.ResponseWriter) (interface{}, error) {
		page, err := fetchPage(db, req.Parameters["id"])
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, &requests.NotFoundError{}
			}
			return nil, err
		}

		if req.Body["name"] != "" {
			page.Name = req.Body["name"]
		}

		_, err = db.Exec("update page set name = ? where _ROWID_ = ?", page.Name, page.ID)
		if err != nil {
			return nil, err
		}

		return page, nil
	}
}

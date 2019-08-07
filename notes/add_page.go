package main

import (
	// "database/sql"
	"encoding/json"
	"net/http"

	requests "./request"
)

func addPage(p *PageRepository) requests.RouteHandler {
	return func(req requests.Request, w http.ResponseWriter) (interface{}, error) {
		page := &Page{}
		err := json.Unmarshal(req.Body, page)

		page, err = p.InsertPage(page)

		if err != nil {
			return nil, err
		}

		return page, nil
	}
}

package main

import (
	"net/http"
	"strconv"

	requests "./request"
)

func updatePage(p *PageRepository) requests.RouteHandler {
	return func(req requests.Request, w http.ResponseWriter) (interface{}, error) {
		id, err := strconv.Atoi(req.Parameters["id"])
		// @TODO throw 403?
		if err != nil {
			return nil, err
		}

		page, error := p.UpdatePage(id, req.Body)

		if error != nil {
			return nil, error
		}

		return page, nil

	}
}

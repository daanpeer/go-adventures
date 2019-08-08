package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	requests "./request"
)

func updatePage(p *PageRepository) requests.RouteHandler {
	return func(req requests.Request, w http.ResponseWriter) (interface{}, error) {
		id, err := strconv.Atoi(req.Parameters["id"])
		if err != nil {
			return nil, &requests.UnprocessableEntity{}
		}

		page := &Page{}
		err = json.Unmarshal(req.Body, page)

		if err != nil {
			return nil, &requests.UnprocessableEntity{}
		}

		newPage, err := p.UpdatePage(id, page)
		if newPage == nil {
			return nil, err
		}

		if err != nil {
			return nil, err
		}

		return newPage, nil
	}
}

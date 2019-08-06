package main

import (
	"net/http"

	requests "./request"
)

func addPage(p *PageRepository) requests.RouteHandler {
	return func(req requests.Request, w http.ResponseWriter) (interface{}, error) {
		if req.Body["name"] == "" {
			return nil, &requests.UnprocessableEntity{}
		}

		page, err := p.InsertPage(req.Body["name"])

		if err != nil {
			return nil, err
		}

		return page, nil
	}
}

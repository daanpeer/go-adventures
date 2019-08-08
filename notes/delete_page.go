package main

import (
	"net/http"
	"strconv"

	requests "./request"
)

func deletePage(p *PageRepository) requests.RouteHandler {
	return func(req requests.Request, w http.ResponseWriter) (interface{}, error) {
		id, err := strconv.Atoi(req.Parameters["id"])

		if err != nil {
			return nil, &requests.UnprocessableEntity{}
		}

		page, err := p.DeletePage(id)
		if err != nil {
			return nil, err
		}

		return page, nil
	}
}

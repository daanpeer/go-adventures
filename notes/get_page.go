package main

import (
	"net/http"
	"strconv"

	requests "./request"
)

func getPage(pr *PageRepository) requests.RouteHandler {
	return func(req requests.Request, w http.ResponseWriter) (interface{}, error) {
		id, err := strconv.Atoi(req.Parameters["id"])

		if err != nil {
			return nil, err
		}

		page, err := pr.FindPageById(id)
		if err != nil {
			return nil, err
		}
		return page, nil
	}
}

package main

import (
	"net/http"
	"strconv"

	requests "./request"
)

func getPages(pr *PageRepository) requests.RouteHandler {
	return func(req requests.Request, w http.ResponseWriter) (interface{}, error) {
		var pages []Page
		var err error

		if req.Parameters["parent_id"] != "" {
			id, err := strconv.Atoi(req.Parameters["parent_id"])
			// @TODO throw 403?
			if err != nil {
				return nil, err
			}
			pages, err = pr.FindPagesByParent(id)
		} else {
			pages, err = pr.FindParent()
		}

		if err != nil {
			return nil, err
		}

		return pages, nil
	}
}

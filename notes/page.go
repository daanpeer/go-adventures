package main

import (
	"database/sql"
	"time"
)

// Page page model
type Page struct {
	ID        int
	Name      string
	CreatedAt time.Time
	ParentID  int
}

func mapPage(rows *sql.Rows) Page {
	page := Page{}
	err := rows.Scan(&page.ID, &page.Name, &page.CreatedAt)
	if err != nil {
		panic(err)
	}
	return page
}

func mapPages(rows *sql.Rows) []Page {
	var pages []Page
	defer rows.Close()
	for rows.Next() {
		pages = append(pages, mapPage(rows))
	}
	return pages
}

func fetchPage(db *sql.DB, id string) (*Page, error) {
	page := Page{}
	err := db.QueryRow(`
		select rowid as id, name, createdAt
		from page
		where _ROWID_ = $1
	`, id).Scan(&page.ID, &page.Name, &page.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &page, err
}

package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	requests "./request"
	_ "github.com/mattn/go-sqlite3"
)

func createDb(db *sql.DB) {
	sqlStmt := `
		create table page
		(
			name varchar,
			created_at datetime,
			updated_at datetime,
			deleted_at datetime,
			content json1,
			parent_id int,
			foreign key (parent_id) references page(_ROWID_)
		);
	`

	_, err := db.Exec(sqlStmt)

	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
}

func insertRows(db *sql.DB) {
	for i := 0; i < 10; i++ {
		res, _ := db.Exec(`insert into page (name, created_at, updated_at, deleted_at, parent_id) values ("test", date('now'), date('now'), null, null)`)
		for j := 0; j < 5; j++ {
			id, _ := res.LastInsertId()
			db.Exec(`insert into page (name, created_at, updated_at, deleted_at, parent_id) values ("sub", date('now'), date('now'), null, $1)`, id)
		}
	}
}

type Page struct {
	Id         int
	Name       string
	Created_at time.Time
	Parent_id  int
}

func mapPage(rows *sql.Rows) Page {
	page := Page{}
	err := rows.Scan(&page.Id, &page.Name, &page.Created_at)
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

func getPages(db *sql.DB) requests.RouteHandler {
	return func(req requests.Request, w http.ResponseWriter) (interface{}, error) {
		var rows *sql.Rows
		var err error
		if req.Parameters["parent_id"] != "" {
			rows, err = db.Query(`select rowid as id, name, created_at from page where parent_id is $1`, req.Parameters["parent_id"])
		} else {
			rows, err = db.Query(`select rowid as id, name, created_at from page where parent_id is null`)
		}

		if err != nil {
			return nil, err
		}

		return mapPages(rows), nil
	}
}

func addPage(db *sql.DB) requests.RouteHandler {
	return func(req requests.Request, w http.ResponseWriter) (interface{}, error) {
		fmt.Println(req.Body, req.Body["name"])
		res, err := db.Exec(`insert into page (
			name,
			created_at,
			updated_at,
			deleted_at
		) values (
			$1,
			date('now'),
			date('now'),
			null
		)`, req.Body["name"])

		if err != nil {
			return nil, err
		}

		id, _ := res.LastInsertId()
		rows, err := db.Query(`select rowid as id, name, created_at from page where id is $1`, id)

		if err != nil {
			return nil, err
		}

		rows.Next()
		defer rows.Close()
		return mapPage(rows), nil
	}
}

func main() {
	os.Remove("./notes.db")

	db, err := sql.Open("sqlite3", "./notes.db")

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createDb(db)
	insertRows(db)

	app := requests.HTTPServer{}
	app.Get("pages/:parent_id", getPages(db))
	app.Get("pages", getPages(db))
	app.Post("pages", addPage(db))
	app.Listen(":8080")
}

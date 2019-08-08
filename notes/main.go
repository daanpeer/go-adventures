package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	requests "./request"
	_ "github.com/mattn/go-sqlite3"
)

func createDb(db *sql.DB) {
	sqlStmt := `
		create table page
		(
			name varchar,
			createdAt datetime,
			updatedAt datetime,
			deletedAt datetime,
			content text,
			parentID int,
			foreign key (parentID) references page(_ROWID_)
		);
	`

	_, err := db.Exec(sqlStmt)

	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
}

// @TODO split content from pagetree

func main() {
	fmt.Println("Removing old db")
	os.Remove("./notes.db")

	fmt.Println("Creating new db")
	db, err := sql.Open("sqlite3", "./notes.db")

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createDb(db)

	go func() {
		fmt.Println("Inserting rows")
		generateDummy(db)
		fmt.Println("Done inserting rows")
	}()

	pr := &PageRepository{db: db}
	app := requests.HTTPServer{}
	app.Get("pages", getPages(pr))
	app.Get("pages/:id", getPage(pr))
	app.Post("pages", addPage(pr))
	app.Patch("pages/:id", updatePage(pr))
	app.Delete("pages/:id", deletePage(pr))
	app.Get("pages/sub/:parent_id", getPages(pr))
	err = app.Listen("localhost:8080")

	if err != nil {
		fmt.Println(err)
	}
}

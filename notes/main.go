package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

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
			content json1,
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

func insertRows(db *sql.DB) {
	for i := 0; i < 10; i++ {
		res, _ := db.Exec(`insert into page (name, createdAt, updatedAt, deletedAt, parentID) values ($1, date('now'), date('now'), null, null)`, "test-"+strconv.Itoa(i))
		ID, _ := res.LastInsertId()
		for j := 0; j < 5; j++ {
			db.Exec(`insert into page (name, createdAt, updatedAt, deletedAt, parentID) values ("sub", date('now'), date('now'), null, $1)`, ID)
		}
	}
}

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
		insertRows(db)
		fmt.Println("Done inserting rows")
	}()

	app := requests.HTTPServer{}
	app.Get("pages", getPages(db))
	app.Get("pages/:id", getPage(db))
	app.Post("pages", addPage(db))
	app.Patch("pages/:id", updatePage(db))
	app.Delete("pages/:id", deletePage(db))
	app.Get("pages/:parent_id", getPages(db))
	err = app.Listen(":8080")

	if err != nil {
		fmt.Println(err)
	}
}

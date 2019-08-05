package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
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

func insertPage(db *sql.DB, parameters map[string]interface{}) int64 {
	data, _ := ioutil.ReadFile("./testpage.md")
	content := string(data)
	res, err := db.Exec(`
		insert into page (
			name,
			content,
			createdAt,
			updatedAt,
			deletedAt,
			parentID
		) values (
			?,
			?,
			date('now'),
			date('now'),
			null,
			?
		)
	`, parameters["name"], content, parameters["parent_id"])
	if err != nil {
		panic(err)
	}
	ID, _ := res.LastInsertId()
	return ID
}

func nest(db *sql.DB, lastID int64, iteration int) {
	if iteration == 5 {
		return
	}

	var name = "sub"
	for index := 0; index < iteration; index++ {
		name += "sub"
	}

	iteration++
	data := map[string]interface{}{
		"name":      name,
		"parent_id": lastID,
	}
	ID := insertPage(db, data)

	lastID = ID
	nest(db, lastID, iteration)
}

func insertRows(db *sql.DB) {
	for i := 0; i < 10; i++ {
		data := map[string]interface{}{
			"name": "test-" + strconv.Atoi(i),
		}
		ID := insertPage(db, data)
		nest(db, ID, 1)
		for j := 0; j < 5; j++ {
			data := map[string]interface{}{
				"name":      "sub-" + string(i),
				"parent_id": ID,
			}

			insertPage(db, data)
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
	app.Get("pages/sub/:parent_id", getPages(db))
	err = app.Listen(":8080")

	if err != nil {
		fmt.Println(err)
	}
}

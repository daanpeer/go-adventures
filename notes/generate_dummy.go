package main

import (
	"database/sql"
	"io/ioutil"
	"strconv"
)

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

func generateDummy(db *sql.DB) {
	for i := 0; i < 10; i++ {
		data := map[string]interface{}{
			"name": "test-" + strconv.Itoa(i),
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

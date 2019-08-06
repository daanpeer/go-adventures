package main

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	requests "./request"
)

// @TODO return pointers to page instead so we can return nil

// Page page model
type Page struct {
	ID        int
	Name      string
	CreatedAt time.Time
	ParentID  int
	Content   string
}

type PageRepository struct {
	db *sql.DB
}

func (p *PageRepository) GetTable() string {
	return "page"
}

func (p *PageRepository) GetColumns() []string {
	return []string{"rowid as id", "name", "createdAt", "content"}
}

func (p *PageRepository) MapPage(rows *sql.Rows) Page {
	page := Page{}
	err := rows.Scan(&page.ID, &page.Name, &page.CreatedAt, &page.Content)
	if err != nil {
		panic(err)
	}
	return page
}

func (p *PageRepository) MapPages(rows *sql.Rows) []Page {
	pages := []Page{}
	defer rows.Close()
	for rows.Next() {
		pages = append(pages, p.MapPage(rows))
	}
	return pages
}

func (p *PageRepository) FindParent() ([]Page, error) {
	rows, err := p.db.Query(fmt.Sprintf(`select %s from %s where parentID is null`, strings.Join(p.GetColumns(), ","), p.GetTable()))

	if err != nil {
		return []Page{}, err
	}

	return p.MapPages(rows), err
}

func (p *PageRepository) FindPageById(id int) (Page, error) {
	rows, err := p.db.Query(fmt.Sprintf(`select %s from %s where _ROWID_ = ?`, strings.Join(p.GetColumns(), ","), p.GetTable()), id)

	if err != nil {
		return Page{}, err
	}

	rows.Next()
	defer rows.Close()
	return p.MapPage(rows), err
}

func (p *PageRepository) DeletePage(id int) (Page, error) {
	page, err := p.FindPageById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return Page{}, &requests.NotFoundError{}
		}
	}

	_, err = p.db.Exec(`
		delete
		from page
		where _ROWID_ = ?
	`, id)

	if err != nil {
		return Page{}, err
	}

	return page, nil
}

func (p *PageRepository) UpdatePage(id int, fields map[string]string) (Page, error) {
	page, err := p.FindPageById(id)
	if err != nil {
		return Page{}, err
	}

	if fields["name"] != "" {
		page.Name = fields["name"]
	}

	_, err = p.db.Exec(fmt.Sprintf("update %s set name = ? where _ROWID_ = ?", p.GetTable()), page.Name, page.ID)
	if err != nil {
		return Page{}, err
	}

	return page, nil
}

func (p *PageRepository) InsertPage(name string) (Page, error) {
	res, err := p.db.Exec(`
	insert into page (
		name,
		createdAt,
		updatedAt,
	deletedAt	
	) values (
		?,
		date('now'),
		date('now'),
		null
	)`, name)

	if err != nil {
		return Page{}, err
	}

	id, err := res.LastInsertId()

	if err != nil {
		return Page{}, nil
	}

	return p.FindPageById(int(id))
}

func (p *PageRepository) FindPagesByParent(id int) ([]Page, error) {
	rows, err := p.db.Query(fmt.Sprintf(`select %s from %s where parentID is $1`, strings.Join(p.GetColumns(), ","), p.GetTable()), id)
	if err != nil {
		return []Page{}, err
	}
	return p.MapPages(rows), err
}

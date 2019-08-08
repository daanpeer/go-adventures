package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	requests "./request"
)

// Page page model
type Page struct {
	ID        int
	Name      string
	CreatedAt time.Time
	ParentID  int
	Content   sql.NullString
}

// MarshalJSON custom marshalJSON for page type
func (page Page) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ID      int    `json:"id"`
		Name    string `json:"name"`
		Content string `json:"content"`
	}{
		ID:      page.ID,
		Name:    page.Name,
		Content: page.Content.String,
	})
}

// UnmarshalJSON unmarshal
func (page *Page) UnmarshalJSON(data []byte) error {
	type Alias Page
	aux := &struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		ParentID int    `json:"parent_id"`
		Content  string `json:"content"`
	}{}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	page.Content = sql.NullString{String: aux.Content}
	page.ID = aux.ID
	page.ParentID = aux.ParentID
	page.Name = aux.Name
	return nil
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
	var pages []Page
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

func (p *PageRepository) FindPageById(id int) (*Page, error) {
	rows, err := p.db.Query(fmt.Sprintf(`select %s from %s where _ROWID_ = ?`, strings.Join(p.GetColumns(), ","), p.GetTable()), id)

	if err != nil {
		return nil, err
	}

	pages := p.MapPages(rows)

	if len(pages) == 0 {
		return nil, nil
	}

	return &pages[0], err
}

func (p *PageRepository) DeletePage(id int) (*Page, error) {
	page, err := p.FindPageById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &requests.NotFoundError{}
		}
		return nil, err
	}

	_, err = p.db.Exec(`
		delete
		from page
		where _ROWID_ = ?
	`, id)

	if err != nil {
		return nil, err
	}

	return page, nil
}

func (p *PageRepository) UpdatePage(id int, newPage *Page) (*Page, error) {
	page, err := p.FindPageById(id)
	if err != nil {
		return nil, err
	}

	if page == nil {
		return nil, &requests.NotFoundError{}
	}

	_, err = p.db.Exec(fmt.Sprintf("update %s set name = ?, content = ? where _ROWID_ = ?", p.GetTable()), newPage.Name, newPage.Content.String, page.ID)
	if err != nil {
		return nil, err
	}

	page.Name = newPage.Name
	page.Content = newPage.Content

	return page, nil
}

func (p *PageRepository) InsertPage(page *Page) (*Page, error) {
	res, err := p.db.Exec(`
	insert into page (
		name,
		createdAt,
		updatedAt,
		deletedAt	
	) values (
		$1,
		date('now'),
		date('now'),
		null
	)`, page.Name)

	var newPage = &Page{}
	if err != nil {
		return newPage, err
	}

	id, err := res.LastInsertId()

	if err != nil {
		return newPage, nil
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

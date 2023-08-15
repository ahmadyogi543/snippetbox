package models

import (
	"database/sql"
	"errors"
	"time"
)

type SnippetModelInterface interface {
	Insert(title string, content string, expires int) (int, error)
	Get(id int) (*Snippet, error)
	Latest() ([]*Snippet, error)
}

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB *sql.DB
}

func (sm *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	query := `
		INSERT INTO snippets (title, content, created, expires)
		VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))
	`

	result, err := sm.DB.Exec(query, title, content, expires)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (sm *SnippetModel) Get(id int) (*Snippet, error) {
	query := `
		SELECT id, title, content, created, expires
		FROM snippets
		WHERE expires > UTC_TIMESTAMP() AND id = ?
	`

	snippet := &Snippet{}
	row := sm.DB.QueryRow(query, id)

	err := row.Scan(
		&snippet.ID,
		&snippet.Title,
		&snippet.Content,
		&snippet.Created,
		&snippet.Expires,
	)
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return snippet, nil
}

func (sm *SnippetModel) Latest() ([]*Snippet, error) {
	query := `
		SELECT id, title, content, created, expires
		FROM snippets
		WHERE expires > UTC_TIMESTAMP()
		ORDER BY id DESC LIMIT 10
	`

	rows, err := sm.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	snippets := []*Snippet{}
	for rows.Next() {

		snippet := &Snippet{}

		err := rows.Scan(
			&snippet.ID,
			&snippet.Title,
			&snippet.Content,
			&snippet.Created,
			&snippet.Expires,
		)
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, snippet)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}

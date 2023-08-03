package models

import (
	"database/sql"
	"errors"
	"time"
)

// Snippet struct field correspond to the field in the DB snippets table.
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// SnippetModel struct wraps a sql.DB connection pool.
type SnippetModel struct {
	DB *sql.DB
}

// insert new snippet into the DB.
func (sm *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	// SQL query with placeholder form title, content and expires.
	query := `
		INSERT INTO snippets (title, content, created, expires)
		VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))
	`

	// use sql.DB.Exec to execute the SQL query along with the value that will inserted.
	result, err := sm.DB.Exec(query, title, content, expires)
	if err != nil {
		return 0, err
	}

	// get the latest inserted row id, convert to int and return it.
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// return specific snippet based on its id.
func (sm *SnippetModel) Get(id int) (*Snippet, error) {
	// SQL query to retrieve a specific snippet with id and should not be expired.
	query := `
		SELECT id, title, content, created, expires
		FROM snippets
		WHERE expires > UTC_TIMESTAMP() AND id = ?
	`

	// use sql.DB.QueryRow to retrieve one row from the database
	// and create snippet struct to store the retrieved value.
	snippet := &Snippet{}
	row := sm.DB.QueryRow(query, id)

	// scan the result and copy it to the snippet field.
	err := row.Scan(
		&snippet.ID,
		&snippet.Title,
		&snippet.Content,
		&snippet.Created,
		&snippet.Expires,
	)
	if err != nil {
		// if error is no rows, return ErrNoRecord (own error implementation).
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return snippet, nil
}

// return the 10 most recently created snippets.
func (sm *SnippetModel) Latest() ([]*Snippet, error) {
	// SQL query to retrieve 10 most recently created snippets.
	query := `
		SELECT id, title, content, created, expires
		FROM snippets
		WHERE expires > UTC_TIMESTAMP()
		ORDER BY id DESC LIMIT 10
	`

	// use sql.DB.Query to retrieve some rows from the database.
	rows, err := sm.DB.Query(query)
	if err != nil {
		return nil, err
	}

	// need to run Close method on rows before exit the Latest method
	// to close DB connection from being used.
	defer rows.Close()

	// snippets slice to store snippets from the database
	// and loop for every row inside the rows that retrieved from the query.
	snippets := []*Snippet{}
	for rows.Next() {
		// represent each snippet
		snippet := &Snippet{}

		// scan the result and copy it to the snippet field.
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

		// append the scanned snippet to the snippets slice
		snippets = append(snippets, snippet)
	}

	// to be safe to check if something goes wrong when iterate
	// through rows above and if there is error, return it.
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}

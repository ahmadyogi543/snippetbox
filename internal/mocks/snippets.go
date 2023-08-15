package mocks

import (
	"time"

	"github.com/ahmadyogi543/snippetbox/internal/models"
)

var mockSnippet = &models.Snippet{
	ID:      1,
	Title:   "A Title",
	Content: "This is a content inside the mock snippet.",
	Created: time.Now(),
	Expires: time.Now(),
}

type SnippetModel struct{}

func (sm *SnippetModel) Insert(title string, content string, expires int) (int, error) {
	return 2, nil
}

func (sm *SnippetModel) Get(id int) (*models.Snippet, error) {
	switch id {
	case 1:
		return mockSnippet, nil
	default:
		return nil, models.ErrNoRecord
	}
}

func (sm *SnippetModel) Latest() ([]*models.Snippet, error) {
	return []*models.Snippet{mockSnippet}, nil
}

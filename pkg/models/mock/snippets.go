package mock

import (
	"time"

	"fredriksiemund/snippetbox/pkg/models"
)

var mockSnippet = &models.Snippet{
	ID:      1,
	Title:   "An old silent pond",
	Content: "An old silent pond...",
	Created: time.Now(),
	Expires: time.Now(),
}

type SnippetRepository struct{}

func (m *SnippetRepository) Insert(title, content, expires string) (int, error) {
	return 2, nil
}

func (m *SnippetRepository) Get(id int) (*models.Snippet, error) {
	switch id {
	case 1:
		return mockSnippet, nil
	default:
		return nil, models.ErrNoRecord
	}
}

func (m *SnippetRepository) Latest() ([]*models.Snippet, error) {
	return []*models.Snippet{mockSnippet}, nil
}

package models

import "time"

type SearchResult struct {
	Title       string
	Description string
	PageID      string
}

type Page struct {
	PageID    string
	Content   string
	Title     string
	CreatedAt time.Time
}

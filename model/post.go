package model

import "time"

type Post struct {
	Id          int64       `json:"id,omitempty"`
	AuthorId    int64       `json:"authorid"`
	ParentId    int64       `json:"parentid"`
	Title       string      `json:"title"`
	MetaTitle   string      `json:"metatitle"`
	Slug        string      `json:"slug"`
	Summary     string      `json:"summary"`
	Published   int16       `json:"published"`
	CreatedAt   time.Time   `json:"createdat"`
	UpdatedAt   time.Time   `json:"updatedat"`
	PublishedAt time.Time   `json:"publishedat"`
	Content     string      `json:"content"`
	Categories  *[]Category `json:"category,omitempty"`
	User        *User       `json:"user,omitempty"`
}

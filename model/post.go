package model

import "time"

/// parsable model
type Post struct {
	Title   string    `json:"title"`
	Href    string    `json:"href"`
	Date    time.Time `json:"date"`
	Author  string    `json:"author"`
	Summary string    `json:"summary"`
}

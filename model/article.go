package model

import "time"

/// parsable model
type Article struct {
	Title       string    `json:"title"`
	Created     time.Time `json:"created"`
	Author      string    `json:"author"`
	Description string    `json:"summary"`
	URL         string    `json:"url"`
}

package model

import "time"

/// parsable model
type Article struct {
	Title       string    `bson:"title"`
	Created     time.Time `bson:"created"`
	Author      string    `bson:"author"`
	Description string    `bson:"summary"`
	URL         string    `bson:"url"`
}

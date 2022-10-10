package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type DBArticle struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Title       string             `bson:"title"`
	Created     time.Time          `bson:"created"`
	Author      string             `bson:"author"`
	Description string             `bson:"summary"`
	URL         string             `bson:"url"`
}

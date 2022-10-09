package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DBArticle struct {
	ID primitive.ObjectID `bson:"_id"`
	Article
}

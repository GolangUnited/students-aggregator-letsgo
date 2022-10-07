package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Articles struct {
	ID primitive.ObjectID `bson:"_id"`
	Article
}

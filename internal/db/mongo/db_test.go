package mongo

import (
	"testing"
	"time"

	"github.com/indikator/aggregator_lets_go/internal/config"
	"github.com/indikator/aggregator_lets_go/model"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestWriteArticle(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	c := config.NewConfig()
	err := c.SetDataFromFile("../../configs/config.yaml")
	if err != nil {
		return
	}
	err = c.Read()
	if err != nil {
		return
	}
	mongoDb := NewDb(c.Database)
	defer mt.Close()

	mt.Run("write article", func(mt *mtest.T) {
		collection = mt.Coll
		id := primitive.NewObjectID()
		mt.AddMockResponses(mtest.CreateSuccessResponse())

		insertedArticle, err := mongoDb.WriteArticle(&model.DBArticle{
			ID:          id,
			Title:       "test_title",
			Created:     time.Date(2022, 1, 1, 1, 1, 1, 1, time.UTC),
			Author:      "mikhailov.mk",
			Description: "test article for db",
			URL:         "test_article.com",
		})

		assert.Nil(t, err)
		assert.Equal(t, &model.DBArticle{
			ID:          id,
			Title:       "test_title",
			Created:     time.Date(2022, 1, 1, 1, 1, 1, 1, time.UTC),
			Author:      "mikhailov.mk",
			Description: "test article for db",
			URL:         "test_article.com",
		}, insertedArticle)
	})

	mt.Run("custom error duplicate", func(mt *mtest.T) {
		collection = mt.Coll
		id := primitive.NewObjectID()
		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index:   1,
			Code:    11000,
			Message: "duplicate key error",
		}))

		insertedArticle, err := mongoDb.WriteArticle(&model.DBArticle{
			ID:          id,
			Title:       "test_title",
			Created:     time.Date(2022, 1, 1, 1, 1, 1, 1, time.UTC),
			Author:      "mikhailov.mk",
			Description: "test article for db",
			URL:         "test_article.com",
		})

		assert.Nil(t, insertedArticle)
		assert.NotNil(t, err)
		assert.True(t, mongo.IsDuplicateKeyError(err))
	})

}

func TestReadArticles(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	c := config.NewConfig()
	err := c.SetDataFromFile("../../tests/configs/mongo/config.yaml")
	if err != nil {
		return
	}
	mongoDb := NewDb(c.Database)
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		collection = mt.Coll
		expectedArticle := model.DBArticle{
			ID:          primitive.NewObjectID(),
			Title:       "test_title",
			Created:     time.Now().AddDate(0, -1, 0),
			Author:      "mikhailov.mk",
			Description: "test article for db",
			URL:         "test_article.com",
		}
		killCursors := mtest.CreateCursorResponse(0, "foo.bar", mtest.NextBatch)
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "foo.bar", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: expectedArticle.ID},
			{Key: "title", Value: expectedArticle.Title},
			{Key: "author", Value: expectedArticle.Author},
			{Key: "created", Value: expectedArticle.Created},
			{Key: "summary", Value: expectedArticle.Description},
			{Key: "url", Value: expectedArticle.URL},
		}), killCursors)

		articleResponse, err := mongoDb.ReadArticles(2)

		assert.Nil(t, err)
		assert.Equal(t, expectedArticle, articleResponse[0])
	})

}

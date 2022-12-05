package mongo

import (
	"fmt"
	"testing"
	"time"

	"github.com/indikator/aggregator_lets_go/internal/config"
	"github.com/indikator/aggregator_lets_go/internal/log/logLevel"
	log "github.com/indikator/aggregator_lets_go/internal/log/stub"
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
	err := c.SetDataFromFile("../../../tests/configs/mongo/config.yaml")
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
	err = c.Read()
	if err != nil {
		t.Errorf("expected nil, got %v", err)

	}
	l := log.NewLog(logLevel.Errors)
	mongoDb := NewDb(c.Database, l)
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
	err := c.SetDataFromFile("../../../tests/configs/mongo/config.yaml")
	if err != nil {
		return
	}
	l := log.NewLog(logLevel.Errors)
	mongoDb := NewDb(c.Database, l)
	defer mt.Close()

	mt.Run("success", func(mt *mtest.T) {
		collection = mt.Coll
		expectedArticle := model.DBArticle{
			ID:          primitive.NewObjectID(),
			Title:       "test_title",
			Created:     time.Date(2022, 1, 1, 1, 1, 1, 0, time.UTC),
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

		articleResponse, err := mongoDb.ReadArticles(7)

		assert.Nil(t, err)
		assert.Equal(t, expectedArticle, articleResponse[0])
	})

	mt.Run("invalid number of days", func(mt *mtest.T) {
		collection = mt.Coll
		mt.AddMockResponses(mtest.CreateWriteErrorsResponse(mtest.WriteError{
			Index:   1,
			Code:    11000,
			Message: "invalid number of days -1",
		}))

		articleResponse, err := mongoDb.ReadArticles(-1)

		assert.Nil(t, articleResponse)
		assert.Equal(t, err, fmt.Errorf("invalid number of days -1"))
	})

}

func TestName(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	c := config.NewConfig()
	err := c.SetDataFromFile("../../../tests/configs/mongo/config.yaml")
	if err != nil {
		return
	}
	c.Read()
	l := log.NewLog(logLevel.Errors)
	mongoDb := NewDb(c.Database, l)
	defer mt.Close()

	name := mongoDb.Name()
	assert.Equal(t, "mongo", name)
}

func TestUrl(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	c := config.NewConfig()
	err := c.SetDataFromFile("../../../tests/configs/mongo/config.yaml")
	if err != nil {
		return
	}
	c.Read()
	l := log.NewLog(logLevel.Errors)
	mongoDb := NewDb(c.Database, l)
	defer mt.Close()

	url := mongoDb.Url()
	assert.Equal(t, "mongodb://mongodb:27017", url)
}

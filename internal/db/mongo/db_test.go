package mongo

import (
	"fmt"
	"github.com/indikator/aggregator_lets_go/internal/config"
	"github.com/indikator/aggregator_lets_go/internal/log/logLevel"
	log "github.com/indikator/aggregator_lets_go/internal/log/stub"
	"github.com/indikator/aggregator_lets_go/model"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"golang.org/x/net/context"
	"testing"
	"time"
)

func mockDecodeResults(curr *mongo.Cursor, articles *[]model.DBArticle, ctx context.Context) error {
	fmt.Println("I am fake")
	err := curr.All(ctx, &articles)
	return err
}

func mockCloseCursor(curr *mongo.Cursor, ctx context.Context) error {
	err := fmt.Errorf("error closing cursor")
	return err
}

func mockCreateClientPass(url string, ctx context.Context) (*mongo.Client, error) {
	fmt.Println(url)
	return nil, nil
}

func mockCreateClientFail(url string, ctx context.Context) (*mongo.Client, error) {
	fmt.Println(url)
	return nil, fmt.Errorf("error creating client")
}

func mockPingClientPass(client *mongo.Client, ctx context.Context) error {
	return nil
}
func mockPingClientFail(client *mongo.Client, ctx context.Context) error {
	return fmt.Errorf("error pinging client")
}

func mockCreateDatabasePass(client *mongo.Client, dbName, collName string) {
	fmt.Println("created database")
}

func mockCreateIndexPass(ctx context.Context) error {
	return nil
}

func mockCreateIndexFail(ctx context.Context) error {
	return fmt.Errorf("error creating index")
}

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
	mt := mtest.New(t, mtest.NewOptions().CreateClient(false).ClientType(mtest.Mock))
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
		DecodeResults_ = DecodeResults
		CloseCursor_ = CloseCursor
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

	mt.Run("no articles found", func(mt *mtest.T) {
		collection = mt.Coll
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    11000,
			Message: "mongo: no documents in result",
		}))

		articleResponse, err := mongoDb.ReadArticles(2)

		assert.Nil(t, articleResponse)
		assert.Equal(t, mtest.CommandError{
			Code:    11000,
			Message: "mongo: no documents in result",
		}.Message, err.Error())
	})

	mt.Run("decode problem", func(mt *mtest.T) {
		originDecodeResults := DecodeResults

		defer func() {
			DecodeResults_ = originDecodeResults
		}()

		DecodeResults_ = mockDecodeResults

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

		_, err := mongoDb.ReadArticles(2)
		fmt.Println(err)
		assert.Equal(t, mtest.CommandError{
			Code:    11000,
			Message: "results argument must be a pointer to a slice, but was a pointer to ptr",
		}.Message, err.Error())
	})

	mt.Run("decode problem", func(mt *mtest.T) {
		originCloseCursor := CloseCursor

		defer func() {
			CloseCursor_ = originCloseCursor
		}()

		CloseCursor_ = mockCloseCursor

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

		_, err := mongoDb.ReadArticles(2)
		fmt.Println(err)
		assert.Equal(t, mtest.CommandError{
			Code:    11000,
			Message: "error closing cursor",
		}.Message, err.Error())
	})

}

func TestInitDb(t *testing.T) {
	c := config.NewConfig()
	err := c.SetDataFromFile("../../../tests/configs/mongo/config.yaml")
	if err != nil {
		return
	}
	c.Read()
	l := log.NewLog(logLevel.Errors)
	mongoDb := NewDb(c.Database, l)

	CreateClient_ = mockCreateClientPass
	PingClient_ = mockPingClientPass
	CreateDatabase_ = mockCreateDatabasePass
	CreateIndex_ = mockCreateIndexPass
	mongoDb.InitDb()

	// all good
	originCreateClient := CreateClient
	originPingClient := PingClient
	originCreateDatabase := CreateDatabase
	originCreateIndex := CreateIndex

	defer func() {
		CreateClient_ = originCreateClient
		PingClient_ = originPingClient
		CreateDatabase_ = originCreateDatabase
		CreateIndex_ = originCreateIndex
	}()

	CreateClient_ = mockCreateClientPass
	PingClient_ = mockPingClientPass
	CreateDatabase_ = mockCreateDatabasePass
	CreateIndex_ = mockCreateIndexPass

	mongoDb.InitDb()

	// client creation failed
	originCreateClient = CreateClient
	defer func() {
		CreateClient_ = originCreateClient
	}()
	CreateClient_ = mockCreateClientFail

	err = mongoDb.InitDb()
	assert.Equal(t, "error creating client", err.Error())

	// client ping failed
	originPingClient = PingClient
	CreateClient_ = mockCreateClientPass
	defer func() {
		PingClient_ = originPingClient
	}()
	PingClient_ = mockPingClientFail

	err = mongoDb.InitDb()
	assert.Equal(t, "error pinging client", err.Error())

	// index creation failed
	originCreateIndex = CreateIndex
	PingClient_ = mockPingClientPass
	defer func() {
		CreateIndex_ = originCreateIndex
	}()
	CreateIndex_ = mockCreateIndexFail

	err = mongoDb.InitDb()
	assert.Equal(t, "error creating index", err.Error())
}

func TestName(t *testing.T) {
	c := config.NewConfig()
	err := c.SetDataFromFile("../../../tests/configs/mongo/config.yaml")
	if err != nil {
		return
	}
	c.Read()
	l := log.NewLog(logLevel.Errors)
	mongoDb := NewDb(c.Database, l)

	name := mongoDb.Name()
	assert.Equal(t, "stub", name)
}

func TestUrl(t *testing.T) {
	c := config.NewConfig()
	err := c.SetDataFromFile("../../../tests/configs/mongo/config.yaml")
	if err != nil {
		return
	}
	c.Read()
	l := log.NewLog(logLevel.Errors)
	mongoDb := NewDb(c.Database, l)

	url := mongoDb.Url()
	assert.Equal(t, "mongodb://mongodb:27017", url)
}

func TestCreateClient(t *testing.T) {
	CreateClient("testUrl", context.TODO())
}

func TestCreateIndex(t *testing.T) {
	CreateIndex(context.TODO())
}

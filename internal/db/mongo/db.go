package mongo

import (
	"context"
	"fmt"
	"github.com/indikator/aggregator_lets_go/internal/log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/indikator/aggregator_lets_go/internal/db"
	"github.com/indikator/aggregator_lets_go/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// instantiate new collection
var collection *mongo.Collection

// init functions
var (
	DecodeResults_  func(curr *mongo.Cursor, articles *[]model.DBArticle, ctx context.Context) error
	CloseCursor_    func(curr *mongo.Cursor, ctx context.Context) error
	CreateClient_   func(url string, ctx context.Context) (*mongo.Client, error)
	PingClient_     func(client *mongo.Client, ctx context.Context) error
	CreateDatabase_ func(client *mongo.Client, dbName, collName string)
	CreateIndex_    func(ctx context.Context) error
)

type database struct {
	name string
	url  string
	log  log.Log
}

// NewDb create an instance of database
func NewDb(c db.Config, l log.Log) db.Db {
	return &database{
		name: c.Name,
		url:  c.Url,
		log:  l,
	}
}

func (db *database) Name() string {
	return db.name
}

func (db *database) Url() string {
	return db.url
}

func (db *database) WriteArticle(article *model.DBArticle) (*model.DBArticle, error) {
	_, err := collection.InsertOne(context.Background(), article)
	if err != nil {
		return nil, err
	}
	return article, nil
}

func (db *database) ReadArticles(nDays int) ([]model.DBArticle, error) {
	//passing bson.D{{}} matches all documents in the collection
	if nDays < 1 {
		return nil, fmt.Errorf("invalid number of days %d", nDays)
	}

	filter := bson.M{"created": bson.M{
		"$gte": primitive.NewDateTimeFromTime(time.Now().AddDate(0, 0, -1*nDays)), //last 7 days
	}}

	ctx := context.Background()
	curr, err := Find(filter, collection, ctx)
	if err != nil {
		return nil, err
	}

	var articles []model.DBArticle

	err = DecodeResults_(curr, &articles, ctx)
	if err != nil {
		CloseCursor_(curr, ctx)
		return nil, err
	}

	err = CloseCursor_(curr, ctx)
	if err != nil {
		return nil, err
	}

	return articles, nil
}

//InitDb creates a new MongoDB client and connect to your running MongoDB server
func (db *database) InitDb() error {

	ctx := context.Background()
	client, err := CreateClient_(db.url, ctx)
	if err != nil {
		return err
	}

	// Next, let’s ensure that your MongoDB server was found and connected to successfully using the Ping method.
	err = PingClient_(client, ctx)
	if err != nil {
		return err
	}

	db.log.WriteInfo("Connected to mongo")

	// create a database
	CreateDatabase_(client, "news", "articles")

	// Declare model for the indexes
	err = CreateIndex_(ctx)
	if err != nil {
		return err
	}

	return nil
}

func Find(filter bson.M, collection *mongo.Collection, ctx context.Context) (*mongo.Cursor, error) {
	curr, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	return curr, nil
}

func DecodeResults(curr *mongo.Cursor, articles *[]model.DBArticle, ctx context.Context) error {
	fmt.Println("I am true")
	err := curr.All(ctx, articles)
	return err
}

func CloseCursor(curr *mongo.Cursor, ctx context.Context) error {
	err := curr.Close(ctx)
	return err
}

func CreateClient(url string, ctx context.Context) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(url)
	client, err := mongo.Connect(ctx, clientOptions)
	return client, err
}

func PingClient(client *mongo.Client, ctx context.Context) error {
	err := client.Ping(ctx, nil)
	return err
}

func CreateDatabase(client *mongo.Client, dbName, collName string) {
	collection = client.Database(dbName).Collection(collName)
}

func CreateIndex(ctx context.Context) error {
	_, err := collection.Indexes().CreateOne(
		ctx,
		mongo.IndexModel{
			Keys: bson.D{
				{Key: "author", Value: 1},
				{Key: "title", Value: 1},
				{Key: "created", Value: 1},
				{Key: "url", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
	)
	return err
}

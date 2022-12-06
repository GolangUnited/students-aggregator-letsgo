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
	articles := []model.DBArticle{}
	cur, err := collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	if err = cur.All(context.Background(), &articles); err != nil {
		return nil, err
	}

	if err = cur.Err(); err != nil {
		return nil, err
	}

	//once exhausted, close the cursor
	err = cur.Close(context.Background())
	if err != nil {
		return nil, err
	}
	
	return articles, nil
}

// InitDb creates a new MongoDB client and connect to your running MongoDB server
func (db *database) InitDb() error {
	clientOptions := options.Client().ApplyURI(db.url)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return err
	}

	// Next, let’s ensure that your MongoDB server was found and connected to successfully using the Ping method.
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return err
	}

	db.log.WriteInfo("Connected to mongo")

	// create a database
	collection = client.Database("news").Collection("articles")

	// Declare model for the indexes
	_, err = collection.Indexes().CreateOne(
		context.Background(),
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

	return nil
}

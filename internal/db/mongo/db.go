package mongo

import (
	"context"
	"fmt"
	"github.com/indikator/aggregator_lets_go/internal/db"
	"github.com/indikator/aggregator_lets_go/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// instantiate new collection
var collection *mongo.Collection

type database struct {
	url string
}

// NewDb create an instance of database
func NewDb(c db.Config) db.Db {
	URL := c.Url
	return &database{
		url: URL,
	}
}

func (db *database) WriteArticle(article *model.DBArticle) (*model.DBArticle, error) {

	_, err := collection.InsertOne(context.Background(), article)
	if err != nil {
		article = nil
	}
	return article, err

}

func (db *database) ReadAllArticles() ([]model.DBArticle, error) {
	//passing bson.D{{}} matches all documents in the collection
	filter := bson.D{{}}
	var articles []model.DBArticle
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

	if len(articles) == 0 {
		return nil, mongo.ErrNoDocuments
	}
	return articles, nil
}

// DBInit creates a new MongoDB client and connect to your running MongoDB server
func (db *database) DBInit() error {
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
	fmt.Println("Connected to mongo")

	// create a database
	collection = client.Database("news").Collection("articles")

	// Declare model for the indexes
	indexName, err := collection.Indexes().CreateOne(
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
	fmt.Println(indexName)

	return nil
}

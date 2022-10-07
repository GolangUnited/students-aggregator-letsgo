package godev

import (
	"context"
	"fmt"
	"github.com/indikator/aggregator_lets_go/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

// instantiate new collection and context
var collection *mongo.Collection
var ctx = context.TODO()

// DBinit creates a new MongoDB client and connect to your running MongoDB server
func DBinit(uri string) {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		log.Fatal(err)
	}
	// Next, let’s ensure that your MongoDB server was found and connected to successfully using the Ping method.

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to mongo")

	// create a database
	collection = client.Database("news").Collection("articles")
}

func WriteArticle(article *model.Article) error {
	_, err := collection.InsertOne(ctx, article)
	return err
}

func ReadAll() ([]*model.Article, error) {
	//passing bson.D{{}} matches all documents in the collection
	filter := bson.D{{}}
	return filterArticles(filter)
}

func filterArticles(filter interface{}) ([]*model.Article, error) {
	// a slice of articles for storing the decoded documents
	var articles []*model.Article
	cur, err := collection.Find(ctx, filter)
	if err != nil {
		return articles, err
	}

	for cur.Next(ctx) {
		var a model.Article
		err := cur.Decode(&a)
		if err != nil {
			return articles, err
		}

		articles = append(articles, &a)
	}

	if err := cur.Err(); err != nil {
		return articles, err
	}

	//once exhausted, close the cursor
	cur.Close(ctx)

	if len(articles) == 0 {
		return articles, mongo.ErrNoDocuments
	}

	return articles, nil
}

package mongo

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
var ctx = context.TODO()

// DBinit creates a new MongoDB client and connect to your running MongoDB server
func DBinit(uri string) *mongo.Collection {
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
	collection := client.Database("news").Collection("articles")
	return collection
}

func WriteArticle(article *model.DBArticle, collection *mongo.Collection) error {
	_, err := collection.InsertOne(ctx, article)
	return err
}

func ReadAllArticles(collection *mongo.Collection) ([]*model.DBArticle, error) {
	//passing bson.D{{}} matches all documents in the collection
	filter := bson.D{{}}
	return filterArticles(filter, collection)
}

func filterArticles(filter interface{}, collection *mongo.Collection) ([]*model.DBArticle, error) {
	// a slice of articles for storing the decoded documents
	var articles []*model.DBArticle
	cur, err := collection.Find(ctx, filter)
	if err != nil {
		return articles, err
	}

	for cur.Next(ctx) {
		var a model.DBArticle
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

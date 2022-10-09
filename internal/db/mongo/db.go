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
//var Ctx = context.TODO()
var collection *mongo.Collection

// DBinit creates a new MongoDB client and connect to your running MongoDB server
func DBinit(uri string) {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}
	// Next, let’s ensure that your MongoDB server was found and connected to successfully using the Ping method.

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to mongo")

	// create a database
	collection = client.Database("news").Collection("articles")
}

func WriteArticle(article *model.DBArticle) (*model.DBArticle, error) {
	_, err := collection.InsertOne(context.Background(), article)
	if err != nil {
		article = nil
	}
	return article, err
}

func ReadAllArticles() ([]model.DBArticle, error) {
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

	if err := cur.Err(); err != nil {
		return nil, err
	}

	//once exhausted, close the cursor
	cur.Close(context.Background())

	if len(articles) == 0 {
		return nil, mongo.ErrNoDocuments
	}
	return articles, nil
}

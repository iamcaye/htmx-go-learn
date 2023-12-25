package repos

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
	"github.com/iamcaye/htmx-go-learn/cmd/models"
)
// Replace the placeholder with your Atlas connection string
const uri = "mongodb://root:example@localhost:27017/?maxPoolSize=10&wmajority=1&wtimeout=30"

func InitMongo() (mongo.Client, error) {
	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to MongoDB!")
	return *client, err
}

func getCollection(client mongo.Client, collectionName string) *mongo.Collection {
	return client.Database("blog").Collection(collectionName)
}

func GetPosts(client mongo.Client) (models.Posts, error) {
	collection := getCollection(client, "posts")
	res, err := collection.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}

	var posts []models.Post
	if res.All(context.Background(), &posts); err != nil {
		return nil, err
	}
	return posts, nil
}

func AddPost(client mongo.Client, post models.Post) error {
	collection := getCollection(client, "posts")
	_, err := collection.InsertOne(context.Background(), post)
	return err
}

package database

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

func DBSet() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb+srv://hungnqptit:Ngohung98vn@flutterlearning.5nu9hmb.mongodb.net/?retryWrites=true&w=majority"))

	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	client.Ping(context.TODO(), nil)
	if err != nil {
		log.Println("Failed to connect db")
		return nil
	}

	fmt.Println("Success connect db")
	return client
}

var Client *mongo.Client = DBSet()

func UserData(client *mongo.Client, collectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database("Ecommerce").Collection(collectionName)
	return collection
}

func ProductData(client *mongo.Client, collectionName string) *mongo.Collection {
	var productionCollection *mongo.Collection = client.Database("Ecommerce").Collection(collectionName)
	return productionCollection
}

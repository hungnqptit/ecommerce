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
	client, err := mongo.NewClient(options.Client().ApplyURI(
		"mongodb://localhost:27017",
		//"mongodb+srv://hungnqptit:dnBSb1VDOBKc8zjT@clusterflutter.1xkdsh1.mongodb.net/"
	))

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

func SavingInfoData(client *mongo.Client) *mongo.Collection {
	var collection *mongo.Collection = client.Database("Ecommerce").Collection("SavingInfo")
	return collection
}

func ProductData(client *mongo.Client, collectionName string) *mongo.Collection {
	var productionCollection *mongo.Collection = client.Database("Ecommerce").Collection(collectionName)
	return productionCollection
}

func ProductDataTest(client *mongo.Client, collectionName string) *mongo.Collection {
	var productionCollection *mongo.Collection = client.Database("Ecommerce").Collection(collectionName)
	return productionCollection
}

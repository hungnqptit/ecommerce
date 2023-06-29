package database

import (
	"context"
	"ecommerce/models"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

var ProductCollection = UserData(Client, "Products")

var (
	ErrorCantFindProduct    = errors.New("can't find the product")
	ErrorCantDecodeProduct  = errors.New("can't find the product")
	ErrorUserIdIsNotValid   = errors.New("this user is not valid")
	ErrorCantUpdateUser     = errors.New("cannot add this product to the cart")
	ErrorCantRemoveItemCart = errors.New("cannot remove this item from the cart")
	ErrorCantGetItem        = errors.New("was unable to get this item from the cart")
	ErrorCantBuyCartItem    = errors.New("cannot update the purchase")
)

func AddProductToCart(ctx context.Context, prodCollection *mongo.Collection, userCollection *mongo.Collection, prodID primitive.ObjectID, userID string) error {
	searchFromDB, err := ProductCollection.Find(ctx, bson.M{"_id": prodID})
	if err != nil {
		return ErrorCantFindProduct
	}
	var productCart []models.ProductUser
	err = searchFromDB.All(ctx, &productCart)
	if err != nil {
		log.Println(err)
		return ErrorCantDecodeProduct
	}
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrorUserIdIsNotValid
	}

	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{
		Key: "userCart",
		Value: bson.D{
			{
				Key:   "$each",
				Value: productCart,
			},
		},
	}}}}
	_, err = userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return ErrorCantUpdateUser
	}

	return nil
}

func RemoveCartItem(ctx context.Context, prodCollection *mongo.Collection, userCollection *mongo.Collection, prodID primitive.ObjectID, userID string) error {
	searchFromDB, err := ProductCollection.Find(ctx, bson.M{"_id": prodID})
	if err != nil {
		return ErrorCantFindProduct
	}
	var productCart []models.ProductUser
	err = searchFromDB.All(ctx, &productCart)
	if err != nil {
		log.Println(err)
		return ErrorCantDecodeProduct
	}
	id, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		log.Println(err)
		return ErrorUserIdIsNotValid
	}

	filter := bson.D{primitive.E{Key: "_id", Value: id}}
	update := bson.D{{Key: "pull", Value: bson.M{"userCart": bson.M{"_id": prodID}}}}
	_, err = userCollection.UpdateMany(ctx, filter, update)
	if err != nil {
		return ErrorCantUpdateUser
	}

	return nil
}

func InstantBuyer(context.Context, *mongo.Collection, *mongo.Collection, primitive.ObjectID, string) error {
	return nil
}
func BuyItemFromCart(context.Context, *mongo.Collection, string) error {
	return nil
}

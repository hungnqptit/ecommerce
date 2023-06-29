package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID            primitive.ObjectID `json:"_id" bson:"_id"`
	FirstName     *string            `json:"firstName" validate:"required,min=2,max=30"`
	LastName      *string            `json:"lastName" validate:"required,min=2,max=30"`
	Password      *string            `json:"password" validate:"required,min=6"`
	Email         *string            `json:"email" validate:"email,required"`
	Phone         *string            `json:"phone" validate:"required"`
	Token         *string            `json:"token" `
	RefreshToken  *string            `json:"refreshToken"`
	CreateAt      time.Time          `json:"createAt"`
	UpdatedAt     time.Time          `json:"updatedAt"`
	UserId        string             `json:"userId"`
	UserCart      []ProductUser      `json:"userCart" bson:"userCart"`
	AddressDetail []Address          `json:"address" bson:"address"`
	OrderStatus   []Order            `json:"orders" bson:"orders"`
}

type Product struct {
	ProductID   primitive.ObjectID `json:"_id" bson:"_id"`
	ProductName *string            `json:"productName" `
	Price       *uint64            `json:"price" `
	Rating      *uint8             `json:"rating" `
	Image       *string            `json:"image" `
}

type ProductUser struct {
	productID   primitive.ObjectID `json:"_id" bson:"_id"`
	productName *string            `json:"productName" `
	price       int                `json:"price" `
	rating      *uint              `json:"rating" `
	image       *string            `json:"image" `
}

type Address struct {
	addressID primitive.ObjectID `bson:"_id"`
	house     *string            `json:"houseName" bson:"houseName"`
	street    *string            `json:"streetName" bson:"streetName"`
	city      *string            `json:"cityName" bson:"cityName"`
	pinCode   *string            `json:"pinCode" bson:"pinCode"`
}

type Order struct {
	orderID       primitive.ObjectID `bson:"_id"`
	orderCart     []ProductUser      `json:"orderList" bson:"orderList"`
	OrderedAt     time.Time          `json:"orderedAt" bson:"orderedAt"`
	price         int                `json:"totalPrice" bson:"totalPrice"`
	discount      *int               `json:"discount" bson:"discount"`
	paymentMethod Payment            `json:"paymentMethod" bson:"paymentMethod"`
}

type Payment struct {
	digital bool
	cod     bool
}

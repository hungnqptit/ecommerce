package tokens

import (
	"context"
	"ecommerce/database"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type SignedDetails struct {
	Email     string
	FirstName string
	LastName  string
	Uid       string
	jwt.StandardClaims
}

var UserData = database.UserData(database.Client, "Users")

var SECRET_KEY = "test"

//os.Getenv("SECRET_KEY")

func TokenGenerator(email string, firstName string, lastName string, uid string) (signedToken string, signedRefreshtoken string, err error) {
	claims := &SignedDetails{Email: email, FirstName: firstName, LastName: lastName, Uid: uid, StandardClaims: jwt.StandardClaims{
		ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
	}}

	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))

	if err != nil {
		return "", "", err
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS384, refreshClaims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Panicln(err.Error())
		return
	}
	return token, refreshToken, err
}

func ValidateToken(signedToken string) (claims *SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(signedToken, &SignedDetails{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})
	if err != nil {
		msg = err.Error()
		return
	}
	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = "Token invalid"
		return
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = "Token is already expired"
		return
	}
	return claims, msg
}

func UpdateAllToken(signedToken string, signedRefreshToken string, userId string) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	var updateObj primitive.D

	updateObj = append(updateObj, bson.E{Key: "token", Value: signedToken})
	updateObj = append(updateObj, bson.E{Key: "refreshToken", Value: signedRefreshToken})
	updateAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{Key: "updatedAt", Value: updateAt})

	upsert := true
	filter := bson.M{"userid": userId}
	opt := options.UpdateOptions{Upsert: &upsert}
	_, err := UserData.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: updateObj}}, &opt)
	defer cancel()
	if err != nil {
		log.Panicln(err)
	}

}

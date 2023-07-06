package controller

import (
	"context"
	"ecommerce/database"
	"ecommerce/models"
	"ecommerce/tokens"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
)

type Application struct {
	prodCollection *mongo.Collection
	userCollection *mongo.Collection
}

func NewApplication(prodCollection, userCollection *mongo.Collection) *Application {
	return &Application{
		prodCollection: prodCollection,
		userCollection: userCollection,
	}
}

var UserCollection = database.UserData(database.Client, "Users")

var ProductCollection = database.UserData(database.Client, "Products")

var SavingInfoCollection = database.SavingInfoData(database.Client)

var Validate = validator.New()

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panicln(err)
	}
	return string(bytes)
}

func VerifyPassword(userPassword string, givenPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(givenPassword), []byte(userPassword))
	valid := true
	msg := ""
	if err != nil {
		msg = "Login or Password is incorrect"
		valid = false
	}
	return valid, msg
}

func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 1000*time.Second)
		defer cancel()

		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationError := Validate.Struct(user)
		if validationError != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationError.Error()})
			return
		}

		count, err := UserCollection.CountDocuments(ctx, bson.M{"phone": user.Email})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			log.Panicln(err)
			return
		}
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
		}

		count, err = UserCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			log.Panicln(err)
			return
		}
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Phone already in use"})
			return
		}
		password := HashPassword(*user.Password)
		user.Password = &password
		user.CreateAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.UserId = user.ID.Hex()
		token, refreshToken, _ := tokens.TokenGenerator(*user.Email, *user.FirstName, *user.LastName, user.UserId)
		user.Token = &token
		user.RefreshToken = &refreshToken
		user.UserCart = make([]models.ProductUser, 0)
		user.AddressDetail = make([]models.Address, 0)
		user.OrderStatus = make([]models.Order, 0)
		_, inserter := UserCollection.InsertOne(ctx, user)
		if inserter != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "The user did not created"})
			return
		}
		c.JSON(http.StatusCreated, "Successfully signed in!")
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 1000*time.Second)
		defer cancel()

		var user models.User
		var foundUser models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := UserCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Login or Password is incorrect"})
			return
		}

		PasswordValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
		defer cancel()

		if !PasswordValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			fmt.Println(msg)
			return
		}

		token, refreshToken, err := tokens.TokenGenerator(*foundUser.Email, *foundUser.FirstName, *foundUser.LastName, foundUser.UserId)
		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			fmt.Println(err.Error())
			return
		}

		tokens.UpdateAllToken(token, refreshToken, foundUser.UserId)
		foundUser.Token = &token
		foundUser.RefreshToken = &refreshToken
		c.JSON(http.StatusFound, foundUser)
	}
}

func SearchProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var productList []models.Product
		var ctx, cancel = context.WithTimeout(context.Background(), 1000*time.Second)
		defer cancel()

		cursor, err := ProductCollection.Find(ctx, bson.D{})
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "Something went wrong, please try again sometime")
			return
		}

		err = cursor.All(ctx, &productList)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		defer cursor.Close(ctx)

		if err := cursor.Err(); err != nil {
			log.Println(err)
			c.IndentedJSON(400, "invalid")
			return
		}
		defer cancel()
		c.IndentedJSON(200, productList)
	}
}

func SearchProductByQuery() gin.HandlerFunc {
	return func(c *gin.Context) {
		var searchProducts []models.Product
		queryParam := c.Query("name")

		if queryParam == "" {
			log.Println("query is empty")
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid search index"})
			c.Abort()
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 1000*time.Second)
		defer cancel()

		searchQueryDb, err := ProductCollection.Find(ctx, bson.M{"productionName": bson.M{"$regex": queryParam}})
		if err != nil {
			c.IndentedJSON(404, "Something went wrong while fetching data")
			return
		}

		err = searchQueryDb.All(ctx, &searchProducts)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(400, "invalid")
			return
		}

		defer searchQueryDb.Close(ctx)

		if err := searchQueryDb.Err(); err != nil {
			log.Println(err)
			c.IndentedJSON(400, "Invalid request")
			return
		}

		defer cancel()
		c.IndentedJSON(200, searchProducts)
	}
}

type AddSavingRequest struct {
	UserID       string  `json:"userID"`
	Amount       float64 `json:"amount"`
	Category     int64   `json:"category"`
	TermDuration int64   `json:"termDuration"`
	SavingName   string  `json:"savingName"`
}

func AddSavingItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		request := &AddSavingRequest{}
		if err := c.BindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if request.UserID == "" {
			fmt.Println("UserID must not be null")
			c.IndentedJSON(http.StatusBadRequest, "UserID must not be null")
			return
		}

		if request.SavingName == "" {
			fmt.Println("Saving name must not be null")
			c.IndentedJSON(http.StatusBadRequest, "Saving name must not be null")
			return
		}

		//if err != nil {
		//	fmt.Println("Input amount must be number")
		//	c.IndentedJSON(http.StatusBadRequest, "Input amount must be number")
		//	return
		//}
		//request.Amount = inputSavingMoney

		if request.Category >= 3 || request.Category < 0 {
			fmt.Println("Category invalid")
			c.IndentedJSON(http.StatusBadRequest, "Category invalid")
			return
		}

		if request.TermDuration < 0 {
			fmt.Println("Term invalid")
			c.IndentedJSON(http.StatusBadRequest, "Category invalid")
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var foundUser models.User
		objectId, err := primitive.ObjectIDFromHex(request.UserID)
		if err != nil {
			log.Println("Invalid id")
			c.IndentedJSON(http.StatusBadRequest, "Invalid user ID")
		}
		errFind := UserCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: objectId}}).Decode(&foundUser)
		if errFind != nil {
			fmt.Println("UserID not true")
			c.IndentedJSON(http.StatusBadRequest, "UserID not true")
			return
		}
		fmt.Println(foundUser)
		if *foundUser.Money < request.Amount {
			fmt.Println("Account money insufficient")
			c.IndentedJSON(http.StatusInternalServerError, "Account money insufficient")
			return
		}
		var foundedRates []models.SearchedInfo
		cursor, errFindRate := SavingInfoCollection.Aggregate(ctx,
			mongo.Pipeline{
				bson.D{{"$unwind", bson.D{{"path", "$rates"}}}},
				bson.D{{"$match", bson.D{{"rates.termDuration", request.TermDuration}}}},
			})
		err = cursor.All(ctx, &foundedRates)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		defer cursor.Close(ctx)
		//errFindCursor.Decode(&foundedRates)
		if errFindRate != nil {
			fmt.Println("Term duration not valid")
			c.IndentedJSON(http.StatusInternalServerError, "Term duration not valid")
			return
		}

		newSaving := models.Saving{
			SavingID:        primitive.NewObjectID(),
			SavingAmount:    &request.Amount,
			SavingTermByDay: &request.TermDuration,
			SavingRate:      foundedRates[0].Rates.Rate,
			SavingName:      &request.SavingName,
			CategoryType:    &request.Category,
		}
		errValidate := Validate.Struct(newSaving)
		if errValidate != nil {
			fmt.Println("Saving object invalid")
			c.IndentedJSON(http.StatusInternalServerError, "Data invalid")
			return
		}
		foundUser.SavingList = append(foundUser.SavingList, newSaving)
		remainingMoney := *foundUser.Money - request.Amount
		foundUser.Money = &remainingMoney
		filter := bson.M{"userid": request.UserID}

		upsert := true
		opt := options.UpdateOptions{Upsert: &upsert}
		_, errUpdate := database.UserData(database.Client, "Users").UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: foundUser}}, &opt)
		defer cancel()
		if errUpdate != nil {
			log.Panicln(err)
		}
		c.IndentedJSON(200, "Saving create successfully")
	}
}

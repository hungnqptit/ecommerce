package main

import (
	controller "ecommerce/controllers"
	"ecommerce/database"
	"ecommerce/routes"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	app := controller.NewApplication(database.ProductData(database.Client, "Products"), database.UserData(database.Client, "Users"))
	router := gin.New()
	router.Use(gin.Logger())
	database.DBSet()
	//db mysql
	db := database.ConnectDb()
	defer db.Close()

	routes.UserRoutes(router)
	//router.Use(middleware.Authentication())
	router.GET("/add_to_cart", app.AddToCart())
	router.GET("/remove_item", app.RemoveItemFromCart())
	router.POST("/saving/add_saving", controller.AddSavingItem())
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Println("Running on port ", port)
	if err := router.Run(":" + port); err != nil {
		log.Panicf("error: %s", err)
	}
}

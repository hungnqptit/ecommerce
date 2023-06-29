package main

import (
	controller "ecommerce/controllers"
	"ecommerce/database"
	"ecommerce/middleware"
	"ecommerce/routes"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	app := controller.NewApplication(database.ProductData(database.Client, "Products"), database.UserData(database.Client, "Users"))

	router := gin.New()
	router.Use(gin.Logger())
	database.DBSet()
	routes.UserRoutes(router)
	router.Use(middleware.Authentication())
	router.GET("/add_to_cart", app.AddToCart())
	router.GET("/remove_item", app.RemoveItemFromCart())
	log.Fatal(router.Run(":" + port))
}

package main

import (
	controller "ecommerce/controllers"
	"ecommerce/database"
	"ecommerce/middleware"
	"ecommerce/routes"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	fmt.Println("Running on port ", port)
	app := controller.NewApplication(database.ProductData(database.Client, "Products"), database.UserData(database.Client, "Users"))
	router := gin.New()
	router.Use(gin.Logger())
	database.DBSet()
	routes.UserRoutes(router)
	router.Use(middleware.Authentication())
	router.GET("/add_to_cart", app.AddToCart())
	router.GET("/remove_item", app.RemoveItemFromCart())
	log.Fatal(router.Run(":" + port))
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

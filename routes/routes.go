package routes

import (
	"ecommerce/controllers"
	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("users/signup", controller.SignUp())
	incomingRoutes.POST("users/login", controller.Login())
	//incomingRoutes.POST("/saving/add", controller.AddSavingItem())
	//incomingRoutes.POST("admin/add_product", controllers.ProductViewerAdmin())
	//incomingRoutes.POST("users/product_view", controllers.SearchProduct())
	//incomingRoutes.POST("users/search", controllers.Search())
}

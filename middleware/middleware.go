package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		ClientToken := c.Request.Header.Get("Authorization")
		fmt.Println(ClientToken)
		//if ClientToken == "" {
		//	c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Method not allowed"})
		//	c.Abort()
		//	return
		//}
		//claims, err := tokens.ValidateToken(strings.Replace("Bearer ", ClientToken, "", 1))
		//if err != "" {
		//	c.JSON(http.StatusMethodNotAllowed, gin.H{"error": err})
		//	c.Abort()
		//	return
		//}
		//
		//c.Set("email", claims.Email)
		//c.Set("uid", claims.Uid)
		c.Next()
	}
}

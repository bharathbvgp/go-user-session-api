package middleware

import (
	"net/http"
	"time"
	"userapp/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(c *gin.Context) {
	cookie , err := c.Cookie("token");
	
	if err != nil {
		// c.JSON(http.StatusUnauthorized, gin.H{"error" : "Authorization token is required"})
		c.JSON(http.StatusOK , gin.H{"IsSessionValid" : "false"})
		c.Abort()
		return 
	}
	
	claims , err := utils.ValidateToken(cookie)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	// checking if the token is expired or not 
	if claims.ExpiresAt < time.Now().Unix() {
		c.JSON(http.StatusUnauthorized , gin.H{"IsSessionValid" : "false"})
		c.Abort()
		return 
	}

	c.Set("userID" , claims.UserID)
	c.Next()
}
package routes

import (
	"userapp/controllers"
	"userapp/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	router.POST("/signup" , controllers.RegisterUser)
	router.POST("/login" , controllers.LoginUser)
	router.GET("/checksession" , middleware.AuthMiddleware , controllers.CheckSession)
	router.POST("/logout" , middleware.AuthMiddleware , controllers.Logout)
}
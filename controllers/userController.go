package controllers

import (
	"net/http"
	"time"
	"userapp/database"
	"userapp/models"
	"userapp/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(c *gin.Context) {
	var newUser models.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest , gin.H{"error" : err.Error()})
		return 
	}
	hashedPassword , err := bcrypt.GenerateFromPassword([]byte(newUser.Password) , bcrypt.DefaultCost) 
	if err != nil {
		c.JSON(http.StatusInternalServerError , gin.H{ "error" : "Failed to hash password"})
		return
	}
	newUser.Password = string(hashedPassword)
	if result := database.DB.Create(&newUser); result.Error != nil {
		c.JSON(http.StatusInternalServerError , gin.H{"error" : result.Error.Error()})
		return 
	}
	c.JSON(http.StatusCreated , newUser)
}

func LoginUser(c *gin.Context) {
	var loginDetails struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&loginDetails); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var user models.User
	// checking whether user exists in database or not 
	if result := database.DB.Where("email = ?" , loginDetails.Email).First(&user); result.Error != nil {
		c.JSON(http.StatusOK, gin.H{"error": "Invalid email or password"})
		return
	}
	// validating password

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginDetails.Password)); err != nil {
		c.JSON(http.StatusUnauthorized,gin.H{"error": "Invalid email or password"})
		return
	}
	// get the token 
	token , err := utils.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	// setting the cookie , so it will be sent from client everytime it sends out a request to the server 

	http.SetCookie(c.Writer, &http.Cookie{
		Name: "token",
		Value : token,
		Expires: time.Now().Add(12 * time.Hour),
		HttpOnly: true,
	})

	c.JSON(http.StatusOK, gin.H{"message" : "Login successfu;" })

}

// If this controller got hit then the session is valid 
// as it passes through the middlware which checks for the expiry of session

func CheckSession(c *gin.Context) {
	c.JSON(http.StatusOK , gin.H{"IsSessionValid" : "true"})
}

func Logout(c *gin.Context) {
	// If we remove the cookie from the headers then automatically the session will be invalidated right 
	http.SetCookie(c.Writer , &http.Cookie{
		Name: "token",
		Value: "",
		Expires: time.Unix(0,0),
		HttpOnly: true,
	})
	c.JSON(http.StatusOK , gin.H{"message" : "Logged out successfully"})
}


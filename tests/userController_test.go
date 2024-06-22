package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"userapp/database"
	"userapp/models"
	"userapp/routes"
	"userapp/utils"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func setupRouter() *gin.Engine {
	router := gin.Default()
	routes.SetupRoutes(router)
	return router
}

func hashPassword(password string) string {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword)
}

func clearDatabase() {
	database.DB.Exec("DELETE FROM users")
}

func TestRegisterUser(t *testing.T) {
	database.SetupDatabase()
	clearDatabase()
	router := setupRouter()

	user := models.User{
		Name:     "Test User",
		Email:    "unique1@example.com", // Use a unique email for this test
		Password: "password",
	}
	userJSON, _ := json.Marshal(user)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(userJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var responseUser models.User
	json.Unmarshal(w.Body.Bytes(), &responseUser)
	assert.Equal(t, "Test User", responseUser.Name)
}

func TestLoginUser(t *testing.T) {
	database.SetupDatabase()
	clearDatabase()
	router := setupRouter()

	user := models.User{
		Name:     "Test User",
		Email:    "unique2@example.com", // Use a unique email for this test
		Password: hashPassword("password"),
	}
	database.DB.Create(&user)

	loginDetails := map[string]string{
		"email":    "unique2@example.com", // Match the unique email
		"password": "password",
	}
	loginJSON, _ := json.Marshal(loginDetails)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(loginJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Check if the token is set in the cookie
	cookies := w.Result().Cookies()
	assert.NotEmpty(t, cookies, "Expected at least one cookie")
	tokenCookie := cookies[0]
	assert.Equal(t, "token", tokenCookie.Name)
	assert.NotEmpty(t, tokenCookie.Value)
}

func TestCheckSession(t *testing.T) {
	database.SetupDatabase()
	clearDatabase()
	router := setupRouter()

	// Simulating a login to get a valid token
	user := models.User{
		Name:     "Test User",
		Email:    "unique3@example.com", // Use a unique email for this test
		Password: hashPassword("password"),
	}
	database.DB.Create(&user)
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/checksession", nil)
	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: token,
	})
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestLogout(t *testing.T) {
	database.SetupDatabase()
	clearDatabase()
	router := setupRouter()

	// Simulating a login to get a valid token
	user := models.User{
		Name:     "Test User",
		Email:    "unique4@example.com", // Use a unique email for this test
		Password: hashPassword("password"),
	}
	database.DB.Create(&user)
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/logout", nil)
	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: token,
	})
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Check if the token is removed from the cookie
	cookies := w.Result().Cookies()
	assert.NotEmpty(t, cookies, "Expected at least one cookie")
	tokenCookie := cookies[0]
	assert.Equal(t, "token", tokenCookie.Name)
	assert.Empty(t, tokenCookie.Value)
	assert.True(t, tokenCookie.Expires.Before(time.Now()))
}

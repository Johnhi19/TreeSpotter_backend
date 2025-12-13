package handlers

import (
	"database/sql"
	"net/http"
	"time"

	database "github.com/Johnhi19/TreeSpotter_backend/db"
	"github.com/Johnhi19/TreeSpotter_backend/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("SUPER_SECRET_KEY")

// ----------------------
// Register
// ----------------------
func Register(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Check if username exists
	var exists string
	err := database.DB.QueryRow("SELECT Username FROM User WHERE Username = ?", user.Username).Scan(&exists)
	if err != sql.ErrNoRows && err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if exists != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already taken"})
		return
	}

	// Hash password
	hashed, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	// Insert into DB
	_, err = database.DB.Exec(
		"INSERT INTO User (Username, Password) VALUES (?, ?)",
		user.Username,
		string(hashed),
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered"})
}

// ----------------------
// Login
// ----------------------
func Login(c *gin.Context) {
	var user models.User
	var stored models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Get user by username
	err := database.DB.QueryRow(
		"SELECT ID, Username, Password FROM User WHERE Username = ?",
		user.Username,
	).Scan(&stored.ID, &stored.Username, &stored.Password)

	if err == sql.ErrNoRows || err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Compare password
	if bcrypt.CompareHashAndPassword([]byte(stored.Password), []byte(user.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": stored.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // token expires after a day
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token creation failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

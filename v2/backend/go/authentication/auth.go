package authentication

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// User model for database
type User struct {
	gorm.Model
	Email        string `gorm:"unique;not null"`
	Password     string `json:"-"`
	Name         string
	Provider     string // "local", "github", or "google"
	ProviderID   string
	RefreshToken string `json:"-"`
}

// Database instance
var db *gorm.DB

// JWT secret key
var jwtSecret []byte

// OAuth configs
var githubOauthConfig *oauth2.Config
var googleOauthConfig *oauth2.Config

// InitDB initializes the database connection
func InitDB() {
	var err error
	db, err = gorm.Open(sqlite.Open("sqlai.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto migrate the User model
	db.AutoMigrate(&User{})

	// Set JWT secret from environment or use default (in production, always use env variable)
	jwtSecret = []byte(GetEnv("JWT_SECRET", "your-default-secret-key-change-in-production"))

	// Initialize OAuth2 configs
	initOAuthConfigs()
}

// InitOAuthConfigs initializes OAuth configurations
func initOAuthConfigs() {
	githubOauthConfig = &oauth2.Config{
		ClientID:     GetEnv("GITHUB_CLIENT_ID", ""),
		ClientSecret: GetEnv("GITHUB_CLIENT_SECRET", ""),
		RedirectURL:  GetEnv("GITHUB_REDIRECT_URI", "http://localhost:8080/api/auth/github/callback"),
		Scopes:       []string{"user:email"},
		Endpoint:     github.Endpoint,
	}

	googleOauthConfig = &oauth2.Config{
		ClientID:     GetEnv("GOOGLE_CLIENT_ID", ""),
		ClientSecret: GetEnv("GOOGLE_CLIENT_SECRET", ""),
		RedirectURL:  GetEnv("GOOGLE_REDIRECT_URI", "http://localhost:8080/api/auth/google/callback"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
}

// RegisterRoutes registers all authentication routes
func RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/register", register)
	router.POST("/login", login)
	router.POST("/refresh", refreshToken)
	router.GET("/user", AuthMiddleware(), getUser)
	
	// OAuth endpoints
	router.GET("/github", githubAuth)
	router.GET("/github/callback", githubCallback)
	router.GET("/google", googleAuth)
	router.GET("/google/callback", googleCallback)
}

// GetEnv retrieves an environment variable with a fallback value
func GetEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// AuthMiddleware returns a middleware for protected routes
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		tokenString := authHeader[7:]

		// Parse and validate the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse claims"})
			c.Abort()
			return
		}

		// Set user ID in context
		c.Set("userID", uint(claims["user_id"].(float64)))
		c.Next()
	}
}

// Register handler - creates a new user
func register(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		Name     string `json:"name"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user already exists
	var existingUser User
	if result := db.Where("email = ?", input.Email).First(&existingUser); result.RowsAffected > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create user
	user := User{
		Email:    input.Email,
		Password: string(hashedPassword),
		Name:     input.Name,
		Provider: "local",
	}

	if result := db.Create(&user); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Generate tokens
	accessToken, refreshToken, err := generateTokens(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
		return
	}

	// Store refresh token in database
	db.Model(&user).Update("refresh_token", refreshToken)

	c.JSON(http.StatusCreated, gin.H{
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
		},
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// Login handler - authenticates a user
func login(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user
	var user User
	if result := db.Where("email = ?", input.Email).First(&user); result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Verify password for local accounts only
	if user.Provider == "local" {
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Please use " + user.Provider + " to login"})
		return
	}

	// Generate tokens
	accessToken, refreshToken, err := generateTokens(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
		return
	}

	// Store refresh token
	db.Model(&user).Update("refresh_token", refreshToken)

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
		},
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// Refresh token handler
func refreshToken(c *gin.Context) {
	var input struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse and validate the refresh token
	token, err := jwt.Parse(input.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse claims"})
		return
	}

	// Get user ID from claims
	userID := uint(claims["user_id"].(float64))

	// Find user
	var user User
	if result := db.First(&user, userID); result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	// Verify the refresh token matches the one in database
	if user.RefreshToken != input.RefreshToken {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	// Generate new tokens
	accessToken, refreshToken, err := generateTokens(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
		return
	}

	// Update refresh token in database
	db.Model(&user).Update("refresh_token", refreshToken)

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// Get user handler - returns the current authenticated user
func getUser(c *gin.Context) {
	userID, _ := c.Get("userID")
	
	var user User
	if result := db.First(&user, userID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
		},
	})
}

// Generate JWT tokens (access and refresh)
func generateTokens(userID uint) (string, string, error) {
	// Create access token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 1).Unix(), // 1 hour
	})
	accessTokenString, err := accessToken.SignedString(jwtSecret)
	if err != nil {
		return "", "", err
	}

	// Create refresh token
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 days
	})
	refreshTokenString, err := refreshToken.SignedString(jwtSecret)
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}
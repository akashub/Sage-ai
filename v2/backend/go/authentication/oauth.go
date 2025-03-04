package authentication

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GitHub OAuth handlers
func githubAuth(c *gin.Context) {
	// Generate a random state parameter
	state := "random-state"
	
	// Redirect to GitHub's authorization endpoint
	url := githubOauthConfig.AuthCodeURL(state)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func githubCallback(c *gin.Context) {
	// Get authorization code from callback
	code := c.Query("code")
	
	// Exchange code for token
	token, err := githubOauthConfig.Exchange(c, code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to exchange token: " + err.Error()})
		return
	}
	
	// Get user info from GitHub API
	client := githubOauthConfig.Client(c, token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info: " + err.Error()})
		return
	}
	defer resp.Body.Close()
	
	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response: " + err.Error()})
		return
	}
	
	// Parse user info
	var githubUser struct {
		ID        int    `json:"id"`
		Login     string `json:"login"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		AvatarURL string `json:"avatar_url"`
	}
	if err := json.Unmarshal(body, &githubUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user info: " + err.Error()})
		return
	}
	
	// If email is not provided, get emails from GitHub API
	if githubUser.Email == "" {
		emails, err := getGithubEmails(client)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user emails: " + err.Error()})
			return
		}
		if len(emails) > 0 {
			githubUser.Email = emails[0]
		}
	}
	
	// Check if user exists
	var user User
	result := db.Where("provider = ? AND provider_id = ?", "github", githubUser.ID).First(&user)
	
	if result.Error != nil {
		// User doesn't exist, create a new one
		user = User{
			Email:      githubUser.Email,
			Name:       githubUser.Name,
			Provider:   "github",
			ProviderID: string(githubUser.ID),
		}
		db.Create(&user)
	}
	
	// Generate tokens
	accessToken, refreshToken, err := generateTokens(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
		return
	}
	
	// Update refresh token
	db.Model(&user).Update("refresh_token", refreshToken)
	
	// Redirect to frontend with tokens in URL parameters
	frontendURL := GetEnv("FRONTEND_URL", "http://localhost:3000")
	c.Redirect(http.StatusTemporaryRedirect, 
		frontendURL + "/auth-callback?access_token=" + accessToken + 
		"&refresh_token=" + refreshToken + 
		"&user_id=" + string(user.ID))
}

// Helper function to get GitHub user emails
func getGithubEmails(client *http.Client) ([]string, error) {
	resp, err := client.Get("https://api.github.com/user/emails")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	var emails []struct {
		Email    string `json:"email"`
		Primary  bool   `json:"primary"`
		Verified bool   `json:"verified"`
	}
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	if err := json.Unmarshal(body, &emails); err != nil {
		return nil, err
	}
	
	// Find primary and verified email
	var primaryEmails []string
	for _, email := range emails {
		if email.Primary && email.Verified {
			return []string{email.Email}
		}
		if email.Verified {
			primaryEmails = append(primaryEmails, email.Email)
		}
	}
	
	return primaryEmails, nil
}

// Google OAuth handlers
func googleAuth(c *gin.Context) {
	// Generate a random state parameter
	state := "random-state"
	
	// Redirect to Google's authorization endpoint
	url := googleOauthConfig.AuthCodeURL(state)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func googleCallback(c *gin.Context) {
	// Get authorization code from callback
	code := c.Query("code")
	
	// Exchange code for token
	token, err := googleOauthConfig.Exchange(c, code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to exchange token: " + err.Error()})
		return
	}
	
	// Get user info from Google API
	client := googleOauthConfig.Client(c, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info: " + err.Error()})
		return
	}
	defer resp.Body.Close()
	
	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response: " + err.Error()})
		return
	}
	
	// Parse user info
	var googleUser struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}
	if err := json.Unmarshal(body, &googleUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user info: " + err.Error()})
		return
	}
	
	// Check if user exists
	var user User
	result := db.Where("provider = ? AND provider_id = ?", "google", googleUser.ID).First(&user)
	
	if result.Error != nil {
		// User doesn't exist, create a new one
		user = User{
			Email:      googleUser.Email,
			Name:       googleUser.Name,
			Provider:   "google",
			ProviderID: googleUser.ID,
		}
		db.Create(&user)
	}
	
	// Generate tokens
	accessToken, refreshToken, err := generateTokens(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
		return
	}
	
	// Update refresh token
	db.Model(&user).Update("refresh_token", refreshToken)
	
	// Redirect to frontend with tokens in URL parameters
	frontendURL := GetEnv("FRONTEND_URL", "http://localhost:3000")
	c.Redirect(http.StatusTemporaryRedirect, 
		frontendURL + "/auth-callback?access_token=" + accessToken + 
		"&refresh_token=" + refreshToken + 
		"&user_id=" + string(user.ID))
}
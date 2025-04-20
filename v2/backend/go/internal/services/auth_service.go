// backend/go/internal/services/auth_service.go
package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
	"strings"
	"net/url"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"sage-ai-v2/internal/models"
	"sage-ai-v2/pkg/logger"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidCredential = errors.New("invalid credentials")
	ErrUserExists        = errors.New("user already exists")
)

type AuthService struct {
	db         *sql.DB
	jwtSecret  []byte
	jwtExpiry  time.Duration
	oauthConfs map[string]models.OAuthConfig
}

func NewAuthService(db *sql.DB, jwtSecret string, jwtExpiry time.Duration, oauthConfs map[string]models.OAuthConfig) *AuthService {
	return &AuthService{
		db:         db,
		jwtSecret:  []byte(jwtSecret),
		jwtExpiry:  jwtExpiry,
		oauthConfs: oauthConfs,
	}
}

// SignIn authenticates a user with email/password
func (s *AuthService) SignIn(ctx context.Context, req models.SignInRequest) (*models.AuthResponse, error) {
	logger.InfoLogger.Printf("Attempting to sign in user: %s", req.Email)

	// Find user by email
	user, err := s.getUserByEmail(ctx, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.ErrorLogger.Printf("Sign in failed for %s: user not found", req.Email)
			return nil, ErrUserNotFound
		}
		logger.ErrorLogger.Printf("Database error during sign in: %v", err)
		return nil, fmt.Errorf("database error: %w", err)
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		logger.ErrorLogger.Printf("Sign in failed for %s: invalid password", req.Email)
		return nil, ErrInvalidCredential
	}

	// Update last login time
	if err := s.updateLastLogin(ctx, user.ID); err != nil {
		logger.ErrorLogger.Printf("Failed to update last login time: %v", err)
		// Non-critical error, continue with sign-in
	}

	// Generate JWT token
	token, err := s.generateJWT(user)
	if err != nil {
		logger.ErrorLogger.Printf("Failed to generate JWT: %v", err)
		return nil, fmt.Errorf("token generation failed: %w", err)
	}

	logger.InfoLogger.Printf("User signed in successfully: %s", req.Email)
	return &models.AuthResponse{
		User:        *user,
		AccessToken: token,
	}, nil
}

// SignUp creates a new user with email/password
func (s *AuthService) SignUp(ctx context.Context, req models.SignUpRequest) (*models.AuthResponse, error) {
	logger.InfoLogger.Printf("Attempting to create new user: %s", req.Email)

	// Check if user already exists
	exists, err := s.userExists(ctx, req.Email)
	if err != nil {
		logger.ErrorLogger.Printf("Database error checking user existence: %v", err)
		return nil, fmt.Errorf("database error: %w", err)
	}

	if exists {
		logger.ErrorLogger.Printf("Sign up failed: user already exists: %s", req.Email)
		return nil, ErrUserExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.ErrorLogger.Printf("Failed to hash password: %v", err)
		return nil, fmt.Errorf("password hashing failed: %w", err)
	}

	// Create user
	user := &models.User{
		ID:           uuid.New().String(),
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Name:         req.Name,
		CreatedAt:    time.Now(),
		LastLoginAt:  time.Now(),
		ProviderType: "email",
	}

	// Insert user into database
	if err := s.createUser(ctx, user); err != nil {
		logger.ErrorLogger.Printf("Failed to create user in database: %v", err)
		return nil, fmt.Errorf("user creation failed: %w", err)
	}

	// Generate JWT token
	token, err := s.generateJWT(user)
	if err != nil {
		logger.ErrorLogger.Printf("Failed to generate JWT: %v", err)
		return nil, fmt.Errorf("token generation failed: %w", err)
	}

	logger.InfoLogger.Printf("User created successfully: %s", req.Email)
	return &models.AuthResponse{
		User:        *user,
		AccessToken: token,
	}, nil
}

// OAuthSignIn handles sign-in/sign-up via OAuth providers (Google, GitHub)
func (s *AuthService) OAuthSignIn(ctx context.Context, provider string, code string, redirectURI string) (*models.AuthResponse, error) {
	logger.InfoLogger.Printf("OAuth sign-in attempt with provider: %s", provider)

	// Get provider config
	conf, ok := s.oauthConfs[provider]
	if !ok {
		logger.ErrorLogger.Printf("Invalid OAuth provider: %s", provider)
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}

	// Override redirect URI if provided
	if redirectURI != "" {
		conf.RedirectURI = redirectURI
	}

	// Exchange code for token (implementation depends on the OAuth provider)
	userInfo, err := s.exchangeCodeForUserInfo(ctx, provider, code, conf)
	if err != nil {
		logger.ErrorLogger.Printf("Failed to exchange code for user info: %v", err)
		return nil, fmt.Errorf("OAuth authentication failed: %w", err)
	}

	// Check if user exists
	user, err := s.getUserByProviderID(ctx, provider, userInfo["id"])
	if err != nil && err != sql.ErrNoRows {
		logger.ErrorLogger.Printf("Database error checking OAuth user: %v", err)
		return nil, fmt.Errorf("database error: %w", err)
	}

	// Create user if not exists
	if err == sql.ErrNoRows {
		logger.InfoLogger.Printf("Creating new user from OAuth provider: %s", provider)
		user = &models.User{
			ID:            uuid.New().String(),
			Email:         userInfo["email"],
			Name:          userInfo["name"],
			CreatedAt:     time.Now(),
			LastLoginAt:   time.Now(),
			ProviderType:  provider,
			ProviderID:    userInfo["id"],
			ProfilePicURL: userInfo["picture"],
		}

		if err := s.createUser(ctx, user); err != nil {
			logger.ErrorLogger.Printf("Failed to create OAuth user in database: %v", err)
			return nil, fmt.Errorf("user creation failed: %w", err)
		}
	} else {
		// Update last login time
		if err := s.updateLastLogin(ctx, user.ID); err != nil {
			logger.ErrorLogger.Printf("Failed to update last login time: %v", err)
			// Non-critical error, continue with sign-in
		}
	}

	// Generate JWT token
	token, err := s.generateJWT(user)
	if err != nil {
		logger.ErrorLogger.Printf("Failed to generate JWT: %v", err)
		return nil, fmt.Errorf("token generation failed: %w", err)
	}

	logger.InfoLogger.Printf("OAuth user signed in successfully: %s via %s", user.Email, provider)
	return &models.AuthResponse{
		User:        *user,
		AccessToken: token,
	}, nil
}

// VerifyToken validates a JWT token and returns the user ID
func (s *AuthService) VerifyToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate signing algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.jwtSecret, nil
	})

	if err != nil {
		return "", fmt.Errorf("token parsing failed: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check expiration
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			return "", errors.New("token expired")
		}
		return claims["sub"].(string), nil
	}

	return "", errors.New("invalid token")
}

// Helper methods

func (s *AuthService) generateJWT(user *models.User) (string, error) {
	expirationTime := time.Now().Add(s.jwtExpiry)
	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": expirationTime.Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *AuthService) getUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `SELECT id, email, password_hash, name, created_at, last_login_at, 
              provider_type, provider_id, refresh_token, profile_pic_url 
              FROM users WHERE email = ?`
	
	user := &models.User{}
	err := s.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.CreatedAt,
		&user.LastLoginAt, &user.ProviderType, &user.ProviderID, &user.RefreshToken,
		&user.ProfilePicURL,
	)
	
	if err != nil {
		return nil, err
	}
	
	return user, nil
}

func (s *AuthService) getUserByProviderID(ctx context.Context, provider, providerID string) (*models.User, error) {
	query := `SELECT id, email, password_hash, name, created_at, last_login_at, 
              provider_type, provider_id, refresh_token, profile_pic_url 
              FROM users WHERE provider_type = ? AND provider_id = ?`
	
	user := &models.User{}
	err := s.db.QueryRowContext(ctx, query, provider, providerID).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.CreatedAt,
		&user.LastLoginAt, &user.ProviderType, &user.ProviderID, &user.RefreshToken,
		&user.ProfilePicURL,
	)
	
	if err != nil {
		return nil, err
	}
	
	return user, nil
}

func (s *AuthService) userExists(ctx context.Context, email string) (bool, error) {
	query := `SELECT COUNT(*) FROM users WHERE email = ?`
	
	var count int
	err := s.db.QueryRowContext(ctx, query, email).Scan(&count)
	if err != nil {
		return false, err
	}
	
	return count > 0, nil
}

func (s *AuthService) createUser(ctx context.Context, user *models.User) error {
	query := `INSERT INTO users (id, email, password_hash, name, created_at, last_login_at, 
              provider_type, provider_id, refresh_token, profile_pic_url) 
              VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	_, err := s.db.ExecContext(ctx, query,
		user.ID, user.Email, user.PasswordHash, user.Name, user.CreatedAt,
		user.LastLoginAt, user.ProviderType, user.ProviderID, user.RefreshToken,
		user.ProfilePicURL,
	)
	
	return err
}

func (s *AuthService) updateLastLogin(ctx context.Context, userID string) error {
	query := `UPDATE users SET last_login_at = ? WHERE id = ?`
	
	_, err := s.db.ExecContext(ctx, query, time.Now(), userID)
	return err
}

func (s *AuthService) exchangeCodeForUserInfo(ctx context.Context, provider, code string, conf models.OAuthConfig) (map[string]string, error) {
	switch provider {
	case "google":
		return s.exchangeGoogleCode(ctx, code, conf)
	case "github":
		return s.exchangeGitHubCode(ctx, code, conf)
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
}

func (s *AuthService) GetOAuthURL(provider, redirectURI string) (string, error) {
	// Get provider config
	conf, ok := s.oauthConfs[provider]
	if !ok {
		return "", fmt.Errorf("unsupported provider: %s", provider)
	}

	// Override redirect URI if provided
	if redirectURI != "" {
		conf.RedirectURI = redirectURI
	}

	// Build authorization URL
	params := url.Values{}
	params.Set("client_id", conf.ClientID)
	params.Set("redirect_uri", conf.RedirectURI)
	params.Set("scope", strings.Join(conf.Scopes, " "))
	params.Set("response_type", "code")
	params.Set("state", generateRandomState())

	// Add provider-specific parameters
	switch provider {
	case "google":
		params.Set("access_type", "offline")
		params.Set("prompt", "consent")
	case "github":
		// GitHub doesn't require additional parameters
	}

	return fmt.Sprintf("%s?%s", conf.AuthURL, params.Encode()), nil
}

// GetUserByID retrieves a user by ID
func (s *AuthService) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
	logger.InfoLogger.Printf("Getting user by ID: %s", userID)

	query := `SELECT id, email, password_hash, name, created_at, last_login_at, 
              provider_type, provider_id, refresh_token, profile_pic_url 
              FROM users WHERE id = ?`
	
	user := &models.User{}
	err := s.db.QueryRowContext(ctx, query, userID).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.CreatedAt,
		&user.LastLoginAt, &user.ProviderType, &user.ProviderID, &user.RefreshToken,
		&user.ProfilePicURL,
	)
	
	if err != nil {
		return nil, err
	}
	
	return user, nil
}

// Helper function to generate a random state for OAuth
func generateRandomState() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
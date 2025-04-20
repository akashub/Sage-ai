// handlers/auth_test.go
package handlers

import (
    "bytes"
    "context"
    "encoding/json"
    "errors"
    "net/http"
    "net/http/httptest"
    "testing"
    "time"

    "sage-ai-v2/internal/models"
)

// Define error constants used in tests
var (
    ErrInternal = errors.New("internal error")
    ErrUserNotFound = errors.New("user not found")
    ErrInvalidCredential = errors.New("invalid credentials")
    ErrUserExists = errors.New("user already exists")
)

// MockAuthService is a mock implementation of the AuthService for testing
type MockAuthService struct {
    // Mock fields to track calls and set return values
    SignInCalled       bool
    SignUpCalled       bool
    OAuthSignInCalled  bool
    GetOAuthURLCalled  bool
    VerifyTokenCalled  bool
    GetUserByIDCalled  bool
    
    // Return values
    SignInResponse     *models.AuthResponse
    SignUpResponse     *models.AuthResponse
    OAuthSignInResponse *models.AuthResponse
    OAuthURL           string
    UserID             string
    User               *models.User
    
    // Error returns
    SignInErr          error
    SignUpErr          error
    OAuthSignInErr     error
    GetOAuthURLErr     error
    VerifyTokenErr     error
    GetUserByIDErr     error
}

// Implement the AuthService interface methods
func (m *MockAuthService) SignIn(ctx context.Context, req models.SignInRequest) (*models.AuthResponse, error) {
    m.SignInCalled = true
    return m.SignInResponse, m.SignInErr
}

func (m *MockAuthService) SignUp(ctx context.Context, req models.SignUpRequest) (*models.AuthResponse, error) {
    m.SignUpCalled = true
    return m.SignUpResponse, m.SignUpErr
}

func (m *MockAuthService) OAuthSignIn(ctx context.Context, provider, code, redirectURI string) (*models.AuthResponse, error) {
    m.OAuthSignInCalled = true
    return m.OAuthSignInResponse, m.OAuthSignInErr
}

func (m *MockAuthService) GetOAuthURL(provider, redirectURI string) (string, error) {
    m.GetOAuthURLCalled = true
    return m.OAuthURL, m.GetOAuthURLErr
}

func (m *MockAuthService) VerifyToken(token string) (string, error) {
    m.VerifyTokenCalled = true
    return m.UserID, m.VerifyTokenErr
}

func (m *MockAuthService) GetUserByID(ctx context.Context, id string) (*models.User, error) {
    m.GetUserByIDCalled = true
    return m.User, m.GetUserByIDErr
}

// TestSignInHandler tests the SignInHandler function
func TestSignInHandler(t *testing.T) {
    // Create test cases
    testCases := []struct {
        name           string
        method         string
        requestBody    models.SignInRequest
        mockResponse   *models.AuthResponse
        mockError      error
        expectedStatus int
    }{
        {
            name:   "Valid sign in",
            method: "POST",
            requestBody: models.SignInRequest{
                Email:    "test@example.com",
                Password: "password123",
            },
            mockResponse: &models.AuthResponse{
				User: models.User{ 
					ID:             "123",
					Email:          "test@example.com",
					Name:           "Test User",
					PasswordHash:   "",
					CreatedAt:      time.Now(),
					LastLoginAt:    time.Now(),
				},
				AccessToken: "test-token",
			},
            mockError:      nil,
            expectedStatus: http.StatusOK,
        },
        {
            name:   "Invalid credentials",
            method: "POST",
            requestBody: models.SignInRequest{
                Email:    "test@example.com",
                Password: "wrongpassword",
            },
            mockResponse:   nil,
            mockError:      ErrInvalidCredential,
            expectedStatus: http.StatusUnauthorized,
        },
        {
            name:   "User not found",
            method: "POST",
            requestBody: models.SignInRequest{
                Email:    "nonexistent@example.com",
                Password: "password123",
            },
            mockResponse:   nil,
            mockError:      ErrUserNotFound,
            expectedStatus: http.StatusNotFound,
        },
        {
            name:   "Internal error",
            method: "POST",
            requestBody: models.SignInRequest{
                Email:    "test@example.com",
                Password: "password123",
            },
            mockResponse:   nil,
            mockError:      ErrInternal,
            expectedStatus: http.StatusInternalServerError,
        },
        {
            name:           "Method not allowed",
            method:         "GET",
            requestBody:    models.SignInRequest{},
            mockResponse:   nil,
            mockError:      nil,
            expectedStatus: http.StatusMethodNotAllowed,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            // Create a mock auth service
            mockService := &MockAuthService{
                SignInResponse: tc.mockResponse,
                SignInErr:      tc.mockError,
            }

            // Create a new AuthHandler with the mock service
            handler := NewAuthHandler(mockService)

            // Create a request
            reqBody, _ := json.Marshal(tc.requestBody)
            req, _ := http.NewRequest(tc.method, "/api/auth/signin", bytes.NewBuffer(reqBody))
            req.Header.Set("Content-Type", "application/json")

            // Create a response recorder
            rr := httptest.NewRecorder()

            // Call the handler
            handler.SignInHandler(rr, req)

            // Check the status code
            if status := rr.Code; status != tc.expectedStatus {
                t.Errorf("handler returned wrong status code: got %v want %v", status, tc.expectedStatus)
            }

            // For successful requests, check the response
            if tc.expectedStatus == http.StatusOK {
                var response models.AuthResponse
                if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
                    t.Errorf("failed to parse response JSON: %v", err)
                }

                // Verify the response
                if response.User.ID != tc.mockResponse.User.ID ||
                   response.User.Email != tc.mockResponse.User.Email ||
                   response.AccessToken != tc.mockResponse.AccessToken {
                    t.Errorf("handler returned unexpected response: got %v want %v", response, tc.mockResponse)
                }

                // Check that cookies were set
                cookies := rr.Result().Cookies()
                var authCookie *http.Cookie
                for _, cookie := range cookies {
                    if cookie.Name == "auth_token" {
                        authCookie = cookie
                        break
                    }
                }

                if authCookie == nil {
                    t.Errorf("auth_token cookie not set")
                } else if authCookie.Value != tc.mockResponse.AccessToken {
                    t.Errorf("auth_token cookie has wrong value: got %v want %v", authCookie.Value, tc.mockResponse.AccessToken)
                }
            }

            // Verify that the service method was called
            if tc.method == "POST" && tc.expectedStatus != http.StatusMethodNotAllowed {
                if !mockService.SignInCalled {
                    t.Error("SignIn method was not called")
                }
            }
        })
    }
}

// TestSignUpHandler tests the SignUpHandler function
func TestSignUpHandler(t *testing.T) {
    testCases := []struct {
        name           string
        method         string
        requestBody    models.SignUpRequest
        mockResponse   *models.AuthResponse
        mockError      error
        expectedStatus int
    }{
        {
            name:   "Valid sign up",
            method: "POST",
            requestBody: models.SignUpRequest{
                Email:    "new@example.com",
                Password: "password123",
                Name:     "New User",
            },  
			mockResponse: &models.AuthResponse{
				User: models.User{ 
					ID:             "123",
					Email:          "new@example.com",
					Name:           "New User",
					PasswordHash:   "",
					CreatedAt:      time.Now(),
					LastLoginAt:    time.Now(),
				},
				AccessToken: "test-token",
			},
            mockError:      nil,
            expectedStatus: http.StatusCreated,
        },
        {
            name:   "User already exists",
            method: "POST",
            requestBody: models.SignUpRequest{
                Email:    "existing@example.com",
                Password: "password123",
                Name:     "Existing User",
            },
            mockResponse:   nil,
            mockError:      ErrUserExists,
            expectedStatus: http.StatusConflict,
        },
        {
            name:   "Internal error",
            method: "POST",
            requestBody: models.SignUpRequest{
                Email:    "test@example.com",
                Password: "password123",
                Name:     "Test User",
            },
            mockResponse:   nil,
            mockError:      ErrInternal,
            expectedStatus: http.StatusInternalServerError,
        },
        {
            name:           "Method not allowed",
            method:         "GET",
            requestBody:    models.SignUpRequest{},
            mockResponse:   nil,
            mockError:      nil,
            expectedStatus: http.StatusMethodNotAllowed,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            // Create a mock auth service
            mockService := &MockAuthService{
                SignUpResponse: tc.mockResponse,
                SignUpErr:      tc.mockError,
            }

            // Create a new AuthHandler with the mock service
            handler := NewAuthHandler(mockService)

            // Create a request
            reqBody, _ := json.Marshal(tc.requestBody)
            req, _ := http.NewRequest(tc.method, "/api/auth/signup", bytes.NewBuffer(reqBody))
            req.Header.Set("Content-Type", "application/json")

            // Create a response recorder
            rr := httptest.NewRecorder()

            // Call the handler
            handler.SignUpHandler(rr, req)

            // Check the status code
            if status := rr.Code; status != tc.expectedStatus {
                t.Errorf("handler returned wrong status code: got %v want %v", status, tc.expectedStatus)
            }

            // For successful requests, check the response
            if tc.expectedStatus == http.StatusCreated {
                var response models.AuthResponse
                if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
                    t.Errorf("failed to parse response JSON: %v", err)
                }

                // Verify the response
                if response.User.ID != tc.mockResponse.User.ID ||
                   response.User.Email != tc.mockResponse.User.Email ||
                   response.AccessToken != tc.mockResponse.AccessToken {
                    t.Errorf("handler returned unexpected response: got %v want %v", response, tc.mockResponse)
                }

                // Check that cookies were set
                cookies := rr.Result().Cookies()
                var authCookie *http.Cookie
                for _, cookie := range cookies {
                    if cookie.Name == "auth_token" {
                        authCookie = cookie
                        break
                    }
                }

                if authCookie == nil {
                    t.Errorf("auth_token cookie not set")
                } else if authCookie.Value != tc.mockResponse.AccessToken {
                    t.Errorf("auth_token cookie has wrong value: got %v want %v", authCookie.Value, tc.mockResponse.AccessToken)
                }
            }

            // Verify that the service method was called
            if tc.method == "POST" && tc.expectedStatus != http.StatusMethodNotAllowed {
                if !mockService.SignUpCalled {
                    t.Error("SignUp method was not called")
                }
            }
        })
    }
}

// TestOAuthURLHandler tests the OAuthURLHandler function
func TestOAuthURLHandler(t *testing.T) {
    testCases := []struct {
        name           string
        method         string
        provider       string
        redirectURI    string
        mockURL        string
        mockError      error
        expectedStatus int
    }{
        {
            name:           "Valid GitHub OAuth URL",
            method:         "GET",
            provider:       "github",
            redirectURI:    "http://localhost:3000/oauth-callback",
            mockURL:        "https://github.com/login/oauth/authorize?client_id=test-id&redirect_uri=http%3A%2F%2Flocalhost%3A3000%2Foauth-callback&response_type=code&scope=read%3Auser+user%3Aemail",
            mockError:      nil,
            expectedStatus: http.StatusOK,
        },
        {
            name:           "Valid Google OAuth URL",
            method:         "GET",
            provider:       "google",
            redirectURI:    "http://localhost:3000/oauth-callback",
            mockURL:        "https://accounts.google.com/o/oauth2/auth?client_id=test-id&redirect_uri=http%3A%2F%2Flocalhost%3A3000%2Foauth-callback&response_type=code&scope=email+profile",
            mockError:      nil,
            expectedStatus: http.StatusOK,
        },
        {
            name:           "Missing redirect URI",
            method:         "GET",
            provider:       "github",
            redirectURI:    "",
            mockURL:        "",
            mockError:      nil,
            expectedStatus: http.StatusBadRequest,
        },
        {
            name:           "Method not allowed",
            method:         "POST",
            provider:       "github",
            redirectURI:    "http://localhost:3000/oauth-callback",
            mockURL:        "",
            mockError:      nil,
            expectedStatus: http.StatusMethodNotAllowed,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            // Create a mock auth service
            mockService := &MockAuthService{
                OAuthURL:      tc.mockURL,
                GetOAuthURLErr: tc.mockError,
            }

            // Create a new AuthHandler with the mock service
            handler := NewAuthHandler(mockService)

            // Create a request
            url := "/api/auth/oauth/url/" + tc.provider
            if tc.redirectURI != "" {
                url += "?redirect_uri=" + tc.redirectURI
            }
            req, _ := http.NewRequest(tc.method, url, nil)

            // Create a response recorder
            rr := httptest.NewRecorder()

            // Call the handler
            handler.OAuthURLHandler(rr, req)

            // Check the status code
            if status := rr.Code; status != tc.expectedStatus {
                t.Errorf("handler returned wrong status code: got %v want %v", status, tc.expectedStatus)
            }

            // For successful requests, check the response
            if tc.expectedStatus == http.StatusOK {
                var response struct {
                    URL string `json:"url"`
                }
                if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
                    t.Errorf("failed to parse response JSON: %v", err)
                }

                // Verify the response
                if response.URL != tc.mockURL {
                    t.Errorf("handler returned unexpected URL: got %v want %v", response.URL, tc.mockURL)
                }
            }

            // Verify that the service method was called when appropriate
            if tc.method == "GET" && tc.redirectURI != "" && tc.expectedStatus != http.StatusMethodNotAllowed {
                if !mockService.GetOAuthURLCalled {
                    t.Error("GetOAuthURL method was not called")
                }
            }
        })
    }
}

// TestOAuthSignInHandler tests the OAuthSignInHandler function
func TestOAuthSignInHandler(t *testing.T) {
    testCases := []struct {
        name           string
        method         string
        provider       string
        requestBody    models.OAuthRequest
        mockResponse   *models.AuthResponse
        mockError      error
        expectedStatus int
    }{
        {
            name:     "Valid GitHub OAuth sign in",
            method:   "POST",
            provider: "github",
            requestBody: models.OAuthRequest{
                Code:        "test-code",
                RedirectURI: "http://localhost:3000/oauth-callback",
            },
			mockResponse: &models.AuthResponse{
				User: models.User{ 
					ID:             "123",
					Email:          "user@github.com",
					Name:           "GitHub User",
					PasswordHash:   "",
					CreatedAt:      time.Now(),
					LastLoginAt:    time.Now(),
				},
				AccessToken: "test-token",
			},
            mockError:      nil,
            expectedStatus: http.StatusOK,
        },
        {
            name:     "Valid Google OAuth sign in",
            method:   "POST",
            provider: "google",
            requestBody: models.OAuthRequest{
                Code:        "test-code",
                RedirectURI: "http://localhost:3000/oauth-callback",
            },
			mockResponse: &models.AuthResponse{
				User: models.User{ 
					ID:             "456",
					Email:          "user@gmail.com",
					Name:           "Goog;e User",
					PasswordHash:   "",
					CreatedAt:      time.Now(),
					LastLoginAt:    time.Now(),
				},
				AccessToken: "test-token",
			},
            mockError:      nil,
            expectedStatus: http.StatusOK,
        },
        {
            name:     "OAuth error",
            method:   "POST",
            provider: "github",
            requestBody: models.OAuthRequest{
                Code:        "invalid-code",
                RedirectURI: "http://localhost:3000/oauth-callback",
            },
            mockResponse:   nil,
            mockError:      errors.New("OAuth authentication failed"),
            expectedStatus: http.StatusInternalServerError,
        },
        {
            name:     "Rate limit error",
            method:   "POST",
            provider: "github",
            requestBody: models.OAuthRequest{
                Code:        "test-code",
                RedirectURI: "http://localhost:3000/oauth-callback",
            },
            mockResponse:   nil,
            mockError:      errors.New("429: rate limit exceeded"),
            expectedStatus: http.StatusTooManyRequests,
        },
        {
            name:           "Method not allowed",
            method:         "GET",
            provider:       "github",
            requestBody:    models.OAuthRequest{},
            mockResponse:   nil,
            mockError:      nil,
            expectedStatus: http.StatusMethodNotAllowed,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            // Create a mock auth service
            mockService := &MockAuthService{
                OAuthSignInResponse: tc.mockResponse,
                OAuthSignInErr:      tc.mockError,
            }

            // Create a new AuthHandler with the mock service
            handler := NewAuthHandler(mockService)

            // Create a request
            reqBody, _ := json.Marshal(tc.requestBody)
            req, _ := http.NewRequest(tc.method, "/api/auth/oauth/"+tc.provider, bytes.NewBuffer(reqBody))
            req.Header.Set("Content-Type", "application/json")

            // Create a response recorder
            rr := httptest.NewRecorder()

            // Call the handler
            handler.OAuthSignInHandler(rr, req)

            // Check the status code
            if status := rr.Code; status != tc.expectedStatus {
                t.Errorf("handler returned wrong status code: got %v want %v", status, tc.expectedStatus)
            }

            // For successful requests, check the response
            if tc.expectedStatus == http.StatusOK {
                var response models.AuthResponse
                if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
                    t.Errorf("failed to parse response JSON: %v", err)
                }

                // Verify the response
                if response.User.ID != tc.mockResponse.User.ID ||
                   response.User.Email != tc.mockResponse.User.Email ||
                   response.AccessToken != tc.mockResponse.AccessToken {
                    t.Errorf("handler returned unexpected response: got %v want %v", response, tc.mockResponse)
                }

                // Check that cookies were set
                cookies := rr.Result().Cookies()
                var authCookie *http.Cookie
                for _, cookie := range cookies {
                    if cookie.Name == "auth_token" {
                        authCookie = cookie
                        break
                    }
                }

                if authCookie == nil {
                    t.Errorf("auth_token cookie not set")
                } else if authCookie.Value != tc.mockResponse.AccessToken {
                    t.Errorf("auth_token cookie has wrong value: got %v want %v", authCookie.Value, tc.mockResponse.AccessToken)
                }
            }

            // Verify that the service method was called
            if tc.method == "POST" && tc.expectedStatus != http.StatusMethodNotAllowed {
                if !mockService.OAuthSignInCalled {
                    t.Error("OAuthSignIn method was not called")
                }
            }
        })
    }
}

// TestSignOutHandler tests the SignOutHandler function
func TestSignOutHandler(t *testing.T) {
    testCases := []struct {
        name           string
        method         string
        expectedStatus int
    }{
        {
            name:           "Valid sign out",
            method:         "POST",
            expectedStatus: http.StatusOK,
        },
        {
            name:           "Method not allowed",
            method:         "GET",
            expectedStatus: http.StatusMethodNotAllowed,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            // Create a mock auth service
            mockService := &MockAuthService{}

            // Create a new AuthHandler with the mock service
            handler := NewAuthHandler(mockService)

            // Create a request
            req, _ := http.NewRequest(tc.method, "/api/auth/signout", nil)

            // Create a response recorder
            rr := httptest.NewRecorder()

            // Call the handler
            handler.SignOutHandler(rr, req)

            // Check the status code
            if status := rr.Code; status != tc.expectedStatus {
                t.Errorf("handler returned wrong status code: got %v want %v", status, tc.expectedStatus)
            }

            // For successful requests, check the response
            if tc.expectedStatus == http.StatusOK {
                var response struct {
                    Success bool `json:"success"`
                }
                if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
                    t.Errorf("failed to parse response JSON: %v", err)
                }

                // Verify the response
                if !response.Success {
                    t.Errorf("handler returned unexpected response: got %v want true", response.Success)
                }

                // Check that cookies were cleared
                cookies := rr.Result().Cookies()
                var authCookie *http.Cookie
                for _, cookie := range cookies {
                    if cookie.Name == "auth_token" {
                        authCookie = cookie
                        break
                    }
                }

                if authCookie == nil {
                    t.Errorf("auth_token cookie not set")
                } else if authCookie.MaxAge != -1 {
                    t.Errorf("auth_token cookie not set to expire: got %v want -1", authCookie.MaxAge)
                }
            }
        })
    }
}

// TestGetUserHandler tests the GetUserHandler function
func TestGetUserHandler(t *testing.T) {
    testCases := []struct {
        name           string
        method         string
        token          string
        mockUserID     string
        mockUser       *models.User
        mockVerifyErr  error
        mockGetUserErr error
        expectedStatus int
    }{
        {
            name:       "Valid token",
            method:     "GET",
            token:      "valid-token",
            mockUserID: "123",
            mockUser: &models.User{
                ID:          "123",
                Email:       "user@example.com",
                Name:        "Test User",
                CreatedAt:   time.Now(),
                LastLoginAt: time.Now(),
            },
            mockVerifyErr:  nil,
            mockGetUserErr: nil,
            expectedStatus: http.StatusOK,
        },
        {
            name:           "Invalid token",
            method:         "GET",
            token:          "invalid-token",
            mockUserID:     "",
            mockUser:       nil,
            mockVerifyErr:  errors.New("invalid token"),
            mockGetUserErr: nil,
            expectedStatus: http.StatusUnauthorized,
        },
        {
            name:           "User not found",
            method:         "GET",
            token:          "valid-token",
            mockUserID:     "456",
            mockUser:       nil,
            mockVerifyErr:  nil,
            mockGetUserErr: errors.New("user not found"),
            expectedStatus: http.StatusInternalServerError,
        },
        {
            name:           "No token",
            method:         "GET",
            token:          "",
            mockUserID:     "",
            mockUser:       nil,
            mockVerifyErr:  nil,
            mockGetUserErr: nil,
            expectedStatus: http.StatusUnauthorized,
        },
        {
            name:           "Method not allowed",
            method:         "POST",
            token:          "valid-token",
            mockUserID:     "",
            mockUser:       nil,
            mockVerifyErr:  nil,
            mockGetUserErr: nil,
            expectedStatus: http.StatusMethodNotAllowed,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            // Create a mock auth service
            mockService := &MockAuthService{
                UserID:        tc.mockUserID,
                User:          tc.mockUser,
                VerifyTokenErr: tc.mockVerifyErr,
                GetUserByIDErr: tc.mockGetUserErr,
            }

            // Create a new AuthHandler with the mock service
            handler := NewAuthHandler(mockService)

            // Create a request
            req, _ := http.NewRequest(tc.method, "/api/auth/user", nil)
            
            // Add token if needed
            if tc.token != "" {
                req.Header.Set("Authorization", "Bearer "+tc.token)
            }

            // Create a response recorder
            rr := httptest.NewRecorder()

            // Call the handler
            handler.GetUserHandler(rr, req)

            // Check the status code
            if status := rr.Code; status != tc.expectedStatus {
                t.Errorf("handler returned wrong status code: got %v want %v", status, tc.expectedStatus)
            }

            // For successful requests, check the response
            if tc.expectedStatus == http.StatusOK {
                var response models.User
                if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
                    t.Errorf("failed to parse response JSON: %v", err)
                }

                // Verify the response
                if response.ID != tc.mockUser.ID ||
                   response.Email != tc.mockUser.Email ||
                   response.Name != tc.mockUser.Name {
                    t.Errorf("handler returned unexpected user: got %v want %v", response, tc.mockUser)
                }
            }

            // Verify that the service methods were called appropriately
            if tc.method == "GET" && tc.token != "" && tc.expectedStatus != http.StatusMethodNotAllowed {
                if !mockService.VerifyTokenCalled {
                    t.Error("VerifyToken method was not called")
                }
                
                if tc.mockVerifyErr == nil && !mockService.GetUserByIDCalled {
                    t.Error("GetUserByID method was not called when token was valid")
                }
            }
        })
    }
}

// TestAuthMiddleware tests the AuthMiddleware function
func TestAuthMiddleware(t *testing.T) {
    testCases := []struct {
        name           string
        token          string
        mockUserID     string
        mockVerifyErr  error
        expectedStatus int
    }{
        {
            name:           "Valid token",
            token:          "valid-token",
            mockUserID:     "123",
            mockVerifyErr:  nil,
            expectedStatus: http.StatusOK,
        },
        {
            name:           "Invalid token",
            token:          "invalid-token",
            mockUserID:     "",
            mockVerifyErr:  errors.New("invalid token"),
            expectedStatus: http.StatusUnauthorized,
        },
        {
            name:           "No token",
            token:          "",
            mockUserID:     "",
            mockVerifyErr:  nil,
            expectedStatus: http.StatusUnauthorized,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            // Create a mock auth service
            mockService := &MockAuthService{
                UserID:        tc.mockUserID,
                VerifyTokenErr: tc.mockVerifyErr,
            }

            // Create a new AuthHandler with the mock service
            handler := NewAuthHandler(mockService)

            // Create a test handler that will be wrapped by the middleware
            testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                // Check if user ID was added to context
                userID := r.Context().Value("userID")
                if userID != tc.mockUserID {
                    t.Errorf("userID not set in context: got %v want %v", userID, tc.mockUserID)
                }
                
                w.WriteHeader(http.StatusOK)
                w.Write([]byte("test handler called"))
            })

            // Create a request
            req, _ := http.NewRequest("GET", "/protected", nil)
            
            // Add token if needed
            if tc.token != "" {
                req.Header.Set("Authorization", "Bearer "+tc.token)
            }

            // Create a response recorder
            rr := httptest.NewRecorder()

            // Apply the middleware to the test handler
            middlewareHandler := handler.AuthMiddleware(testHandler)
            
            // Call the handler with middleware
            middlewareHandler.ServeHTTP(rr, req)

            // Check the status code
            if status := rr.Code; status != tc.expectedStatus {
                t.Errorf("middleware returned wrong status code: got %v want %v", status, tc.expectedStatus)
            }

            // Verify that the service method was called if a token was provided
            if tc.token != "" && !mockService.VerifyTokenCalled {
                t.Error("VerifyToken method was not called")
            }
        })
    }
}
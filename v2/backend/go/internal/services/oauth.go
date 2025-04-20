// backend/go/internal/services/oauth.go
package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	
	"sage-ai-v2/internal/models"
	"sage-ai-v2/pkg/logger"
)

// Implement the real OAuth exchange for Google
func (s *AuthService) exchangeGoogleCode(ctx context.Context, code string, conf models.OAuthConfig) (map[string]string, error) {
	logger.InfoLogger.Printf("Exchanging Google OAuth code for token")
	
	// Exchange authorization code for token
	data := url.Values{}
	data.Set("code", code)
	data.Set("client_id", conf.ClientID)
	data.Set("client_secret", conf.ClientSecret)
	data.Set("redirect_uri", conf.RedirectURI)
	data.Set("grant_type", "authorization_code")
	
	req, err := http.NewRequestWithContext(ctx, "POST", conf.TokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create token request: %w", err)
	}
	
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("token request failed: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("token request returned %d: %s", resp.StatusCode, string(body))
	}
	
	var tokenResp struct {
		AccessToken string `json:"access_token"`
		IDToken     string `json:"id_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}
	
	// Use the access token to get user info
	userReq, err := http.NewRequestWithContext(ctx, "GET", conf.UserInfoURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create user info request: %w", err)
	}
	
	userReq.Header.Add("Authorization", "Bearer "+tokenResp.AccessToken)
	
	userResp, err := client.Do(userReq)
	if err != nil {
		return nil, fmt.Errorf("user info request failed: %w", err)
	}
	defer userResp.Body.Close()
	
	if userResp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(userResp.Body)
		return nil, fmt.Errorf("user info request returned %d: %s", userResp.StatusCode, string(body))
	}
	
	var userInfo struct {
		ID            string `json:"id"`
		Email         string `json:"email"`
		VerifiedEmail bool   `json:"verified_email"`
		Name          string `json:"name"`
		GivenName     string `json:"given_name"`
		FamilyName    string `json:"family_name"`
		Picture       string `json:"picture"`
		Locale        string `json:"locale"`
	}
	
	if err := json.NewDecoder(userResp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to decode user info response: %w", err)
	}
	
	// Return user info in a standardized format
	return map[string]string{
		"id":      userInfo.ID,
		"email":   userInfo.Email,
		"name":    userInfo.Name,
		"picture": userInfo.Picture,
	}, nil
}

// Implement the real OAuth exchange for GitHub
func (s *AuthService) exchangeGitHubCode(ctx context.Context, code string, conf models.OAuthConfig) (map[string]string, error) {
	logger.InfoLogger.Printf("Exchanging GitHub OAuth code for token")
	
	// Exchange authorization code for token
	data := url.Values{}
	data.Set("code", code)
	data.Set("client_id", conf.ClientID)
	data.Set("client_secret", conf.ClientSecret)
	data.Set("redirect_uri", conf.RedirectURI)
	
	req, err := http.NewRequestWithContext(ctx, "POST", conf.TokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create token request: %w", err)
	}
	
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/json")
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("token request failed: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("token request returned %d: %s", resp.StatusCode, string(body))
	}
	
	var tokenResp struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		Scope       string `json:"scope"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}
	
	// Use the access token to get user info
	userReq, err := http.NewRequestWithContext(ctx, "GET", conf.UserInfoURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create user info request: %w", err)
	}
	
	userReq.Header.Add("Authorization", "token "+tokenResp.AccessToken)
	userReq.Header.Add("Accept", "application/json")
	
	userResp, err := client.Do(userReq)
	if err != nil {
		return nil, fmt.Errorf("user info request failed: %w", err)
	}
	defer userResp.Body.Close()
	
	if userResp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(userResp.Body)
		return nil, fmt.Errorf("user info request returned %d: %s", userResp.StatusCode, string(body))
	}
	
	var userInfo struct {
		ID        int    `json:"id"`
		Login     string `json:"login"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		AvatarURL string `json:"avatar_url"`
	}
	
	if err := json.NewDecoder(userResp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to decode user info response: %w", err)
	}
	
	// If email is not public, make an additional request to get emails
	if userInfo.Email == "" {
		emailReq, err := http.NewRequestWithContext(ctx, "GET", "https://api.github.com/user/emails", nil)
		if err != nil {
			return nil, fmt.Errorf("failed to create email request: %w", err)
		}
		
		emailReq.Header.Add("Authorization", "token "+tokenResp.AccessToken)
		emailReq.Header.Add("Accept", "application/json")
		
		emailResp, err := client.Do(emailReq)
		if err != nil {
			return nil, fmt.Errorf("email request failed: %w", err)
		}
		defer emailResp.Body.Close()
		
		var emails []struct {
			Email    string `json:"email"`
			Primary  bool   `json:"primary"`
			Verified bool   `json:"verified"`
		}
		
		if err := json.NewDecoder(emailResp.Body).Decode(&emails); err != nil {
			return nil, fmt.Errorf("failed to decode email response: %w", err)
		}
		
		// Find primary email
		for _, email := range emails {
			if email.Primary && email.Verified {
				userInfo.Email = email.Email
				break
			}
		}
		
		// If no primary email, use the first verified one
		if userInfo.Email == "" && len(emails) > 0 {
			for _, email := range emails {
				if email.Verified {
					userInfo.Email = email.Email
					break
				}
			}
		}
	}
	
	// Return user info in a standardized format
	return map[string]string{
		"id":      fmt.Sprintf("%d", userInfo.ID),
		"email":   userInfo.Email,
		"name":    userInfo.Name,
		"picture": userInfo.AvatarURL,
	}, nil
}
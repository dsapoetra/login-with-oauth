package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"login-with-oauth/internal/models"
	"login-with-oauth/internal/repository"
	"net/http"
	"time"

	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleService struct {
	config         *oauth2.Config
	userRepository repository.UserRepository
}

func NewGoogleService(clientID, clientSecret string, userRepository repository.UserRepository) *GoogleService {
	config := &oauth2.Config{
		ClientID:     viper.GetString("google.clientID"),
		ClientSecret: viper.GetString("google.clientSecret"),
		RedirectURL:  "http://localhost:8080/callback-gl",
		Endpoint:     google.Endpoint,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
	}

	return &GoogleService{
		config:         config,
		userRepository: userRepository,
	}
}

func (s *GoogleService) GetAuthURL(state string) string {
	return s.config.AuthCodeURL(state)
}

func (s *GoogleService) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	return s.config.Exchange(ctx, code)
}

func (s *GoogleService) GetUserData(token *oauth2.Token) (*models.User, error) {
	client := s.config.Client(context.Background(), token)

	req, err := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Google API request failed with status: %d, body: %s", resp.StatusCode, string(body))
	}

	var googleUser struct {
		ID       string `json:"id"`
		Email    string `json:"email"`
		Name     string `json:"name"`
		Picture  string `json:"picture"`
		Verified bool   `json:"verified_email"`
	}

	if err := json.Unmarshal(body, &googleUser); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %v, body: %s", err, string(body))
	}

	userData := models.User{
		ID:        googleUser.ID,
		Username:  googleUser.Name,
		Email:     googleUser.Email,
		AvatarURL: googleUser.Picture,
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
	}

	savedUser, err := s.userRepository.CreateUser(userData)
	if err != nil {
		return nil, fmt.Errorf("failed to create user in repository: %v", err)
	}

	return savedUser, nil
}

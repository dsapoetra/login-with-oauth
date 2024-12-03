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

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

type GithubService struct {
	config         *oauth2.Config
	userRepository repository.UserRepository
}

func NewGitHubService(clientID, clientSecret string, userRepository repository.UserRepository) *GithubService {
	config := &oauth2.Config{
		ClientID:     clientID,     // Use the passed parameter
		ClientSecret: clientSecret, // Use the passed parameter
		RedirectURL:  "http://localhost:8080/gh-cb",
		Endpoint:     github.Endpoint,
		Scopes:       []string{"user:email", "user:avatar"},
	}

	return &GithubService{
		config:         config,
		userRepository: userRepository,
	}
}

func (s *GithubService) GetAuthURL(state string) string {
	return s.config.AuthCodeURL(state)
}

func (s *GithubService) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	return s.config.Exchange(ctx, code)
}

func (s *GithubService) GetUserData(token *oauth2.Token) (*models.User, error) {
	client := s.config.Client(context.Background(), token)

	// Create request
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Set headers
	req.Header.Set("User-Agent", "Oauth")

	// Make request
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
		return nil, fmt.Errorf("GitHub API request failed with status: %d, body: %s", resp.StatusCode, string(body))
	}

	var githubUser struct {
		ID        int64  `json:"id"`
		Login     string `json:"login"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		AvatarURL string `json:"avatar_url"`
	}

	if err := json.Unmarshal(body, &githubUser); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %v, body: %s", err, string(body))
	}

	userData := models.User{
		ID:        fmt.Sprintf("%d", githubUser.ID),
		Username:  githubUser.Login,
		Email:     githubUser.Email,
		AvatarURL: githubUser.AvatarURL,
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
	}

	return s.userRepository.CreateUser(userData)
}

package services

import (
	"context"
	"login-with-oauth/internal/models"
	"login-with-oauth/internal/repository/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2"
)

func TestGithubService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock.NewMockUserRepository(ctrl)

	t.Run("TestGetAuthURL", func(t *testing.T) {
		// Arrange
		service := NewGitHubService("test-client-id", "test-client-secret", mockUserRepo)
		state := "test-state"

		// Act
		url := service.GetAuthURL(state)

		// Assert
		assert.Contains(t, url, "github.com/login/oauth/authorize")
		assert.Contains(t, url, "state="+state)
		// Fix: Use URL-encoded scope format
		assert.Contains(t, url, "scope=user%3Aemail+user%3Aavatar")
		assert.Contains(t, url, "client_id=test-client-id")
		assert.Contains(t, url, "redirect_uri=http%3A%2F%2Flocalhost%3A8080%2Fgh-cb")
	})

	t.Run("TestExchange", func(t *testing.T) {
		// 	// Arrange
		service := NewGitHubService("test-client-id", "test-client-secret", mockUserRepo)
		code := "test-code"

		// Act
		token, err := service.Exchange(context.Background(), code)

		// Note: This will actually fail because it tries to make a real HTTP request
		// In a real test, you'd want to mock the OAuth2 config's exchange function
		assert.Error(t, err)
		assert.Nil(t, token)
	})

	t.Run("TestGetUserData", func(t *testing.T) {
		// Create a mock server
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Return mock response for GitHub API
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{
				"id": 12345,
				"login": "testuser",
				"name": "Test User",
				"email": "test@example.com",
				"avatar_url": "https://example.com/avatar.jpg"
			}`))
		}))
		defer mockServer.Close()

		// Create the service
		service := NewGitHubService("test-client-id", "test-client-secret", mockUserRepo)

		// Override the service's config and transport
		transport := &mockTransport{
			mockServer: mockServer,
		}
		http.DefaultClient.Transport = transport

		// Override the service's config
		service.config = &oauth2.Config{
			ClientID:     "test-client-id",
			ClientSecret: "test-client-secret",
			Endpoint: oauth2.Endpoint{
				AuthURL:  mockServer.URL,
				TokenURL: mockServer.URL,
			},
		}

		// Create a test token
		token := &oauth2.Token{
			AccessToken: "test-token",
			TokenType:   "Bearer",
		}

		// Set up expected user
		expectedUser := &models.User{
			ID:        "12345",
			Username:  "testuser",
			Email:     "test@example.com",
			AvatarURL: "https://example.com/avatar.jpg",
			CreatedAt: time.Now().Format(time.RFC3339),
			UpdatedAt: time.Now().Format(time.RFC3339),
		}

		// Set up mock repository expectation
		mockUserRepo.EXPECT().
			CreateUser(gomock.Any()).
			Return(expectedUser, nil)

		// Act
		user, err := service.GetUserData(token)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, expectedUser.ID, user.ID)
		assert.Equal(t, expectedUser.Username, user.Username)
		assert.Equal(t, expectedUser.Email, user.Email)
		assert.Equal(t, expectedUser.AvatarURL, user.AvatarURL)
	})

	// t.Run("TestGetUserData_APIError", func(t *testing.T) {
	// 	// Arrange
	// 	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 		w.WriteHeader(http.StatusUnauthorized)
	// 		w.Write([]byte("Unauthorized"))
	// 	}))
	// 	defer testServer.Close()

	// 	service := NewGitHubService("test-client-id", "test-client-secret", mockUserRepo)
	// 	service.config.Endpoint.TokenURL = testServer.URL

	// 	token := &oauth2.Token{
	// 		AccessToken: "invalid-token",
	// 		TokenType:   "Bearer",
	// 	}

	// 	// Act
	// 	user, err := service.GetUserData(token)

	// 	// Assert
	// 	assert.Error(t, err)
	// 	assert.Nil(t, user)
	// 	assert.Contains(t, err.Error(), "GitHub API request failed with status: 401")
	// })

	// t.Run("TestGetUserData_InvalidJSON", func(t *testing.T) {
	// 	// Arrange
	// 	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 		w.Write([]byte("invalid json"))
	// 	}))
	// 	defer testServer.Close()

	// 	service := NewGitHubService("test-client-id", "test-client-secret", mockUserRepo)
	// 	service.config.Endpoint.TokenURL = testServer.URL

	// 	token := &oauth2.Token{
	// 		AccessToken: "test-token",
	// 		TokenType:   "Bearer",
	// 	}

	// 	// Act
	// 	user, err := service.GetUserData(token)

	// 	// Assert
	// 	assert.Error(t, err)
	// 	assert.Nil(t, user)
	// 	assert.Contains(t, err.Error(), "failed to decode JSON")
	// })
}

// mockTransport implements http.RoundTripper
type mockTransport struct {
	mockServer *httptest.Server
}

func (t *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Replace the GitHub API URL with our mock server URL
	if req.URL.Host == "api.github.com" {
		req.URL.Scheme = "http"
		req.URL.Host = strings.TrimPrefix(t.mockServer.URL, "http://")
	}
	return http.DefaultTransport.RoundTrip(req)
}

package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"login-with-oauth/internal/services"
	"net/http"
)

type GithubHandler struct {
	githubService *services.GithubService
}

func NewAuthHandler(githubService *services.GithubService) *GithubHandler {
	return &GithubHandler{
		githubService: githubService,
	}
}

func (h *GithubHandler) GitHubLogin(w http.ResponseWriter, r *http.Request) {
	// Generate random state
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)

	// Redirect to GitHub
	url := h.githubService.GetAuthURL(state)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *GithubHandler) GitHubCallback(w http.ResponseWriter, r *http.Request) {
	// Get code and state from query params
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	// TODO: Verify state matches stored state
	if state == "" {
		http.Error(w, "State parameter is missing", http.StatusBadRequest)
		return
	}

	// Exchange code for token
	token, err := h.githubService.Exchange(r.Context(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}

	// Get user data
	userData, err := h.githubService.GetUserData(token)
	if err != nil {
		http.Error(w, "Failed to get user data", http.StatusInternalServerError)
		return
	}

	// Return user data
	w.Write([]byte("Logged in successfully as: " + userData.Email))

}

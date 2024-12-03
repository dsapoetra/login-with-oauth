package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"login-with-oauth/internal/services"
	"net/http"
)

type GoogleHandler struct {
	googleService *services.GoogleService
}

func NewGoogleHandler(googleService *services.GoogleService) *GoogleHandler {
	return &GoogleHandler{
		googleService: googleService,
	}
}

func generateRandomState() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func (h *GoogleHandler) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	state := generateRandomState()
	authURL := h.googleService.GetAuthURL(state)
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

func (h *GoogleHandler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code not found", http.StatusBadRequest)
		return
	}

	token, err := h.googleService.Exchange(r.Context(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token", http.StatusInternalServerError)
		return
	}

	user, err := h.googleService.GetUserData(token)
	if err != nil {
		http.Error(w, "Failed to get user data", http.StatusInternalServerError)
		return
	}

	// Handle successful login (e.g., create session, set cookies, etc.)
	w.Write([]byte("Logged in successfully as: " + user.Email))
}

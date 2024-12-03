package services

import (
	mocks "login-with-oauth/internal/repository/mock"
	"testing"

	"github.com/golang/mock/gomock" // Changed from go.uber.org/mock/gomock
	"github.com/stretchr/testify/assert"
)

func TestNewGoogleService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	service := NewGoogleService("test-client-id", "test-client-secret", mockRepo)

	assert.NotNil(t, service)
	assert.NotNil(t, service.config)
	assert.Equal(t, mockRepo, service.userRepository)
}

func TestGetAuthURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockUserRepository(ctrl)
	service := NewGoogleService("test-client-id", "test-client-secret", mockRepo)

	url := service.GetAuthURL("test-state")
	assert.Contains(t, url, "accounts.google.com/o/oauth2/auth")
	assert.Contains(t, url, "test-state")
}

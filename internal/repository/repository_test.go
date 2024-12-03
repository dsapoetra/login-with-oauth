package repository_test

import (
	"errors"
	"login-with-oauth/internal/models"
	"login-with-oauth/internal/repository/mock"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock repository
	mockRepo := mock.NewMockUserRepository(ctrl)

	testUser := models.User{
		ID:        "123",
		Username:  "testuser",
		Email:     "test@example.com",
		AvatarURL: "https://example.com/avatar.jpg",
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
	}

	// Set expectations
	mockRepo.EXPECT().
		CreateUser(gomock.Any()).
		Return(&testUser, nil).
		Times(1)

	// Test the mock
	savedUser, err := mockRepo.CreateUser(testUser)
	assert.NoError(t, err)
	assert.Equal(t, testUser.ID, savedUser.ID)
}

// Test with error case
func TestCreateUser_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockUserRepository(ctrl)

	testUser := models.User{
		ID:        "123",
		Username:  "testuser",
		Email:     "test@example.com",
		AvatarURL: "https://example.com/avatar.jpg",
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
	}

	// Expect error
	mockRepo.EXPECT().
		CreateUser(gomock.Any()).
		Return(nil, errors.New("database error")).
		Times(1)

	savedUser, err := mockRepo.CreateUser(testUser)
	assert.Error(t, err)
	assert.Nil(t, savedUser)
}

package mocks

import (
	"context"

	"github.com/slodkiadrianek/octopus/internal/DTO"
	"github.com/slodkiadrianek/octopus/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) InsertUserToDB(ctx context.Context, user DTO.CreateUser, password string) error {
	args := m.Called(ctx, user, password)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateUser(ctx context.Context, user DTO.CreateUser, userID int) error {
	args := m.Called(ctx, user, userID)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateUserNotifications(ctx context.Context, userID int, userNotifications DTO.UpdateUserNotificationsSettings) error {
	args := m.Called(ctx, userNotifications, userID)
	return args.Error(0)
}

func (m *MockUserRepository) DeleteUser(ctx context.Context, password string, userID int) error {
	args := m.Called(ctx, userID, password)
	return args.Error(0)
}

func (m *MockUserRepository) FindUserByID(ctx context.Context, userID int) (models.User, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *MockUserRepository) FindUserByEmail(ctx context.Context, email string) (models.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *MockUserRepository) ChangeUserPassword(ctx context.Context, userID int, newPassword string) error {
	args := m.Called(ctx, userID, newPassword)
	return args.Error(0)
}

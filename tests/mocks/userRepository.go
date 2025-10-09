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

func (m *MockUserRepository) InsertUserToDb(ctx context.Context, user DTO.CreateUser, password string) error {
	args := m.Called(ctx, user, password)
	return args.Error(0)
}
func (m *MockUserRepository) UpdateUser(ctx context.Context, user DTO.CreateUser, userId int) error {
	args := m.Called(ctx, user, userId)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateUserNotifications(ctx context.Context, userId int, userNotifications DTO.UpdateUserNotifications) error {
	args := m.Called(ctx, userNotifications, userId)
	return args.Error(0)
}
func (m *MockUserRepository) DeleteUser(ctx context.Context, password string, userId int) error {
	args := m.Called(ctx, userId, password)
	return args.Error(0)
}
func (m *MockUserRepository) FindUserById(ctx context.Context, userId int) (models.User, error) {
	args := m.Called(ctx, userId)
	return args.Get(0).(models.User), args.Error(1)
}
func (m *MockUserRepository) FindUserByEmail(ctx context.Context, email string) (models.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *MockUserRepository) ChangeUserPassword(ctx context.Context, userId int, newPassword string) error {
	args := m.Called(ctx, userId, newPassword)
	return args.Error(0)
}

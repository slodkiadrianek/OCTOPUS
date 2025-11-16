package user

import (
	"context"
	"errors"
	"testing"

	"github.com/slodkiadrianek/octopus/internal/DTO"
	"github.com/slodkiadrianek/octopus/internal/models"
	"github.com/slodkiadrianek/octopus/internal/services/interfaces"
	"github.com/slodkiadrianek/octopus/tests"
	"github.com/slodkiadrianek/octopus/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserService_InsertUserToDb(t *testing.T) {
	type args struct {
		name          string
		expectedError *string
		password      string
		setupMock     func() interfaces.UserRepository
	}
	testsScenarios := []args{
		{
			name:          "Proper data",
			expectedError: nil,
			password:      "fdEW4$#f4303",
			setupMock: func() interfaces.UserRepository {
				m := new(mocks.MockUserRepository)
				m.On("FindUserByEmail", mock.Anything, mock.Anything).Return(
					models.User{}, nil)
				m.On("InsertUserToDb", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				return m
			},
		},
		{
			name:          "User not found",
			expectedError: tests.Ptr("User not found"),
			password:      "fdEW4$#f4303",
			setupMock: func() interfaces.UserRepository {
				m := new(mocks.MockUserRepository)
				m.On("FindUserByEmail", mock.Anything, mock.Anything).Return(
					models.User{}, errors.New("user not found"))
				return m
			},
		},
		{
			name:          "Failed to find user by email",
			expectedError: tests.Ptr("Failed to execute query"),
			password:      "fdEW4$#f4303",
			setupMock: func() interfaces.UserRepository {
				m := new(mocks.MockUserRepository)
				m.On("FindUserByEmail", mock.Anything, mock.Anything).Return(
					models.User{}, errors.New("failed to execute query"))
				return m
			},
		},
		{
			name:          "User already exists",
			expectedError: tests.Ptr("User with this email already exists"),
			password:      "fdEW4$#f4303",
			setupMock: func() interfaces.UserRepository {
				m := new(mocks.MockUserRepository)
				m.On("FindUserByEmail", mock.Anything, mock.Anything).Return(
					models.User{
						ID: 1,
					}, nil)
				return m
			},
		},
		{
			name:          "Failed to hash password",
			expectedError: tests.Ptr("bcrypt: password length exceeds 72 bytes"),
			password:      "fdEW4$#f4303r3er236575467nfw7f9348htx0f94378xfh349fxyh349xf8@#34RFDFE42423423423",
			setupMock: func() interfaces.UserRepository {
				m := new(mocks.MockUserRepository)
				m.On("FindUserByEmail", mock.Anything, mock.Anything).Return(
					models.User{
						ID: 0,
					}, nil)
				return m
			},
		},
		{
			name:          "Failed to insert user to db",
			expectedError: tests.Ptr("Failed to insert user to db"),
			password:      "fdEW4$#f4303",
			setupMock: func() interfaces.UserRepository {
				m := new(mocks.MockUserRepository)
				m.On("FindUserByEmail", mock.Anything, mock.Anything).Return(
					models.User{
						ID: 0,
					}, nil)

				m.On("InsertUserToDb", mock.Anything, mock.Anything,
					mock.Anything).Return(errors.New("failed to insert user to db"))
				return m
			},
		},
	}

	for _, testScenario := range testsScenarios {
		t.Run(testScenario.name, func(t *testing.T) {
			ctx := context.Background()
			loggerService := tests.CreateLogger()
			userRepository := testScenario.setupMock()
			cacheService := tests.CreateCacheService(loggerService)
			userService := NewUserService(loggerService, userRepository, cacheService)
			user := DTO.NewCreateUser("adikurek@cos.com", "Adrian", "Kurek")
			err := userService.InsertUserToDb(ctx, *user, testScenario.password)
			if testScenario.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), *testScenario.expectedError)
			}
		})
	}
}

func TestUserService_UpdateUser(t *testing.T) {
	type args struct {
		name          string
		expectedError *string
		setupMock     func() interfaces.UserRepository
	}
	testsScenarios := []args{
		{
			name:          "Proper data",
			expectedError: nil,
			setupMock: func() interfaces.UserRepository {
				m := new(mocks.MockUserRepository)
				m.On("UpdateUser", mock.Anything, mock.Anything, mock.Anything).Return(
					nil)
				return m
			},
		},
		{
			name:          "Failed to update an user",
			expectedError: tests.Ptr("Failed to update an user"),
			setupMock: func() interfaces.UserRepository {
				m := new(mocks.MockUserRepository)
				m.On("UpdateUser", mock.Anything, mock.Anything, mock.Anything).Return(
					errors.New("failed to update an user"))
				return m
			},
		},
	}
	for _, testScenario := range testsScenarios {
		t.Run(testScenario.name, func(t *testing.T) {
			ctx := context.Background()
			loggerService := tests.CreateLogger()
			userRepository := testScenario.setupMock()
			cacheService := tests.CreateCacheService(loggerService)
			userService := NewUserService(loggerService, userRepository, cacheService)
			user := DTO.NewCreateUser("adikurek@cos.com", "Adrina", "Kurek")
			err := userService.UpdateUser(ctx, *user, 3)
			if testScenario.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), *testScenario.expectedError)
			}
		})
	}
}

func TestUserService_UpdateUserNotifications(t *testing.T) {
	type args struct {
		name          string
		expectedError *string
		setupMock     func() interfaces.UserRepository
	}
	testsScenarios := []args{
		{
			name:          "Proper data",
			expectedError: nil,
			setupMock: func() interfaces.UserRepository {
				m := new(mocks.MockUserRepository)
				m.On("UpdateUserNotifications", mock.Anything, mock.Anything, mock.Anything).Return(
					nil)
				return m
			},
		},
		{
			name:          "Failed to update an user",
			expectedError: tests.Ptr("Failed to update an user"),
			setupMock: func() interfaces.UserRepository {
				m := new(mocks.MockUserRepository)
				m.On("UpdateUserNotifications", mock.Anything, mock.Anything, mock.Anything).Return(
					errors.New("failed to update an user"))
				return m
			},
		},
	}
	for _, testScenario := range testsScenarios {
		t.Run(testScenario.name, func(t *testing.T) {
			ctx := context.Background()
			loggerService := tests.CreateLogger()
			userRepository := testScenario.setupMock()
			cacheService := tests.CreateCacheService(loggerService)

			userService := NewUserService(loggerService, userRepository, cacheService)
			user := DTO.UpdateUserNotificationsSettings{
				DiscordNotificationsSettings: false,
				SlackNotificationsSettings:   false,
				EmailNotificationsSettings:   false,
			}
			err := userService.UpdateUserNotifications(ctx, 3, user)
			if testScenario.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), *testScenario.expectedError)
			}
		})
	}
}

func TestUserService_DeleteUser(t *testing.T) {
	type args struct {
		name          string
		expectedError *string
		password      string
		setupMock     func() interfaces.UserRepository
	}
	testsScenarios := []args{
		{
			name:          "Proper data",
			expectedError: nil,
			password:      "ci$#fm43980faz",
			setupMock: func() interfaces.UserRepository {
				m := new(mocks.MockUserRepository)
				m.On("FindUserById", mock.Anything, mock.Anything).Return(
					models.User{ID: 1, Password: "$2a$10$0f4BED0dDgYCE8xVREwhUeyjpKTtBIm4eO.xrPC/H8kvsBVM2gpdq"}, nil)
				m.On("DeleteUser", mock.Anything, mock.Anything, mock.Anything).Return(
					nil)
				return m
			},
		},
		{
			name:          "User does not exist",
			expectedError: tests.Ptr("User with this id does not exist"),
			password:      "ci$#fm43980faz",
			setupMock: func() interfaces.UserRepository {
				m := new(mocks.MockUserRepository)
				m.On("FindUserById", mock.Anything, mock.Anything).Return(
					models.User{ID: 0}, nil)
				return m
			},
		},
		{
			name:          "Wrong password provided",
			expectedError: tests.Ptr("Wrong password provided"),
			password:      "ci$#fm43980faz2",
			setupMock: func() interfaces.UserRepository {
				m := new(mocks.MockUserRepository)
				m.On("FindUserById", mock.Anything, mock.Anything).Return(
					models.User{ID: 1, Password: "$2a$10$0f4BED0dDgYCE8xVREwhUeyjpKTtBIm4eO.xrPC/H8kvsBVM2gpdq"}, nil)
				return m
			},
		},
		{
			name:          "Failed to find a user",
			expectedError: tests.Ptr("Failed to find a user"),
			password:      "ci$#fm43980faz2",
			setupMock: func() interfaces.UserRepository {
				m := new(mocks.MockUserRepository)
				m.On("FindUserById", mock.Anything, mock.Anything).Return(models.User{},
					errors.New("failed to find a user"))
				return m
			},
		},
		{
			name:          "Failed to delete a user",
			expectedError: tests.Ptr("Failed to delete a user"),
			password:      "ci$#fm43980faz",
			setupMock: func() interfaces.UserRepository {
				m := new(mocks.MockUserRepository)
				m.On("FindUserById", mock.Anything, mock.Anything).Return(
					models.User{ID: 1, Password: "$2a$10$0f4BED0dDgYCE8xVREwhUeyjpKTtBIm4eO.xrPC/H8kvsBVM2gpdq"}, nil)
				m.On("DeleteUser", mock.Anything, mock.Anything, mock.Anything).Return(
					errors.New("failed to delete a user"))
				return m
			},
		},
	}
	for _, testScenario := range testsScenarios {
		t.Run(testScenario.name, func(t *testing.T) {
			ctx := context.Background()
			loggerService := tests.CreateLogger()
			userRepository := testScenario.setupMock()
			cacheService := tests.CreateCacheService(loggerService)

			userService := NewUserService(loggerService, userRepository, cacheService)
			err := userService.DeleteUser(ctx, 3, testScenario.password)
			if testScenario.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), *testScenario.expectedError)
			}
		})
	}
}

func TestUserService_ChangeUserPassword(t *testing.T) {
	type args struct {
		name          string
		expectedError *string
		password      string
		newPassword   string
		setupMock     func() interfaces.UserRepository
	}
	testsScenarios := []args{
		{
			name:          "Proper data",
			expectedError: nil,
			password:      "ci$#fm43980faz",
			newPassword:   "ci$#fm43980faz",
			setupMock: func() interfaces.UserRepository {
				m := new(mocks.MockUserRepository)
				m.On("FindUserById", mock.Anything, mock.Anything).Return(
					models.User{ID: 1, Password: "$2a$10$0f4BED0dDgYCE8xVREwhUeyjpKTtBIm4eO.xrPC/H8kvsBVM2gpdq"}, nil)
				m.On("ChangeUserPassword", mock.Anything, mock.Anything, mock.Anything).Return(
					nil)
				return m
			},
		},
		{
			name:          "User with this id does not exist",
			expectedError: tests.Ptr("User with this id does not exist"),
			password:      "ci$#fm43980faz",
			newPassword:   "ci$#fm43980faz",
			setupMock: func() interfaces.UserRepository {
				m := new(mocks.MockUserRepository)
				m.On("FindUserById", mock.Anything, mock.Anything).Return(
					models.User{ID: 0}, nil)
				return m
			},
		},
		{
			name:          "Wrong password provided",
			expectedError: tests.Ptr("Wrong current password provided"),
			password:      "ci$#fm43980faz2",
			newPassword:   "ci$#fm43980faz",
			setupMock: func() interfaces.UserRepository {
				m := new(mocks.MockUserRepository)
				m.On("FindUserById", mock.Anything, mock.Anything).Return(
					models.User{ID: 1, Password: "$2a$10$0f4BED0dDgYCE8xVREwhUeyjpKTtBIm4eO.xrPC/H8kvsBVM2gpdq"}, nil)
				return m
			},
		},
		{
			name:          "Wrong password provid3ed",
			expectedError: tests.Ptr("bcrypt: password length exceeds 72 bytes"),
			password:      "ci$#fm43980faz",
			newPassword:   "ci$#fm4398d432-89m52348-$#@rt43t43tZ#4t43t39ty4343324fn87634c8-t734tct43t430faz",
			setupMock: func() interfaces.UserRepository {
				m := new(mocks.MockUserRepository)
				m.On("FindUserById", mock.Anything, mock.Anything).Return(
					models.User{ID: 1, Password: "$2a$10$0f4BED0dDgYCE8xVREwhUeyjpKTtBIm4eO.xrPC/H8kvsBVM2gpdq"}, nil)
				return m
			},
		},
		{
			name:          "Failed to find a user",
			expectedError: tests.Ptr("Failed to find a user"),
			password:      "ci$#fm43980faz2",
			newPassword:   "ci$#fm43980faz",
			setupMock: func() interfaces.UserRepository {
				m := new(mocks.MockUserRepository)
				m.On("FindUserById", mock.Anything, mock.Anything).Return(models.User{},
					errors.New("failed to find a user"))
				return m
			},
		},
		{
			name:          "Failed to change password",
			expectedError: tests.Ptr("Failed to change password"),
			password:      "ci$#fm43980faz",
			newPassword:   "ci$#fm43980faz",
			setupMock: func() interfaces.UserRepository {
				m := new(mocks.MockUserRepository)
				m.On("FindUserById", mock.Anything, mock.Anything).Return(
					models.User{ID: 1, Password: "$2a$10$0f4BED0dDgYCE8xVREwhUeyjpKTtBIm4eO.xrPC/H8kvsBVM2gpdq"}, nil)
				m.On("ChangeUserPassword", mock.Anything, mock.Anything, mock.Anything).Return(
					errors.New("failed to change password"))
				return m
			},
		},
	}
	for _, testScenario := range testsScenarios {
		t.Run(testScenario.name, func(t *testing.T) {
			ctx := context.Background()
			loggerService := tests.CreateLogger()
			userRepository := testScenario.setupMock()
			cacheService := tests.CreateCacheService(loggerService)

			userService := NewUserService(loggerService, userRepository, cacheService)
			err := userService.ChangeUserPassword(ctx, 3, testScenario.password, testScenario.newPassword)
			if testScenario.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), *testScenario.expectedError)
			}
		})
	}
}

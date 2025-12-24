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
		expectedError error
		password      string
		setupMock     func() interfaces.UserRepository
	}
	testsScenarios := []args{
		{
			name:          "Proper data",
			expectedError: nil,
			password:      "fdEW4$#f4303",
			setupMock: func() interfaces.UserRepository {
				mUserRepository := new(mocks.MockUserRepository)
				mUserRepository.On("FindUserByEmail", mock.Anything, mock.Anything).Return(
					models.User{}, nil)
				mUserRepository.On("InsertUserToDb", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				return mUserRepository
			},
		},
		{
			name:          "User not found",
			expectedError: errors.New("user not found"),
			password:      "fdEW4$#f4303",
			setupMock: func() interfaces.UserRepository {
				mUserRepository := new(mocks.MockUserRepository)
				mUserRepository.On("FindUserByEmail", mock.Anything, mock.Anything).Return(
					models.User{}, errors.New("user not found"))
				return mUserRepository
			},
		},
		{
			name:          "Failed to find user by email",
			expectedError: errors.New("failed to execute query"),
			password:      "fdEW4$#f4303",
			setupMock: func() interfaces.UserRepository {
				mUserRepository := new(mocks.MockUserRepository)
				mUserRepository.On("FindUserByEmail", mock.Anything, mock.Anything).Return(
					models.User{}, errors.New("failed to execute query"))
				return mUserRepository
			},
		},
		{
			name:          "User already exists",
			expectedError: errors.New("user with this email already exists"),
			password:      "fdEW4$#f4303",
			setupMock: func() interfaces.UserRepository {
				mUserRepository := new(mocks.MockUserRepository)
				mUserRepository.On("FindUserByEmail", mock.Anything, mock.Anything).Return(
					models.User{
						ID: 1,
					}, nil)
				return mUserRepository
			},
		},
		{
			name:          "Failed to hash password",
			expectedError: errors.New("bcrypt: password length exceeds 72 bytes"),
			password:      "fdEW4$#f4303r3er236575467nfw7f9348htx0f94378xfh349fxyh349xf8@#34RFDFE42423423423",
			setupMock: func() interfaces.UserRepository {
				mUserRepository := new(mocks.MockUserRepository)
				mUserRepository.On("FindUserByEmail", mock.Anything, mock.Anything).Return(
					models.User{
						ID: 0,
					}, nil)
				return mUserRepository
			},
		},
		{
			name:          "Failed to insert user to db",
			expectedError: errors.New("failed to insert user to db"),
			password:      "fdEW4$#f4303",
			setupMock: func() interfaces.UserRepository {
				mUserRepository := new(mocks.MockUserRepository)
				mUserRepository.On("FindUserByEmail", mock.Anything, mock.Anything).Return(
					models.User{
						ID: 0,
					}, nil)

				mUserRepository.On("InsertUserToDb", mock.Anything, mock.Anything,
					mock.Anything).Return(errors.New("failed to insert user to db"))
				return mUserRepository
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
			err := userService.InsertUserToDB(ctx, *user, testScenario.password)
			if testScenario.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), testScenario.expectedError.Error())
			}
		})
	}
}

func TestUserService_UpdateUser(t *testing.T) {
	type args struct {
		name          string
		expectedError error
		setupMock     func() interfaces.UserRepository
	}
	testsScenarios := []args{
		{
			name:          "Proper data",
			expectedError: nil,
			setupMock: func() interfaces.UserRepository {
				mUserRepository := new(mocks.MockUserRepository)
				mUserRepository.On("UpdateUser", mock.Anything, mock.Anything, mock.Anything).Return(
					nil)
				return mUserRepository
			},
		},
		{
			name:          "Failed to update an user",
			expectedError: errors.New("failed to update an user"),
			setupMock: func() interfaces.UserRepository {
				mUserRepository := new(mocks.MockUserRepository)
				mUserRepository.On("UpdateUser", mock.Anything, mock.Anything, mock.Anything).Return(
					errors.New("failed to update an user"))
				return mUserRepository
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
				assert.Contains(t, err.Error(), testScenario.expectedError.Error())
			}
		})
	}
}

func TestUserService_UpdateUserNotifications(t *testing.T) {
	type args struct {
		name          string
		expectedError error
		setupMock     func() interfaces.UserRepository
	}
	testsScenarios := []args{
		{
			name:          "Proper data",
			expectedError: nil,
			setupMock: func() interfaces.UserRepository {
				mUserRepository := new(mocks.MockUserRepository)
				mUserRepository.On("UpdateUserNotifications", mock.Anything, mock.Anything, mock.Anything).Return(
					nil)
				return mUserRepository
			},
		},
		{
			name:          "Failed to update an user",
			expectedError: errors.New("failed to update an user"),
			setupMock: func() interfaces.UserRepository {
				mUserRepository := new(mocks.MockUserRepository)
				mUserRepository.On("UpdateUserNotifications", mock.Anything, mock.Anything, mock.Anything).Return(
					errors.New("failed to update an user"))
				return mUserRepository
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
				assert.Contains(t, err.Error(), testScenario.expectedError.Error())
			}
		})
	}
}

func TestUserService_GetUser(t *testing.T) {
	type args struct {
		name          string
		expectedError error
		setupMock     func() (interfaces.UserRepository, interfaces.CacheService)
	}
	testsScenarios := []args{
		{
			name:          "Proper data without saved user in cache",
			expectedError: nil,
			setupMock: func() (interfaces.UserRepository, interfaces.CacheService) {
				mUserRepository := new(mocks.MockUserRepository)
				mUserRepository.On("FindUserById", mock.Anything, mock.Anything).Return(
					models.User{ID: 1, Password: "$2a$10$0f4BED0dDgYCE8xVREwhUeyjpKTtBIm4eO.xrPC/H8kvsBVM2gpdq"}, nil)
				mCacheService := new(mocks.MockCacheService)
				mCacheService.On("ExistsData", mock.Anything, mock.Anything).Return(int64(0), nil)
				mCacheService.On("SetData", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
				return mUserRepository, mCacheService
			},
		},
		{
			name:          "Proper data with saved user in cache",
			expectedError: nil,
			setupMock: func() (interfaces.UserRepository, interfaces.CacheService) {
				mUserRepository := new(mocks.MockUserRepository)
				mUserRepository.On("FindUserById", mock.Anything, mock.Anything).Return(
					models.User{ID: 1, Password: "$2a$10$0f4BED0dDgYCE8xVREwhUeyjpKTtBIm4eO.xrPC/H8kvsBVM2gpdq"}, nil)
				mCacheService := new(mocks.MockCacheService)
				mCacheService.On("ExistsData", mock.Anything, mock.Anything).Return(int64(1), nil)
				mCacheService.On("GetData", mock.Anything, mock.Anything).Return(`{
				  "id": 1,
				  "email": "joedoe@email.com",
				  "name": "Joe",
				  "surname": "Doe",
				  "password": "$2a$10$0f4BED0dDgYCE8xVREwhUeyjpKTtBIm4eO.xrPC/H8kvsBVM2gpdq",
				  "discord_notifications": false,
				  "email_notifications_settings": true,
				  "slack_notifications_settings": false,
				  "created_at": "2023-01-01T00:00:00Z",
				  "updated_at": "2023-01-01T00:00:00Z"
				}`, nil)
				mCacheService.On("SetData", mock.Anything, mock.Anything, mock.Anything,
					mock.Anything).Return(nil)
				return mUserRepository, mCacheService
			},
		},
		{
			name:          "Failed to read user from cache",
			expectedError: errors.New("failed to read data from cache"),
			setupMock: func() (interfaces.UserRepository, interfaces.CacheService) {
				mUserRepository := new(mocks.MockUserRepository)
				mUserRepository.On("FindUserById", mock.Anything, mock.Anything).Return(
					models.User{ID: 1, Password: "$2a$10$0f4BED0dDgYCE8xVREwhUeyjpKTtBIm4eO.xrPC/H8kvsBVM2gpdq"}, nil)
				mCacheService := new(mocks.MockCacheService)
				mCacheService.On("ExistsData", mock.Anything, mock.Anything).Return(int64(1), nil)
				mCacheService.On("GetData", mock.Anything, mock.Anything).Return(``, errors.New("failed to read data from cache"))
				mCacheService.On("SetData", mock.Anything, mock.Anything, mock.Anything,
					mock.Anything).Return(nil)
				return mUserRepository, mCacheService
			},
		},
		{
			name:          "ExistsData failed",
			expectedError: errors.New("failed to check exists data"),
			setupMock: func() (interfaces.UserRepository, interfaces.CacheService) {
				mUserRepository := new(mocks.MockUserRepository)
				mCacheService := new(mocks.MockCacheService)
				mCacheService.On("ExistsData", mock.Anything, mock.Anything).Return(int64(0), errors.New("failed to check exists data"))
				return mUserRepository, mCacheService
			},
		},
		{
			name:          "Failed to set data in cache",
			expectedError: errors.New("failed to set data in cache"),
			setupMock: func() (interfaces.UserRepository, interfaces.CacheService) {
				mUserRepository := new(mocks.MockUserRepository)
				mUserRepository.On("FindUserById", mock.Anything, mock.Anything).Return(
					models.User{ID: 1, Password: "$2a$10$0f4BED0dDgYCE8xVREwhUeyjpKTtBIm4eO.xrPC/H8kvsBVM2gpdq"}, nil)
				mCacheService := new(mocks.MockCacheService)
				mCacheService.On("ExistsData", mock.Anything, mock.Anything).Return(int64(0), nil)
				mCacheService.On("SetData", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("failed to set data in cache"))
				return mUserRepository, mCacheService
			},
		},
	}
	for _, testScenario := range testsScenarios {
		t.Run(testScenario.name, func(t *testing.T) {
			ctx := context.Background()
			loggerService := tests.CreateLogger()
			userRepository, cacheService := testScenario.setupMock()

			userService := NewUserService(loggerService, userRepository, cacheService)
			user, err := userService.GetUser(ctx, 3)
			if testScenario.expectedError == nil {
				assert.NoError(t, err)
				assert.NotEmpty(t, user)
			} else {
				assert.Error(t, err)
				assert.Empty(t, user)
				assert.Contains(t, err.Error(), testScenario.expectedError.Error())
			}
		})
	}
}

func TestUserService_DeleteUser(t *testing.T) {
	type args struct {
		name          string
		expectedError error
		password      string
		setupMock     func() (interfaces.UserRepository, interfaces.CacheService)
	}
	testsScenarios := []args{
		{
			name:          "Proper data without user data saved in the cache",
			expectedError: nil,
			password:      "ci$#fm43980faz",
			setupMock: func() (interfaces.UserRepository, interfaces.CacheService) {
				mUserRepository := new(mocks.MockUserRepository)
				mUserRepository.On("FindUserById", mock.Anything, mock.Anything).Return(
					models.User{ID: 1, Password: "$2a$10$0f4BED0dDgYCE8xVREwhUeyjpKTtBIm4eO.xrPC/H8kvsBVM2gpdq"}, nil)
				mUserRepository.On("DeleteUser", mock.Anything, mock.Anything, mock.Anything).Return(
					nil)
				mCacheService := new(mocks.MockCacheService)
				mCacheService.On("ExistsData", mock.Anything, mock.Anything).Return(int64(0), nil)
				mCacheService.On("SetData", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
				mCacheService.On("DeleteData", mock.Anything, mock.Anything).Return(nil)
				return mUserRepository, mCacheService
			},
		},
		{
			name:          "Proper data with user data saved in the cache",
			expectedError: nil,
			password:      "ci$#fm43980faz",
			setupMock: func() (interfaces.UserRepository, interfaces.CacheService) {
				mUserRepository := new(mocks.MockUserRepository)
				mUserRepository.On("DeleteUser", mock.Anything, mock.Anything, mock.Anything).Return(
					nil)
				mCacheService := new(mocks.MockCacheService)
				mCacheService.On("ExistsData", mock.Anything, mock.Anything).Return(int64(1), nil)
				mCacheService.On("GetData", mock.Anything, mock.Anything).Return(`{
				  "id": 1,
				  "email": "joedoe@email.com",
				  "name": "Joe",
				  "surname": "Doe",
				  "password": "$2a$10$0f4BED0dDgYCE8xVREwhUeyjpKTtBIm4eO.xrPC/H8kvsBVM2gpdq",
				  "discord_notifications": false,
				  "email_notifications_settings": true,
				  "slack_notifications_settings": false,
				  "created_at": "2023-01-01T00:00:00Z",
				  "updated_at": "2023-01-01T00:00:00Z"
				}`, nil)
				mCacheService.On("DeleteData", mock.Anything, mock.Anything).Return(nil)
				return mUserRepository, mCacheService
			},
		},
		{
			name:          "Failed to get data from cache",
			expectedError: errors.New("failed to get data from cache"),
			password:      "ci$#fm43980faz",
			setupMock: func() (interfaces.UserRepository, interfaces.CacheService) {
				mUserRepository := new(mocks.MockUserRepository)
				mCacheService := new(mocks.MockCacheService)
				mCacheService.On("ExistsData", mock.Anything, mock.Anything).Return(int64(1),
					nil)
				mCacheService.On("GetData", mock.Anything, mock.Anything).Return(``,
					errors.New("failed to get data from cache"))
				return mUserRepository, mCacheService
			},
		},
		{
			name:          "ExistsData failed",
			expectedError: errors.New("failed to check exists data"),
			password:      "ci$#fm43980faz",
			setupMock: func() (interfaces.UserRepository, interfaces.CacheService) {
				mUserRepository := new(mocks.MockUserRepository)
				mCacheService := new(mocks.MockCacheService)
				mCacheService.On("ExistsData", mock.Anything, mock.Anything).Return(int64(0), errors.New("failed to check exists data"))
				return mUserRepository, mCacheService
			},
		},
		{
			name:          "User does not exist",
			expectedError: errors.New("user with this id does not exist"),
			password:      "ci$#fm43980faz",
			setupMock: func() (interfaces.UserRepository, interfaces.CacheService) {
				mUserRepository := new(mocks.MockUserRepository)
				mUserRepository.On("FindUserById", mock.Anything, mock.Anything).Return(
					models.User{ID: 0}, nil)
				mCacheService := new(mocks.MockCacheService)
				mCacheService.On("ExistsData", mock.Anything, mock.Anything).Return(int64(0), nil)
				return mUserRepository, mCacheService
			},
		},
		{
			name:          "Wrong password provided",
			expectedError: errors.New("wrong password provided"),
			password:      "ci$#fm43980faz2",
			setupMock: func() (interfaces.UserRepository, interfaces.CacheService) {
				mUserRepository := new(mocks.MockUserRepository)
				mUserRepository.On("FindUserById", mock.Anything, mock.Anything).Return(
					models.User{ID: 1, Password: "$2a$10$0f4BED0dDgYCE8xVREwhUeyjpKTtBIm4eO.xrPC/H8kvsBVM2gpdq"}, nil)
				mCacheService := new(mocks.MockCacheService)
				mCacheService.On("ExistsData", mock.Anything, mock.Anything).Return(int64(0), nil)
				mCacheService.On("SetData", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
				return mUserRepository, mCacheService
			},
		},
		{
			name:          "Failed to find a user",
			expectedError: errors.New("failed to find a user"),
			password:      "ci$#fm43980faz2",
			setupMock: func() (interfaces.UserRepository, interfaces.CacheService) {
				mUserRepository := new(mocks.MockUserRepository)
				mUserRepository.On("FindUserById", mock.Anything, mock.Anything).Return(
					models.User{}, errors.New("failed to find a user"))
				mCacheService := new(mocks.MockCacheService)
				mCacheService.On("ExistsData", mock.Anything, mock.Anything).Return(int64(0), nil)
				return mUserRepository, mCacheService
			},
		},

		{
			name:          "Failed to delete a user from database",
			expectedError: errors.New("failed to delete a user from database"),
			password:      "ci$#fm43980faz",
			setupMock: func() (interfaces.UserRepository, interfaces.CacheService) {
				mUserRepository := new(mocks.MockUserRepository)
				mUserRepository.On("FindUserById", mock.Anything, mock.Anything).Return(
					models.User{ID: 1, Password: "$2a$10$0f4BED0dDgYCE8xVREwhUeyjpKTtBIm4eO.xrPC/H8kvsBVM2gpdq"}, nil)
				mUserRepository.On("DeleteUser", mock.Anything, mock.Anything, mock.Anything).Return(
					errors.New("failed to delete a user from database"))
				mCacheService := new(mocks.MockCacheService)
				mCacheService.On("ExistsData", mock.Anything, mock.Anything).Return(int64(0), nil)
				mCacheService.On("SetData", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
				return mUserRepository, mCacheService
			},
		},
		{
			name:          "Failed to delete a user from cache",
			expectedError: errors.New("failed to delete a user from cache"),
			password:      "ci$#fm43980faz",
			setupMock: func() (interfaces.UserRepository, interfaces.CacheService) {
				mUserRepository := new(mocks.MockUserRepository)
				mUserRepository.On("FindUserById", mock.Anything, mock.Anything).Return(
					models.User{ID: 1, Password: "$2a$10$0f4BED0dDgYCE8xVREwhUeyjpKTtBIm4eO.xrPC/H8kvsBVM2gpdq"}, nil)
				mUserRepository.On("DeleteUser", mock.Anything, mock.Anything, mock.Anything).Return(
					nil)
				mCacheService := new(mocks.MockCacheService)
				mCacheService.On("ExistsData", mock.Anything, mock.Anything).Return(int64(0), nil)
				mCacheService.On("SetData", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
				mCacheService.On("DeleteData", mock.Anything, mock.Anything).Return(errors.New(
					"failed to delete a user from cache"))

				return mUserRepository, mCacheService
			},
		},
	}
	for _, testScenario := range testsScenarios {
		t.Run(testScenario.name, func(t *testing.T) {
			ctx := context.Background()
			loggerService := tests.CreateLogger()
			userRepository, cacheService := testScenario.setupMock()
			userService := NewUserService(loggerService, userRepository, cacheService)
			err := userService.DeleteUser(ctx, 3, testScenario.password)
			if testScenario.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), testScenario.expectedError.Error())
			}
		})
	}
}

func TestUserService_ChangeUserPassword(t *testing.T) {
	type args struct {
		name          string
		expectedError error
		password      string
		newPassword   string
		setupMock     func() (interfaces.UserRepository, interfaces.CacheService)
	}
	testsScenarios := []args{
		{
			name:          "Proper data",
			expectedError: nil,
			password:      "ci$#fm43980faz",
			newPassword:   "ci$#fm43980faz",
			setupMock: func() (interfaces.UserRepository, interfaces.CacheService) {
				mUserRepository := new(mocks.MockUserRepository)
				mUserRepository.On("FindUserById", mock.Anything, mock.Anything).Return(
					models.User{ID: 1, Password: "$2a$10$0f4BED0dDgYCE8xVREwhUeyjpKTtBIm4eO.xrPC/H8kvsBVM2gpdq"}, nil)
				mUserRepository.On("ChangeUserPassword", mock.Anything, mock.Anything, mock.Anything).Return(
					nil)
				mCacheService := new(mocks.MockCacheService)
				mCacheService.On("ExistsData", mock.Anything, mock.Anything).Return(int64(0), nil)
				mCacheService.On("SetData", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

				return mUserRepository, mCacheService
			},
		},

		{
			name:          "User with this id does not exist",
			expectedError: errors.New("user with this id does not exist"),
			password:      "ci$#fm43980faz",
			newPassword:   "ci$#fm43980faz",
			setupMock: func() (interfaces.UserRepository, interfaces.CacheService) {
				mUserRepository := new(mocks.MockUserRepository)
				mUserRepository.On("FindUserById", mock.Anything, mock.Anything).Return(
					models.User{ID: 0}, nil)
				mCacheService := new(mocks.MockCacheService)
				mCacheService.On("ExistsData", mock.Anything, mock.Anything).Return(int64(0), nil)
				return mUserRepository, mCacheService
			},
		},
		{
			name:          "Wrong password provided",
			expectedError: errors.New("wrong current password provided"),
			password:      "ci$#fm43980faz2",
			newPassword:   "ci$#fm43980faz",
			setupMock: func() (interfaces.UserRepository, interfaces.CacheService) {
				mUserRepository := new(mocks.MockUserRepository)
				mUserRepository.On("FindUserById", mock.Anything, mock.Anything).Return(
					models.User{ID: 1, Password: "$2a$10$0f4BED0dDgYCE8xVREwhUeyjpKTtBIm4eO.xrPC/H8kvsBVM2gpdq"}, nil)
				mCacheService := new(mocks.MockCacheService)
				mCacheService.On("ExistsData", mock.Anything, mock.Anything).Return(int64(0), nil)
				mCacheService.On("SetData", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
				return mUserRepository, mCacheService
			},
		},
		{
			name:          "Wrong password provided",
			expectedError: errors.New("bcrypt: password length exceeds 72 bytes"),
			password:      "ci$#fm43980faz",
			newPassword:   "ci$#fm4398d432-89m52348-$#@rt43t43tZ#4t43t39ty4343324fn87634c8-t734tct43t430faz",
			setupMock: func() (interfaces.UserRepository, interfaces.CacheService) {
				mUserRepository := new(mocks.MockUserRepository)
				mUserRepository.On("FindUserById", mock.Anything, mock.Anything).Return(
					models.User{ID: 1, Password: "$2a$10$0f4BED0dDgYCE8xVREwhUeyjpKTtBIm4eO.xrPC/H8kvsBVM2gpdq"}, nil)
				mCacheService := new(mocks.MockCacheService)
				mCacheService.On("ExistsData", mock.Anything, mock.Anything).Return(int64(0), nil)
				mCacheService.On("SetData", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
				return mUserRepository, mCacheService
			},
		},
		{
			name:          "Failed to find a user",
			expectedError: errors.New("failed to find a user"),
			password:      "ci$#fm43980faz2",
			newPassword:   "ci$#fm43980faz",
			setupMock: func() (interfaces.UserRepository, interfaces.CacheService) {
				mUserRepository := new(mocks.MockUserRepository)
				mUserRepository.On("FindUserById", mock.Anything, mock.Anything).Return(models.User{},
					errors.New("failed to find a user"))
				mCacheService := new(mocks.MockCacheService)
				mCacheService.On("ExistsData", mock.Anything, mock.Anything).Return(int64(0), nil)
				mCacheService.On("SetData", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
				return mUserRepository, mCacheService
			},
		},
		{
			name:          "Failed to change password",
			expectedError: errors.New("failed to change password"),
			password:      "ci$#fm43980faz",
			newPassword:   "ci$#fm43980faz",
			setupMock: func() (interfaces.UserRepository, interfaces.CacheService) {
				mUserRepository := new(mocks.MockUserRepository)
				mUserRepository.On("FindUserById", mock.Anything, mock.Anything).Return(
					models.User{ID: 1, Password: "$2a$10$0f4BED0dDgYCE8xVREwhUeyjpKTtBIm4eO.xrPC/H8kvsBVM2gpdq"}, nil)
				mUserRepository.On("ChangeUserPassword", mock.Anything, mock.Anything, mock.Anything).Return(
					errors.New("failed to change password"))
				mCacheService := new(mocks.MockCacheService)
				mCacheService.On("ExistsData", mock.Anything, mock.Anything).Return(int64(0), nil)
				mCacheService.On("SetData", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
				return mUserRepository, mCacheService
			},
		},
		{
			name:          "ExistsData failed",
			expectedError: errors.New("failed to check exists data"),
			password:      "ci$#fm43980faz",
			setupMock: func() (interfaces.UserRepository, interfaces.CacheService) {
				mUserRepository := new(mocks.MockUserRepository)
				mCacheService := new(mocks.MockCacheService)
				mCacheService.On("ExistsData", mock.Anything, mock.Anything).Return(int64(0), errors.New("failed to check exists data"))
				return mUserRepository, mCacheService
			},
		},
		{
			name:          "Failed to get data from cache",
			expectedError: errors.New("failed to get data from cache"),
			password:      "ci$#fm43980faz",
			setupMock: func() (interfaces.UserRepository, interfaces.CacheService) {
				mUserRepository := new(mocks.MockUserRepository)
				mCacheService := new(mocks.MockCacheService)
				mCacheService.On("ExistsData", mock.Anything, mock.Anything).Return(int64(1),
					nil)
				mCacheService.On("GetData", mock.Anything, mock.Anything).Return(``,
					errors.New("failed to get data from cache"))
				return mUserRepository, mCacheService
			},
		},
	}
	for _, testScenario := range testsScenarios {
		t.Run(testScenario.name, func(t *testing.T) {
			ctx := context.Background()
			loggerService := tests.CreateLogger()
			userRepository, cacheService := testScenario.setupMock()

			userService := NewUserService(loggerService, userRepository, cacheService)
			err := userService.ChangeUserPassword(ctx, 3, testScenario.password, testScenario.newPassword)
			if testScenario.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), testScenario.expectedError.Error())
			}
		})
	}
}

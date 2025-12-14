package user

import (
	"context"
	"errors"
	"testing"

	"github.com/slodkiadrianek/octopus/internal/DTO"
	"github.com/slodkiadrianek/octopus/internal/config"
	"github.com/slodkiadrianek/octopus/internal/middleware"
	"github.com/slodkiadrianek/octopus/internal/models"
	"github.com/slodkiadrianek/octopus/internal/services/interfaces"
	"github.com/slodkiadrianek/octopus/tests"
	"github.com/slodkiadrianek/octopus/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthService_LoginUser(t *testing.T) {
	type args struct {
		name          string
		password      string
		expectedError error
		setupMock     func() interfaces.UserRepository
	}
	env, err := config.SetConfig(tests.EnvFileLocationForServices)

	if err != nil {
		panic(err)
	}
	testsScenarios := []args{
		{
			name:          "Proper data to login user",
			password:      "ci$#fm43980faz",
			expectedError: nil,
			setupMock: func() interfaces.UserRepository {
				mUserRepository := new(mocks.MockUserRepository)
				mUserRepository.On("FindUserByEmail", mock.Anything, mock.Anything).Return(
					models.User{
						ID:       1,
						Password: "$2a$10$0f4BED0dDgYCE8xVREwhUeyjpKTtBIm4eO.xrPC/H8kvsBVM2gpdq",
					}, nil)
				return mUserRepository
			},
		},
		{
			name:          "Failed to find user by email",
			password:      "ci$#fm43980faz",
			expectedError: errors.New("failed to find user by email"),
			setupMock: func() interfaces.UserRepository {
				mUserRepository := new(mocks.MockUserRepository)
				mUserRepository.On("FindUserByEmail", mock.Anything, mock.Anything).Return(
					models.User{}, errors.New("failed to find user by email"))
				return mUserRepository
			},
		},
		{
			name:          "User with this email does not exist",
			password:      "ci$#fm43980faz",
			expectedError: errors.New("User with this email does not exist"),
			setupMock: func() interfaces.UserRepository {
				mUserRepository := new(mocks.MockUserRepository)
				mUserRepository.On("FindUserByEmail", mock.Anything, mock.Anything).Return(
					models.User{ID: 0}, nil)
				return mUserRepository
			},
		},
		{
			name:          "Wrong password provided",
			password:      "ci$#fm43980faz",
			expectedError: errors.New("Wrong password provided"),
			setupMock: func() interfaces.UserRepository {
				mUserRepository := new(mocks.MockUserRepository)
				mUserRepository.On("FindUserByEmail", mock.Anything, mock.Anything).Return(
					models.User{
						ID:       1,
						Password: "$2a$10$333430f4BED0dDgYCE8xVREwhUeyjpKTtBIm4eO.xrPC/H8kvsBVM2gpdq",
					}, nil)
				return mUserRepository
			},
		},
	}
	for _, testScenario := range testsScenarios {
		t.Run(testScenario.name, func(t *testing.T) {
			ctx := context.Background()
			loggerService := tests.CreateLogger()
			cacheService := tests.CreateCacheService(loggerService)
			userRepository := testScenario.setupMock()
			jwt := middleware.NewJWT(env.JWTSecret, loggerService, cacheService)
			loginData := DTO.LoginUser{Email: "asdfjsdf8932@gmail.com", Password: testScenario.password}
			authService := NewAuthService(loggerService, userRepository, jwt)
			token, err := authService.LoginUser(ctx, loginData)
			if testScenario.expectedError == nil {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
			} else {
				assert.Error(t, err)
				assert.Empty(t, token)
				assert.Contains(t, err.Error(), testScenario.expectedError.Error())
			}
		})
	}
}

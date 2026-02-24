package auth

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/AliceOrlandini/Auto-Light-Pi/internal/auth/mocks"
	"github.com/AliceOrlandini/Auto-Light-Pi/internal/refresh_token"
	"github.com/AliceOrlandini/Auto-Light-Pi/internal/user"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

func TestService_Register(t *testing.T) {
	tests := []struct {
		name          string
		username      string
		email         string
		password      string
		userName      string
		surname       string
		setupMock     func(*mocks.MockuserRepository)
		expectedError error
	}{
		{
			name:     "success",
			username: "mario",
			email:    "mariorossi@gmail.com",
			password: "Testtest123",
			userName: "mario",
			surname:  "rossi",
			setupMock: func(m *mocks.MockuserRepository) {
				m.EXPECT().GetOneByEmail(gomock.Any(), "mariorossi@gmail.com").Return(nil, nil)
				m.EXPECT().GetOneByUsername(gomock.Any(), "mario").Return(nil, nil)
				m.EXPECT().CreateOne(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:     "email_already_exists",
			username: "mario",
			email:    "mariorossi@gmail.com",
			password: "Testtest123",
			userName: "mario",
			surname:  "rossi",
			setupMock: func(m *mocks.MockuserRepository) {
				m.EXPECT().GetOneByEmail(gomock.Any(), "mariorossi@gmail.com").Return(&user.User{ID: "1"}, nil)
			},
			expectedError: ErrUserAlreadyExists,
		},
		{
			name:     "username_already_exists",
			username: "mario",
			email:    "mariorossi@gmail.com",
			password: "Testtest123",
			userName: "mario",
			surname:  "rossi",
			setupMock: func(m *mocks.MockuserRepository) {
				m.EXPECT().GetOneByEmail(gomock.Any(), "mariorossi@gmail.com").Return(nil, nil)
				m.EXPECT().GetOneByUsername(gomock.Any(), "mario").Return(&user.User{ID: "1"}, nil)
			},
			expectedError: ErrUserAlreadyExists,
		},
		{
			name:     "db_error_on_email_check",
			username: "mario",
			email:    "mariorossi@gmail.com",
			password: "Testtest123",
			userName: "mario",
			surname:  "rossi",
			setupMock: func(m *mocks.MockuserRepository) {
				m.EXPECT().GetOneByEmail(gomock.Any(), "mariorossi@gmail.com").Return(nil, errors.New("db error"))
			},
			expectedError: errors.New("db error"),
		},
		{
			name:     "db_error_on_username_check",
			username: "mario",
			email:    "mariorossi@gmail.com",
			password: "Testtest123",
			userName: "mario",
			surname:  "rossi",
			setupMock: func(m *mocks.MockuserRepository) {
				m.EXPECT().GetOneByEmail(gomock.Any(), "mariorossi@gmail.com").Return(nil, nil)
				m.EXPECT().GetOneByUsername(gomock.Any(), "mario").Return(nil, errors.New("db error"))
			},
			expectedError: errors.New("db error"),
		},
		{
			name:     "db_error_on_create",
			username: "mario",
			email:    "mariorossi@gmail.com",
			password: "Testtest123",
			userName: "mario",
			surname:  "rossi",
			setupMock: func(m *mocks.MockuserRepository) {
				m.EXPECT().GetOneByEmail(gomock.Any(), "mariorossi@gmail.com").Return(nil, nil)
				m.EXPECT().GetOneByUsername(gomock.Any(), "mario").Return(nil, nil)
				m.EXPECT().CreateOne(gomock.Any(), gomock.Any()).Return(errors.New("db error"))
			},
			expectedError: errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUserRepo := mocks.NewMockuserRepository(ctrl)
			mockTokenRepo := mocks.NewMockrefreshTokenRepository(ctrl)
			tt.setupMock(mockUserRepo)

			s := NewAuthService(mockUserRepo, mockTokenRepo)
			err := s.Register(context.Background(), tt.username, tt.email, tt.password, tt.userName, tt.surname)

			if tt.expectedError != nil {
				if err == nil {
					t.Errorf("expected error %v, got nil", tt.expectedError)
				} else {
					if !errors.Is(err, tt.expectedError) && err.Error() != tt.expectedError.Error() {
						t.Errorf("expected error %v, got %v", tt.expectedError, err)
					}
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
			}
		})
	}
}

func TestService_LoginByUsername(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("Testtest123"), 14)

	tests := []struct {
		name          string
		username      string
		password      string
		setupMock     func(*mocks.MockuserRepository)
		expectedUser  *user.User
		expectedError error
	}{
		{
			name:     "success",
			username: "mario",
			password: "Testtest123",
			setupMock: func(m *mocks.MockuserRepository) {
				m.EXPECT().GetOneByUsername(gomock.Any(), "mario").Return(&user.User{
					Username: "mario",
					Password: string(hashedPassword),
				}, nil)
			},
			expectedUser: &user.User{
				Username: "mario",
				Password: string(hashedPassword),
			},
			expectedError: nil,
		},
		{
			name:     "user_not_exists",
			username: "mario",
			password: "Testtest123",
			setupMock: func(m *mocks.MockuserRepository) {
				m.EXPECT().GetOneByUsername(gomock.Any(), "mario").Return(nil, nil)
			},
			expectedUser:  nil,
			expectedError: ErrUserNotExists,
		},
		{
			name:     "db_error",
			username: "mario",
			password: "Testtest123",
			setupMock: func(m *mocks.MockuserRepository) {
				m.EXPECT().GetOneByUsername(gomock.Any(), "mario").Return(nil, errors.New("db error"))
			},
			expectedUser:  nil,
			expectedError: errors.New("db error"),
		},
		{
			name:     "invalid_password",
			username: "mario",
			password: "WrongPassword",
			setupMock: func(m *mocks.MockuserRepository) {
				m.EXPECT().GetOneByUsername(gomock.Any(), "mario").Return(&user.User{
					Username: "mario",
					Password: string(hashedPassword),
				}, nil)
			},
			expectedUser:  nil,
			expectedError: ErrInvalidPassword,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUserRepo := mocks.NewMockuserRepository(ctrl)
			mockTokenRepo := mocks.NewMockrefreshTokenRepository(ctrl)
			tt.setupMock(mockUserRepo)

			s := NewAuthService(mockUserRepo, mockTokenRepo)
			user, err := s.LoginByUsername(context.Background(), tt.username, tt.password)

			if tt.expectedError != nil {
				if err == nil {
					t.Errorf("expected error %v, got nil", tt.expectedError)
				} else {
					if !errors.Is(err, tt.expectedError) && err.Error() != tt.expectedError.Error() {
						t.Errorf("expected error %v, got %v", tt.expectedError, err)
					}
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
				if user.Username != tt.expectedUser.Username {
					t.Errorf("expected user %v, got %v", tt.expectedUser, user)
				}
			}
		})
	}
}

func TestService_LoginByEmail(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("Testtest123"), 14)

	tests := []struct {
		name          string
		email         string
		password      string
		setupMock     func(*mocks.MockuserRepository)
		expectedUser  *user.User
		expectedError error
	}{
		{
			name:     "success",
			email:    "mariorossi@gmail.com",
			password: "Testtest123",
			setupMock: func(m *mocks.MockuserRepository) {
				m.EXPECT().GetOneByEmail(gomock.Any(), "mariorossi@gmail.com").Return(&user.User{
					Email:    "mariorossi@gmail.com",
					Password: string(hashedPassword),
				}, nil)
			},
			expectedUser: &user.User{
				Email:    "mariorossi@gmail.com",
				Password: string(hashedPassword),
			},
			expectedError: nil,
		},
		{
			name:     "user_not_exists",
			email:    "mariorossi@gmail.com",
			password: "Testtest123",
			setupMock: func(m *mocks.MockuserRepository) {
				m.EXPECT().GetOneByEmail(gomock.Any(), "mariorossi@gmail.com").Return(nil, nil)
			},
			expectedUser:  nil,
			expectedError: ErrUserNotExists,
		},
		{
			name:     "db_error",
			email:    "mariorossi@gmail.com",
			password: "Testtest123",
			setupMock: func(m *mocks.MockuserRepository) {
				m.EXPECT().GetOneByEmail(gomock.Any(), "mariorossi@gmail.com").Return(nil, errors.New("db error"))
			},
			expectedUser:  nil,
			expectedError: errors.New("db error"),
		},
		{
			name:     "invalid_password",
			email:    "mariorossi@gmail.com",
			password: "WrongPassword",
			setupMock: func(m *mocks.MockuserRepository) {
				m.EXPECT().GetOneByEmail(gomock.Any(), "mariorossi@gmail.com").Return(&user.User{
					Email:    "mariorossi@gmail.com",
					Password: string(hashedPassword),
				}, nil)
			},
			expectedUser:  nil,
			expectedError: ErrInvalidPassword,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUserRepo := mocks.NewMockuserRepository(ctrl)
			mockTokenRepo := mocks.NewMockrefreshTokenRepository(ctrl)
			tt.setupMock(mockUserRepo)

			s := NewAuthService(mockUserRepo, mockTokenRepo)
			user, err := s.LoginByEmail(context.Background(), tt.email, tt.password)

			if tt.expectedError != nil {
				if err == nil {
					t.Errorf("expected error %v, got nil", tt.expectedError)
				} else {
					if !errors.Is(err, tt.expectedError) && err.Error() != tt.expectedError.Error() {
						t.Errorf("expected error %v, got %v", tt.expectedError, err)
					}
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
				if user.Email != tt.expectedUser.Email {
					t.Errorf("expected user %v, got %v", tt.expectedUser, user)
				}
			}
		})
	}
}

func TestService_GenerateJWT(t *testing.T) {
	os.Setenv("JWT_SECRET", "secret")
	os.Setenv("APPLICATION_NAME", "app")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockuserRepository(ctrl)
	mockTokenRepo := mocks.NewMockrefreshTokenRepository(ctrl)

	s := NewAuthService(mockUserRepo, mockTokenRepo)
	token, err := s.GenerateJWT("UserID")

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if token == "" {
		t.Errorf("expected token, got empty string")
	}
}

func TestService_GenerateRefreshToken(t *testing.T) {
	tests := []struct {
		name          string
		userID        string
		setupMock     func(*mocks.MockrefreshTokenRepository)
		expectedError error
	}{
		{
			name:   "success",
			userID: "UserID",
			setupMock: func(m *mocks.MockrefreshTokenRepository) {
				m.EXPECT().CreateOne(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:   "db_error",
			userID: "UserID",
			setupMock: func(m *mocks.MockrefreshTokenRepository) {
				m.EXPECT().CreateOne(gomock.Any(), gomock.Any()).Return(errors.New("db error"))
			},
			expectedError: errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUserRepo := mocks.NewMockuserRepository(ctrl)
			mockTokenRepo := mocks.NewMockrefreshTokenRepository(ctrl)
			tt.setupMock(mockTokenRepo)

			s := NewAuthService(mockUserRepo, mockTokenRepo)
			token, err := s.GenerateRefreshToken(context.Background(), tt.userID)

			if tt.expectedError != nil {
				if err == nil {
					t.Errorf("expected error %v, got nil", tt.expectedError)
				} else {
					if err.Error() != tt.expectedError.Error() {
						t.Errorf("expected error %v, got %v", tt.expectedError, err)
					}
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
				if token == nil {
					t.Errorf("expected token, got nil")
				}
			}
		})
	}
}

func TestService_ValidateRefreshToken(t *testing.T) {
	validToken := "rt1.validtokenopaque"

	tests := []struct {
		name           string
		token          string
		setupMock      func(*mocks.MockrefreshTokenRepository)
		expectedUserID string
		expectedError  error
	}{
		{
			name:  "success",
			token: validToken,
			setupMock: func(m *mocks.MockrefreshTokenRepository) {
				m.EXPECT().GetOneUserIDByTokenHash(gomock.Any(), gomock.Any()).Return("UserID", nil)
			},
			expectedUserID: "UserID",
			expectedError:  nil,
		},
		{
			name:  "invalid_prefix",
			token: "invalid.token",
			setupMock: func(m *mocks.MockrefreshTokenRepository) {
			},
			expectedUserID: "",
			expectedError:  ErrInvalidToken,
		},
		{
			name:  "token_not_found",
			token: validToken,
			setupMock: func(m *mocks.MockrefreshTokenRepository) {
				m.EXPECT().GetOneUserIDByTokenHash(gomock.Any(), gomock.Any()).Return("", refresh_token.ErrTokenHashNotFound)
			},
			expectedUserID: "",
			expectedError:  ErrInvalidToken,
		},
		{
			name:  "db_error",
			token: validToken,
			setupMock: func(m *mocks.MockrefreshTokenRepository) {
				m.EXPECT().GetOneUserIDByTokenHash(gomock.Any(), gomock.Any()).Return("", errors.New("db error"))
			},
			expectedUserID: "",
			expectedError:  errors.New("db error"),
		},
		{
			name:  "empty_user_id",
			token: validToken,
			setupMock: func(m *mocks.MockrefreshTokenRepository) {
				m.EXPECT().GetOneUserIDByTokenHash(gomock.Any(), gomock.Any()).Return("", nil)
			},
			expectedUserID: "",
			expectedError:  ErrInvalidToken,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUserRepo := mocks.NewMockuserRepository(ctrl)
			mockTokenRepo := mocks.NewMockrefreshTokenRepository(ctrl)
			tt.setupMock(mockTokenRepo)

			s := NewAuthService(mockUserRepo, mockTokenRepo)
			userID, err := s.ValidateRefreshToken(context.Background(), tt.token)

			if tt.expectedError != nil {
				if err == nil {
					t.Errorf("expected error %v, got nil", tt.expectedError)
				} else {
					if !errors.Is(err, tt.expectedError) && err.Error() != tt.expectedError.Error() {
						t.Errorf("expected error %v, got %v", tt.expectedError, err)
					}
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
				if userID != tt.expectedUserID {
					t.Errorf("expected userID %v, got %v", tt.expectedUserID, userID)
				}
			}
		})
	}
}

func TestService_RotateRefreshToken(t *testing.T) {
	tests := []struct {
		name          string
		userID        string
		setupMock     func(*mocks.MockrefreshTokenRepository)
		expectedError error
	}{
		{
			name:   "success",
			userID: "UserID",
			setupMock: func(m *mocks.MockrefreshTokenRepository) {
				m.EXPECT().DeleteOneByUserID(gomock.Any(), "UserID").Return(nil)
				m.EXPECT().CreateOne(gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:   "delete_error",
			userID: "UserID",
			setupMock: func(m *mocks.MockrefreshTokenRepository) {
				m.EXPECT().DeleteOneByUserID(gomock.Any(), "UserID").Return(errors.New("db error"))
			},
			expectedError: errors.New("db error"),
		},
		{
			name:   "create_error",
			userID: "UserID",
			setupMock: func(m *mocks.MockrefreshTokenRepository) {
				m.EXPECT().DeleteOneByUserID(gomock.Any(), "UserID").Return(nil)
				m.EXPECT().CreateOne(gomock.Any(), gomock.Any()).Return(errors.New("db error"))
			},
			expectedError: errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUserRepo := mocks.NewMockuserRepository(ctrl)
			mockTokenRepo := mocks.NewMockrefreshTokenRepository(ctrl)
			tt.setupMock(mockTokenRepo)

			s := NewAuthService(mockUserRepo, mockTokenRepo)
			token, err := s.RotateRefreshToken(context.Background(), tt.userID)

			if tt.expectedError != nil {
				if err == nil {
					t.Errorf("expected error %v, got nil", tt.expectedError)
				} else {
					if err.Error() != tt.expectedError.Error() {
						t.Errorf("expected error %v, got %v", tt.expectedError, err)
					}
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
				if token == nil {
					t.Errorf("expected token, got nil")
				}
			}
		})
	}
}
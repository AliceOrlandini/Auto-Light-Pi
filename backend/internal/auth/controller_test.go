package auth

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AliceOrlandini/Auto-Light-Pi/internal/auth/mocks"
	"github.com/AliceOrlandini/Auto-Light-Pi/internal/refresh_token"
	"github.com/AliceOrlandini/Auto-Light-Pi/internal/user"
	"github.com/gin-gonic/gin"
	"go.uber.org/mock/gomock"
)

func init() { gin.SetMode(gin.TestMode) }

/**
 * newTestContext is a helper function to create a new gin.Context and httptest.ResponseRecorder
 * This is necessary since the function that we are testing requires a gin.Context.
 */
func newTestContext(method, path string, body []byte) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	req := httptest.NewRequest(method, path, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	c.Request = req
	return c, w
}

func TestController_Register(t *testing.T) {
	tests := []struct {
		name         string
		body         string
		expectedCode int
		setupMock    func(*mocks.MockauthService)
	}{
		{
			name:         "success",
			body:         `{"username":"mario","email":"mariorossi@gmail.com","password":"Testtest123","name":"mario","surname":"rossi"}`,
			expectedCode: http.StatusCreated,
			// We mock the Service.Register method to return nil, which means that the registration was successful.
			// In the other test cases, we don't need to mock the Service.Register method because the validation will fail before calling it.
			setupMock: func(m *mocks.MockauthService) {
				m.EXPECT().
					Register(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
			},
		},
		{
			name: "user_already_exists",
			body: `{"username":"mario","email":"mariorossi@gmail.com","password":"Testtest123","name":"mario","surname":"rossi"}`,
			expectedCode: http.StatusBadRequest,
			setupMock: func(m *mocks.MockauthService) {
				m.EXPECT().
					Register(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(ErrUserAlreadyExists)
			},
		},
		{
			name: "context_cancelled",
			body: `{"username":"mario","email":"mariorossi@gmail.com","password":"Testtest123","name":"mario","surname":"rossi"}`,
			// no expected code, since the context is cancelled before calling the Register method
			expectedCode: http.StatusOK,
			setupMock: func(m *mocks.MockauthService) {
				m.EXPECT().
					Register(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(context.Canceled)
			},
		},
		{
			name: "deadline_exceeded",
			body: `{"username":"mario","email":"mariorossi@gmail.com","password":"Testtest123","name":"mario","surname":"rossi"}`,
			// no expected code, since the context is cancelled before calling the Register method
			expectedCode: http.StatusRequestTimeout,
			setupMock: func(m *mocks.MockauthService) {
				m.EXPECT().
					Register(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(context.DeadlineExceeded)
			},
		},
		{
			name: "internal_server_error",
			body: `{"username":"mario","email":"mariorossi@gmail.com","password":"Testtest123","name":"mario","surname":"rossi"}`,
			expectedCode: http.StatusInternalServerError,
			setupMock: func(m *mocks.MockauthService) {
				m.EXPECT().
					Register(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(fmt.Errorf("some internal error"))
			},
		},
		{
			name:         "missing_username",
			body:         `{"email":"mariorossi@gmail.com","password":"Testtest123","name":"mario","surname":"rossi"}`,
			expectedCode: http.StatusBadRequest,
			setupMock:    func(m *mocks.MockauthService) {},
		},
		{
			name:         "missing_email",
			body:         `{"username":"mario","password":"Testtest123","name":"mario","surname":"rossi"}`,
			expectedCode: http.StatusBadRequest,
			setupMock:    func(m *mocks.MockauthService) {},
		},
		{
			name:         "missing_password",
			body:         `{"username":"mario","email":"mariorossi@gmail.com","name":"mario","surname":"rossi"}`,
			expectedCode: http.StatusBadRequest,
			setupMock:    func(m *mocks.MockauthService) {},
		},
		{
			name:         "missing_name",
			body:         `{"username":"mario","email":"mariorossi@gmail.com","password":"Testtest123","surname":"rossi"}`,
			expectedCode: http.StatusBadRequest,
			setupMock:    func(m *mocks.MockauthService) {},
		},
		{
			name:         "missing_surname",
			body:         `{"username":"mario","email":"mariorossi@gmail.com","password":"Testtest123","name":"mario"}`,
			expectedCode: http.StatusBadRequest,
			setupMock:    func(m *mocks.MockauthService) {},
		},
		{
			name:         "empty_username",
			body:         `{"username":"","email":"mariorossi@gmail.com","password":"Testtest123","name":"mario","surname":"rossi"}`,
			expectedCode: http.StatusBadRequest,
			setupMock:    func(m *mocks.MockauthService) {},
		},
		{
			name:         "empty_email",
			body:         `{"username":"mario","email":"","password":"Testtest123","name":"mario","surname":"rossi"}`,
			expectedCode: http.StatusBadRequest,
			setupMock:    func(m *mocks.MockauthService) {},
		},
		{
			name:         "empty_password",
			body:         `{"username":"mario","email":"mariorossi@gmail.com","password":"","name":"mario","surname":"rossi"}`,
			expectedCode: http.StatusBadRequest,
			setupMock:    func(m *mocks.MockauthService) {},
		},
		{
			name:         "empty_name",
			body:         `{"username":"mario","email":"mariorossi@gmail.com","password":"Testtest123","name":"","surname":"rossi"}`,
			expectedCode: http.StatusBadRequest,
			setupMock:    func(m *mocks.MockauthService) {},
		},
		{
			name:         "empty_surname",
			body:         `{"username":"mario","email":"mariorossi@gmail.com","password":"Testtest123","name":"mario","surname":""}`,
			expectedCode: http.StatusBadRequest,
			setupMock:    func(m *mocks.MockauthService) {},
		},
		{
			name:         "invalid_json",
			body:         `{"username":123}`,
			expectedCode: http.StatusBadRequest,
			setupMock:    func(m *mocks.MockauthService) {},
		},
		{
			name:         "short_password",
			body:         `{"username":"mario","email":"mariorossi@gmail.com","password":"short","name":"mario","surname":"rossi"}`,
			expectedCode: http.StatusBadRequest,
			setupMock:    func(m *mocks.MockauthService) {},
		},
		{
			name:         "long_password",
			body:         `{"username":"mario","email":"mariorossi@gmail.com","password":"thisisaverylongpasswordthatexceedsthelimit","name":"mario","surname":"rossi"}`,
			expectedCode: http.StatusBadRequest,
			setupMock:    func(m *mocks.MockauthService) {},
		},
		{
			name:         "no_number",
			body:         `{"username":"mario","email":"mariorossi@gmail.com","password":"PasswordWithoutNumber","name":"mario","surname":"rossi"}`,
			expectedCode: http.StatusBadRequest,
			setupMock:    func(m *mocks.MockauthService) {},
		},
		{
			name:         "no_uppercase",
			body:         `{"username":"mario","email":"mariorossi@gmail.com","password":"passwordwithoutuppercase1","name":"mario","surname":"rossi"}`,
			expectedCode: http.StatusBadRequest,
			setupMock:    func(m *mocks.MockauthService) {},
		},
		{
			name:         "no_lowercase",
			body:         `{"username":"mario","email":"mariorossi@gmail.com","password":"PASSWORDWITHOUTLOWERCASE1","name":"mario","surname":"rossi"}`,
			expectedCode: http.StatusBadRequest,
			setupMock:    func(m *mocks.MockauthService) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAuthService := mocks.NewMockauthService(ctrl)
			if tt.setupMock != nil {
				tt.setupMock(mockAuthService)
			}

			uc := &Controller{service: mockAuthService}

			body := []byte(tt.body)
			c, w := newTestContext(http.MethodPost, "/register", body)

			uc.Register(c)

			if w.Code != tt.expectedCode {
				t.Fatalf("got %d want %d; body=%s", w.Code, tt.expectedCode, w.Body.String())
			}
		})
	}
}

func TestController_LoginByUsername(t *testing.T) {
	tests := []struct {
		name string
		body string
		expectedCode int
		setupMock func(*mocks.MockauthService)
	}{
		{
			name: "success",
			body: `{"username":"mario","password":"Testtest123"}`,
			expectedCode: http.StatusOK,
			setupMock: func(m *mocks.MockauthService) {
				dummyUser := &user.User{
					ID:       "1",
					Username: "mario",
					Email:    "mario@example.com",
					Password: "hashedpassword",
					Name:     "Mario",
					Surname:  "Rossi",
				}
				m.EXPECT().LoginByUsername(gomock.Any(), "mario", "Testtest123").Return(dummyUser, nil)
				m.EXPECT().GenerateJWT("1").Return("dummy_jwt_token", nil)
				dummyRefreshToken := &refresh_token.RefreshToken{
					RefreshToken: "dummy_refresh_token",
					UserID:       "1",
					TTL:          time.Now().Add(24 * time.Hour),
				}
				m.EXPECT().GenerateRefreshToken(gomock.Any(), "1").Return(dummyRefreshToken, nil)
			},
		},
		{
			name: "missing_username",
			body: `{"password":"Testtest123"}`,
			expectedCode: http.StatusBadRequest,
			setupMock: func(m *mocks.MockauthService) {},
		},
		{
			name: "missing_password",
			body: `{"username":"mario"}`,
			expectedCode: http.StatusBadRequest,
			setupMock: func(m *mocks.MockauthService) {},
		},
		{
			name: "invalid_json",
			body: `{"username":123}`,
			expectedCode: http.StatusBadRequest,
			setupMock: func(m *mocks.MockauthService) {},
		},
		{
			name: "password_not_valid",
			body: `{"username":"mario","password":"invalidpassword"}`,
			expectedCode: http.StatusBadRequest,
			setupMock: func(m *mocks.MockauthService) {},
		},
		{
			name: "login_unsuccessful",
			body: `{"username":"mario","password":"Testtest123"}`,
			expectedCode: http.StatusUnauthorized,
			setupMock: func(m *mocks.MockauthService) {
				m.EXPECT().
				LoginByUsername(gomock.Any(), "mario", "Testtest123").
				Return(nil, ErrUserNotExists)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAuthService := mocks.NewMockauthService(ctrl)
			if tt.setupMock != nil {
				tt.setupMock(mockAuthService)
			}

			uc := NewAuthController(mockAuthService)

			body := []byte(tt.body)
			c, w := newTestContext(http.MethodPost, "/login/username", body)

			uc.LoginByUsername(c)

			if w.Code != tt.expectedCode {
				t.Fatalf("got %d want %d; body=%s", w.Code, tt.expectedCode, w.Body.String())
			}
		})
	}
}

func TestController_LoginByEmail(t *testing.T) {
	tests := []struct {
		name string
		body string
		expectedCode int
		setupMock func(*mocks.MockauthService)
	}{
		{
			name: "success",
			body: `{"email":"mario@example.com","password":"Testtest123"}`,
			expectedCode: http.StatusOK,
			setupMock: func(m *mocks.MockauthService) {
				dummyUser := &user.User{
					ID:       "1",
					Username: "mario",
					Email:    "mario@example.com",
					Password: "hashedpassword",
					Name:     "Mario",
					Surname:  "Rossi",
				}
				dummyRefreshToken := &refresh_token.RefreshToken{
					RefreshToken: "dummy_refresh_token",
					UserID:       "1",
					TTL:          time.Now().Add(24 * time.Hour),
				}
				m.EXPECT().
					LoginByEmail(gomock.Any(), "mario@example.com", "Testtest123").
					Return(dummyUser, nil)
				m.EXPECT().
					GenerateJWT("1").
					Return("dummy_jwt_token", nil)
				m.EXPECT().
					GenerateRefreshToken(gomock.Any(), "1").
					Return(dummyRefreshToken, nil)
			},
		},
		{
			name: "missing_email",
			body: `{"password":"Testtest123"}`,
			expectedCode: http.StatusBadRequest,
			setupMock: func(m *mocks.MockauthService) {},
		},
		{
			name: "missing_password",
			body: `{"email":"mario@example.com"}`,
			expectedCode: http.StatusBadRequest,
			setupMock: func(m *mocks.MockauthService) {},
		},
		{
			name: "invalid_json",
			body: `{"email":123}`,
			expectedCode: http.StatusBadRequest,
			setupMock: func(m *mocks.MockauthService) {},
		},
		{
			name: "password_not_valid",
			body: `{"email":"mario@example.com","password":"invalidpassword"}`,
			expectedCode: http.StatusBadRequest,
			setupMock: func(m *mocks.MockauthService) {},
		},
		{
			name: "login_unsuccessful",
			body: `{"email":"mario@example.com","password":"Testtest123"}`,
			expectedCode: http.StatusUnauthorized,
			setupMock: func(m *mocks.MockauthService) {
				m.EXPECT().
				LoginByEmail(gomock.Any(), "mario@example.com", "Testtest123").
				Return(nil, ErrUserNotExists)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAuthService := mocks.NewMockauthService(ctrl)
			if tt.setupMock != nil {
				tt.setupMock(mockAuthService)
			}

			uc := NewAuthController(mockAuthService)

			body := []byte(tt.body)
			c, w := newTestContext(http.MethodPost, "/login/email", body)

			uc.LoginByEmail(c)

			if w.Code != tt.expectedCode {
				t.Fatalf("got %d want %d; body=%s", w.Code, tt.expectedCode, w.Body.String())
			}
		})
	}
}

func TestController_RefreshToken(t *testing.T) {
	tests := []struct {
		name         string
		cookies      []*http.Cookie
		expectedCode int
		setupMock    func(*mocks.MockauthService)
	}{
		{
			name:         "success",
			cookies:      []*http.Cookie{{Name: "__Host-refresh_token", Value: "dummy_refresh_token"}},
			expectedCode: http.StatusOK,
			setupMock: func(m *mocks.MockauthService) {
				dummyNewRefreshToken := &refresh_token.RefreshToken{
					RefreshToken: "new_dummy_refresh_token",
					UserID:       "1",
					TTL:          time.Now().Add(24 * time.Hour),
				}

				m.EXPECT().
					ValidateRefreshToken(gomock.Any(), "dummy_refresh_token").
					Return("1", nil)

				m.EXPECT().
					RotateRefreshToken(gomock.Any(), "1").
					Return(dummyNewRefreshToken, nil)

				m.EXPECT().
					GenerateJWT("1").
					Return("dummy_jwt_token", nil)
			},
		},
		{
			name:         "missing_refresh_token_cookie",
			cookies:      []*http.Cookie{},
			expectedCode: http.StatusUnauthorized,
			setupMock:    func(m *mocks.MockauthService) {},
		},
		{
			name: 			 "invalid_refresh_token",
			cookies:      []*http.Cookie{{Name: "__Host-refresh_token", Value: "invalid_refresh_token"}},
			expectedCode: http.StatusUnauthorized,
			setupMock: func(m *mocks.MockauthService) {
				m.EXPECT().
					ValidateRefreshToken(gomock.Any(), "invalid_refresh_token").
					Return("", ErrInvalidToken)
			},
		},
		{
			name: 			 "validate_refresh_token_context_cancelled",
			cookies:      []*http.Cookie{{Name: "__Host-refresh_token", Value: "dummy_refresh_token"}},
			expectedCode: http.StatusOK,
			setupMock: func(m *mocks.MockauthService) {
				m.EXPECT().
					ValidateRefreshToken(gomock.Any(), "dummy_refresh_token").
					Return("", context.Canceled)
			},
		},
		{
			name: 			 "validate_refresh_token_deadline_exceeded",
			cookies:      []*http.Cookie{{Name: "__Host-refresh_token", Value: "dummy_refresh_token"}},
			expectedCode: http.StatusRequestTimeout,
			setupMock: func(m *mocks.MockauthService) {
				m.EXPECT().
					ValidateRefreshToken(gomock.Any(), "dummy_refresh_token").
					Return("", context.DeadlineExceeded)
			},
		},
		{
			name: 			 "validate_refresh_token_internal_error",
			cookies:      []*http.Cookie{{Name: "__Host-refresh_token", Value: "dummy_refresh_token"}},
			expectedCode: http.StatusInternalServerError,
			setupMock: func(m *mocks.MockauthService) {
				m.EXPECT().
					ValidateRefreshToken(gomock.Any(), "dummy_refresh_token").
					Return("", fmt.Errorf("some internal error"))
			},
		},
		{
			name: 			 "rotate_refresh_token_internal_error",
			cookies:      []*http.Cookie{{Name: "__Host-refresh_token", Value: "dummy_refresh_token"}},
			expectedCode: http.StatusInternalServerError,
			setupMock: func(m *mocks.MockauthService) {
				m.EXPECT().
					ValidateRefreshToken(gomock.Any(), "dummy_refresh_token").
					Return("1", nil)

				m.EXPECT().
					RotateRefreshToken(gomock.Any(), "1").
					Return(nil, fmt.Errorf("some error"))
			},
		},
		{
			name: 			 "rotate_refresh_token_context_cancelled",
			cookies:      []*http.Cookie{{Name: "__Host-refresh_token", Value: "dummy_refresh_token"}},
			expectedCode: http.StatusOK,
			setupMock: func(m *mocks.MockauthService) {
				m.EXPECT().
					ValidateRefreshToken(gomock.Any(), "dummy_refresh_token").
					Return("1", nil)

				m.EXPECT().
					RotateRefreshToken(gomock.Any(), "1").
					Return(nil, context.Canceled)
			},
		},
		{
			name: 			 "rotate_refresh_token_deadline_exceeded",
			cookies:      []*http.Cookie{{Name: "__Host-refresh_token", Value: "dummy_refresh_token"}},
			expectedCode: http.StatusRequestTimeout,
			setupMock: func(m *mocks.MockauthService) {
				m.EXPECT().
					ValidateRefreshToken(gomock.Any(), "dummy_refresh_token").
					Return("1", nil)

				m.EXPECT().
					RotateRefreshToken(gomock.Any(), "1").
					Return(nil, context.DeadlineExceeded)
			},
		},
		{
			name: 			 "generate_jwt_error",
			cookies:      []*http.Cookie{{Name: "__Host-refresh_token", Value: "dummy_refresh_token"}},
			expectedCode: http.StatusInternalServerError,
			setupMock: func(m *mocks.MockauthService) {
				m.EXPECT().
					ValidateRefreshToken(gomock.Any(), "dummy_refresh_token").
					Return("1", nil)

				dummyNewRefreshToken := &refresh_token.RefreshToken{
					RefreshToken: "new_dummy_refresh_token",
					UserID:       "1",
					TTL:          time.Now().Add(24 * time.Hour),
				}

				m.EXPECT().
					RotateRefreshToken(gomock.Any(), "1").
					Return(dummyNewRefreshToken, nil)

				m.EXPECT().
					GenerateJWT("1").
					Return("", fmt.Errorf("some error"))	
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAuthService := mocks.NewMockauthService(ctrl)
			if tt.setupMock != nil {
				tt.setupMock(mockAuthService)
			}

			uc := NewAuthController(mockAuthService)

			c, w := newTestContext(http.MethodPost, "/refresh-token", nil)

			if len(tt.cookies) > 0 {
				for _, cookie := range tt.cookies {
					c.Request.AddCookie(cookie)
				}
			}

			uc.RefreshToken(c)

			if w.Code != tt.expectedCode {
				t.Fatalf("got %d want %d; body=%s", w.Code, tt.expectedCode, w.Body.String())
			}
		})
	}
}

func TestController_HandleLoginError(t *testing.T) {
	tests := []struct {
		name         string
		err          error
		expectedCode int
	}{
		{
			name:         "user_not_exists",
			err:          ErrUserNotExists,
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "invalid_password",
			err:          ErrInvalidPassword,
			expectedCode: http.StatusUnauthorized,
		},
		{
			name:         "context_cancelled",
			err:          context.Canceled,
			expectedCode: http.StatusOK,
		},
		{
			name:         "deadline_exceeded",
			err:          context.DeadlineExceeded,
			expectedCode: http.StatusRequestTimeout,
		},
		{
			name:         "internal_server_error",
			err:          fmt.Errorf("some internal error"),
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAuthService := mocks.NewMockauthService(ctrl)

			uc := NewAuthController(mockAuthService)

			c, w := newTestContext(http.MethodPost, "/login/username", nil)

			uc.handleLoginError(c, tt.err)

			if w.Code != tt.expectedCode {
				t.Fatalf("got %d want %d; body=%s", w.Code, tt.expectedCode, w.Body.String())
			}
		})
	}
}

func TestController_HandleSuccessfulLogin(t *testing.T) {
	tests := []struct {
		name				 string
		user 				 *user.User
		expectedCode int
		setupMock    func(*mocks.MockauthService)
	}{
		{
			name:         "success",
			expectedCode: http.StatusOK,
			user: &user.User{
				ID:       "1",
				Username: "mario",
				Email:    "mario@example.com",
				Password: "hashedpassword",
				Name:     "Mario",
				Surname:  "Rossi",
			},
			setupMock: func(m *mocks.MockauthService) {
				dummyRefreshToken := &refresh_token.RefreshToken{
					RefreshToken: "dummy_refresh_token",
					UserID:       "1",
					TTL:          time.Now().Add(24 * time.Hour),
				}
				m.EXPECT().
					GenerateJWT("1").
					Return("dummy_jwt_token", nil)
				m.EXPECT().
					GenerateRefreshToken(gomock.Any(), "1").
					Return(dummyRefreshToken, nil)
			},
		},
		{
			name: 			 "generate_jwt_error",
			expectedCode: http.StatusInternalServerError,
			user: &user.User{
				ID:			 "1",
				Username: "mario",
				Email:    "mario@example.com",
				Password: "hashedpassword",
				Name:     "Mario",
				Surname:  "Rossi",
			},
			setupMock: func(m *mocks.MockauthService) {
				m.EXPECT().
					GenerateJWT("1").
					Return("", fmt.Errorf("some error"))
			},
		},
		{
			name: 			 "generate_refresh_token_internal_error",
			expectedCode: http.StatusInternalServerError,
			user: &user.User{
				ID:			 "1",
				Username: "mario",
				Email:    "mario@example.com",
				Password: "hashedpassword",
				Name:     "Mario",
				Surname:  "Rossi",
			},
			setupMock: func(m *mocks.MockauthService) {
				m.EXPECT().
					GenerateJWT("1").
					Return("dummy_jwt_token", nil)
				m.EXPECT().
					GenerateRefreshToken(gomock.Any(), "1").
					Return(nil, fmt.Errorf("some error"))
			},
		},
		{
			name: 			 "generate_refresh_token_context_cancelled",
			expectedCode: http.StatusOK,
			user: &user.User{
				ID:			 "1",
				Username: "mario",
				Email:    "mario@example.com",
				Password: "hashedpassword",
				Name:     "Mario",
				Surname:  "Rossi",
			},
			setupMock: func(m *mocks.MockauthService) {
				m.EXPECT().
					GenerateJWT("1").
					Return("dummy_jwt_token", nil)
				m.EXPECT().
					GenerateRefreshToken(gomock.Any(), "1").
					Return(nil, context.Canceled)
			},
		},
		{
			name: 			 "generate_refresh_token_deadline_exceeded",
			expectedCode: http.StatusRequestTimeout,
			user: &user.User{
				ID:			 "1",
				Username: "mario",
				Email:    "mario@example.com",
				Password: "hashedpassword",
				Name:     "Mario",
				Surname:  "Rossi",
			},
			setupMock: func(m *mocks.MockauthService) {
				m.EXPECT().
					GenerateJWT("1").
					Return("dummy_jwt_token", nil)
				m.EXPECT().
					GenerateRefreshToken(gomock.Any(), "1").
					Return(nil, context.DeadlineExceeded)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockAuthService := mocks.NewMockauthService(ctrl)
			if tt.setupMock != nil {
				tt.setupMock(mockAuthService)
			}

			uc := NewAuthController(mockAuthService)

			c, w := newTestContext(http.MethodPost, "/login/username", nil)

			uc.handleSuccessfulLogin(c, tt.user)

			if w.Code != tt.expectedCode {
				t.Fatalf("got %d want %d; body=%s", w.Code, tt.expectedCode, w.Body.String())
			}
		})
	}
}
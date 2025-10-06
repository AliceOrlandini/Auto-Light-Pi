package auth

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AliceOrlandini/Auto-Light-Pi/internal/auth/mocks"
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

func TestRegister_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mocks.NewMockauthService(ctrl)
	mockAuthService.
		EXPECT().
		Register(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(nil)

	uc := &Controller{service: mockAuthService}
	
	body := []byte(`{"username":"mario","email":"mariorossi@gmail.com","password":"Testtest123","name":"mario","surname":"rossi"}`)
	c, w := newTestContext(http.MethodPost, "/register", body)

	uc.Register(c)

	if w.Code != http.StatusCreated {
		t.Fatalf("got %d want %d; body=%s", w.Code, http.StatusCreated, w.Body.String())
	}
}

func TestRegister_BadPayloads(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mocks.NewMockauthService(ctrl)
	uc := &Controller{service: mockAuthService}

	tests := []struct {
		name string
		body string
	}{
		{"missing_username", `{"email":"mariorossi@gmail.com","password":"Testtest123","name":"mario","surname":"rossi"}`},
		{"missing_email",    `{"username":"mario","password":"Testtest123","name":"mario","surname":"rossi"}`},
		{"missing_password", `{"username":"mario","email":"mariorossi@gmail.com","name":"mario","surname":"rossi"}`},
		{"missing_name",     `{"username":"mario","email":"mariorossi@gmail.com","password":"Testtest123","surname":"rossi"}`},
		{"missing_surname",  `{"username":"mario","email":"mariorossi@gmail.com","password":"Testtest123","name":"mario"}`},
		{"empty_username",   `{"username":"","email":"mariorossi@gmail.com","password":"Testtest123","name":"mario","surname":"rossi"}`},
		{"empty_email",      `{"username":"mario","email":"","password":"Testtest123","name":"mario","surname":"rossi"}`},
		{"empty_password",   `{"username":"mario","email":"mariorossi@gmail.com","password":"","name":"mario","surname":"rossi"}`},
		{"empty_name",       `{"username":"mario","email":"mariorossi@gmail.com","password":"Testtest123","name":"","surname":"rossi"}`},
		{"empty_surname",    `{"username":"mario","email":"mariorossi@gmail.com","password":"Testtest123","name":"mario","surname":""}`},
		{"invalid_json",     `{"username":123}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := []byte(tt.body)
			c, w := newTestContext(http.MethodPost, "/register", body)

			uc.Register(c)

			if w.Code != http.StatusBadRequest {
				t.Fatalf("got %d want %d; body=%s", w.Code, http.StatusBadRequest, w.Body.String())
			}
		})
	}
}

func TestRegister_InvalidPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mocks.NewMockauthService(ctrl)
	uc := &Controller{service: mockAuthService}

	tests := []struct {
		name     string
		password string
	}{
		{"short_password", "short"},
		{"long_password", "thisisaverylongpasswordthatexceedsthelimit"},
		{"no_number", "PasswordWithoutNumber"},
		{"no_uppercase", "passwordwithoutuppercase1"},
		{"no_lowercase", "PASSWORDWITHOUTLOWERCASE1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := []byte(`{"username":"mario","email":"mariorossi@gmail.com","password":"` + tt.password + `","name":"mario","surname":"rossi"}`)
			c, w := newTestContext(http.MethodPost, "/register", body)

			uc.Register(c)

			if w.Code != http.StatusBadRequest {
				t.Fatalf("got %d want %d; body=%s", w.Code, http.StatusBadRequest, w.Body.String())
			}
		})
	}
}
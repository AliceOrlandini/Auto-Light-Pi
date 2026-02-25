package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func init() { gin.SetMode(gin.TestMode) }

func TestAuthMiddleware(t *testing.T) {
	// Set the environment variable for the secret
	os.Setenv("JWT_SECRET", "supersecret")

	// Helper function to create a token for testing
	createToken := func(userID string, secret string, expired bool, method jwt.SigningMethod) string {
		claims := jwt.MapClaims{
			"sub": userID,
			"exp": time.Now().Add(time.Hour).Unix(),
		}
		if expired {
			claims["exp"] = time.Now().Add(-time.Hour).Unix()
		}

		token := jwt.NewWithClaims(method, claims)
		signedString, _ := token.SignedString([]byte(secret))
		return signedString
	}

	tests := []struct {
		name           string
		authHeader     string
		setupHeader    func() string // Function to dynamically generate header (e.g., using createToken)
		expectedStatus int
		expectedUserID string
	}{
		{
			name: "success",
			setupHeader: func() string {
				return "Bearer " + createToken("user123", "supersecret", false, jwt.SigningMethodHS256)
			},
			expectedStatus: http.StatusOK,
			expectedUserID: "user123",
		},
		{
			name:           "missing_auth_header",
			authHeader:     "", // Empty header
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "invalid_auth_header_format_no_bearer",
			authHeader:     "Basic user:pass",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "invalid_auth_header_format_malformed_bearer",
			authHeader:     "Bearer", // Missing token part
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "invalid_token_string",
			authHeader:     "Bearer invalid.token.string",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "wrong_secret",
			setupHeader: func() string {
				return "Bearer " + createToken("user123", "wrongsecret", false, jwt.SigningMethodHS256)
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "expired_token",
			setupHeader: func() string {
				return "Bearer " + createToken("user123", "supersecret", true, jwt.SigningMethodHS256)
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "wrong_signing_method_none",
			setupHeader: func() string {
				// Header: {"alg":"RS256","typ":"JWT"}
				header := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9" 
				// Payload: {"sub":"user123"}
				payload := "eyJzdWIiOiJ1c2VyMTIzIn0"
				// Signature: garbage
				signature := "garbage"
				return "Bearer " + header + "." + payload + "." + signature
			},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize the middleware again to pick up the env var if needed, 
			// though typical implementation reads it once.
			// If AuthMiddleware reads os.Getenv outside the returned function, 
			// it must be called after Setenv.
			middleware := AuthMiddleware()

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			req, _ := http.NewRequest("GET", "/", nil)

			headerValue := tt.authHeader
			if tt.setupHeader != nil {
				headerValue = tt.setupHeader()
			}

			if headerValue != "" {
				req.Header.Set("Authorization", headerValue)
			}
			c.Request = req

			middleware(c)

			if w.Code != tt.expectedStatus {
				t.Errorf("test %s: expected status %d, got %d", tt.name, tt.expectedStatus, w.Code)
			}

			if tt.expectedStatus == http.StatusOK {
				userID, exists := c.Get("userID")
				if !exists {
					t.Errorf("test %s: expected userID in context, got none", tt.name)
				}
				if userID != tt.expectedUserID {
					t.Errorf("test %s: expected userID %s, got %v", tt.name, tt.expectedUserID, userID)
				}
			}
		})
	}
}



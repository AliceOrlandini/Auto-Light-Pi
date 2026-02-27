package routes

import (
	"bytes"
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/AliceOrlandini/Auto-Light-Pi/internal/auth"
	"github.com/AliceOrlandini/Auto-Light-Pi/internal/refresh_token"
	"github.com/AliceOrlandini/Auto-Light-Pi/internal/testutils"
	"github.com/AliceOrlandini/Auto-Light-Pi/internal/user"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

func init() { gin.SetMode(gin.TestMode) }

var testPostgresDB *sql.DB
var testRedisDB *redis.Client

func TestMain(m *testing.M) {
	pgConnectionStr := testutils.SetupPostgres()
	testPostgresDB, _ = sql.Open("postgres", pgConnectionStr)
	redisConnectionStr := testutils.SetupRedis()
	opt, _ := redis.ParseURL(redisConnectionStr)
	testRedisDB = redis.NewClient(opt)
	os.Exit(m.Run())
}

func TestIntegrationRoutes(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()

	// Setup env vars for JWT
	// Note: We must set this before SetupRoutes is called, because AuthMiddleware reads 
	// the environment variable during initialization.
	os.Setenv("JWT_SECRET", "supersecret")
	os.Setenv("APPLICATION_NAME", "TestApp")

	// Initialize real repositories and services
	userRepo := user.NewUserRepository(testPostgresDB)
	rtRepo := refresh_token.NewRefreshTokenRepository(testRedisDB)
	authService := auth.NewAuthService(userRepo, rtRepo)
	authController := auth.NewAuthController(authService)

	router := SetupRoutes(authController)

	// Helper to create valid token for auth middleware tests
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
		method         string
		path           string
		body           string
		setupData      func(context.Context) error
		setupRequest   func(*http.Request)
		expectedStatus int
		verifyResponse func(*testing.T, *httptest.ResponseRecorder)
		checkDBDataPresence func(context.Context) (bool, error)
	}{
		{
			name:   "register_route",
			method: "POST",
			path:   "/api/register",
			body:   `{"username":"mario","email":"mariorossi@gmail.com","password":"Testtest123","name":"mario","surname":"rossi"}`,
			setupData: func(ctx context.Context) error {
				// This ensure user is fresh
				_, err := testPostgresDB.ExecContext(ctx, "DELETE FROM user_account WHERE email = $1", "mariorossi@gmail.com")
				return err
			},
			expectedStatus: http.StatusCreated,
			checkDBDataPresence: func(ctx context.Context) (bool, error) {
				var exists bool
				err := testPostgresDB.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM user_account WHERE email = $1)", "mariorossi@gmail.com").Scan(&exists)
				return exists, err
			},
		},
		{
			name:   "login_email_route",
			method: "POST",
			path:   "/api/login/email",
			body:   `{"email":"luigi@gmail.com","password":"Testtest123"}`,
			setupData: func(ctx context.Context) error {
				// Ensure clean state
				testPostgresDB.ExecContext(ctx, "DELETE FROM user_account WHERE email = $1", "luigi@gmail.com")
				// Create user via service to ensure password hashing
				return authService.Register(ctx, "luigi", "luigi@gmail.com", "Testtest123", "luigi", "verdi")
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "login_username_route",
			method: "POST",
			path:   "/api/login/username",
			body:   `{"username":"peach","password":"Testtest123"}`,
			setupData: func(ctx context.Context) error {
				testPostgresDB.ExecContext(ctx, "DELETE FROM user_account WHERE username = $1", "peach")
				return authService.Register(ctx, "peach", "peach@gmail.com", "Testtest123", "peach", "princess")
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "ping_route_unauthorized",
			method:         "GET",
			path:           "/api/ping",
			setupRequest: 	func(req *http.Request) {
				req.Header.Set("Authorization", "Bearer " + createToken("user-uuid-123", "supersecret", true, jwt.SigningMethodHS256))
			},
			setupData:      func(ctx context.Context) error { return nil },
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:   "ping_route_authorized",
			method: "GET",
			path:   "/api/ping",
			setupData: func(ctx context.Context) error { return nil },
			setupRequest: func(req *http.Request) {
				req.Header.Set("Authorization", "Bearer " + createToken("user-uuid-123", "supersecret", false, jwt.SigningMethodHS256))
			},
			expectedStatus: http.StatusOK,
			verifyResponse: func(t *testing.T, w *httptest.ResponseRecorder) {
				if !strings.Contains(w.Body.String(), "Hello, user user-uuid-123") {
					t.Errorf("expected body to contain 'Hello, user user-uuid-123', got %s", w.Body.String())
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupData != nil {
				if err := tt.setupData(ctx); err != nil {
					t.Fatalf("setupData failed: %v", err)
				}
			}

			w := httptest.NewRecorder()
			var body *bytes.Buffer
			if tt.body != "" {
				body = bytes.NewBufferString(tt.body)
			} else {
				body = bytes.NewBuffer(nil)
			}
			req, _ := http.NewRequest(tt.method, tt.path, body)
			req.Header.Set("Content-Type", "application/json")

			if tt.setupRequest != nil {
				tt.setupRequest(req)
			}

			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d. Body: %s", tt.expectedStatus, w.Code, w.Body.String())
			}

			if tt.checkDBDataPresence != nil {
				exists, err := tt.checkDBDataPresence(ctx)
				if err != nil {
					t.Fatalf("checkDBDataPresence failed: %v", err)
				}
				if !exists {
					t.Errorf("expected data presence check to be true, got false")
				}
			}

			if tt.verifyResponse != nil {
				tt.verifyResponse(t, w)
			}
		})
	}
}

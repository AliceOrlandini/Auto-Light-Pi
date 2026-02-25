package routes

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/AliceOrlandini/Auto-Light-Pi/internal/auth"
	"github.com/AliceOrlandini/Auto-Light-Pi/internal/refresh_token"
	"github.com/AliceOrlandini/Auto-Light-Pi/internal/user"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	redisContainer "github.com/testcontainers/testcontainers-go/modules/redis"
)

func init() { gin.SetMode(gin.TestMode) }

func getRootPath() string {
    _, b, _, _ := runtime.Caller(0)
    return filepath.Join(filepath.Dir(b), "..", "..") 
}

func setupContainers(ctx context.Context) (*sql.DB, *redis.Client, func(), error) {
	// Postgres Container
	pgContainer, err := postgres.Run(ctx,
		"postgres:18-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("user"),
		postgres.WithPassword("password"),
		postgres.WithInitScripts(filepath.Join(getRootPath(), "testdata", "schema.sql")),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to start postgres container: %w", err)
	}

	pgConnectionString, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		pgContainer.Terminate(ctx)
		return nil, nil, nil, fmt.Errorf("failed to get postgres connection string: %w", err)
	}

	pgDB, err := sql.Open("postgres", pgConnectionString)
	if err != nil {
		pgContainer.Terminate(ctx)
		return nil, nil, nil, fmt.Errorf("failed to open postgres connection: %w", err)
	}

	// Redis Container
	rContainer, err := redisContainer.Run(ctx,
		"redis:alpine",
		redisContainer.WithSnapshotting(0, 0),
		redisContainer.WithLogLevel(redisContainer.LogLevelVerbose),
	)
	if err != nil {
		pgDB.Close()
		pgContainer.Terminate(ctx)
		return nil, nil, nil, fmt.Errorf("failed to start redis container: %w", err)
	}

	rConnectionString, err := rContainer.ConnectionString(ctx)
	if err != nil {
		pgDB.Close()
		pgContainer.Terminate(ctx)
		rContainer.Terminate(ctx)
		return nil, nil, nil, fmt.Errorf("failed to get redis connection string: %w", err)
	}

	opt, err := redis.ParseURL(rConnectionString)
	if err != nil {
		pgDB.Close()
		pgContainer.Terminate(ctx)
		rContainer.Terminate(ctx)
		return nil, nil, nil, fmt.Errorf("failed to parse redis url: %w", err)
	}
	rdb := redis.NewClient(opt)

	cleanup := func() {
		pgDB.Close()
		pgContainer.Terminate(ctx)
		rdb.Close()
		rContainer.Terminate(ctx)
	}

	return pgDB, rdb, cleanup, nil
}

func TestIntegrationRoutes(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()
	pgDB, rdb, cleanup, err := setupContainers(ctx)
	if err != nil {
		t.Fatalf("failed to setup containers: %v", err)
	}
	defer cleanup()

	// Setup env vars for JWT
	// Note: We must set this before SetupRoutes is called, because AuthMiddleware reads 
	// the environment variable during initialization.
	os.Setenv("JWT_SECRET", "supersecret")
	os.Setenv("APPLICATION_NAME", "TestApp")

	// Initialize real repositories and services
	userRepo := user.NewUserRepository(pgDB)
	rtRepo := refresh_token.NewRefreshTokenRepository(rdb)
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
				_, err := pgDB.ExecContext(ctx, "DELETE FROM user_account WHERE email = $1", "mariorossi@gmail.com")
				return err
			},
			expectedStatus: http.StatusCreated,
			checkDBDataPresence: func(ctx context.Context) (bool, error) {
				var exists bool
				err := pgDB.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM user_account WHERE email = $1)", "mariorossi@gmail.com").Scan(&exists)
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
				pgDB.ExecContext(ctx, "DELETE FROM user_account WHERE email = $1", "luigi@gmail.com")
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
				pgDB.ExecContext(ctx, "DELETE FROM user_account WHERE username = $1", "peach")
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

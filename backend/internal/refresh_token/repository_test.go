package refresh_token

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/AliceOrlandini/Auto-Light-Pi/internal/testutils"
	"github.com/redis/go-redis/v9"
)

var testRedisDB *redis.Client

func TestMain(m *testing.M) {
	redisConnectionStr := testutils.SetupRedis()
	opt, _ := redis.ParseURL(redisConnectionStr)
	testRedisDB = redis.NewClient(opt)

	os.Exit(m.Run())
}

func TestRepository_CreateOne(t *testing.T) {
	ctx := context.Background()
	repo := NewRefreshTokenRepository(testRedisDB)

	tests := []struct {
		name         string
		setupFunc    func(ctx context.Context, rdb *redis.Client) error
		refreshToken RefreshToken
		verifyFunc   func(ctx context.Context, rdb *redis.Client, rt RefreshToken) error
		wantErr      bool
	}{
		{
			name: "create_new_token_success",
			refreshToken: RefreshToken{
				RefreshToken:     "hash_success",
				RefreshTokenHash: "hash_success",
				UserID:           "user_success",
				TTL:              time.Now().Add(1 * time.Hour),
			},
			verifyFunc: func(ctx context.Context, rdb *redis.Client, rt RefreshToken) error {
				// Verify lookup hash -> userID
				userID, err := rdb.Get(ctx, "rth:"+rt.RefreshTokenHash).Result()
				if err != nil {
					return fmt.Errorf("lookup by hash failed: %w", err)
				}
				if userID != rt.UserID {
					return fmt.Errorf("expected userID %s, got %s", rt.UserID, userID)
				}

				// Verify lookup userID -> hash
				hash, err := rdb.Get(ctx, "rtu:"+rt.UserID).Result()
				if err != nil {
					return fmt.Errorf("lookup by userID failed: %w", err)
				}
				if hash != rt.RefreshTokenHash {
					return fmt.Errorf("expected hash %s, got %s", rt.RefreshTokenHash, hash)
				}
				return nil
			},
			wantErr: false,
		},
		{
			name: "rotate_token_removes_old_one",
			setupFunc: func(ctx context.Context, rdb *redis.Client) error {
				// Pre-populate an old token for this user
				if err := rdb.Set(ctx, "rtu:user_rotate", "old_hash", time.Hour).Err(); err != nil {
					return err
				}
				if err := rdb.Set(ctx, "rth:old_hash", "user_rotate", time.Hour).Err(); err != nil {
					return err
				}
				return nil
			},
			refreshToken: RefreshToken{
				RefreshToken:     "new_hash",
				RefreshTokenHash: "new_hash",
				UserID:           "user_rotate",
				TTL:              time.Now().Add(1 * time.Hour),
			},
			verifyFunc: func(ctx context.Context, rdb *redis.Client, rt RefreshToken) error {
				// New token should exist
				newHash, err := rdb.Get(ctx, "rtu:"+rt.UserID).Result()
				if err != nil {
					return fmt.Errorf("failed to get new token: %w", err)
				}
				if newHash != rt.RefreshTokenHash {
					return fmt.Errorf("expected new hash %s, got %s", rt.RefreshTokenHash, newHash)
				}

				// Old token should be gone
				_, err = rdb.Get(ctx, "rth:old_hash").Result()
				if err != redis.Nil {
					return fmt.Errorf("expected old token hash key to be deleted, but got: %v", err)
				}
				return nil
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean DB before each test
			testRedisDB.FlushDB(ctx)

			if tt.setupFunc != nil {
				if err := tt.setupFunc(ctx, testRedisDB); err != nil {
					t.Fatalf("setupFunc failed: %v", err)
				}
			}

			err := repo.CreateOne(ctx, &tt.refreshToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateOne() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.verifyFunc != nil && err == nil {
				if err := tt.verifyFunc(ctx, testRedisDB, tt.refreshToken); err != nil {
					t.Errorf("verification failed: %v", err)
				}
			}
		})
	}
}

func TestRepository_GetOneUserIDByTokenHash(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	ctx := context.Background()
	repo := NewRefreshTokenRepository(testRedisDB)

	tests := []struct {
		name      string
		tokenHash string
		setupFunc func()
		want      string
		wantErr   error
	}{
		{
			name:      "found",
			tokenHash: "hash_found",
			setupFunc: func() {
				testRedisDB.Set(ctx, "rth:hash_found", "user_found", time.Hour)
			},
			want:    "user_found",
			wantErr: nil,
		},
		{
			name:      "not_found",
			tokenHash: "hash_missing",
			setupFunc: func() {},
			want:      "",
			wantErr:   ErrTokenHashNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testRedisDB.FlushDB(ctx)
			tt.setupFunc()

			got, err := repo.GetOneUserIDByTokenHash(ctx, tt.tokenHash)
			if err != tt.wantErr {
				t.Errorf("GetOneUserIDByTokenHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetOneUserIDByTokenHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepository_GetOneTokenHashByUserID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	ctx := context.Background()
	repo := NewRefreshTokenRepository(testRedisDB)

	tests := []struct {
		name      string
		userID    string
		setupFunc func()
		want      string
		wantErr   error
	}{
		{
			name:   "found",
			userID: "user_found",
			setupFunc: func() {
				testRedisDB.Set(ctx, "rtu:user_found", "hash_found", time.Hour)
			},
			want:    "hash_found",
			wantErr: nil,
		},
		{
			name:      "not_found",
			userID:    "user_missing",
			setupFunc: func() {},
			want:      "",
			wantErr:   ErrUserIDNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testRedisDB.FlushDB(ctx)
			tt.setupFunc()

			got, err := repo.GetOneTokenHashByUserID(ctx, tt.userID)
			if err != tt.wantErr {
				t.Errorf("GetOneTokenHashByUserID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetOneTokenHashByUserID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepository_DeleteOneByUserID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	ctx := context.Background()
	repo := NewRefreshTokenRepository(testRedisDB)

	tests := []struct {
		name       string
		userID     string
		setupFunc  func()
		verifyFunc func() error
	}{
		{
			name:   "delete_success",
			userID: "user_delete",
			setupFunc: func() {
				// Setup a full valid token pair
				testRedisDB.Set(ctx, "rtu:user_delete", "hash_delete", time.Hour)
				testRedisDB.Set(ctx, "rth:hash_delete", "user_delete", time.Hour)
			},
			verifyFunc: func() error {
				// Verify user mapping is gone
				_, err := testRedisDB.Get(ctx, "rtu:user_delete").Result()
				if err != redis.Nil {
					return fmt.Errorf("expected rtu key to be deleted, got %v", err)
				}
				// Verify hash mapping is gone
				_, err = testRedisDB.Get(ctx, "rth:hash_delete").Result()
				if err != redis.Nil {
					return fmt.Errorf("expected rth key to be deleted, got %v", err)
				}
				return nil
			},
		},
		{
			name:      "delete_non_existent",
			userID:    "user_ghost",
			setupFunc: func() {},
			verifyFunc: func() error {
				return nil // Should just succeed without error
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testRedisDB.FlushDB(ctx)
			tt.setupFunc()

			if err := repo.DeleteOneByUserID(ctx, tt.userID); err != nil {
				t.Errorf("DeleteOneByUserID() error = %v", err)
			}

			if tt.verifyFunc != nil {
				if err := tt.verifyFunc(); err != nil {
					t.Errorf("verification failed: %v", err)
				}
			}
		})
	}
}

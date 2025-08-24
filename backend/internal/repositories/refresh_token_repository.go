package repositories

import (
	"context"
	"errors"

	"github.com/AliceOrlandini/Auto-Light-Pi/internal/models"
	"github.com/redis/go-redis/v9"
)

var (
	ErrUserIDNotFound  = errors.New("user ID not found")
	ErrTokenHashNotFound = errors.New("token hash not found")
)

type RefreshTokenRepository struct {
	db *redis.Client
}

func NewRefreshTokenRepository(db *redis.Client) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

func (r *RefreshTokenRepository) CreateOne(ctx context.Context, refreshToken *models.RefreshToken) error {
	// we will save two records:
	// 1. rth:{tokenHash} -> userId used as lookup when a new request comes in
	// 2. rtu:{userId} -> {tokenHash} used for the revocation of the token

	rtKey := "rth:" + refreshToken.RefreshTokenHash
	rtuKey := "rtu:" + refreshToken.UserID

	// we use a Lua script since we need to perform multiple operations atomically
	lua := `
		local rtKey = KEYS[1]
		local ruKey = KEYS[2]
		local userId = ARGV[1]
		local tokenHash = ARGV[2]
		local pxat   = tonumber(ARGV[3])

		-- We need first to retreive the old token by lookup for the userId
		local oldTokenHash = redis.call("GET", ruKey)
		-- If there is an old token, we need to delete it
		if oldTokenHash then
			local oldrtKey = "rth:" .. oldTokenHash
			redis.call("DEL", oldrtKey)
			redis.call("DEL", ruKey)
		end

		-- Insert new values
		redis.call("SET", rtKey, userId, "PXAT", pxat)
		redis.call("SET", ruKey, tokenHash, "PXAT", pxat)

		return 1
	`

	pxat := refreshToken.TTL.UnixMilli()
	_, err := r.db.Eval(ctx, lua, []string{rtKey, rtuKey}, refreshToken.UserID, refreshToken.RefreshTokenHash, pxat).Int64()
	if err != nil {
		// if there is a Redis error, we consider it as an internal server error
		return err
	}
	return nil
}

func (r *RefreshTokenRepository) GetOneUserIDByTokenHash(ctx context.Context, tokenHash string) (string, error) {
	rtKey := "rth:" + tokenHash
	userId, err := r.db.Get(ctx, rtKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", ErrTokenHashNotFound
		}
		return "", err
	}
	return userId, nil
}

func (r *RefreshTokenRepository) GetOneTokenHashByUserID(ctx context.Context, userID string) (string, error) {
	rtuKey := "rtu:" + userID
	tokenHash, err := r.db.Get(ctx, rtuKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", ErrUserIDNotFound
		}
		return "", err
	}
	return tokenHash, nil
}

func (r *RefreshTokenRepository) RevokeOneByUserID(ctx context.Context, userID string) error {
	// we need to remove two records, 
	// the second is obtained using the value from the first
	rtuKey := "rtu:" + userID

	lua := `
		local rtuKey = KEYS[1]
		local tokenHash = redis.call("GET", rtuKey)
	
		-- If tokenHash exists than we generate rtKey and delete it
		if tokenHash then
				local rtKey = "rth:" .. tokenHash
				redis.call("DEL", rtKey)
		end

		-- Then we delete the rtuKey
		redis.call("DEL", rtuKey)

		-- This function is idempotent, so it is always a success
		return 1
	`

	_, err := r.db.Eval(ctx, lua, []string{rtuKey}, userID).Result()
	if err != nil {
		return err
	}
	
	return nil
}
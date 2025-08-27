package refresh_token

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ErrUserIDNotFound  = errors.New("user ID not found")
	ErrTokenHashNotFound = errors.New("token hash not found")
)

type refreshTokenEntity struct {
	RefreshTokenHash string
	UserID           string
	TTL              time.Time
}

type repository struct {
	db *redis.Client
}

func NewRefreshTokenRepository(db *redis.Client) *repository {
	return &repository{db: db}
}

func (r *repository) CreateOne(ctx context.Context, refreshToken *RefreshToken) error {
	// we will save two records:
	// 1. rth:{tokenHash} -> userID used as lookup when a new request comes in
	// 2. rtu:{userID} -> {tokenHash} used for the revocation of the token

	refreshTokenEntity := toEntity(refreshToken)

	rtKey := "rth:" + refreshTokenEntity.RefreshTokenHash
	rtuKey := "rtu:" + refreshTokenEntity.UserID

	// we use a Lua script since we need to perform multiple operations atomically
	lua := `
		local rtKey = KEYS[1]
		local ruKey = KEYS[2]
		local userID = ARGV[1]
		local tokenHash = ARGV[2]
		local pxat   = tonumber(ARGV[3])

		-- We need first to retreive the old token by lookup for the userID
		local oldTokenHash = redis.call("GET", ruKey)
		-- If there is an old token, we need to delete it
		if oldTokenHash then
			local oldrtKey = "rth:" .. oldTokenHash
			redis.call("DEL", oldrtKey)
			redis.call("DEL", ruKey)
		end

		-- Insert new values
		redis.call("SET", rtKey, userID, "PXAT", pxat)
		redis.call("SET", ruKey, tokenHash, "PXAT", pxat)

		return 1
	`

	pxat := refreshTokenEntity.TTL.UnixMilli()
	_, err := r.db.Eval(ctx, lua, []string{rtKey, rtuKey}, refreshTokenEntity.UserID, refreshTokenEntity.RefreshTokenHash, pxat).Int64()
	if err != nil {
		// if there is a Redis error, we consider it as an internal server error
		return err
	}
	return nil
}

func (r *repository) GetOneUserIDByTokenHash(ctx context.Context, tokenHash string) (string, error) {
	rtKey := "rth:" + tokenHash
	userID, err := r.db.Get(ctx, rtKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", ErrTokenHashNotFound
		}
		return "", err
	}
	return userID, nil
}

func (r *repository) GetOneTokenHashByUserID(ctx context.Context, userID string) (string, error) {
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

func (r *repository) DeleteOneByUserID(ctx context.Context, userID string) error {
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

func toEntity(rt *RefreshToken) *refreshTokenEntity {
	return &refreshTokenEntity{
		RefreshTokenHash: rt.RefreshTokenHash,
		UserID:           rt.UserID,
		TTL:              rt.TTL,
	}
}
package refresh_token

import (
	"time"
)

type RefreshToken struct {
	RefreshToken     string
	RefreshTokenHash string
	UserID           string
	TTL              time.Time
}
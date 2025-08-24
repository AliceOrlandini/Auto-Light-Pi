package models

import (
	"time"
)

type RefreshToken struct {
	RefreshToken     string    `json:"refresh_token"`
	RefreshTokenHash string    `json:"-"`
	UserID           string    `json:"-"`
	TTL              time.Time `json:"-"`
	CreatedAt        time.Time `json:"-"`
}
package models

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"time"
)

const (
	ScopeAuthentication = "authentication"
)

// Token is the type for authentication tokens
type Token struct {
	PlaintText string    `json:"token"`
	UserID     int64     `json:"-"`
	Hash       []byte    `json:"-"`
	Expiry     time.Time `json:"expiry"`
	Scope      string    `json:"-"`
}

// GenerateToken generates a token that for ttl and returns it
func GenerateToken(userID int, ttl time.Duration, scope string) (*Token, error) {
	token := &Token{
		UserID: int64(userID),
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}

	randomBytes := make([]byte, 16)

	_, err := rand.Read(randomBytes)

	if err != nil {
		return nil, err
	}

	token.PlaintText = base64.StdEncoding.WithPadding(base64.NoPadding).EncodeToString(randomBytes)

	hash := sha256.Sum256([]byte(token.PlaintText))

	token.Hash = hash[:]

	return token, nil
}

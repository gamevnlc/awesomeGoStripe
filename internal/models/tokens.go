package models

import (
	"context"
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

func (m *DBModel) InsertToken(t *Token, u User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var newTokenId int
	//goland:noinspection ALL
	stmt := `
		insert into tokens (user_id, name, email, token_hash, created_at, updated_at)
		values ($1, $2, $3, $4, $5, $6) returning id
	`
	err := m.DB.QueryRowContext(ctx, stmt,
		u.ID,
		u.LastName,
		u.Email,
		t.Hash,
		time.Now(),
		time.Now(),
	).Scan(&newTokenId)
	if err != nil {
		return err
	}

	return nil
}

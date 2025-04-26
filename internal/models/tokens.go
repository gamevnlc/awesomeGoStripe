package models

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
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

	token.PlaintText = base32.StdEncoding.WithPadding(base64.NoPadding).EncodeToString(randomBytes)

	hash := sha256.Sum256([]byte(token.PlaintText))

	token.Hash = hash[:]

	return token, nil
}

func (m *DBModel) InsertToken(t *Token, u User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	//delete existing tokens

	//goland:noinspection ALL
	query := `delete from tokens where user_id = $1`
	_, err := m.DB.ExecContext(ctx, query, u.ID)
	if err != nil {
		return err
	}

	var newTokenId int
	//goland:noinspection ALL
	stmt := `
		insert into tokens (user_id, name, email, token_hash, expiry, created_at, updated_at)
		values ($1, $2, $3, $4, $5, $6, $7) returning id
	`
	err = m.DB.QueryRowContext(ctx, stmt,
		u.ID,
		u.LastName,
		u.Email,
		t.Hash,
		t.Expiry,
		time.Now(),
		time.Now(),
	).Scan(&newTokenId)
	if err != nil {
		return err
	}

	return nil
}

func (m *DBModel) GetUserForToken(token string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tokenHash := sha256.Sum256([]byte(token))

	var user User

	//goland:noinspection ALL
	query := `
		select 
		    u.id, u.first_name, u.last_name, u.email
		from users u
		inner join tokens t on u.id = t.user_id
		where 
		    t.token_hash = $1
			and t.expiry > $2
	`

	err := m.DB.QueryRowContext(ctx, query, tokenHash[:], time.Now()).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

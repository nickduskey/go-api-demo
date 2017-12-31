package models

import (
	"database/sql"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/nickduskey/api-demo/utils"
)

// TokenClaims represents a jwt object with claims
// https://tools.ietf.org/html/rfc7519
type TokenClaims struct {
	jwt.StandardClaims
	Role string `json:"role"`
	Csrf string `json:"csrf"`
}

// RefreshToken represents a refresh token stored in the db
type RefreshToken struct {
	ID  int
	Jti string
}

// RefreshTokenValidTime is the expiration time of refresh tokens
const RefreshTokenValidTime = time.Hour * 72

// AuthTokenValidTime is the expiration time of auth tokens
const AuthTokenValidTime = time.Minute * 15

// GenerateCSRFSecret generates a random 32 character string
// to be used as a CSRF secret in the JWT
func GenerateCSRFSecret() (string, error) {
	return utils.GenerateRandomString(32)
}

// GetRefreshTokenByJTI retrieves a token by its JTI
func (t *RefreshToken) GetRefreshTokenByJTI(db *sql.DB) error {
	return db.QueryRow("SELECT jti FROM refresh_tokens WHERE jti=$1",
		t.Jti).Scan(&t.Jti)
}

// StoreRefreshToken generates a random JTI and then stores it in the db
func StoreRefreshToken(db *sql.DB) (jti string, err error) {
	jti, err = utils.GenerateRandomString(32)
	if err != nil {
		return jti, err
	}

	// check to make sure our jti is unique
	// TODO

	// insert refreshToken into db
	t := RefreshToken{Jti: jti}
	err = t.CreateRefreshToken(db)

	if err != nil {
		return jti, err
	}

	return jti, nil
}

// CreateRefreshToken stores a token in the db
func (t *RefreshToken) CreateRefreshToken(db *sql.DB) error {
	return db.QueryRow(
		"INSERT INTO refresh_tokens(jti) VALUES($1) RETURNING id",
		t.Jti).Scan(&t.ID)
}

// DeleteRefreshToken revokes a token from the db
func (t *RefreshToken) DeleteRefreshToken(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM refresh_tokens WHERE id=$1", t.ID)

	return err
}

// CheckRefreshToken checks the refresh token
func CheckRefreshToken(jti string, db *sql.DB) bool {
	// Fetch Refresh Token from the db
	t := RefreshToken{Jti: jti}
	if err := t.GetRefreshTokenByJTI(db); err != nil {
		return false
	}

	return t.Jti != ""
}

// RevokeRefreshToken revokes a refresh token stored in the db
func RevokeRefreshToken(jti string, db *sql.DB) error {
	t := RefreshToken{Jti: jti}
	err := t.DeleteRefreshToken(db)
	if err != nil {
		return err
	}
	return nil
}

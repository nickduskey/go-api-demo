package auth

import (
	"crypto/rsa"
	"database/sql"
	"errors"
	"io/ioutil"
	"log"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/nickduskey/api-demo/models"
)

const (
	privKeyPath = "keys/app.rsa"
	pubKeyPath  = "keys/app.rsa.pub"
)

var (
	verifyKey *rsa.PublicKey
	signKey   *rsa.PrivateKey
)

// InitAuth initializes the signing keys
func InitAuth() error {
	signBytes, err := ioutil.ReadFile(privKeyPath)
	if err != nil {
		return err
	}

	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		return err
	}

	verifyBytes, err := ioutil.ReadFile(pubKeyPath)
	if err != nil {
		return err
	}

	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		return err
	}

	return nil
}

// CreateNewTokens generates new jwt
func CreateNewTokens(uuid, role string, db *sql.DB) (authTokenString, refreshTokenString, csrfSecret string, err error) {
	// generate the csrf secret
	csrfSecret, err = models.GenerateCSRFSecret() // TODO implement this
	if err != nil {
		return
	}

	// generate the refresh token
	refreshTokenString, err = createRefreshTokenString(uuid, role, csrfSecret, db)

	// generate the auth token
	authTokenString, err = createAuthTokenString(uuid, role, csrfSecret)
	if err != nil {
		return
	}

	return
}

// CheckAndRefreshTokens checks status and refreshes tokens
func CheckAndRefreshTokens(oldAuthTokenString, oldRefreshTokenString, oldCsrfSecret string, db *sql.DB) (newAuthTokenString, newRefreshTokenString, newCsrfSecret string, err error) {
	// first check that a csrf token was provided
	if oldCsrfSecret == "" {
		log.Println("No CSRF token!")
		err = errors.New("Unauthorized")
		return
	}

	// then check that it matches what is in the auth token claims
	authToken, err := jwt.ParseWithClaims(oldAuthTokenString, &models.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return verifyKey, nil
	})
	authTokenClaims, ok := authToken.Claims.(*models.TokenClaims)
	if !ok {
		return
	}
	if oldCsrfSecret != authTokenClaims.Csrf {
		log.Println("CSRF token doesn't match jwt!")
		err = errors.New("Unauthorized")
		return
	}

	// next, check the auth token in a stateless manner
	if authToken.Valid {
		log.Println("Auth token is valid")
		// auth token hasn't expired
		// we need to return the csrf secret because that is what the function calls for
		newCsrfSecret = authTokenClaims.Csrf

		// update the exp of the refresh token string, but don't save it to the db
		// we don't need to check if our refresh token is valid here
		// because we aren't renewing the auth token, the auth token is already valid
		newRefreshTokenString, err = updateRefreshTokenExp(oldRefreshTokenString)
		newAuthTokenString = oldAuthTokenString
		return
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		log.Println("Auth token is not valid")
		if ve.Errors&(jwt.ValidationErrorExpired) != 0 {
			log.Println("Auth token is expired")
			newAuthTokenString, newCsrfSecret, err = updateAuthTokenString(oldRefreshTokenString, oldAuthTokenString, db)
			if err != nil {
				return
			}

			// update the exp of the refresh token string
			newRefreshTokenString, err = updateRefreshTokenExp(oldRefreshTokenString)
			if err != nil {
				return
			}

			// update the csrf string of the refresh token
			newRefreshTokenString, err = updateRefreshTokenCsrf(newRefreshTokenString, newCsrfSecret)
			return
		}
		log.Println("Error in auth token")
		err = errors.New("Error in auth token")
		return
	} else {
		log.Println("Error in auth token")
		err = errors.New("Error in auth token")
		return
	}
}

func createAuthTokenString(uuid, role, csrfSecret string) (authTokenString string, err error) {
	authTokenExp := time.Now().Add(models.AuthTokenValidTime).Unix()
	authClaims := models.TokenClaims{
		jwt.StandardClaims{
			Subject:   uuid,
			ExpiresAt: authTokenExp,
		},
		role,
		csrfSecret,
	}

	// create a signer for rsa 256
	authJwt := jwt.NewWithClaims(jwt.GetSigningMethod("RS256"), authClaims)

	// generate the auth token string
	authTokenString, err = authJwt.SignedString(signKey)
	return
}

func createRefreshTokenString(uuid, role, csrfString string, db *sql.DB) (refreshTokenString string, err error) {
	refreshTokenExp := time.Now().Add(models.RefreshTokenValidTime).Unix()
	refreshJti, err := models.StoreRefreshToken(db) // pass this function the db
	if err != nil {
		return
	}

	refreshClaims := models.TokenClaims{
		jwt.StandardClaims{
			Id:        refreshJti, // jti
			Subject:   uuid,
			ExpiresAt: refreshTokenExp,
		},
		role,
		csrfString,
	}

	// create a signer for rsa 256
	refreshJwt := jwt.NewWithClaims(jwt.GetSigningMethod("RS256"), refreshClaims)

	// generate the refresh token string
	refreshTokenString, err = refreshJwt.SignedString(signKey)
	return
}

func updateRefreshTokenExp(oldRefreshTokenString string) (newRefreshTokenString string, err error) {
	refreshToken, err := jwt.ParseWithClaims(oldRefreshTokenString, &models.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return verifyKey, nil
	})

	oldRefreshTokenClaims, ok := refreshToken.Claims.(*models.TokenClaims)
	if !ok {
		return
	}

	refreshTokenExp := time.Now().Add(models.RefreshTokenValidTime).Unix()

	refreshClaims := models.TokenClaims{
		jwt.StandardClaims{
			Id:        oldRefreshTokenClaims.StandardClaims.Id, // jti
			Subject:   oldRefreshTokenClaims.StandardClaims.Subject,
			ExpiresAt: refreshTokenExp,
		},
		oldRefreshTokenClaims.Role,
		oldRefreshTokenClaims.Csrf,
	}

	// create a signer for rsa 256
	refreshJwt := jwt.NewWithClaims(jwt.GetSigningMethod("RS256"), refreshClaims)

	// generate the refresh token string
	newRefreshTokenString, err = refreshJwt.SignedString(signKey)
	return
}

func updateAuthTokenString(refreshTokenString, oldAuthTokenSTring string, db *sql.DB) (newAuthTokenString, csrfSecret string, err error) {
	refreshToken, err := jwt.ParseWithClaims(refreshTokenString, &models.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return verifyKey, nil
	})
	refreshTokenClaims, ok := refreshToken.Claims.(*models.TokenClaims)
	if !ok {
		err = errors.New("Error reading jwt claims")
		return
	}

	// check if the refresh token has been revoked
	if models.CheckRefreshToken(refreshTokenClaims.StandardClaims.Id, db) {
		// the refresh token has not been revoked
		// has it expired?
		if refreshToken.Valid {
			// nope the refresh token has not expired
			// issue a new auth token
			authToken, _ := jwt.ParseWithClaims(oldAuthTokenSTring, &models.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
				return verifyKey, nil
			})

			oldAuthTokenClaims, ok := authToken.Claims.(*models.TokenClaims)
			if !ok {
				err = errors.New("Error reading jwt claims")
				return
			}

			// our policy is to regenerate the csrf secret for each new auth token
			csrfSecret, err = models.GenerateCSRFSecret()
			if err != nil {
				return
			}

			newAuthTokenString, err = createAuthTokenString(oldAuthTokenClaims.StandardClaims.Subject, oldAuthTokenClaims.Role, csrfSecret)

			return
		}
		log.Println("Refresh token has expired!")
		// the refresh token has expired
		// Revoke the token in our db and require the user to login again
		models.RevokeRefreshToken(refreshTokenClaims.StandardClaims.Id, db)

		err = errors.New("Unauthorized")
		return
	}
	log.Println("Refresh token has been revoked!")
	err = errors.New("Unauthorized")
	return
}

// RevokeRefreshToken revokes a refresh token
func RevokeRefreshToken(refreshTokenString string, db *sql.DB) error {
	refreshToken, err := jwt.ParseWithClaims(refreshTokenString, &models.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return verifyKey, nil
	})
	if err != nil {
		return errors.New("Could not parse refresh token with claims")
	}

	refreshTokenClaims, ok := refreshToken.Claims.(*models.TokenClaims)
	if !ok {
		return errors.New("Could not read refresh token claims")
	}

	models.RevokeRefreshToken(refreshTokenClaims.StandardClaims.Id, db) // TODO implement this

	return nil
}

func updateRefreshTokenCsrf(oldRefreshTokenString string, newCsrfString string) (newRefreshTokenString string, err error) {
	refreshToken, err := jwt.ParseWithClaims(oldRefreshTokenString, &models.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return verifyKey, nil
	})

	oldRefreshTokenClaims, ok := refreshToken.Claims.(*models.TokenClaims)
	if !ok {
		return
	}

	refreshClaims := models.TokenClaims{
		jwt.StandardClaims{
			Id:        oldRefreshTokenClaims.StandardClaims.Id, // jti
			Subject:   oldRefreshTokenClaims.StandardClaims.Subject,
			ExpiresAt: oldRefreshTokenClaims.StandardClaims.ExpiresAt,
		},
		oldRefreshTokenClaims.Role,
		newCsrfString,
	}

	// create a signer for rsa 256
	refreshJwt := jwt.NewWithClaims(jwt.GetSigningMethod("RS256"), refreshClaims)

	// generate the refresh token string
	newRefreshTokenString, err = refreshJwt.SignedString(signKey)
	return
}

// GrabUUID grabs uuid
func GrabUUID(authTokenString string) (string, error) {
	authToken, _ := jwt.ParseWithClaims(authTokenString, &models.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return "", errors.New("Error fetching claims")
	})

	authTokenClaims, ok := authToken.Claims.(*models.TokenClaims)
	if !ok {
		return "", errors.New("Error fetching claims")
	}

	return authTokenClaims.StandardClaims.Subject, nil
}

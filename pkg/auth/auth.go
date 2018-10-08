package auth

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

// IAMConfig IAM config
type IAMConfig struct {
	SecretKey              string        `default:"OpenPitrix-lC4LipAXPYsuqw5F"`
	ExpireTime             time.Duration `default:"2h"`
	RefreshTokenExpireTime time.Duration `default:"336h"` // default is 2 week
}

// ErrExpired error expired
var ErrExpired = fmt.Errorf("access token expired")

const (
	//RequesterKey requester key
	RequesterKey = "requester"
	//TokenType token type
	TokenType = "Bearer"
)

//Requester requester info
type Requester struct {
	UserID string `json:"user_id,omitempty"`
	Role   string `json:"role,omitempty"`
}

//ToJSON format requester to json string
func (r *Requester) ToJSON() string {
	json, _ := json.Marshal(r)
	return string(json)
}

func trimKey(k string) []byte {
	return []byte(strings.TrimSpace(k))
}

//Validate jwt token
func Validate(k, s string) (*Requester, error) {
	tok, err := jwt.ParseSigned(s)
	if err != nil {
		return nil, err
	}
	c := &jwt.Claims{}
	requester := &Requester{}
	err = tok.Claims(trimKey(k), c, requester)
	if err != nil {
		return nil, err
	}
	if c.Expiry.Time().Unix() < time.Now().Unix() {
		return nil, ErrExpired
	}
	requester.UserID = c.Subject
	return requester, nil
}

//Generate jwt token
func Generate(k string, expire time.Duration, userID, role string) (string, error) {
	// TODO: use RS512 or ES512 to encrypt token
	// https://auth0.com/blog/brute-forcing-hs256-is-possible-the-importance-of-using-strong-keys-to-sign-jwts/

	signer, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.HS512, Key: trimKey(k)}, nil)
	if err != nil {
		return "", err
	}
	requester := &Requester{
		Role: role,
	}
	now := time.Now()
	c := &jwt.Claims{
		IssuedAt: jwt.NewNumericDate(now),
		Expiry:   jwt.NewNumericDate(now.Add(expire)),
		// TODO: add jti
		Subject: userID,
	}
	return jwt.Signed(signer).Claims(requester).Claims(c).CompactSerialize()
}

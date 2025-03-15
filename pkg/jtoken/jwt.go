package jtoken

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

const (
	expirationJWT = time.Hour * 5
	HeaderJWT     = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"
	/*
		{
			"alg": "HS256",
			"typ": "JWT"
		}
	*/

	errEmptyToken  = "empty token"
	errInvalidFmt  = "invalid token format"
	errInvalidHdr  = "invalid header format"
	errInvalidEnc  = "invalid payload encoding"
	errInvalidSig  = "invalid signature"
	errInvalidJSON = "invalid payload format"
	errTokenExp    = "token is expired"
	errNoUsername  = "username is required"
	errEmptyPL     = "empty payload"
	errSignFail    = "unable to generate access token"
)

// Payload represents JWT claims
type Payload struct {
	IsAdmin   bool      `json:"is_admin"`
	Username  string    `json:"username"`
	ExpiresAt time.Time `json:"expires_at"`
}

// Validate checks required fields in the Payload
func (p *Payload) Validate() error {
	if p == nil {
		return errors.New(errEmptyPL)
	}
	if p.Username == "" {
		return errors.New(errNoUsername)
	}
	return nil
}

// GenerateAccessToken creates a signed JWT token
func GenerateAccessToken(payload *Payload, secretKey string) (string, error) {
	if secretKey == "" {
		return "", errors.New(errSignFail)
	}

	if err := payload.Validate(); err != nil {
		return "", err
	}

	// payload.ExpiresAt = time.Now().Add(expirationJWT)
	if payload.ExpiresAt.Before(time.Now()) {
		return "", errors.New(errTokenExp)
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return "", errors.New(errInvalidJSON)
	}

	encodedPayload := base64.RawURLEncoding.EncodeToString(data)
	signature, err := SignHS256(HeaderJWT+"."+encodedPayload, secretKey)
	if err != nil {
		return "", fmt.Errorf("%s: %v", errSignFail, err)
	}

	return HeaderJWT + "." + encodedPayload + "." + signature, nil
}

// VerifyJWT checks token format, signature, and expiration
func VerifyJWT(token, secretKey string) (*Payload, error) {
	if token == "" {
		return nil, errors.New(errEmptyToken)
	}

	if secretKey == "" {
		return nil, errors.New(errSignFail)
	}

	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, errors.New(errInvalidFmt)
	}
	if parts[0] != HeaderJWT {
		return nil, errors.New(errInvalidFmt)
	}

	// Decode payload
	payload, err := DecodePayload(token)
	if err != nil {
		return nil, err
	}

	// Check expiration
	if payload.ExpiresAt.Before(time.Now()) {
		return nil, errors.New(errTokenExp)
	}

	// Verify signature
	expectedSignature, err := SignHS256(parts[0]+"."+parts[1], secretKey)
	if err != nil {
		return nil, err
	}
	if parts[2] != expectedSignature {
		return nil, errors.New(errInvalidSig)
	}

	return payload, nil
}

// SignHS256 generates an HMAC-SHA256 signature
func SignHS256(data, secretKey string) (string, error) {
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(data))
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil)), nil
}

// decodePayload decodes a Base64-encoded JSON payload into a Payload struct
func DecodePayload(jtoken string) (*Payload, error) {
	parts := strings.Split(jtoken, ".")
	if len(parts) != 3 {
		return nil, errors.New(errInvalidFmt)
	}

	data, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, errors.New(errInvalidEnc)
	}

	var payload Payload
	if err := json.Unmarshal(data, &payload); err != nil {
		return nil, errors.New(errInvalidJSON)
	}

	return &payload, nil
}

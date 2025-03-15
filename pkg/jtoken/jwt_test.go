package jtoken_test

import (
	"cafeteria/pkg/jtoken"
	"testing"
	"time"
)

func TestGenerateAccessToken(t *testing.T) {
	testCases := []struct {
		name      string
		payload   *jtoken.Payload
		secret    string
		wantError string
	}{
		{
			name:      "Missing username",
			payload:   &jtoken.Payload{},
			secret:    "very_secret",
			wantError: "username is required",
		},
		{
			name:      "Empty payload",
			payload:   nil,
			secret:    "very_secret",
			wantError: "empty payload",
		},
		{
			name: "Valid payload",
			payload: &jtoken.Payload{
				Username:  "test_user",
				IsAdmin:   false,
				ExpiresAt: time.Now().Add(time.Hour),
			},
			secret:    "very_secret",
			wantError: "",
		},
		{
			name: "Expired token",
			payload: &jtoken.Payload{
				Username:  "expired_user",
				IsAdmin:   false,
				ExpiresAt: time.Now().Add(-time.Hour), // Expired
			},
			secret:    "very_secret",
			wantError: "token is expired", // It generates a token, but should fail in VerifyJWT
		},
		{
			name:      "Empty secret key",
			payload:   &jtoken.Payload{Username: "test_user", ExpiresAt: time.Now().Add(time.Hour)},
			secret:    "",
			wantError: "unable to generate access token",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			token, got := jtoken.GenerateAccessToken(tc.payload, tc.secret)

			if (got != nil && got.Error() != tc.wantError) || (got == nil && tc.wantError != "") {
				t.Errorf("Test %q failed: got error %q, wanted %q", tc.name, got, tc.wantError)
			}

			// If token is generated, verify it to check correctness
			if got == nil && token != "" && tc.wantError == "" {
				_, verifyErr := jtoken.VerifyJWT(token, tc.secret)
				if verifyErr != nil {
					t.Errorf("Verification failed for test %q: %v", tc.name, verifyErr)
				}
			}
		})
	}
}

func TestVerifyJWT(t *testing.T) {
	secret := "very_secret"

	validPayload := &jtoken.Payload{
		Username:  "valid_user",
		IsAdmin:   false,
		ExpiresAt: time.Now().Add(time.Hour),
	}

	validToken, _ := jtoken.GenerateAccessToken(validPayload, secret)
	tamperedToken := validToken[:len(validToken)-2] + "xx" // Simulating tampering

	testCases := []struct {
		name      string
		token     string
		secret    string
		wantError string
	}{
		{
			name:      "Valid token",
			token:     validToken,
			secret:    secret,
			wantError: "",
		},
		{
			name:      "Invalid token format",
			token:     "invalid.token.format",
			secret:    secret,
			wantError: "invalid token format",
		},
		{
			name:      "Invalid signature",
			token:     tamperedToken,
			secret:    secret,
			wantError: "invalid signature",
		},
		{
			name:      "Invalid secret key",
			token:     validToken,
			secret:    "wrong_secret",
			wantError: "invalid signature",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, got := jtoken.VerifyJWT(tc.token, tc.secret)

			if (got != nil && got.Error() != tc.wantError) || (got == nil && tc.wantError != "") {
				t.Errorf("Test %q failed: got %q, wanted %q", tc.name, got, tc.wantError)
			}
		})
	}
}

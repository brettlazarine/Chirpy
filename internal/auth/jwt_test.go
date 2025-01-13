package auth

import (
	"testing"
	"time"

	"net/http"

	"github.com/google/uuid"
)

func TestJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecretPass := "secretive"
	tokenSecretFail := "NotSoSecretive"
	expiresInHour := 1 * time.Hour
	expiresInNow := 0 * time.Hour

	testPass := []struct {
		name        string
		userID      uuid.UUID
		tokenSecret string
		expiresIn   time.Duration
		wantError   bool
	}{
		{
			name:        "TestMakeJWT",
			userID:      userID,
			tokenSecret: tokenSecretPass,
			expiresIn:   expiresInHour,
			wantError:   false,
		},
		{
			name:        "TestValidateJWT",
			userID:      userID,
			tokenSecret: tokenSecretPass,
			expiresIn:   expiresInHour,
			wantError:   false,
		},
	}

	testFailSecret := []struct {
		name            string
		userID          uuid.UUID
		tokenSecretOK   string
		tokenSecretFail string
		expiresIn       time.Duration
		wantError       bool
	}{
		{
			name:            "TestValidateJWT_WrongSecret",
			userID:          userID,
			tokenSecretOK:   tokenSecretFail,
			tokenSecretFail: tokenSecretPass,
			expiresIn:       expiresInHour,
			wantError:       true,
		},
	}

	testFailExpired := []struct {
		name        string
		userID      uuid.UUID
		tokenSecret string
		expiresIn   time.Duration
		wantError   bool
	}{
		{
			name:        "TestValidateJWT_Expired",
			userID:      userID,
			tokenSecret: tokenSecretPass,
			expiresIn:   expiresInNow,
			wantError:   true,
		},
	}

	for _, tt := range testPass {
		t.Run(tt.name, func(t *testing.T) {
			_, err := MakeJWT(tt.userID, tt.tokenSecret, tt.expiresIn)
			if (err != nil) != tt.wantError {
				t.Errorf("JWT *%v* error = %v, wantError %v", tt.name, err, tt.wantError)
			}
		})
	}
	for _, tt := range testFailSecret {
		t.Run(tt.name, func(t *testing.T) {
			token, _ := MakeJWT(tt.userID, tt.tokenSecretOK, tt.expiresIn)
			_, err := ValidateJWT(token, tt.tokenSecretFail)
			if (err != nil) != tt.wantError {
				t.Errorf("JWT *%v* error = %v, wantError %v", tt.name, err, tt.wantError)
			}
		})
	}
	for _, tt := range testFailExpired {
		t.Run(tt.name, func(t *testing.T) {
			token, _ := MakeJWT(tt.userID, tt.tokenSecret, tt.expiresIn)
			time.Sleep(2 * time.Second)
			_, err := ValidateJWT(token, tt.tokenSecret)
			if (err != nil) != tt.wantError {
				t.Errorf("JWT *%v* error = %v, wantError %v", tt.name, err, tt.wantError)
			}
		})
	}
}

func TestValidateJWT_BootDev(t *testing.T) {
	userID := uuid.New()
	validToken, _ := MakeJWT(userID, "secret", time.Hour)

	tests := []struct {
		name        string
		tokenString string
		tokenSecret string
		wantUserID  uuid.UUID
		wantErr     bool
	}{
		{
			name:        "Valid token",
			tokenString: validToken,
			tokenSecret: "secret",
			wantUserID:  userID,
			wantErr:     false,
		},
		{
			name:        "Invalid token",
			tokenString: "invalid.token.string",
			tokenSecret: "secret",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
		{
			name:        "Wrong secret",
			tokenString: validToken,
			tokenSecret: "wrong_secret",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUserID, err := ValidateJWT(tt.tokenString, tt.tokenSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotUserID != tt.wantUserID {
				t.Errorf("ValidateJWT() gotUserID = %v, want %v", gotUserID, tt.wantUserID)
			}
		})
	}
}

func TestGetBearerToken(t *testing.T) {
	validToken := "Bearer valid.token.string"
	tokenString := "valid.token.string"
	invalidToken := "invalid.token.string"

	tests := []struct {
		name      string
		header    string
		wantToken string
		wantErr   bool
	}{
		{
			name:      "Valid token",
			header:    validToken,
			wantToken: tokenString,
			wantErr:   false,
		},
		{
			name:      "Invalid token",
			header:    invalidToken,
			wantToken: "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			headers := make(map[string][]string)
			headers["Authorization"] = []string{tt.header}
			gotToken, err := GetBearerToken(headers)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBearerToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotToken != tt.wantToken {
				t.Errorf("GetBearerToken() gotToken = %v, want %v", gotToken, tt.wantToken)
			}
		})
	}
}

func TestGetBearerToken_BootDev(t *testing.T) {
	tests := []struct {
		name      string
		headers   http.Header
		wantToken string
		wantErr   bool
	}{
		{
			name: "Valid Bearer token",
			headers: http.Header{
				"Authorization": []string{"Bearer valid_token"},
			},
			wantToken: "valid_token",
			wantErr:   false,
		},
		{
			name:      "Missing Authorization header",
			headers:   http.Header{},
			wantToken: "",
			wantErr:   true,
		},
		{
			name: "Malformed Authorization header",
			headers: http.Header{
				"Authorization": []string{"InvalidBearer token"},
			},
			wantToken: "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotToken, err := GetBearerToken(tt.headers)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBearerToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotToken != tt.wantToken {
				t.Errorf("GetBearerToken() gotToken = %v, want %v", gotToken, tt.wantToken)
			}
		})
	}
}

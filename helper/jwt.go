package helper

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"google.golang.org/api/idtoken"
)

type TokenPayload struct {
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Iss           string `json:"iss"`
	Aud           string `json:"aud"`
	Exp           int64  `json:"exp"`
	Iat           int64  `json:"iat"`
}

func DecodeIdToken(idToken string) (*TokenPayload, error) {
	parts := strings.Split(idToken, ".")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid token")
	}

	payload := parts[1]
	decoded, err := base64.RawURLEncoding.DecodeString(payload)
	if err != nil {
		return nil, err
	}

	tokenPayload := TokenPayload{}
	err = json.Unmarshal(decoded, &tokenPayload)
	if err != nil {
		return nil, err
	}

	return &tokenPayload, nil
}

func VerifyGoogleIdToken(ctx context.Context, idToken string, audience string) (*idtoken.Payload, error) {
	payload, err := idtoken.Validate(ctx, idToken, audience)
	if err != nil {
		return nil, fmt.Errorf("token verification failed: %w", err)
	}
	return payload, nil
}

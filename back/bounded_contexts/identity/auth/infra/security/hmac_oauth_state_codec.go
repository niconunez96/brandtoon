package security

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"brandtoonapi/bounded_contexts/identity/auth/domain"
)

const maxOAuthStateAge = 10 * time.Minute

type HMACOAuthStateCodec struct {
	secret []byte
}

type serializedOAuthState struct {
	IssuedAt   string `json:"issuedAt"`
	Nonce      string `json:"nonce"`
	RedirectTo string `json:"redirectTo"`
}

func NewHMACOAuthStateCodec(secret string) *HMACOAuthStateCodec {
	return &HMACOAuthStateCodec{secret: []byte(secret)}
}

func (c *HMACOAuthStateCodec) Decode(rawState string) (*domain.OAuthState, error) {
	parts := strings.Split(rawState, ".")
	if len(parts) != 2 {
		return nil, errors.New("invalid oauth state format")
	}

	payloadSegment := parts[0]
	signatureSegment := parts[1]
	if !hmac.Equal([]byte(signatureSegment), []byte(c.sign(payloadSegment))) {
		return nil, errors.New("invalid oauth state signature")
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(payloadSegment)
	if err != nil {
		return nil, err
	}

	var payload serializedOAuthState
	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return nil, err
	}

	issuedAt, err := time.Parse(time.RFC3339Nano, payload.IssuedAt)
	if err != nil {
		return nil, err
	}

	if time.Since(issuedAt.UTC()) > maxOAuthStateAge {
		return nil, errors.New("oauth state expired")
	}

	return &domain.OAuthState{
		IssuedAt:   issuedAt.UTC(),
		Nonce:      payload.Nonce,
		RedirectTo: payload.RedirectTo,
	}, nil
}

func (c *HMACOAuthStateCodec) Encode(state domain.OAuthState) (string, error) {
	payload := serializedOAuthState{
		IssuedAt:   state.IssuedAt.UTC().Format(time.RFC3339Nano),
		Nonce:      state.Nonce,
		RedirectTo: state.RedirectTo,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	payloadSegment := base64.RawURLEncoding.EncodeToString(payloadBytes)
	return fmt.Sprintf("%s.%s", payloadSegment, c.sign(payloadSegment)), nil
}

func (c *HMACOAuthStateCodec) sign(payload string) string {
	hash := hmac.New(sha256.New, c.secret)
	hash.Write([]byte(payload))
	return base64.RawURLEncoding.EncodeToString(hash.Sum(nil))
}

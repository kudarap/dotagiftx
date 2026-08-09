package github

import (
	"bytes"
	"context"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type (
	Config struct {
		AppID          string
		PrivateKey     string
		InstallationID string
		Owner          string
		Repository     string
	}

	Client struct {
		http   *http.Client
		config Config
	}
)

func New(conf Config) *Client {
	return &Client{
		http:   &http.Client{},
		config: conf,
	}
}

func (c *Client) CreateIssue(ctx context.Context, title string, body string) (string, error) {
	if title == "" || body == "" {
		return "", fmt.Errorf("title and body required")
	}

	payload := struct {
		Title string `json:"title"`
		Body  string `json:"body"`
	}{
		title,
		body,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	token, err := c.getInstallationAccessToken(ctx)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues", c.config.Owner, c.config.Repository)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-GitHub-Api-Version", "2026-03-10")

	res, err := c.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("req err: %w", err)
	}

	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		body, _ := io.ReadAll(res.Body)
		return "", fmt.Errorf("github status %d: %s", res.StatusCode, body)
	}

	var result struct {
		HTMLURL string `json:"html_url"`
	}

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.HTMLURL, nil
}

func (c *Client) getInstallationAccessToken(ctx context.Context) (string, error) {
	url := fmt.Sprintf("https://api.github.com/app/installations/%s/access_tokens", c.config.InstallationID)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return "", err
	}

	appJWT, err := c.appJWT()
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", appJWT))
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-Github-Api-Version", "2026-03-10")

	res, err := c.http.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		body, _ := io.ReadAll(res.Body)
		return "", fmt.Errorf("installation github status %d: %s", res.StatusCode, body)
	}

	var token struct {
		Token     string    `json:"token"`
		ExpiresAt time.Time `json:"expires_at"`
	}

	if err := json.NewDecoder(res.Body).Decode(&token); err != nil {
		return "", err
	}

	return token.Token, nil
}

// Generates a JWT based on the private key provided from Github
func (c *Client) appJWT() (string, error) {
	// Decode the key as its encoded to base64
	pembytes, err := base64.StdEncoding.DecodeString(strings.TrimSpace(c.config.PrivateKey))
	if err != nil {
		return "", err
	}

	block, _ := pem.Decode(pembytes)
	if block == nil {
		return "", errors.New("invalid or empty private key")
	}

	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("parse RSA private key: %w", err)
	}

	now := time.Now().Unix()
	claims := jwt.MapClaims{
		"iat": now - 60,
		"exp": now + 9*60,
		"iss": c.config.AppID,
	}
	return jwt.NewWithClaims(jwt.SigningMethodRS256, claims).
		SignedString(key)
}

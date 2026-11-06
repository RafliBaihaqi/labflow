package turnstile

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

const siteverifyURL = "https://challenges.cloudflare.com/turnstile/v0/siteverify"

type verifyResponse struct {
	Success bool     `json:"success"`
	Codes   []string `json:"error-codes"`
}

func Verify(token, remoteIP string) error {
	env := os.Getenv("ENV")
	if env == "" {
		env = os.Getenv("GO_ENV")
	}
	isProd := env == "production"

	// Dev-only bypass: explicit flag, not allowed in production
	if os.Getenv("BYPASS_TURNSTILE") == "true" {
		if isProd {
			return fmt.Errorf("turnstile: BYPASS_TURNSTILE=true not allowed in production")
		}
		fmt.Println("[turnstile] Dev Mode: BYPASS_TURNSTILE is active — skipping verification")
		return nil
	}

	// Cloudflare test token that always passes - allow only in non-prod
	if token == "1x00000000000000000000AA" {
		if isProd {
			return fmt.Errorf("turnstile: test token not allowed in production")
		}
		fmt.Println("[turnstile] Dev Mode: Bypassing siteverify for test/dummy token")
		return nil
	}

	secret := os.Getenv("TURNSTILE_SECRET")
	if secret == "" {
		// Fail closed in prod, allow skip with warning in dev
		if isProd {
			return fmt.Errorf("turnstile: TURNSTILE_SECRET not set - failing closed (prod must enforce)")
		}
		fmt.Println("[turnstile] WARNING: TURNSTILE_SECRET not set — skipping verification (dev only)")
		return nil
	}

	fields := url.Values{
		"secret":   {secret},
		"response": {token},
	}
	if remoteIP != "" {
		fields.Set("remoteip", remoteIP)
	}

	resp, err := http.PostForm(siteverifyURL, fields)
	if err != nil {
		errStr := fmt.Errorf("turnstile: siteverify request failed: %w", err)
		fmt.Printf("[turnstile] ERROR: %v\n", errStr)
		return errStr
	}
	defer resp.Body.Close()

	var result verifyResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		errStr := fmt.Errorf("turnstile: failed to decode siteverify response: %w", err)
		fmt.Printf("[turnstile] ERROR: %v\n", errStr)
		return errStr
	}

	if !result.Success {
		errStr := fmt.Errorf("turnstile: challenge failed: %v", result.Codes)
		fmt.Printf("[turnstile] ERROR: %v (used secret: %s...)\n", errStr, secret[:min(len(secret), 6)])
		return errStr
	}
	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

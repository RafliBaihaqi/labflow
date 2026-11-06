package bootstrap

import (
	"fmt"
	"os"
)

func CookieKey() string {
	cookieKey := os.Getenv("COOKIE_SECRET")
	if cookieKey == "" {
		panic("COOKIE_SECRET environment variable is not set")
	}
	return cookieKey
}

// RequireEnv returns the value of the named environment variable, or an error
// if it is unset or empty. Use for configuration that must be present at boot.
func RequireEnv(key string) (string, error) {
	v := os.Getenv(key)
	if v == "" {
		return "", fmt.Errorf("required environment variable %s is not set", key)
	}
	return v, nil
}

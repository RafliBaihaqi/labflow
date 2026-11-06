package bootstrap

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadEnvFile(t *testing.T) {
	t.Run("missing file returns nil and does not error", func(t *testing.T) {
		dir := t.TempDir()
		if err := loadEnvFile(filepath.Join(dir, ".env")); err != nil {
			t.Fatalf("missing .env should not error, got %v", err)
		}
	})

	t.Run("valid file loads variables", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, ".env")
		if err := os.WriteFile(path, []byte("LABFLOW_TEST_KEY=hello\n"), 0o600); err != nil {
			t.Fatalf("write temp env: %v", err)
		}
		// Register the var for automatic cleanup, then clear it: godotenv does
		// not overwrite variables that are already present in the environment.
		t.Setenv("LABFLOW_TEST_KEY", "placeholder")
		os.Unsetenv("LABFLOW_TEST_KEY")
		if err := loadEnvFile(path); err != nil {
			t.Fatalf("valid .env should load, got %v", err)
		}
		if got := os.Getenv("LABFLOW_TEST_KEY"); got != "hello" {
			t.Fatalf("expected LABFLOW_TEST_KEY=hello, got %q", got)
		}
	})

	t.Run("malformed file returns error", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, ".env")
		if err := os.WriteFile(path, []byte("THIS LINE HAS NO EQUALS SIGN\n"), 0o600); err != nil {
			t.Fatalf("write temp env: %v", err)
		}
		if err := loadEnvFile(path); err == nil {
			t.Fatal("malformed .env should return an error, got nil")
		}
	})
}

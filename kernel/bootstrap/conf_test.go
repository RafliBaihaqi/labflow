package bootstrap

import (
	"testing"
)

func TestRequireEnv(t *testing.T) {
	t.Run("returns value when set", func(t *testing.T) {
		t.Setenv("LABFLOW_REQ_KEY", "certificates")
		got, err := RequireEnv("LABFLOW_REQ_KEY")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != "certificates" {
			t.Fatalf("expected 'certificates', got %q", got)
		}
	})

	t.Run("returns error when empty", func(t *testing.T) {
		t.Setenv("LABFLOW_REQ_KEY", "")
		if _, err := RequireEnv("LABFLOW_REQ_KEY"); err == nil {
			t.Fatal("expected error for empty env var, got nil")
		}
	})
}

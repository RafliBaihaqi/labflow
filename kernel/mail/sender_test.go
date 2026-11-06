package mail

import (
	"strings"
	"testing"
)

func TestRenderVerificationEmail(t *testing.T) {
	body, err := renderVerificationEmail("123456")
	if err != nil {
		t.Fatalf("render should succeed with embedded template, got %v", err)
	}
	if !strings.Contains(string(body), "123456") {
		t.Fatalf("rendered body should contain the OTP, got: %s", string(body))
	}
}

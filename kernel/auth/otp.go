package auth

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// GenerateOTP returns a cryptographically secure 6-digit numeric code.
func GenerateOTP() string {
	maxVal := big.NewInt(1000000)
	n, err := rand.Int(rand.Reader, maxVal)
	if err != nil {
		panic(fmt.Sprintf("failed to generate OTP: %v", err))
	}
	return fmt.Sprintf("%06d", n.Int64())
}

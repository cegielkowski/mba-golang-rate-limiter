package pkg

import (
	"context"
	"mba-golang-rate-limiter/pkg"
	"testing"
	"time"
)

func TestRateLimiter(t *testing.T) {
	rl, err := pkg.NewRateLimiter()
	if err != nil {
		t.Fatalf("Could not create rate limiter: %v", err)
	}

	ctx := context.Background()

	// Test IP rate limiting
	for i := 0; i < rl.GetDefaultIPLimit(); i++ {
		allowed, err := rl.AllowIP(ctx, "192.168.1.1")
		if err != nil || !allowed {
			t.Fatalf("IP rate limit failed on iteration %d: %v, allowed: %v", i, err, allowed)
		}
	}

	allowed, err := rl.AllowIP(ctx, "192.168.1.1")
	if err != nil {
		t.Fatalf("IP rate limit error: %v", err)
	}
	if allowed {
		t.Fatal("IP rate limit did not block request as expected")
	}

	// Log the sleep duration for debugging
	t.Log("Sleeping for 6 seconds to allow rate limit to reset")
	time.Sleep(6 * time.Second) // Increase sleep time to ensure expiration

	allowed, err = rl.AllowIP(ctx, "192.168.1.1")
	if err != nil || !allowed {
		t.Fatalf("IP rate limit did not reset after expiration: %v, allowed: %v", err, allowed)
	} else {
		t.Log("IP rate limit successfully reset after expiration")
	}

	// Test Token rate limiting
	for i := 0; i < rl.GetDefaultTokenLimit(); i++ {
		allowed, err := rl.AllowToken(ctx, "abc123")
		if err != nil || !allowed {
			t.Fatalf("Token rate limit failed on iteration %d: %v, allowed: %v", i, err, allowed)
		}
	}

	allowed, err = rl.AllowToken(ctx, "abc123")
	if err != nil {
		t.Fatalf("Token rate limit error: %v", err)
	}
	if allowed {
		t.Fatal("Token rate limit did not block request as expected")
	}

	// Log the sleep duration for debugging
	t.Log("Sleeping for 6 seconds to allow rate limit to reset")
	time.Sleep(6 * time.Second) // Increase sleep time to ensure expiration

	allowed, err = rl.AllowToken(ctx, "abc123")
	if err != nil || !allowed {
		t.Fatalf("Token rate limit did not reset after expiration: %v, allowed: %v", err, allowed)
	} else {
		t.Log("Token rate limit successfully reset after expiration")
	}
}

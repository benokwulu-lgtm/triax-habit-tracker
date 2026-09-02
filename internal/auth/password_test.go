package auth

import "testing"

func TestHashAndCheckPassword(t *testing.T) {
	password := "correct-horse-battery-staple"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}

	if !CheckPassword(password, hash) {
		t.Error("expected CheckPassword to return true for the correct password")
	}

	if CheckPassword("wrong-password", hash) {
		t.Error("expected CheckPassword to return false for an incorrect password")
	}
}

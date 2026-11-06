package security

import "testing"

func TestPasswordHashVerification(t *testing.T) {
	hash, err := HashPassword("correct horse battery staple")
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}
	if hash == "correct horse battery staple" {
		t.Fatal("password was stored in plain text")
	}
	if !VerifyPassword(hash, "correct horse battery staple") {
		t.Fatal("VerifyPassword() rejected the correct password")
	}
	if VerifyPassword(hash, "incorrect password") {
		t.Fatal("VerifyPassword() accepted an incorrect password")
	}
}

func TestVerifyPasswordRejectsMalformedHash(t *testing.T) {
	if VerifyPassword("not-a-hash", "password") {
		t.Fatal("VerifyPassword() accepted a malformed hash")
	}
}

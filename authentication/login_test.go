package authentication

import "testing"

func TestAuthService_Login(t *testing.T) {
	svc := NewAuthService()
	if _, err := svc.CreateUser(&UserReq{Username: "jane", Password: "s3cret"}); err != nil {
		t.Fatalf("create user: %v", err)
	}

	if err := svc.Login("preeti", "admin123"); err != nil {
		t.Errorf("seeded user should log in, got %v", err)
	}
	if err := svc.Login("jane", "s3cret"); err != nil {
		t.Errorf("created user should log in, got %v", err)
	}
	if err := svc.Login("jane", "wrong"); err == nil {
		t.Error("wrong password should fail")
	}
}

package dto

import "testing"

func TestRegisterRequest_Validate(t *testing.T) {
	tests := []struct {
		name     string
		req      RegisterRequest
		wantErr  bool
		errField string
	}{
		{"valid", RegisterRequest{Name: "Alice", Email: "a@b.com", Password: "secret123"}, false, ""},
		{"missing name", RegisterRequest{Email: "a@b.com", Password: "secret123"}, true, "name"},
		{"missing email", RegisterRequest{Name: "Alice", Password: "secret123"}, true, "email"},
		{"invalid email", RegisterRequest{Name: "Alice", Email: "not-an-email", Password: "secret123"}, true, "email"},
		{"short password", RegisterRequest{Name: "Alice", Email: "a@b.com", Password: "abc"}, true, "password"},
		{"all empty", RegisterRequest{}, true, "name"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := tt.req.Validate()
			if tt.wantErr && errs == nil {
				t.Fatalf("expected validation error for field %q, got nil", tt.errField)
			}
			if !tt.wantErr && errs != nil {
				t.Fatalf("expected no error, got: %v", errs)
			}
			if tt.wantErr && tt.errField != "" {
				if _, ok := errs[tt.errField]; !ok {
					t.Errorf("expected error on field %q, got fields: %v", tt.errField, errs)
				}
			}
		})
	}
}

func TestLoginRequest_Validate(t *testing.T) {
	tests := []struct {
		name     string
		req      LoginRequest
		wantErr  bool
		errField string
	}{
		{"valid", LoginRequest{Email: "a@b.com", Password: "secret"}, false, ""},
		{"missing email", LoginRequest{Password: "secret"}, true, "email"},
		{"missing password", LoginRequest{Email: "a@b.com"}, true, "password"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := tt.req.Validate()
			if tt.wantErr && errs == nil {
				t.Fatalf("expected validation error for field %q, got nil", tt.errField)
			}
			if !tt.wantErr && errs != nil {
				t.Fatalf("expected no error, got: %v", errs)
			}
			if tt.wantErr && tt.errField != "" {
				if _, ok := errs[tt.errField]; !ok {
					t.Errorf("expected error on field %q, got fields: %v", tt.errField, errs)
				}
			}
		})
	}
}

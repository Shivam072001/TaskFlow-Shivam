package dto

import "testing"

func TestCreateTaskRequest_Validate(t *testing.T) {
	tests := []struct {
		name     string
		req      CreateTaskRequest
		wantErr  bool
		errField string
	}{
		{"valid minimal", CreateTaskRequest{Title: "Fix bug"}, false, ""},
		{"valid with priority", CreateTaskRequest{Title: "Fix bug", Priority: "high"}, false, ""},
		{"missing title", CreateTaskRequest{Priority: "medium"}, true, "title"},
		{"invalid priority", CreateTaskRequest{Title: "Fix bug", Priority: "critical"}, true, "priority"},
		{
			"valid dates",
			CreateTaskRequest{Title: "Task", StartDate: strPtr("2025-01-01"), DueDate: strPtr("2025-12-31")},
			false, "",
		},
		{
			"bad start_date format",
			CreateTaskRequest{Title: "Task", StartDate: strPtr("not-a-date")},
			true, "start_date",
		},
		{
			"bad due_date format",
			CreateTaskRequest{Title: "Task", DueDate: strPtr("31/12/2025")},
			true, "due_date",
		},
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

func TestUpdateTaskRequest_Validate(t *testing.T) {
	tests := []struct {
		name     string
		req      UpdateTaskRequest
		wantErr  bool
		errField string
	}{
		{"empty update (no fields)", UpdateTaskRequest{}, false, ""},
		{"valid status", UpdateTaskRequest{Status: strPtr("done")}, false, ""},
		{"invalid status", UpdateTaskRequest{Status: strPtr("cancelled")}, true, "status"},
		{"blocked without reason", UpdateTaskRequest{Status: strPtr("blocked")}, true, "blocked_reason"},
		{
			"blocked with reason",
			UpdateTaskRequest{Status: strPtr("blocked"), BlockedReason: strPtr("waiting for API")},
			false, "",
		},
		{
			"blocked with blocked_by_task",
			UpdateTaskRequest{Status: strPtr("blocked"), BlockedByTask: strPtr("PROJ-1")},
			false, "",
		},
		{"empty title", UpdateTaskRequest{Title: strPtr("")}, true, "title"},
		{"invalid priority", UpdateTaskRequest{Priority: strPtr("urgent")}, true, "priority"},
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

func TestParsePagination_Defaults(t *testing.T) {
	// ParsePagination requires an *http.Request — tested through the parseInt helper logic.
	// Direct validation: parseInt for valid/invalid inputs.
	if v := parseInt("5"); v != 5 {
		t.Errorf("parseInt(\"5\") = %d, want 5", v)
	}
	if v := parseInt("0"); v != 0 {
		t.Errorf("parseInt(\"0\") = %d, want 0", v)
	}
	if v := parseInt("abc"); v != 0 {
		t.Errorf("parseInt(\"abc\") = %d, want 0", v)
	}
	if v := parseInt(""); v != 0 {
		t.Errorf("parseInt(\"\") = %d, want 0", v)
	}
	if v := parseInt("100"); v != 100 {
		t.Errorf("parseInt(\"100\") = %d, want 100", v)
	}
}

func strPtr(s string) *string { return &s }

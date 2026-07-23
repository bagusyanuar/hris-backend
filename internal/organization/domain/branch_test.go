package domain

import (
	"errors"
	"testing"
)

func TestNewBranch(t *testing.T) {
	city := "Jakarta"

	tests := []struct {
		name      string
		companyID string
		code      string
		branch    string
		city      *string
		isMain    bool
		wantErr   error
	}{
		{
			name:      "happy path main branch",
			companyID: "co-1",
			code:      "JKT",
			branch:    "Kantor Pusat Jakarta",
			city:      &city,
			isMain:    true,
			wantErr:   nil,
		},
		{
			name:      "happy path without city (nullable)",
			companyID: "co-1",
			code:      "SBY",
			branch:    "Cabang Surabaya",
			city:      nil,
			isMain:    false,
			wantErr:   nil,
		},
		{
			name:      "empty companyID is invalid (scope wajib)",
			companyID: "",
			code:      "JKT",
			branch:    "Kantor Pusat Jakarta",
			wantErr:   ErrInvalidInput,
		},
		{
			name:      "empty code is invalid",
			companyID: "co-1",
			code:      "",
			branch:    "Kantor Pusat Jakarta",
			wantErr:   ErrInvalidInput,
		},
		{
			name:      "empty name is invalid",
			companyID: "co-1",
			code:      "JKT",
			branch:    "",
			wantErr:   ErrInvalidInput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			branch, err := NewBranch(tt.companyID, tt.code, tt.branch, tt.city, tt.isMain)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected error %v, got %v", tt.wantErr, err)
				}
				if branch != nil {
					t.Fatalf("expected nil branch on error, got %+v", branch)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if branch.ID == "" {
				t.Fatal("expected generated ID, got empty string")
			}
			if branch.CompanyID != tt.companyID {
				t.Errorf("expected companyID %q, got %q", tt.companyID, branch.CompanyID)
			}
			if branch.IsMain != tt.isMain {
				t.Errorf("expected isMain %v, got %v", tt.isMain, branch.IsMain)
			}
			if !branch.IsActive {
				t.Error("expected new branch to default IsActive=true")
			}
			if branch.CreatedAt.IsZero() || branch.UpdatedAt.IsZero() {
				t.Error("expected CreatedAt/UpdatedAt to be set")
			}
		})
	}
}

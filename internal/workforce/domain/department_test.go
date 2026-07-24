package domain

import (
	"errors"
	"testing"
)

func TestNewDepartment(t *testing.T) {
	parentID := "11111111-1111-1111-1111-111111111111"

	tests := []struct {
		name      string
		companyID string
		code      string
		deptName  string
		parentID  *string
		wantErr   error
	}{
		{
			name:      "happy path root department",
			companyID: "co-1",
			code:      "DIR",
			deptName:  "Direksi",
			parentID:  nil,
			wantErr:   nil,
		},
		{
			name:      "happy path child department",
			companyID: "co-1",
			code:      "TI",
			deptName:  "Divisi Teknologi Informasi",
			parentID:  &parentID,
			wantErr:   nil,
		},
		{
			name:      "empty company id is invalid",
			companyID: "",
			code:      "TI",
			deptName:  "Divisi TI",
			wantErr:   ErrInvalidInput,
		},
		{
			name:      "empty code is invalid",
			companyID: "co-1",
			code:      "",
			deptName:  "Divisi TI",
			wantErr:   ErrInvalidInput,
		},
		{
			name:      "empty name is invalid",
			companyID: "co-1",
			code:      "TI",
			deptName:  "",
			wantErr:   ErrInvalidInput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			department, err := NewDepartment(tt.companyID, tt.code, tt.deptName, tt.parentID)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected error %v, got %v", tt.wantErr, err)
				}
				if department != nil {
					t.Fatalf("expected nil department on error, got %+v", department)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if department.ID == "" {
				t.Fatal("expected generated ID, got empty string")
			}
			if department.CompanyID != tt.companyID {
				t.Errorf("expected company id %q, got %q", tt.companyID, department.CompanyID)
			}
			if department.Code != tt.code {
				t.Errorf("expected code %q, got %q", tt.code, department.Code)
			}
			if department.Name != tt.deptName {
				t.Errorf("expected name %q, got %q", tt.deptName, department.Name)
			}
			if department.ParentID != tt.parentID {
				t.Errorf("expected parent id %v, got %v", tt.parentID, department.ParentID)
			}
			if !department.IsActive {
				t.Error("expected new department to default IsActive=true")
			}
			if department.CreatedAt.IsZero() || department.UpdatedAt.IsZero() {
				t.Error("expected CreatedAt/UpdatedAt to be set")
			}
		})
	}
}

func TestNewDepartment_GeneratesUniqueIDs(t *testing.T) {
	d1, err := NewDepartment("co-1", "DIR", "Direksi", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	d2, err := NewDepartment("co-1", "OPR", "Operasional", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d1.ID == d2.ID {
		t.Fatal("expected different IDs for different departments")
	}
}

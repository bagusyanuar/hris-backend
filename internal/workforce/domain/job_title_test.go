package domain

import (
	"errors"
	"testing"
)

func TestNewJobTitle(t *testing.T) {
	tests := []struct {
		name       string
		companyID  string
		code       string
		titleName  string
		gradeLevel int
		wantErr    error
	}{
		{
			name:       "happy path",
			companyID:  "co-1",
			code:       "MGR",
			titleName:  "Manajer",
			gradeLevel: 7,
			wantErr:    nil,
		},
		{
			name:       "empty company id is invalid",
			companyID:  "",
			code:       "MGR",
			titleName:  "Manajer",
			gradeLevel: 7,
			wantErr:    ErrInvalidInput,
		},
		{
			name:       "empty code is invalid",
			companyID:  "co-1",
			code:       "",
			titleName:  "Manajer",
			gradeLevel: 7,
			wantErr:    ErrInvalidInput,
		},
		{
			name:       "empty name is invalid",
			companyID:  "co-1",
			code:       "MGR",
			titleName:  "",
			gradeLevel: 7,
			wantErr:    ErrInvalidInput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jobTitle, err := NewJobTitle(tt.companyID, tt.code, tt.titleName, tt.gradeLevel)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected error %v, got %v", tt.wantErr, err)
				}
				if jobTitle != nil {
					t.Fatalf("expected nil job title on error, got %+v", jobTitle)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if jobTitle.ID == "" {
				t.Fatal("expected generated ID, got empty string")
			}
			if jobTitle.CompanyID != tt.companyID {
				t.Errorf("expected company id %q, got %q", tt.companyID, jobTitle.CompanyID)
			}
			if jobTitle.GradeLevel != tt.gradeLevel {
				t.Errorf("expected grade level %d, got %d", tt.gradeLevel, jobTitle.GradeLevel)
			}
			if !jobTitle.IsActive {
				t.Error("expected new job title to default IsActive=true")
			}
			if jobTitle.CreatedAt.IsZero() || jobTitle.UpdatedAt.IsZero() {
				t.Error("expected CreatedAt/UpdatedAt to be set")
			}
		})
	}
}

func TestNewJobTitle_GeneratesUniqueIDs(t *testing.T) {
	t1, err := NewJobTitle("co-1", "MGR", "Manajer", 7)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	t2, err := NewJobTitle("co-1", "SPV", "Supervisor", 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if t1.ID == t2.ID {
		t.Fatal("expected different IDs for different job titles")
	}
}

package domain

import (
	"errors"
	"testing"
)

func TestNewJobPosition(t *testing.T) {
	reportsTo := "22222222-2222-2222-2222-222222222222"

	tests := []struct {
		name           string
		companyID      string
		departmentID   string
		jobTitleID     string
		posName        string
		reportsToID    *string
		headcountQuota int
		wantErr        error
		wantQuota      int
	}{
		{
			name:           "happy path with explicit quota",
			companyID:      "co-1",
			departmentID:   "dept-1",
			jobTitleID:     "title-1",
			posName:        "Manajer Pengembangan",
			reportsToID:    &reportsTo,
			headcountQuota: 3,
			wantErr:        nil,
			wantQuota:      3,
		},
		{
			name:           "quota defaults to 1 when zero",
			companyID:      "co-1",
			departmentID:   "dept-1",
			jobTitleID:     "title-1",
			posName:        "Direktur Utama",
			reportsToID:    nil,
			headcountQuota: 0,
			wantErr:        nil,
			wantQuota:      1,
		},
		{
			name:           "quota defaults to 1 when negative",
			companyID:      "co-1",
			departmentID:   "dept-1",
			jobTitleID:     "title-1",
			posName:        "Direktur Utama",
			reportsToID:    nil,
			headcountQuota: -5,
			wantErr:        nil,
			wantQuota:      1,
		},
		{
			name:         "empty company id is invalid",
			companyID:    "",
			departmentID: "dept-1",
			jobTitleID:   "title-1",
			posName:      "Staf",
			wantErr:      ErrInvalidInput,
		},
		{
			name:         "empty department id is invalid",
			companyID:    "co-1",
			departmentID: "",
			jobTitleID:   "title-1",
			posName:      "Staf",
			wantErr:      ErrInvalidInput,
		},
		{
			name:         "empty job title id is invalid",
			companyID:    "co-1",
			departmentID: "dept-1",
			jobTitleID:   "",
			posName:      "Staf",
			wantErr:      ErrInvalidInput,
		},
		{
			name:         "empty name is invalid",
			companyID:    "co-1",
			departmentID: "dept-1",
			jobTitleID:   "title-1",
			posName:      "",
			wantErr:      ErrInvalidInput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jobPosition, err := NewJobPosition(tt.companyID, tt.departmentID, tt.jobTitleID, tt.posName, tt.reportsToID, tt.headcountQuota)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected error %v, got %v", tt.wantErr, err)
				}
				if jobPosition != nil {
					t.Fatalf("expected nil job position on error, got %+v", jobPosition)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if jobPosition.ID == "" {
				t.Fatal("expected generated ID, got empty string")
			}
			if jobPosition.HeadcountQuota != tt.wantQuota {
				t.Errorf("expected headcount quota %d, got %d", tt.wantQuota, jobPosition.HeadcountQuota)
			}
			if !jobPosition.IsActive {
				t.Error("expected new job position to default IsActive=true")
			}
			if jobPosition.CreatedAt.IsZero() || jobPosition.UpdatedAt.IsZero() {
				t.Error("expected CreatedAt/UpdatedAt to be set")
			}
		})
	}
}

func TestNewJobPosition_GeneratesUniqueIDs(t *testing.T) {
	p1, err := NewJobPosition("co-1", "dept-1", "title-1", "Staf A", nil, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	p2, err := NewJobPosition("co-1", "dept-1", "title-1", "Staf B", nil, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p1.ID == p2.ID {
		t.Fatal("expected different IDs for different job positions")
	}
}

package domain

import (
	"errors"
	"testing"
)

func TestNewCompany(t *testing.T) {
	npwp := "01.234.567.8-901.000"
	bpjs := "JKN-0001"

	tests := []struct {
		name      string
		code      string
		legalName string
		npwp      *string
		bpjsNo    *string
		wantErr   error
	}{
		{
			name:      "happy path with npwp and bpjs",
			code:      "PTA",
			legalName: "PT Alpha Nusantara",
			npwp:      &npwp,
			bpjsNo:    &bpjs,
			wantErr:   nil,
		},
		{
			name:      "happy path without npwp and bpjs (nullable)",
			code:      "PTB",
			legalName: "PT Beta Sejahtera",
			npwp:      nil,
			bpjsNo:    nil,
			wantErr:   nil,
		},
		{
			name:      "empty code is invalid",
			code:      "",
			legalName: "PT Alpha Nusantara",
			wantErr:   ErrInvalidInput,
		},
		{
			name:      "empty legal name is invalid",
			code:      "PTA",
			legalName: "",
			wantErr:   ErrInvalidInput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			company, err := NewCompany(tt.code, tt.legalName, tt.npwp, tt.bpjsNo)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected error %v, got %v", tt.wantErr, err)
				}
				if company != nil {
					t.Fatalf("expected nil company on error, got %+v", company)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if company.ID == "" {
				t.Fatal("expected generated ID, got empty string")
			}
			if company.Code != tt.code {
				t.Errorf("expected code %q, got %q", tt.code, company.Code)
			}
			if company.LegalName != tt.legalName {
				t.Errorf("expected legal name %q, got %q", tt.legalName, company.LegalName)
			}
			if !company.IsActive {
				t.Error("expected new company to default IsActive=true")
			}
			if company.CreatedAt.IsZero() || company.UpdatedAt.IsZero() {
				t.Error("expected CreatedAt/UpdatedAt to be set")
			}
		})
	}
}

func TestNewCompany_GeneratesUniqueIDs(t *testing.T) {
	c1, err := NewCompany("PTA", "PT Alpha Nusantara", nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	c2, err := NewCompany("PTB", "PT Beta Sejahtera", nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c1.ID == c2.ID {
		t.Fatal("expected different IDs for different companies")
	}
}

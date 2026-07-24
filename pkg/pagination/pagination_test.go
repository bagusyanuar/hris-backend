package pagination

import "testing"

func TestRequest_SearchClause(t *testing.T) {
	tests := []struct {
		name       string
		search     string
		columns    []string
		wantClause string
		wantArgs   int
	}{
		{
			name:       "empty search returns empty clause",
			search:     "",
			columns:    []string{"code", "name"},
			wantClause: "",
			wantArgs:   0,
		},
		{
			name:       "empty columns returns empty clause",
			search:     "manajer",
			columns:    nil,
			wantClause: "",
			wantArgs:   0,
		},
		{
			name:       "single column",
			search:     "manajer",
			columns:    []string{"name"},
			wantClause: "(name ILIKE ?)",
			wantArgs:   1,
		},
		{
			name:       "multiple columns",
			search:     "TI",
			columns:    []string{"code", "name"},
			wantClause: "(code ILIKE ? OR name ILIKE ?)",
			wantArgs:   2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := Request{Search: tt.search}
			clause, args := req.SearchClause(tt.columns...)

			if clause != tt.wantClause {
				t.Errorf("expected clause %q, got %q", tt.wantClause, clause)
			}
			if len(args) != tt.wantArgs {
				t.Fatalf("expected %d args, got %d", tt.wantArgs, len(args))
			}
			for _, a := range args {
				if a != "%"+tt.search+"%" {
					t.Errorf("expected arg %q, got %q", "%"+tt.search+"%", a)
				}
			}
		})
	}
}

func TestNewMeta_TotalPages(t *testing.T) {
	tests := []struct {
		name  string
		limit int
		total int64
		want  int64
	}{
		{name: "exact division", limit: 20, total: 40, want: 2},
		{name: "remainder rounds up", limit: 20, total: 41, want: 3},
		{name: "total less than limit", limit: 20, total: 5, want: 1},
		{name: "total zero", limit: 20, total: 0, want: 0},
		{name: "limit zero guarded", limit: 0, total: 10, want: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			meta := NewMeta(Request{Page: 1, Limit: tt.limit}, tt.total)
			if meta.TotalPages != tt.want {
				t.Errorf("expected TotalPages %d, got %d", tt.want, meta.TotalPages)
			}
		})
	}
}

package domain

import (
	"context"
	"errors"
	"testing"
)

// fakeParentLookup simulates repo.FindParentID pakai map in-memory: id -> parentID (nil = root).
func fakeParentLookup(chain map[string]*string) func(ctx context.Context, id string) (*string, error) {
	return func(ctx context.Context, id string) (*string, error) {
		parent, ok := chain[id]
		if !ok {
			return nil, errors.New("not found")
		}
		return parent, nil
	}
}

func strPtr(s string) *string { return &s }

func TestDetectCycle(t *testing.T) {
	// Chain: A(root) <- B <- C <- D
	a, b, c, d := "A", "B", "C", "D"
	chain := map[string]*string{
		a: nil,
		b: &a,
		c: &b,
		d: &c,
	}

	t.Run("no cycle when assigning valid deeper parent", func(t *testing.T) {
		// Assign a new node E to report to D — E doesn't exist in chain yet, no cycle.
		err := DetectCycle(context.Background(), "E", d, fakeParentLookup(chain))
		if err != nil {
			t.Fatalf("expected no cycle, got %v", err)
		}
	})

	t.Run("no cycle when new parent is root", func(t *testing.T) {
		err := DetectCycle(context.Background(), "E", a, fakeParentLookup(chain))
		if err != nil {
			t.Fatalf("expected no cycle, got %v", err)
		}
	})

	t.Run("direct self-reference is a cycle", func(t *testing.T) {
		err := DetectCycle(context.Background(), b, b, fakeParentLookup(chain))
		if !errors.Is(err, ErrHierarchyCycle) {
			t.Fatalf("expected ErrHierarchyCycle, got %v", err)
		}
	})

	t.Run("indirect cycle detected (re-parent ancestor to its own descendant)", func(t *testing.T) {
		// Try to make B report to D, but D is already a descendant of B (B<-C<-D) — cycle.
		err := DetectCycle(context.Background(), b, d, fakeParentLookup(chain))
		if !errors.Is(err, ErrHierarchyCycle) {
			t.Fatalf("expected ErrHierarchyCycle, got %v", err)
		}
	})

	t.Run("re-parenting a leaf to a sibling is not a cycle", func(t *testing.T) {
		// D currently reports to C; re-parent D to report to A instead — not a cycle.
		err := DetectCycle(context.Background(), d, a, fakeParentLookup(chain))
		if err != nil {
			t.Fatalf("expected no cycle, got %v", err)
		}
	})

	t.Run("propagates lookup error", func(t *testing.T) {
		err := DetectCycle(context.Background(), "X", "unknown-node", fakeParentLookup(chain))
		if err == nil || errors.Is(err, ErrHierarchyCycle) {
			t.Fatalf("expected lookup error, got %v", err)
		}
	})

	t.Run("guards against runaway depth on corrupted data", func(t *testing.T) {
		// root points to itself — infinite loop without the maxHierarchyDepth guard,
		// simulates corrupted data where the walk never reaches a nil root.
		selfLoop := map[string]*string{"root": strPtr("root")}
		err := DetectCycle(context.Background(), "never-matches", "root", fakeParentLookup(selfLoop))
		if !errors.Is(err, ErrHierarchyCycle) {
			t.Fatalf("expected ErrHierarchyCycle from depth guard, got %v", err)
		}
	})
}

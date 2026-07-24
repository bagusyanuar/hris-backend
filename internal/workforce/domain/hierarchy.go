package domain

import (
	"context"
	"errors"
)

var ErrHierarchyCycle = errors.New("hierarchy cycle detected")

// maxHierarchyDepth adalah jaring pengaman kalau data korup bikin walk-up
// jalan tak berujung (bug di tempat lain) — org chart realistis jauh lebih
// dangkal dari ini (decision-log.md ADR-002).
const maxHierarchyDepth = 100

// DetectCycle dipakai bareng Department.ParentID & JobPosition.ReportsToID
// (decision-log.md ADR-002) — walk-up parent chain di application layer,
// bukan recursive CTE. candidateID = ID entity yang sedang diubah; newParentID
// = parent/reports_to baru yang mau di-assign; findParentID = repo.FindParentID
// entity terkait (Department atau JobPosition, tergantung pemanggil).
func DetectCycle(ctx context.Context, candidateID, newParentID string, findParentID func(ctx context.Context, id string) (*string, error)) error {
	current := &newParentID
	for depth := 0; depth < maxHierarchyDepth; depth++ {
		if current == nil {
			return nil // sampai root, tidak ada cycle
		}
		if *current == candidateID {
			return ErrHierarchyCycle
		}
		parent, err := findParentID(ctx, *current)
		if err != nil {
			return err
		}
		current = parent
	}
	return ErrHierarchyCycle // depth exceeded — anggap cycle/data korup
}

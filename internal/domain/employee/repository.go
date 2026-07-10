package employee

import "context"

type Repository interface {
	Save(ctx context.Context, emp *Employee) error
	FindByID(ctx context.Context, id string) (*Employee, error)
	FindAll(ctx context.Context) ([]*Employee, error)
	Update(ctx context.Context, emp *Employee) error
	Delete(ctx context.Context, id string) error

	// ExecuteInTx executes the provided function within a database transaction.
	// The transaction context must be passed down to Repository methods inside the function.
	ExecuteInTx(ctx context.Context, fn func(txCtx context.Context) error) error
}

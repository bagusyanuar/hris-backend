package employee

import "context"

type Repository interface {
	SaveCore(ctx context.Context, emp *Employee) error
	FindByID(ctx context.Context, id string) (*Employee, error)

	SavePersonalData(ctx context.Context, data *PersonalData) error
	FindByKTP(ctx context.Context, ktpNumber string) (*PersonalData, error)

	SaveContact(ctx context.Context, contact *Contact) error

	SaveBanks(ctx context.Context, employeeID string, banks []*Bank) error
	SaveEducations(ctx context.Context, employeeID string, educations []*Education) error
	SaveDocument(ctx context.Context, doc *Document) error
}

//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	"gorm.io/gorm"

	"github.com/bagusyanuar/hris-backend/internal/application/auth"
	"github.com/bagusyanuar/hris-backend/internal/application/employee"
	"github.com/bagusyanuar/hris-backend/internal/application/organization"

	"github.com/bagusyanuar/hris-backend/internal/infrastructure/repository"

	httpAuth "github.com/bagusyanuar/hris-backend/internal/interfaces/http/auth"
	httpEmployee "github.com/bagusyanuar/hris-backend/internal/interfaces/http/employee"
	httpOrg "github.com/bagusyanuar/hris-backend/internal/interfaces/http/organization"

	domainAuth "github.com/bagusyanuar/hris-backend/internal/domain/auth"
	userInfra "github.com/bagusyanuar/hris-backend/internal/user/infrastructure"
)

var RepositorySet = wire.NewSet(
	userInfra.NewUserRepository,
	repository.NewOrganizationRepository,
	repository.NewEmployeeRepository,
)

var ServiceSet = wire.NewSet(
	auth.NewService,
	organization.NewService,
	employee.NewService,
)

var HandlerSet = wire.NewSet(
	httpAuth.NewHandler,
	httpOrg.NewHandler,
	httpEmployee.NewHandler,
	wire.Struct(new(APIHandlers), "*"),
)

func InitializeAPI(db *gorm.DB, tokenGen domainAuth.TokenGenerator) (*APIHandlers, error) {
	wire.Build(
		RepositorySet,
		ServiceSet,
		HandlerSet,
	)
	return &APIHandlers{}, nil
}

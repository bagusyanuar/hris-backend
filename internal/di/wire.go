//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	"gorm.io/gorm"

	"github.com/bagusyanuar/hris-backend/internal/application/employee"
	"github.com/bagusyanuar/hris-backend/internal/application/organization"
	authApp "github.com/bagusyanuar/hris-backend/internal/auth/application"

	"github.com/bagusyanuar/hris-backend/internal/infrastructure/repository"

	authHTTP "github.com/bagusyanuar/hris-backend/internal/auth/transport/http"
	httpEmployee "github.com/bagusyanuar/hris-backend/internal/interfaces/http/employee"
	httpOrg "github.com/bagusyanuar/hris-backend/internal/interfaces/http/organization"

	authDomain "github.com/bagusyanuar/hris-backend/internal/auth/domain"
	userInfra "github.com/bagusyanuar/hris-backend/internal/user/infrastructure"
)

var RepositorySet = wire.NewSet(
	userInfra.NewUserRepository,
	repository.NewOrganizationRepository,
	repository.NewEmployeeRepository,
)

var ServiceSet = wire.NewSet(
	authApp.NewService,
	organization.NewService,
	employee.NewService,
)

var HandlerSet = wire.NewSet(
	authHTTP.NewHandler,
	httpOrg.NewHandler,
	httpEmployee.NewHandler,
	wire.Struct(new(APIHandlers), "*"),
)

func InitializeAPI(db *gorm.DB, tokenGen authDomain.TokenGenerator) (*APIHandlers, error) {
	wire.Build(
		RepositorySet,
		ServiceSet,
		HandlerSet,
	)
	return &APIHandlers{}, nil
}

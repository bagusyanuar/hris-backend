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
	userAdapter "github.com/bagusyanuar/hris-backend/internal/user/adapter"

	orgAdapter "github.com/bagusyanuar/hris-backend/internal/organization/adapter"
	orgApp "github.com/bagusyanuar/hris-backend/internal/organization/application"
	orgHTTP "github.com/bagusyanuar/hris-backend/internal/organization/transport/http"
)

var RepositorySet = wire.NewSet(
	userAdapter.NewUserRepository,
	repository.NewOrganizationRepository,
	repository.NewEmployeeRepository,
	orgAdapter.NewCompanyRepository,
	orgAdapter.NewBranchRepository,
	orgAdapter.NewGormTxManager,
)

var ServiceSet = wire.NewSet(
	authApp.NewService,
	organization.NewService,
	employee.NewService,
	orgApp.NewService,
)

var HandlerSet = wire.NewSet(
	authHTTP.NewHandler,
	httpOrg.NewHandler,
	httpEmployee.NewHandler,
	orgHTTP.NewHandler,
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

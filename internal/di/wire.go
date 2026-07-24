//go:build wireinject
// +build wireinject

package di

import (
	"github.com/google/wire"
	"gorm.io/gorm"

	"github.com/bagusyanuar/hris-backend/internal/application/employee"
	authApp "github.com/bagusyanuar/hris-backend/internal/auth/application"

	"github.com/bagusyanuar/hris-backend/internal/infrastructure/repository"

	authHTTP "github.com/bagusyanuar/hris-backend/internal/auth/transport/http"
	httpEmployee "github.com/bagusyanuar/hris-backend/internal/interfaces/http/employee"

	authDomain "github.com/bagusyanuar/hris-backend/internal/auth/domain"
	userAdapter "github.com/bagusyanuar/hris-backend/internal/user/adapter"

	orgAdapter "github.com/bagusyanuar/hris-backend/internal/organization/adapter"
	orgApp "github.com/bagusyanuar/hris-backend/internal/organization/application"
	orgHTTP "github.com/bagusyanuar/hris-backend/internal/organization/transport/http"

	workforceAdapter "github.com/bagusyanuar/hris-backend/internal/workforce/adapter"
	workforceApp "github.com/bagusyanuar/hris-backend/internal/workforce/application"
	workforceHTTP "github.com/bagusyanuar/hris-backend/internal/workforce/transport/http"
)

var RepositorySet = wire.NewSet(
	userAdapter.NewUserRepository,
	repository.NewEmployeeRepository,
	orgAdapter.NewCompanyRepository,
	orgAdapter.NewBranchRepository,
	orgAdapter.NewGormTxManager,
	workforceAdapter.NewDepartmentRepository,
	workforceAdapter.NewJobTitleRepository,
	workforceAdapter.NewJobPositionRepository,
	workforceAdapter.NewGormTxManager,
)

var ServiceSet = wire.NewSet(
	authApp.NewService,
	employee.NewService,
	orgApp.NewService,
	workforceApp.NewService,
)

var HandlerSet = wire.NewSet(
	authHTTP.NewHandler,
	httpEmployee.NewHandler,
	orgHTTP.NewHandler,
	workforceHTTP.NewHandler,
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

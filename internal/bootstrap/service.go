package bootstrap

import (
	"go-playground/internal/domain"
	"go-playground/internal/service"
)

// Services holds all service instances
type Services struct {
	UserService              *service.UserService
	AuthService              *service.AuthService
	PointsService            domain.PointsService
	TransactionService       *service.TransactionService
	RewardsService           *service.RewardsService
	RedemptionService        *service.RedemptionService
	MerchantService          *service.MerchantService
	MerchantCustomersService *service.MerchantCustomersService
	ProgramService           *service.ProgramService
	ProgramRuleService       *service.ProgramRulesService
}

// InitializeServices initializes all services
func InitializeServices(repos *Repositories) *Services {
	pointsService := service.NewPointsService(repos.PointsRepo, repos.EventRepo)
	legacyPointsService := service.NewLegacyPointsService(pointsService)
	merchantService := service.NewMerchantService(repos.MerchantRepo)
	transactionService := service.NewTransactionService(
		repos.TransactionRepo,
		pointsService,
		repos.EventRepo,
		repos.MerchantCustomersRepo,
	)

	return &Services{
		UserService:        service.NewUserService(repos.UserRepo, repos.CacheRepo),
		AuthService:        service.NewAuthService(repos.UserRepo, repos.AuthRepo, repos.SessionRepo),
		PointsService:      legacyPointsService,
		TransactionService: transactionService,
		RewardsService:     service.NewRewardsService(repos.RewardsRepo),
		RedemptionService: service.NewRedemptionService(
			repos.RedemptionRepo,
			repos.RewardsRepo,
			pointsService,
			transactionService,
			repos.EventRepo,
		),
		MerchantService:          merchantService,
		MerchantCustomersService: service.NewMerchantCustomersService(repos.MerchantCustomersRepo),
		ProgramService:           service.NewProgramService(repos.ProgramRepo),
		ProgramRuleService:       service.NewProgramRulesService(repos.ProgramRuleRepo),
	}
}

package bootstrap

import (
	"go-playground/internal/service"
)

// Services holds all service instances
type Services struct {
	UserService              *service.UserService
	AuthService              *service.AuthService
	PointsService            *service.PointsService
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
	merchantService := service.NewMerchantService(repos.MerchantRepo)
	eventLoggerService := service.NewEventLoggerService(repos.EventRepo)
	transactionService := service.NewTransactionService(
		repos.TransactionRepo,
		pointsService,
		eventLoggerService,
		repos.MerchantCustomersRepo,
	)
	redemptionService := service.NewRedemptionService(
		repos.RedemptionRepo,
		repos.RewardsRepo,
		pointsService,
		transactionService,
		eventLoggerService,
	)

	return &Services{
		UserService: service.NewUserService(
			repos.UserRepo,
			repos.CacheRepo,
		),
		AuthService: service.NewAuthService(
			repos.UserRepo,
			repos.AuthRepo,
			repos.SessionRepo,
		),
		PointsService:            pointsService,
		TransactionService:       transactionService,
		RewardsService:           service.NewRewardsService(repos.RewardsRepo),
		RedemptionService:        redemptionService,
		MerchantService:          merchantService,
		MerchantCustomersService: service.NewMerchantCustomersService(repos.MerchantCustomersRepo),
		ProgramService:           service.NewProgramService(repos.ProgramRepo),
		ProgramRuleService:       service.NewProgramRulesService(repos.ProgramRuleRepo, repos.ProgramRepo),
	}
}

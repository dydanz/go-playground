package bootstrap

import (
	"go-playground/internal/domain"
	"go-playground/internal/service"
)

// Services holds all service instances
type Services struct {
	UserService        *service.UserService
	AuthService        *service.AuthService
	PointsService      domain.PointsService
	TransactionService *service.TransactionService
	RewardsService     *service.RewardsService
	RedemptionService  *service.RedemptionService
	MerchantService    *service.MerchantService
	ProgramService     *service.ProgramService
	ProgramRuleService *service.ProgramRulesService
}

// InitializeServices initializes all services
func InitializeServices(repos *Repositories) *Services {
	pointsService := service.NewPointsService(repos.PointsRepo, repos.EventRepo)
	legacyPointsService := service.NewLegacyPointsService(pointsService)

	return &Services{
		UserService:        service.NewUserService(repos.UserRepo, repos.CacheRepo),
		AuthService:        service.NewAuthService(repos.UserRepo, repos.AuthRepo, repos.SessionRepo),
		PointsService:      legacyPointsService,
		TransactionService: service.NewTransactionService(repos.TransactionRepo, pointsService, repos.EventRepo),
		RewardsService:     service.NewRewardsService(repos.RewardsRepo),
		RedemptionService:  service.NewRedemptionService(repos.RedemptionRepo, repos.RewardsRepo, pointsService, repos.EventRepo),
		MerchantService:    service.NewMerchantService(repos.MerchantRepo),
		ProgramService:     service.NewProgramService(repos.ProgramRepo),
		ProgramRuleService: service.NewProgramRulesService(repos.ProgramRuleRepo),
	}
}

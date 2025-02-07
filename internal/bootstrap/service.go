package bootstrap

import (
	"go-playground/internal/service"
)

// Services holds all service instances
type Services struct {
	UserService        *service.UserService
	AuthService        *service.AuthService
	PointsService      *service.PointsService
	TransactionService *service.TransactionService
	RewardsService     *service.RewardsService
	RedemptionService  *service.RedemptionService
	MerchantService    *service.MerchantService
	ProgramService     *service.ProgramService
	ProgramRuleService *service.ProgramRulesService
}

// InitializeServices initializes all services
func InitializeServices(repos *Repositories) *Services {
	return &Services{
		UserService:        service.NewUserService(repos.UserRepo, repos.CacheRepo),
		AuthService:        service.NewAuthService(repos.UserRepo, repos.AuthRepo, repos.SessionRepo),
		PointsService:      service.NewPointsService(repos.PointsRepo, repos.EventRepo),
		TransactionService: service.NewTransactionService(repos.TransactionRepo, nil, repos.EventRepo), // PointsService will be set after initialization
		RewardsService:     service.NewRewardsService(repos.RewardsRepo),
		RedemptionService:  service.NewRedemptionService(repos.RedemptionRepo, repos.RewardsRepo, nil, repos.EventRepo), // PointsService will be set after initialization
		MerchantService:    service.NewMerchantService(repos.MerchantRepo),
		ProgramService:     service.NewProgramService(repos.ProgramRepo),
		ProgramRuleService: service.NewProgramRulesService(repos.ProgramRuleRepo),
	}
}

// SetupServiceDependencies sets up dependencies between services that couldn't be set during initialization
func (s *Services) SetupServiceDependencies() {
	// Set PointsService in services that depend on it
	// s.TransactionService.SetPointsService(s.PointsService)
	// s.RedemptionService.SetPointsService(s.PointsService)
}

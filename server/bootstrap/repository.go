package bootstrap

import (
	"database/sql"
	"go-playground/server/config"
	"go-playground/server/repository/postgres"
	"go-playground/server/repository/redis"

	redislib "github.com/go-redis/redis/v8"
)

// Repositories holds all repository instances
type Repositories struct {
	UserRepo              *postgres.UserRepository
	CacheRepo             *redis.CacheRepository
	AuthRepo              *postgres.AuthRepository
	PointsRepo            *postgres.PointsRepository
	TransactionRepo       *postgres.TransactionRepository
	RewardsRepo           *postgres.RewardsRepository
	RedemptionRepo        *postgres.RedemptionRepository
	EventRepo             *postgres.EventLogRepository
	MerchantRepo          *postgres.MerchantRepository
	MerchantCustomersRepo *postgres.MerchantCustomersRepository
	ProgramRepo           *postgres.ProgramsRepository
	SessionRepo           redis.SessionRepository
	ProgramRuleRepo       *postgres.ProgramRuleRepository
}

// InitializeRepositories initializes all repositories
func InitializeRepositories(db *sql.DB, dbConn *config.DbConnection, rdb *redislib.Client, cfg *config.Config) *Repositories {
	return &Repositories{
		UserRepo:              postgres.NewUserRepository(db),
		CacheRepo:             redis.NewCacheRepository(rdb),
		AuthRepo:              postgres.NewAuthRepository(db, &cfg.Auth),
		PointsRepo:            postgres.NewPointsRepository(db),
		TransactionRepo:       postgres.NewTransactionRepository(*dbConn),
		RewardsRepo:           postgres.NewRewardsRepository(db),
		RedemptionRepo:        postgres.NewRedemptionRepository(db),
		EventRepo:             postgres.NewEventLogRepository(db),
		MerchantRepo:          postgres.NewMerchantRepository(db),
		MerchantCustomersRepo: postgres.NewMerchantCustomersRepository(db),
		ProgramRepo:           postgres.NewProgramsRepository(db),
		SessionRepo:           redis.NewSessionRepository(rdb),
		ProgramRuleRepo:       postgres.NewProgramRuleRepository(*dbConn),
	}
}

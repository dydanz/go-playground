package bootstrap

import (
	"database/sql"
	"go-playground/internal/config"
	"go-playground/internal/repository/postgres"
	"go-playground/internal/repository/redis"

	redislib "github.com/go-redis/redis/v8"
)

// Repositories holds all repository instances
type Repositories struct {
	UserRepo        *postgres.UserRepository
	CacheRepo       *redis.CacheRepository
	AuthRepo        *postgres.AuthRepository
	PointsRepo      *postgres.PointsRepository
	TransactionRepo *postgres.TransactionRepository
	RewardsRepo     *postgres.RewardsRepository
	RedemptionRepo  *postgres.RedemptionRepository
	EventRepo       *postgres.EventLogRepository
	MerchantRepo    *postgres.MerchantRepository
	ProgramRepo     *postgres.ProgramsRepository
	SessionRepo     redis.SessionRepository
}

// InitializeRepositories initializes all repositories
func InitializeRepositories(db *sql.DB, dbConn *config.DbConnection, rdb *redislib.Client, cfg *config.Config) *Repositories {
	return &Repositories{
		UserRepo:        postgres.NewUserRepository(db),
		CacheRepo:       redis.NewCacheRepository(rdb),
		AuthRepo:        postgres.NewAuthRepository(db, &cfg.Auth),
		PointsRepo:      postgres.NewPointsRepository(db),
		TransactionRepo: postgres.NewTransactionRepository(*dbConn),
		RewardsRepo:     postgres.NewRewardsRepository(db),
		RedemptionRepo:  postgres.NewRedemptionRepository(db),
		EventRepo:       postgres.NewEventLogRepository(db),
		MerchantRepo:    postgres.NewMerchantRepository(db),
		ProgramRepo:     postgres.NewProgramsRepository(db),
		SessionRepo:     redis.NewSessionRepository(rdb),
	}
}

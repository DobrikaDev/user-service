package sql

import (
	"DobrikaDev/user-service/internal/storage/deps"
	"DobrikaDev/user-service/utils/config"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type SqlStorage struct {
	trf deps.TransactionFactory
	deps.TransactionManager
	logger *zap.Logger
}

func NewStorage(trf deps.TransactionFactory, trm deps.TransactionManager, logger *zap.Logger) *SqlStorage {
	return &SqlStorage{
		trf:                trf,
		TransactionManager: trm,
		logger:             logger,
	}
}

func buildDSN(cfg *config.Config) string {
	return fmt.Sprintf(
		"user=%s host=%s port=%d password=%s dbname=%s sslmode=disable",
		cfg.SQL.Username, cfg.SQL.Host, cfg.SQL.Port, cfg.SQL.Password, cfg.SQL.Name,
	)
}

func NewPostgresDB(cfg *config.Config) (*sqlx.DB, error) {
	dsn := buildDSN(cfg)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	db.SetConnMaxIdleTime(time.Second * 30)
	db.SetConnMaxLifetime(time.Minute * 15)
	sqlxDb := sqlx.NewDb(db, "pgx")

	return sqlxDb, nil
}

func MustCreateDB(cfg *config.Config) *sqlx.DB {
	db, err := NewPostgresDB(cfg)
	if err != nil {
		panic(err)
	}

	return db
}

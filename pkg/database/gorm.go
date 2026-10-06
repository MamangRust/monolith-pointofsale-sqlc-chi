package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// NewGormClient creates a *gorm.DB connection from the viper DB_*
// configuration. It is the single database client for all services.
func NewGormClient(logger logger.LoggerInterface) (*gorm.DB, error) {
	dbDriver := viper.GetString("DB_DRIVER")
	if dbDriver != "postgres" && dbDriver != "" {
		logger.Error("GORM postgres driver only supports PostgreSQL", zap.String("DB_DRIVER", dbDriver))
		return nil, fmt.Errorf("gorm postgres driver only supports PostgreSQL, got: %s", dbDriver)
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		viper.GetString("DB_HOST"),
		viper.GetString("DB_PORT"),
		viper.GetString("DB_USERNAME"),
		viper.GetString("DB_NAME"),
		viper.GetString("DB_PASSWORD"),
	)

	maxOpenConns := viper.GetInt("DB_MAX_OPEN_CONNS")
	if maxOpenConns <= 0 {
		maxOpenConns = 100
	}

	maxIdleConns := viper.GetInt("DB_MIN_IDLE_CONNS")
	if maxIdleConns <= 0 {
		maxIdleConns = 50
	}

	connMaxLifetime := viper.GetDuration("DB_CONN_MAX_LIFETIME")
	if connMaxLifetime == 0 {
		connMaxLifetime = time.Hour
	}

	connMaxIdleTime := viper.GetDuration("DB_CONN_MAX_IDLE_TIME")
	if connMaxIdleTime == 0 {
		connMaxIdleTime = 30 * time.Minute
	}

	slowThreshold := viper.GetDuration("DB_SLOW_THRESHOLD")
	if slowThreshold <= 0 {
		slowThreshold = 200 * time.Millisecond
	}

	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:         newZapGormLogger(logger, slowThreshold),
		TranslateError: true,
	})
	if err != nil {
		logger.Error("Failed to open GORM connection", zap.Error(err))
		return nil, fmt.Errorf("failed to open GORM connection: %w", err)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		logger.Error("Failed to get underlying sql.DB from GORM", zap.Error(err))
		return nil, fmt.Errorf("failed to get underlying sql.DB from GORM: %w", err)
	}

	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetConnMaxLifetime(connMaxLifetime)
	sqlDB.SetConnMaxIdleTime(connMaxIdleTime)

	if err := sqlDB.Ping(); err != nil {
		logger.Error("Failed to ping database via GORM", zap.Error(err))
		return nil, fmt.Errorf("failed to ping database via GORM: %w", err)
	}

	logger.Debug("GORM database connection established successfully",
		zap.Int("MaxOpenConns", maxOpenConns),
		zap.Int("MaxIdleConns", maxIdleConns),
		zap.Duration("ConnMaxLifetime", connMaxLifetime),
		zap.Duration("ConnMaxIdleTime", connMaxIdleTime),
		zap.Duration("SlowThreshold", slowThreshold),
	)

	return gormDB, nil
}

// zapGormLogger bridges GORM's logger interface onto the shared zap logger.
// Successful fast queries stay silent; slow queries log as warnings and
// failures as errors, keeping service logs readable while surfacing drift.
type zapGormLogger struct {
	log           logger.LoggerInterface
	slowThreshold time.Duration
}

func newZapGormLogger(log logger.LoggerInterface, slowThreshold time.Duration) gormlogger.Interface {
	return &zapGormLogger{log: log, slowThreshold: slowThreshold}
}

func (l *zapGormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	return l
}

func (l *zapGormLogger) Info(ctx context.Context, message string, args ...interface{}) {
	l.log.Info(fmt.Sprintf(message, args...))
}

func (l *zapGormLogger) Warn(ctx context.Context, message string, args ...interface{}) {
	l.log.Warn(fmt.Sprintf(message, args...))
}

func (l *zapGormLogger) Error(ctx context.Context, message string, args ...interface{}) {
	l.log.Error(fmt.Sprintf(message, args...))
}

func (l *zapGormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)

	switch {
	case err != nil && !errors.Is(err, gorm.ErrRecordNotFound):
		sql, rows := fc()
		l.log.Error("gorm query failed",
			zap.Error(err),
			zap.Duration("elapsed", elapsed),
			zap.Int64("rows", rows),
			zap.String("sql", sql),
		)
	case l.slowThreshold > 0 && elapsed > l.slowThreshold:
		sql, rows := fc()
		l.log.Warn("gorm slow query",
			zap.Duration("elapsed", elapsed),
			zap.Duration("threshold", l.slowThreshold),
			zap.Int64("rows", rows),
			zap.String("sql", sql),
		)
	}
}

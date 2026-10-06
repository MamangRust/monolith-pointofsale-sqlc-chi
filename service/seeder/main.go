package main

import (
	"context"
	"log"

	"github.com/MamangRust/monolith-point-of-sale-pkg/database"
	"github.com/MamangRust/monolith-point-of-sale-pkg/database/seeder"
	"github.com/MamangRust/monolith-point-of-sale-pkg/dotenv"
	"github.com/MamangRust/monolith-point-of-sale-pkg/hash"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"

	"go.uber.org/zap"
)

func main() {
	if err := dotenv.Viper(); err != nil {
		log.Fatalf("Error loading environment variables: %v", err)
	}

	l, err := logger.NewLogger("seeder", nil)
	if err != nil {
		log.Fatalf("Error creating logger: %v", err)
	}

	gormDB, err := database.NewGormClient(l)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer func() {
		if sqlDB, err := gormDB.DB(); err == nil {
			if err := sqlDB.Close(); err != nil {
				l.Error("Failed to close GORM database connection", zap.Error(err))
			} else {
				l.Info("GORM database connection closed")
			}
		}
	}()

	ctx := context.Background()

	s := seeder.NewSeeder(seeder.Deps{
		Db:     gormDB,
		Ctx:    ctx,
		Logger: l,
		Hash:   hash.NewHashingPassword(),
	})

	if err := s.Run(); err != nil {
		log.Fatalf("Seeding failed: %v", err)
	}

	l.Info("Seeding completed successfully.")
}

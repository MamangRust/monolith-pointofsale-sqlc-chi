package seeder

import (
	"context"
	"fmt"

	"github.com/MamangRust/monolith-point-of-sale-pkg/database/models"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type roleSeeder struct {
	db     *gorm.DB
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewRoleSeeder(db *gorm.DB, ctx context.Context, logger logger.LoggerInterface) *roleSeeder {
	return &roleSeeder{
		db:     db,
		ctx:    ctx,
		logger: logger,
	}
}

func (r *roleSeeder) Seed() error {
	randomRoles := []string{
		"ROLE_ADMIN",
		"Admin Access 1",
		"Super Admin",
		"Admin",
		"Store Manager",
		"Cashier",
		"Inventory Staff",
		"Support",
		"Auditor",
		"Viewer",
	}

	totalRoles := len(randomRoles)

	for i, roleName := range randomRoles {
		role := &models.Role{
			RoleName: roleName,
		}

		if err := r.db.WithContext(r.ctx).Create(role).Error; err != nil {
			r.logger.Error("failed to seed role", zap.Int("role", i+1), zap.String("roleName", roleName), zap.Error(err))
			return fmt.Errorf("failed to seed role %d (%s): %w", i+1, roleName, err)
		}
	}

	r.logger.Info("role seeded successfully", zap.Int("totalRoles", totalRoles))
	return nil
}

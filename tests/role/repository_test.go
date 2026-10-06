package role_test

import (
	"context"
	"testing"

	"github.com/MamangRust/monolith-point-of-sale-role/repository"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	tests "github.com/MamangRust/monolith-point-of-sale-test"

	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type RoleRepositoryTestSuite struct {
	suite.Suite
	ts     *tests.TestSuite
	repo   repository.Repositories
	gormDB *gorm.DB
	roleID int
}

func (s *RoleRepositoryTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	gormDB, err := ts.GormDB()
	s.Require().NoError(err)
	s.gormDB = gormDB

	s.repo = repository.NewRepositories(gormDB)
}

func (s *RoleRepositoryTestSuite) TearDownSuite() {
	if sqlDB, err := s.gormDB.DB(); err == nil {
		sqlDB.Close()
	}
	s.ts.Teardown()
}

func (s *RoleRepositoryTestSuite) Test1_CreateRole() {
	ctx := context.Background()
	req := &requests.CreateRoleRequest{
		Name: "Test Role",
	}

	res, err := s.repo.CreateRole(ctx, req)
	s.NoError(err)
	s.NotNil(res)
	s.Equal(req.Name, res.RoleName)
	s.roleID = int(res.RoleID)
}

func (s *RoleRepositoryTestSuite) Test2_FindById() {
	s.Require().NotZero(s.roleID)
	ctx := context.Background()

	found, err := s.repo.FindById(ctx, s.roleID)
	s.NoError(err)
	s.NotNil(found)
	s.Equal(s.roleID, int(found.RoleID))
	s.Equal("Test Role", found.RoleName)
}

func (s *RoleRepositoryTestSuite) Test3_UpdateRole() {
	s.Require().NotZero(s.roleID)
	ctx := context.Background()

	req := &requests.UpdateRoleRequest{
		ID:   &s.roleID,
		Name: "Updated Role",
	}

	res, err := s.repo.UpdateRole(ctx, req)
	s.NoError(err)
	s.NotNil(res)
	s.Equal("Updated Role", res.RoleName)
}

func (s *RoleRepositoryTestSuite) Test4_TrashAndRestore() {
	s.Require().NotZero(s.roleID)
	ctx := context.Background()

	// Trash
	trashed, err := s.repo.TrashedRole(ctx, s.roleID)
	s.NoError(err)
	s.NotNil(trashed)

	// Restore
	restored, err := s.repo.RestoreRole(ctx, s.roleID)
	s.NoError(err)
	s.NotNil(restored)

	// Verify restored
	found, err := s.repo.FindById(ctx, s.roleID)
	s.NoError(err)
	s.NotNil(found)
}

func (s *RoleRepositoryTestSuite) Test5_DeletePermanent() {
	s.Require().NotZero(s.roleID)
	ctx := context.Background()

	// Must be trashed first for permanent delete
	_, err := s.repo.TrashedRole(ctx, s.roleID)
	s.NoError(err)

	success, err := s.repo.DeleteRolePermanent(ctx, s.roleID)
	s.NoError(err)
	s.True(success)
}

func TestRoleRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(RoleRepositoryTestSuite))
}

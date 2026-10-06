package category_test

import (
	"context"
	"testing"

	"github.com/MamangRust/monolith-point-of-sale-category/repository"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	tests "github.com/MamangRust/monolith-point-of-sale-test"

	"github.com/stretchr/testify/suite"
)

type CategoryRepositoryTestSuite struct {
	suite.Suite
	ts         *tests.TestSuite
	repo       repository.Repositories
	categoryID int
}

func (s *CategoryRepositoryTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	gormDB, err := s.ts.GormDB()
	s.Require().NoError(err)

	s.repo = repository.NewRepositories(gormDB)
}

func (s *CategoryRepositoryTestSuite) TearDownSuite() {
	s.ts.Teardown()
}

func (s *CategoryRepositoryTestSuite) Test1_CreateCategory() {
	ctx := context.Background()

	req := &requests.CreateCategoryRequest{
		Name:        "Electronics",
		Description: "Electronic items and gadgets",
	}

	res, err := s.repo.CreateCategory(ctx, req)
	s.NoError(err)
	s.NotNil(res)
	s.Equal(req.Name, res.Name)
	s.categoryID = int(res.CategoryID)
}

func (s *CategoryRepositoryTestSuite) Test2_FindById() {
	s.Require().NotZero(s.categoryID)
	ctx := context.Background()

	found, err := s.repo.FindById(ctx, s.categoryID)
	s.NoError(err)
	s.NotNil(found)
	s.Equal(s.categoryID, int(found.CategoryID))
	s.Equal("Electronics", found.Name)
}

func (s *CategoryRepositoryTestSuite) Test3_UpdateCategory() {
	s.Require().NotZero(s.categoryID)
	ctx := context.Background()

	req := &requests.UpdateCategoryRequest{
		CategoryID:  &s.categoryID,
		Name:        "Updated Electronics",
		Description: "Updated description",
	}

	res, err := s.repo.UpdateCategory(ctx, req)
	s.NoError(err)
	s.NotNil(res)
	s.Equal("Updated Electronics", res.Name)
}

func (s *CategoryRepositoryTestSuite) Test4_TrashAndRestore() {
	s.Require().NotZero(s.categoryID)
	ctx := context.Background()

	// Trash
	trashed, err := s.repo.TrashedCategory(ctx, s.categoryID)
	s.NoError(err)
	s.NotNil(trashed)

	// Verify trashed via query
	trashedFound, err := s.repo.FindByTrashed(ctx, &requests.FindAllCategory{Search: "", Page: 1, PageSize: 10})
	s.NoError(err)
	s.NotEmpty(trashedFound)

	// Restore
	restored, err := s.repo.RestoreCategory(ctx, s.categoryID)
	s.NoError(err)
	s.NotNil(restored)

	// Verify restored
	found, err := s.repo.FindById(ctx, s.categoryID)
	s.NoError(err)
	s.NotNil(found)
}

func (s *CategoryRepositoryTestSuite) Test5_DeletePermanent() {
	s.Require().NotZero(s.categoryID)
	ctx := context.Background()

	// Must be trashed first for permanent delete
	_, err := s.repo.TrashedCategory(ctx, s.categoryID)
	s.NoError(err)

	success, err := s.repo.DeleteCategoryPermanently(ctx, s.categoryID)
	s.NoError(err)
	s.True(success)
}

func TestCategoryRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(CategoryRepositoryTestSuite))
}

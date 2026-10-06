package role_test

import (
	"context"
	pbroles "github.com/MamangRust/monolith-point-of-sale-pb/roles"
	"testing"

	role_cache "github.com/MamangRust/monolith-point-of-sale-role/cache"
	role_handler "github.com/MamangRust/monolith-point-of-sale-role/handler"
	"github.com/MamangRust/monolith-point-of-sale-role/repository"
	"github.com/MamangRust/monolith-point-of-sale-role/service"
	"github.com/MamangRust/monolith-point-of-sale-shared/cache"
	"github.com/MamangRust/monolith-point-of-sale-shared/observability"
	tests "github.com/MamangRust/monolith-point-of-sale-test"

	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type RoleGapiTestSuite struct {
	tests.BaseTestSuite
	client        pbroles.RoleQueryServiceClient
	clientCommand pbroles.RoleCommandServiceClient
}

func (s *RoleGapiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()

	// Infrastructure
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(s.RedisClient(), s.Log, cacheMetrics)
	gormDB, err := s.GormDB()
	s.Require().NoError(err)

	// Role dependencies
	mencache := role_cache.NewMencache(cacheStore)
	repos := repository.NewRepositories(gormDB)
	svc := service.NewService(&service.Deps{
		Repositories:  repos,
		Logger:        s.Log,
		Mencache:      mencache,
		Observability: s.Obs,
	})

	// Handler
	handler := role_handler.NewHandler(&role_handler.Deps{
		Service: svc,
		Logger:  s.Log,
	})

	// Server
	server := grpc.NewServer()
	pbroles.RegisterRoleQueryServiceServer(server, handler.Role)
	pbroles.RegisterRoleCommandServiceServer(server, handler.RoleCommand)

	addr := s.RegisterServer(server)
	conn := s.GetConnection(addr)

	s.client = pbroles.NewRoleQueryServiceClient(conn)
	s.clientCommand = pbroles.NewRoleCommandServiceClient(conn)
}

func (s *RoleGapiTestSuite) TestRoleGapiLifecycle() {
	ctx := context.Background()

	// 1. Create
	createRes, err := s.clientCommand.CreateRole(ctx, &pbroles.CreateRoleRequest{
		Name: "Gapi Role",
	})
	s.Require().NoError(err)
	s.Require().NotNil(createRes)
	roleID := createRes.Data.Id

	// 2. FindById
	getRes, err := s.client.FindByIdRole(ctx, &pbroles.FindByIdRoleRequest{RoleId: roleID})
	s.Require().NoError(err)
	s.Equal("Gapi Role", getRes.Data.Name)

	// 3. FindAll
	allRes, err := s.client.FindAllRole(ctx, &pbroles.FindAllRoleRequest{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.NotEmpty(allRes.Data)

	// 4. FindByActive
	activeRes, err := s.client.FindByActive(ctx, &pbroles.FindAllRoleRequest{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.NotEmpty(activeRes.Data)

	// 5. Update
	updateRes, err := s.clientCommand.UpdateRole(ctx, &pbroles.UpdateRoleRequest{
		Id:   roleID,
		Name: "Gapi Role Updated",
	})
	s.Require().NoError(err)
	s.Equal("Gapi Role Updated", updateRes.Data.Name)

	// 6. Trash
	_, err = s.clientCommand.TrashedRole(ctx, &pbroles.FindByIdRoleRequest{RoleId: roleID})
	s.Require().NoError(err)

	// 7. FindByTrashed
	trashedRes, err := s.client.FindByTrashed(ctx, &pbroles.FindAllRoleRequest{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.NotEmpty(trashedRes.Data)

	// 8. Restore
	_, err = s.clientCommand.RestoreRole(ctx, &pbroles.FindByIdRoleRequest{RoleId: roleID})
	s.Require().NoError(err)

	// 9. DeletePermanent
	_, _ = s.clientCommand.TrashedRole(ctx, &pbroles.FindByIdRoleRequest{RoleId: roleID})
	_, err = s.clientCommand.DeleteRolePermanent(ctx, &pbroles.FindByIdRoleRequest{RoleId: roleID})
	s.Require().NoError(err)

	// 10. RestoreAll
	_, err = s.clientCommand.RestoreAllRole(ctx, &emptypb.Empty{})
	s.Require().NoError(err)

	// 11. DeleteAll
	_, err = s.clientCommand.DeleteAllRolePermanent(ctx, &emptypb.Empty{})
	s.Require().NoError(err)
}

func TestRoleGapiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(RoleGapiTestSuite))
}

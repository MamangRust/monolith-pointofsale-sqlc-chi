package user_test

import (
	"context"
	pbusers "github.com/MamangRust/monolith-point-of-sale-pb/users"
	"testing"

	"github.com/MamangRust/monolith-point-of-sale-pkg/hash"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	tests "github.com/MamangRust/monolith-point-of-sale-test"
	gapi "github.com/MamangRust/monolith-point-of-sale-user/handler"
	"github.com/MamangRust/monolith-point-of-sale-user/service"

	user_cache "github.com/MamangRust/monolith-point-of-sale-user/cache"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserGapiTestSuite struct {
	tests.BaseTestSuite
	client      pbusers.UserCommandServiceClient
	queryClient pbusers.UserQueryServiceClient
	userID      int
}

func (s *UserGapiTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()

	s.SetupRoleService()

	gormDB, err := s.GormDB()
	s.Require().NoError(err)

	// Seed default role required by user service CreateUser
	err = gormDB.WithContext(s.Ctx).Exec(
		`INSERT INTO roles (role_name, created_at, updated_at)
		 VALUES ('Admin Access 1', current_timestamp, current_timestamp)
		 ON CONFLICT (role_name) DO NOTHING`).Error
	s.Require().NoError(err)

	repos := s.UserRepositories()

	log, _ := logger.NewLogger("test", nil)
	hasher := hash.NewHashingPassword()
	cacheStore := s.GetCacheStore()
	mencache := user_cache.NewMencache(cacheStore)

	userService := service.NewService(&service.Deps{
		Repositories:  repos,
		Logger:        log,
		Hash:          hasher,
		Mencache:      mencache,
		Observability: s.Obs,
	})

	// Start gRPC Server
	userHandler := gapi.NewHandler(&gapi.Deps{
		Service: userService,
		Logger:  log,
	})
	server := grpc.NewServer()
	pbusers.RegisterUserQueryServiceServer(server, userHandler.User)
	pbusers.RegisterUserCommandServiceServer(server, userHandler.UserCommand)

	addr := s.RegisterServer(server)
	conn := s.GetConnection(addr)

	s.client = pbusers.NewUserCommandServiceClient(conn)
	s.queryClient = pbusers.NewUserQueryServiceClient(conn)
}

func (s *UserGapiTestSuite) TestUserGapiLifecycle() {
	ctx := context.Background()

	// 1. Create
	createReq := &pbusers.CreateUserRequest{
		Firstname:       "Gapi",
		Lastname:        "User",
		Email:           "gapi.user@example.com",
		Password:        "password123",
		ConfirmPassword: "password123",
	}
	res, err := s.client.Create(ctx, createReq)
	s.Require().NoError(err)
	s.Equal(createReq.Email, res.Data.Email)
	userID := res.Data.Id

	// 2. FindById
	getRes, err := s.queryClient.FindById(ctx, &pbusers.FindByIdUserRequest{Id: userID})
	s.Require().NoError(err)
	s.Equal(userID, getRes.Data.Id)

	// 3. FindAll
	allRes, err := s.queryClient.FindAll(ctx, &pbusers.FindAllUserRequest{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.NotEmpty(allRes.Data)

	// 4. FindByActive
	activeRes, err := s.queryClient.FindByActive(ctx, &pbusers.FindAllUserRequest{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.NotEmpty(activeRes.Data)

	// 5. Update
	updateRes, err := s.client.Update(ctx, &pbusers.UpdateUserRequest{
		Id:              userID,
		Firstname:       "GapiUpdated",
		Lastname:        "UserUpdated",
		Email:           "gapi.updated@example.com",
		Password:        "password123",
		ConfirmPassword: "password123",
	})
	s.Require().NoError(err)
	s.Equal("GapiUpdated", updateRes.Data.Firstname)

	// 6. Trash
	_, err = s.client.TrashedUser(ctx, &pbusers.FindByIdUserRequest{Id: userID})
	s.Require().NoError(err)

	// 7. FindByTrashed
	trashedRes, err := s.queryClient.FindByTrashed(ctx, &pbusers.FindAllUserRequest{Page: 1, PageSize: 10})
	s.Require().NoError(err)
	s.NotEmpty(trashedRes.Data)

	// 8. Restore
	_, err = s.client.RestoreUser(ctx, &pbusers.FindByIdUserRequest{Id: userID})
	s.Require().NoError(err)

	// 9. DeletePermanent
	_, _ = s.client.TrashedUser(ctx, &pbusers.FindByIdUserRequest{Id: userID})
	_, err = s.client.DeleteUserPermanent(ctx, &pbusers.FindByIdUserRequest{Id: userID})
	s.Require().NoError(err)

	// 10. RestoreAll
	_, err = s.client.RestoreAllUser(ctx, &emptypb.Empty{})
	s.Require().NoError(err)

	// 11. DeleteAll
	_, err = s.client.DeleteAllUserPermanent(ctx, &emptypb.Empty{})
	s.Require().NoError(err)
}

func TestUserGapiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(UserGapiTestSuite))
}

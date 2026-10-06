package role_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	pbroles "github.com/MamangRust/monolith-point-of-sale-pb/roles"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	apierror "github.com/MamangRust/monolith-point-of-sale-apigateway/apierror"
	apicache "github.com/MamangRust/monolith-point-of-sale-apigateway/cache/gateway_cache"
	rolehandler "github.com/MamangRust/monolith-point-of-sale-apigateway/handler"
	httpx "github.com/MamangRust/monolith-point-of-sale-apigateway/httpx"
	pbcompat "github.com/MamangRust/monolith-point-of-sale-apigateway/pbcompat"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	role_cache "github.com/MamangRust/monolith-point-of-sale-role/cache"
	role_handler "github.com/MamangRust/monolith-point-of-sale-role/handler"
	"github.com/MamangRust/monolith-point-of-sale-role/repository"
	"github.com/MamangRust/monolith-point-of-sale-role/service"
	"github.com/MamangRust/monolith-point-of-sale-shared/cache"
	"github.com/MamangRust/monolith-point-of-sale-shared/domain/requests"
	response_api "github.com/MamangRust/monolith-point-of-sale-shared/mapper/response/api"
	"github.com/MamangRust/monolith-point-of-sale-shared/observability"
	tests "github.com/MamangRust/monolith-point-of-sale-test"

	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/gorm"
)

type RoleApiTestSuite struct {
	suite.Suite
	ts          *tests.TestSuite
	gormDB      *gorm.DB
	redisClient *redis.Client
	router      chi.Router
	grpcServer  *grpc.Server
	conn        *grpc.ClientConn
	roleID      int
}

func (s *RoleApiTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	gormDB, err := ts.GormDB()
	s.Require().NoError(err)
	s.gormDB = gormDB

	opts, err := redis.ParseURL(s.ts.RedisURL)
	s.Require().NoError(err)
	s.redisClient = redis.NewClient(opts)

	repos := repository.NewRepositories(gormDB)

	log, _ := logger.NewLogger("test", nil)
	obs, _ := observability.NewObservability("test", log)
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(s.redisClient, log, cacheMetrics)
	mencache := role_cache.NewMencache(cacheStore)

	roleService := service.NewService(&service.Deps{
		Repositories:  repos,
		Logger:        log,
		Mencache:      mencache,
		Observability: obs,
	})

	// Start internal gRPC Server for Role module
	roleHandlerGrpc := role_handler.NewHandler(&role_handler.Deps{
		Service: roleService,
		Logger:  log,
	})
	server := grpc.NewServer()
	pbroles.RegisterRoleQueryServiceServer(server, roleHandlerGrpc.Role)
	pbroles.RegisterRoleCommandServiceServer(server, roleHandlerGrpc.RoleCommand)
	s.grpcServer = server

	lis, err := net.Listen("tcp", "localhost:0")
	s.Require().NoError(err)

	go func() {
		_ = server.Serve(lis)
	}()

	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)
	s.conn = conn

	// Setup Echo and API Handler
	e := chi.NewRouter()
	s.router = e

	// Bypass auth middleware by setting user_id and seeding roles in Redis
	e.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r = r.WithContext(httpx.SetValue(r.Context(), "userID", "1"))
			next.ServeHTTP(w, r)
		})
	})

	roles := []string{"Admin_Role_10", "Admin_Admin_14"}
	cache.SetToCache(s.ts.Ctx, cacheStore, "user_roles:1", &roles, 5*time.Minute)

	apiErrorHandler := apierror.NewApiHandler(obs, log)
	roleMapper := response_api.NewRoleResponseMapper()
	gatewayCache := apicache.NewGatewayCache(cacheStore)
	client := pbcompat.NewRoleServiceClient(conn)
	rolehandler.NewHandlerRole(
		e,
		client,
		log,
		roleMapper,
		apiErrorHandler,
		gatewayCache,
	)
}

func (s *RoleApiTestSuite) TearDownSuite() {
	s.conn.Close()
	s.grpcServer.Stop()
	s.redisClient.Close()
	if sqlDB, err := s.gormDB.DB(); err == nil {
		sqlDB.Close()
	}
	s.ts.Teardown()
}

func (s *RoleApiTestSuite) TestRoleApiLifecycle() {
	// 1. Create
	reqBody := requests.CreateRoleRequest{
		Name: "API Role",
	}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/role", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	s.router.ServeHTTP(rec, req)
	s.Require().Equal(http.StatusOK, rec.Code, rec.Body.String())
	var res map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &res)
	data := res["data"].(map[string]interface{})
	s.Equal(reqBody.Name, data["name"])
	s.roleID = int(data["id"].(float64))

	// 2. FindAll
	req = httptest.NewRequest(http.MethodGet, "/api/role", nil)
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 3. FindById
	s.Require().NotZero(s.roleID)
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/role/%d", s.roleID), nil)
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)
	json.Unmarshal(rec.Body.Bytes(), &res)
	data = res["data"].(map[string]interface{})
	s.Equal(float64(s.roleID), data["id"])

	// 4. FindByActive
	req = httptest.NewRequest(http.MethodGet, "/api/role/active", nil)
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 5. FindByTrashed
	req = httptest.NewRequest(http.MethodGet, "/api/role/trashed", nil)
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 6. FindByUserId
	req = httptest.NewRequest(http.MethodGet, "/api/role/user/1", nil)
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 7. Update
	updateBody := requests.UpdateRoleRequest{
		Name: "Updated API Role",
	}
	body, _ = json.Marshal(updateBody)
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/role/update/%d", s.roleID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 8. Restore
	req = httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/role/restore/%d", s.roleID), nil)
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 9. DeletePermanent
	req = httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/role/permanent/%d", s.roleID), nil)
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 10. RestoreAll
	req = httptest.NewRequest(http.MethodPost, "/api/role/restore/all", nil)
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)

	// 11. DeleteAll
	req = httptest.NewRequest(http.MethodDelete, "/api/role/permanent-all", nil)
	rec = httptest.NewRecorder()
	s.router.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)
}

func TestRoleApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(RoleApiTestSuite))
}

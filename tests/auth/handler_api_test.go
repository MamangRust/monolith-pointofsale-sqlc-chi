package auth_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/MamangRust/monolith-point-of-sale-pb"
	pbroles "github.com/MamangRust/monolith-point-of-sale-pb/roles"
	pbuserrole "github.com/MamangRust/monolith-point-of-sale-pb/user_role"
	pbusers "github.com/MamangRust/monolith-point-of-sale-pb/users"
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	apierror "github.com/MamangRust/monolith-point-of-sale-apigateway/apierror"
	auth_cache_api "github.com/MamangRust/monolith-point-of-sale-apigateway/cache/auth"
	apigateway "github.com/MamangRust/monolith-point-of-sale-apigateway/handler"
	httpx "github.com/MamangRust/monolith-point-of-sale-apigateway/httpx"
	pbcompat "github.com/MamangRust/monolith-point-of-sale-apigateway/pbcompat"
	auth_cache "github.com/MamangRust/monolith-point-of-sale-auth/cache"
	"github.com/MamangRust/monolith-point-of-sale-auth/handler"
	"github.com/MamangRust/monolith-point-of-sale-auth/repository"
	"github.com/MamangRust/monolith-point-of-sale-auth/service"
	"github.com/MamangRust/monolith-point-of-sale-pkg/auth"
	"github.com/MamangRust/monolith-point-of-sale-pkg/hash"
	"github.com/MamangRust/monolith-point-of-sale-pkg/logger"
	role_cache "github.com/MamangRust/monolith-point-of-sale-role/cache"
	role_handler "github.com/MamangRust/monolith-point-of-sale-role/handler"
	role_repo "github.com/MamangRust/monolith-point-of-sale-role/repository"
	role_service "github.com/MamangRust/monolith-point-of-sale-role/service"
	"github.com/MamangRust/monolith-point-of-sale-shared/cache"
	response_api "github.com/MamangRust/monolith-point-of-sale-shared/mapper/response/api"
	"github.com/MamangRust/monolith-point-of-sale-shared/observability"
	tests "github.com/MamangRust/monolith-point-of-sale-test"
	user_cache "github.com/MamangRust/monolith-point-of-sale-user/cache"
	user_handler "github.com/MamangRust/monolith-point-of-sale-user/handler"
	user_repo "github.com/MamangRust/monolith-point-of-sale-user/repository"
	user_service "github.com/MamangRust/monolith-point-of-sale-user/service"

	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/gorm"
)

type AuthHandlerApiTestSuite struct {
	suite.Suite
	ts          *tests.TestSuite
	gormDB      *gorm.DB
	redisClient *redis.Client
	server      chi.Router
	email       string
	password    string
	accessToken string
	userID      int
}

func (s *AuthHandlerApiTestSuite) SetupSuite() {
	ts, err := tests.SetupTestSuite()
	s.Require().NoError(err)
	s.ts = ts

	gormDB, err := ts.GormDB()
	s.Require().NoError(err)
	s.gormDB = gormDB

	opts, err := redis.ParseURL(s.ts.RedisURL)
	s.Require().NoError(err)
	s.redisClient = redis.NewClient(opts)
	s.redisClient.FlushAll(context.Background())

	log, _ := logger.NewLogger("test", nil)
	hasher := hash.NewHashingPassword()
	cacheMetrics, _ := observability.NewCacheMetrics("test")
	cacheStore := cache.NewCacheStore(s.redisClient, log, cacheMetrics)
	obs, _ := observability.NewObservability("test", log)

	// 1. Setup Role Service & gRPC Server
	roleMencache := role_cache.NewMencache(cacheStore)
	roleRepos := role_repo.NewRepositories(gormDB)
	roleSvc := role_service.NewService(&role_service.Deps{
		Repositories:  roleRepos,
		Logger:        log,
		Mencache:      roleMencache,
		Observability: obs,
	})
	roleGapi := role_handler.NewHandler(&role_handler.Deps{
		Service: roleSvc,
		Logger:  log,
	})
	roleServer := grpc.NewServer()
	pbroles.RegisterRoleQueryServiceServer(roleServer, roleGapi.Role)
	pbroles.RegisterRoleCommandServiceServer(roleServer, roleGapi.RoleCommand)
	pbuserrole.RegisterUserRoleServiceServer(roleServer, roleGapi.UserRole)
	roleLis, _ := net.Listen("tcp", "localhost:0")
	go roleServer.Serve(roleLis)
	roleConn, _ := grpc.NewClient(roleLis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))

	// 2. Setup User Service & gRPC Server
	userMencache := user_cache.NewMencache(cacheStore)
	userRepos := user_repo.NewRepositories(&user_repo.Deps{
		Db:       gormDB,
		Role:     pbroles.NewRoleQueryServiceClient(roleConn),
		UserRole: pbuserrole.NewUserRoleServiceClient(roleConn),
	})
	userSvc := user_service.NewService(&user_service.Deps{
		Repositories:  userRepos,
		Logger:        log,
		Hash:          hasher,
		Mencache:      userMencache,
		Observability: obs,
	})
	userGapi := user_handler.NewHandler(&user_handler.Deps{
		Service: userSvc,
		Logger:  log,
	})
	userServer := grpc.NewServer()
	pbusers.RegisterUserQueryServiceServer(userServer, userGapi.User)
	pbusers.RegisterUserCommandServiceServer(userServer, userGapi.UserCommand)
	userLis, _ := net.Listen("tcp", "localhost:0")
	go userServer.Serve(userLis)
	userConn, _ := grpc.NewClient(userLis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))

	// 3. Setup Auth Service
	roleCommandClient := pbroles.NewRoleCommandServiceClient(roleConn)
	roleQueryClient := pbroles.NewRoleQueryServiceClient(roleConn)
	userQueryClient := pbusers.NewUserQueryServiceClient(userConn)
	userCommandClient := pbusers.NewUserCommandServiceClient(userConn)
	repos := repository.NewRepositories(gormDB, userQueryClient, userCommandClient, roleQueryClient, roleCommandClient, pbuserrole.NewUserRoleServiceClient(roleConn))

	tokenManager, _ := auth.NewManager("mysecret")
	mencache := auth_cache.NewMencache(cacheStore)
	apiAuthCache := auth_cache_api.NewMencache(cacheStore)

	svc := service.NewService(&service.Deps{
		Repositories:  repos,
		Logger:        log,
		Mencache:      mencache,
		Token:         tokenManager,
		Hash:          hasher,
		Kafka:         nil,
		Observability: obs,
	})

	h := handler.NewAuthHandleGrpc(svc, log)

	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, h)

	lis, err := net.Listen("tcp", "localhost:0")
	s.Require().NoError(err)

	go func() {
		_ = grpcServer.Serve(lis)
	}()

	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)

	s.server = chi.NewRouter()
	apiHandler := apierror.NewApiHandler(obs, log)
	mapper := response_api.NewAuthResponseMapper()

	// Auth bypass middleware for /api/auth/me
	s.server.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if s.userID != 0 {
				r = r.WithContext(httpx.SetValue(r.Context(), "userID", strconv.Itoa(s.userID)))
			}
			next.ServeHTTP(w, r)
		})
	})

	fmt.Println("Calling RegisterAuthHandler...")
	apigateway.NewHandlerAuth(s.server, pbcompat.NewAuthServiceClient(conn), log, mapper, apiHandler, apiAuthCache)
	fmt.Println("RegisterAuthHandler called.")

	s.server.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("pong"))
	})

	s.email = "auth.handler.api.test@example.com"
	s.password = "password123"

	// ROLE_ADMIN sudah di-seed oleh tests.SetupTestSuite (SeedMinimalRoles).
	if _, err := roleCommandClient.CreateRole(context.Background(), &pbroles.CreateRoleRequest{
		Name: "ROLE_ADMIN",
	}); err != nil {
		fmt.Printf("DEBUG: Seed ROLE_ADMIN gRPC (expected conflict): %v\n", err)
	}
}

func (s *AuthHandlerApiTestSuite) TearDownSuite() {
	if s.redisClient != nil {
		s.redisClient.Close()
	}
	s.ts.Teardown()
}

func (s *AuthHandlerApiTestSuite) Test0_Ping() {
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	rec := httptest.NewRecorder()
	s.server.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)
	s.Equal("pong", rec.Body.String())
}

func (s *AuthHandlerApiTestSuite) Test0_Hello() {
	req := httptest.NewRequest(http.MethodGet, "/api/auth/hello", nil)
	rec := httptest.NewRecorder()
	s.server.ServeHTTP(rec, req)
	s.Equal(http.StatusOK, rec.Code)
	s.Equal("Hello", rec.Body.String())
}

// verifyUser simulates the email-verification step (login only accepts
// is_verified = true users).
func (s *AuthHandlerApiTestSuite) verifyUser(email string) {
	err := s.gormDB.Exec("UPDATE users SET is_verified = true WHERE email = ?", email).Error
	s.Require().NoError(err)
}

func (s *AuthHandlerApiTestSuite) Test1_Register() {
	body := map[string]string{
		"firstname":        "Auth",
		"lastname":         "API",
		"email":            s.email,
		"password":         s.password,
		"confirm_password": s.password,
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	s.server.ServeHTTP(rec, req)

	s.Equal(http.StatusOK, rec.Code, "Expected StatusOK, got %d. Body: %s", rec.Code, rec.Body.String())

	var res map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &res)
	s.NoError(err)

	data, ok := res["data"].(map[string]interface{})
	s.True(ok, "Expected 'data' to be a map, got %T. Body: %s", res["data"], rec.Body.String())
	s.userID = int(data["id"].(float64))
}

func (s *AuthHandlerApiTestSuite) Test2_Login() {
	s.verifyUser(s.email)

	body := map[string]string{
		"email":    s.email,
		"password": s.password,
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	s.server.ServeHTTP(rec, req)

	s.Equal(http.StatusOK, rec.Code, "Expected StatusOK, got %d. Body: %s", rec.Code, rec.Body.String())

	var res map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &res)
	s.NoError(err)

	data, ok := res["data"].(map[string]interface{})
	s.True(ok, "Expected 'data' to be a map, got %T. Body: %s", res["data"], rec.Body.String())
	s.accessToken = data["access_token"].(string)
}

func (s *AuthHandlerApiTestSuite) Test4_LoginLockout() {
	email := "locked.api@example.com"
	password := "wrongpassword"

	// Register user first
	regBody := map[string]string{
		"firstname":        "Locked",
		"lastname":         "API",
		"email":            email,
		"password":         "correctpassword",
		"confirm_password": "correctpassword",
	}
	jsonRegBody, _ := json.Marshal(regBody)
	regReq := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBuffer(jsonRegBody))
	regReq.Header.Set("Content-Type", "application/json")
	regRec := httptest.NewRecorder()
	s.server.ServeHTTP(regRec, regReq)
	s.Equal(http.StatusOK, regRec.Code)
	s.verifyUser(email)

	loginBody := map[string]string{
		"email":    email,
		"password": password,
	}
	jsonLoginBody, _ := json.Marshal(loginBody)

	// Fail login 5 times (total 5) — password mismatch is 400 Bad Request
	for i := 0; i < 5; i++ {
		loginReq := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(jsonLoginBody))
		loginReq.Header.Set("Content-Type", "application/json")
		loginRec := httptest.NewRecorder()
		s.server.ServeHTTP(loginRec, loginReq)
		s.Equal(http.StatusBadRequest, loginRec.Code)
	}

	// 6th attempt should return 400 Bad Request (ErrAccountLocked)
	lockedReq := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBuffer(jsonLoginBody))
	lockedReq.Header.Set("Content-Type", "application/json")
	lockedRec := httptest.NewRecorder()
	s.server.ServeHTTP(lockedRec, lockedReq)
	s.Equal(http.StatusBadRequest, lockedRec.Code)
}

func (s *AuthHandlerApiTestSuite) Test3_GetMe() {
	s.Require().NotZero(s.userID)
	s.Require().NotEmpty(s.accessToken)

	req := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	req.Header.Set("Authorization", "Bearer "+s.accessToken)
	rec := httptest.NewRecorder()

	s.server.ServeHTTP(rec, req)

	s.Equal(http.StatusOK, rec.Code, "Expected StatusOK, got %d. Body: %s", rec.Code, rec.Body.String())

	var res map[string]interface{}
	err := json.Unmarshal(rec.Body.Bytes(), &res)
	s.NoError(err)

	data, ok := res["data"].(map[string]interface{})
	s.True(ok, "Expected 'data' to be a map, got %T. Body: %s", res["data"], rec.Body.String())
	s.Equal(s.email, data["email"])
}

func TestAuthHandlerApiSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(AuthHandlerApiTestSuite))
}

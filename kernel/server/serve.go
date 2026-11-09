package server

import (
	"labflow/auth"
	"labflow/bootstrap"
	"labflow/handler"
	"labflow/httpx"
	"labflow/mail"
	"labflow/middleware"
	"labflow/queue"
	"labflow/repo"
	"labflow/service"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
)

func Serve(cookieKey string, db *pgxpool.Pool, redisClient *redis.Client, mailer *mail.Sender, enqueuer *queue.Enqueuer, s3client *s3.Client) {
	addr := ":8080"
	if port := os.Getenv("PORT"); port != "" {
		addr = ":" + port
	}
	srv := &http.Server{
		Addr:              addr,
		Handler:           NewHandler(cookieKey, db, redisClient, enqueuer, s3client),
		ReadHeaderTimeout: 10 * time.Second,
	}
	log.Printf("listening on %s", addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Printf("server stopped: %v", err)
	}
}

// NewHandler wires repositories, services and handlers into the API's routes
// and global middleware.
func NewHandler(cookieKey string, db *pgxpool.Pool, redisClient *redis.Client, enqueuer *queue.Enqueuer, s3client *s3.Client) http.Handler {
	authRepo := repo.NewAuthRepo(db)
	reservationRepo := repo.NewReservationRepo(db)
	testListRepo := repo.NewTestListRepo(db)
	userRepo := repo.NewUserRepo(db)
	testRepo := repo.NewTestRepo(db)
	quotationRepo := repo.NewQuotationRepo(db)
	tokenUtil := auth.NewJwtUtils(cookieKey, redisClient)
	redisStore := bootstrap.NewRedisStore(redisClient)

	fileObjectService := service.NewFileObjectService(s3client)
	authService := service.NewService(authRepo, redisStore, enqueuer, tokenUtil, repo.RoleCustomer)
	reservationService := service.NewReservationService(reservationRepo, redisStore)
	testListService := service.NewTestListService(testListRepo)
	uploadService := service.NewUploadService(s3client, bootstrap.BucketReservations, bootstrap.BucketCertificates, bootstrap.BucketUsers)
	accessService := service.NewAccessService(authRepo)
	userService := service.NewUserService(userRepo)
	testService := service.NewTestService(testRepo, authRepo)
	quotationService := service.NewQuotationService(quotationRepo, enqueuer)

	authHandler := handler.NewAuthHandler(authService)
	reservationHandler := handler.NewReservationHandler(reservationService)
	testlisthandler := handler.NewTestListHandler(testListService)
	uploadHandler := handler.NewUploadHandler(uploadService)
	userHandler := handler.NewUserHandler(userService)
	fileHandler := handler.NewFileHandler(fileObjectService, bootstrap.BucketReservations, bootstrap.BucketUsers, bootstrap.BucketCertificates)
	testHandler := handler.NewTestHandler(testService)
	quotationHandler := handler.NewQuotationHandler(quotationService)
	lifecycleHandler := handler.NewLifecycleHandler()

	// Strict rate limiter for auth (IP-based, 10 requests per minute)
	authLimiter := middleware.RateLimiter(redisClient, middleware.RateLimitConfig{
		Namespace: "auth",
		Limit:     10,
		Window:    time.Minute,
	})
	// User/IP-aware API rate limiter (60 requests per minute)
	apiLimiter := middleware.RateLimiter(redisClient, middleware.RateLimitConfig{
		Namespace: "api",
		Limit:     60,
		Window:    time.Minute,
	})
	jwt := middleware.JWTCookieMiddleware(tokenUtil)

	mux := http.NewServeMux()
	// authRoute: rate limited by IP, no session required.
	authRoute := func(pattern string, h http.HandlerFunc) {
		mux.Handle(pattern, httpx.Chain(h, authLimiter))
	}
	// protected: requires a valid session cookie.
	protected := func(pattern string, h http.HandlerFunc) {
		mux.Handle(pattern, httpx.Chain(h, jwt, apiLimiter))
	}
	// requires: a valid session cookie and the named permission.
	requires := func(permission, pattern string, h http.HandlerFunc) {
		mux.Handle(pattern, httpx.Chain(h, jwt, apiLimiter, middleware.RBACMiddleware(accessService, permission)))
	}

	authRoute("POST /auth/register", authHandler.Register)
	authRoute("POST /auth/login", authHandler.LoginWithCookie)
	authRoute("POST /auth/logout", authHandler.Logout)
	authRoute("POST /auth/verify-otp", authHandler.VerifyOtp)
	authRoute("POST /auth/resend-otp", authHandler.ResendOtp)
	authRoute("POST /auth/forgot-password", authHandler.ForgotPassword)
	authRoute("POST /auth/verify-reset-otp", authHandler.VerifyResetOtp)
	authRoute("POST /auth/reset-password", authHandler.ResetPassword)

	mux.HandleFunc("POST /api/quotations", quotationHandler.CreateQuotation)
	requires("read_quotation", "GET /api/quotations", quotationHandler.ListQuotations)
	requires("update_quotation", "PATCH /api/quotations/{xid}", quotationHandler.UpdateQuotationStatus)

	requires("create_reservation", "POST /api/reservations", reservationHandler.CreateReservation)
	requires("read_reservation", "GET /api/reservations", reservationHandler.ListReservations)
	requires("read_reservation", "GET /api/reservations/stats", reservationHandler.GetReservationStats)
	requires("update_reservation", "PATCH /api/reservations/{xid}", reservationHandler.UpdateReservationsByXid)
	requires("read_reservation", "GET /api/reservations/{xid}", reservationHandler.GetReservationByXid)
	requires("update_reservation_status", "PUT /api/reservations/{xid}/status", reservationHandler.UpdateReservationStatus)

	requires("read_test", "GET /api/tests/catalog/{categoryid}", testlisthandler.GetTestListByCategory)
	requires("update_test", "PATCH /api/tests/{xid}", testHandler.UpdateTestByXid)

	requires("read_test", "GET /api/test-lists", testlisthandler.ListTestLists)
	requires("create_test", "POST /api/test-lists", testlisthandler.CreateTestList)
	requires("update_test", "PUT /api/test-lists/{xid}", testlisthandler.UpdateTestList)
	requires("delete_test", "DELETE /api/test-lists/{xid}", testlisthandler.DeactivateTestList)

	protected("GET /api/auth/me", authHandler.GetMe)
	requires("read_lifecycle", "GET /api/lifecycle/transitions", lifecycleHandler.GetTransitions)

	requires("read_reservation", "GET /api/files/reservations/{object}", fileHandler.ServeReservationFile)
	protected("GET /api/files/users/{object}", fileHandler.ServeUserFile)
	requires("read_certificate", "GET /api/files/certificates/{object}", fileHandler.ServeCertificateFile)

	requires("update_reservation_status", "POST /api/uploads/reservation", uploadHandler.UploadSamplePhoto)
	requires("create_certificate", "POST /api/uploads/certificate", uploadHandler.UploadCertificate)
	requires("update_user_avatar", "POST /api/uploads/user", uploadHandler.UploadUserAvatar)

	requires("update_user", "PATCH /api/users/{xid}", userHandler.UpdateUser)

	mux.Handle("GET /metrics", promhttp.Handler())
	// GET patterns also match HEAD requests.
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		httpx.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	return httpx.Chain(mux,
		middleware.AccessLog(os.Stdout),
		middleware.Recover,
		middleware.Metrics,
		middleware.CORS(middleware.CORSConfig{
			AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000"},
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
		}),
		// Global safety limiter (IP-based, 100 requests per minute)
		middleware.RateLimiter(redisClient, middleware.RateLimitConfig{
			Namespace: "global",
			Limit:     100,
			Window:    time.Minute,
		}),
	)
}

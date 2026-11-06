package server

import (
	"labflow/auth"
	"labflow/bootstrap"
	"labflow/handler"
	"labflow/mail"
	"labflow/middleware"
	"labflow/queue"
	"labflow/repo"
	"labflow/service"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func Serve(cookieKey string, db *gorm.DB, redisClient *redis.Client, mailer *mail.Sender, enqueuer *queue.Enqueuer, s3client *s3.Client) {
	gin.SetMode(gin.ReleaseMode)
	ginServer := gin.Default()
	ginServer.Use(middleware.Metrics())

	sessionStore := cookie.NewStore([]byte(cookieKey))
	sessionStore.Options(sessions.Options{
		Path:     "/",
		Secure:   false,
		HttpOnly: true,
	})
	ginServer.Use(sessions.Sessions("labflow-session", sessionStore))
	ginServer.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	// Global safety limiter (IP-based, 100 requests per minute)
	ginServer.Use(middleware.RateLimiter(redisClient, middleware.RateLimitConfig{
		Namespace: "global",
		Limit:     100,
		Window:    time.Minute,
	}))

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

	authentication := ginServer.Group("/auth")
	// Strict rate limiter for auth (IP-based, 10 requests per minute)
	authentication.Use(middleware.RateLimiter(redisClient, middleware.RateLimitConfig{
		Namespace: "auth",
		Limit:     10,
		Window:    time.Minute,
	}))
	{
		authentication.POST("/register", authHandler.Register)
		authentication.POST("/login", authHandler.LoginWithCookie)
		authentication.POST("/logout", authHandler.Logout)
		authentication.POST("/verify-otp", authHandler.VerifyOtp)
		authentication.POST("/resend-otp", authHandler.ResendOtp)
		authentication.POST("/forgot-password", authHandler.ForgotPassword)
		authentication.POST("/verify-reset-otp", authHandler.VerifyResetOtp)
		authentication.POST("/reset-password", authHandler.ResetPassword)
	}
	ginServer.POST("/api/quotations", quotationHandler.CreateQuotation)
	// Protected routes
	apiGroup := ginServer.Group("/api")
	apiGroup.Use(middleware.JWTCookieMiddleware(tokenUtil))
	// User/IP-aware API rate limiter (60 requests per minute)
	apiGroup.Use(middleware.RateLimiter(redisClient, middleware.RateLimitConfig{
		Namespace: "api",
		Limit:     60,
		Window:    time.Minute,
	}))
	{
		quotations := apiGroup.Group("/quotations")
		{
			quotations.GET("", middleware.RBACMiddleware(accessService, "read_quotation"), quotationHandler.ListQuotations)
			quotations.PATCH("/:xid", middleware.RBACMiddleware(accessService, "update_quotation"), quotationHandler.UpdateQuotationStatus)
		}

		reservations := apiGroup.Group("/reservations")
		{
			reservations.POST("", middleware.RBACMiddleware(accessService, "create_reservation"), reservationHandler.CreateReservation)
			reservations.GET("", middleware.RBACMiddleware(accessService, "read_reservation"), reservationHandler.ListReservations)
			reservations.GET("/stats", middleware.RBACMiddleware(accessService, "read_reservation"), reservationHandler.GetReservationStats)
			reservations.PATCH("/:xid", middleware.RBACMiddleware(accessService, "update_reservation"), reservationHandler.UpdateReservationsByXid)
			reservations.GET("/:xid", middleware.RBACMiddleware(accessService, "read_reservation"), reservationHandler.GetReservationByXid)
			reservations.PUT("/:xid/status", middleware.RBACMiddleware(accessService, "update_reservation_status"), reservationHandler.UpdateReservationStatus)
		}

		tests := apiGroup.Group("/tests")
		{
			tests.GET("/catalog/:categoryid", middleware.RBACMiddleware(accessService, "read_test"), testlisthandler.GetTestListByCategory)
			tests.PATCH("/:xid", middleware.RBACMiddleware(accessService, "update_test"), testHandler.UpdateTestByXid)

		}
		testLists := apiGroup.Group("/test-lists")
		{
			testLists.GET("", middleware.RBACMiddleware(accessService, "read_test"), testlisthandler.ListTestLists)
			testLists.POST("", middleware.RBACMiddleware(accessService, "create_test"), testlisthandler.CreateTestList)
			testLists.PUT("/:xid", middleware.RBACMiddleware(accessService, "update_test"), testlisthandler.UpdateTestList)
			testLists.DELETE("/:xid", middleware.RBACMiddleware(accessService, "delete_test"), testlisthandler.DeactivateTestList)
		}
		apiGroup.GET("/auth/me", authHandler.GetMe)
		apiGroup.GET("/lifecycle/transitions", middleware.RBACMiddleware(accessService, "read_lifecycle"), lifecycleHandler.GetTransitions)
		files := apiGroup.Group("/files")
		{
			files.GET("/reservations/:object", middleware.RBACMiddleware(accessService, "read_reservation"), fileHandler.ServeReservationFile)
			files.GET("/users/:object", fileHandler.ServeUserFile)
			files.GET("/certificates/:object", middleware.RBACMiddleware(accessService, "read_certificate"), fileHandler.ServeCertificateFile)
		}

		apiGroup.POST("/uploads/reservation", middleware.RBACMiddleware(accessService, "update_reservation_status"), uploadHandler.UploadSamplePhoto)
		apiGroup.POST("/uploads/certificate", middleware.RBACMiddleware(accessService, "create_certificate"), uploadHandler.UploadCertificate)
		apiGroup.POST("/uploads/user", middleware.RBACMiddleware(accessService, "update_user_avatar"), uploadHandler.UploadUserAvatar)

		apiGroup.PATCH("/users/:xid", middleware.RBACMiddleware(accessService, "update_user"), userHandler.UpdateUser)
	}

	healthHandler := func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	}
	ginServer.GET("/metrics", gin.WrapH(promhttp.Handler()))
	ginServer.GET("/health", healthHandler)
	ginServer.HEAD("/health", healthHandler)

	err := ginServer.Run()
	if err != nil {
		return
	}
}

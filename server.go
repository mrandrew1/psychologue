package main

import (
	"github.com/mrandrew1/psychologue/cache"
	"github.com/mrandrew1/psychologue/config"
	"github.com/mrandrew1/psychologue/controller"
	"github.com/mrandrew1/psychologue/middleware"
	"github.com/mrandrew1/psychologue/repository"
	"github.com/mrandrew1/psychologue/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

var (
	db                   *gorm.DB                        = config.SetupDatabaseConnection()
	userRepository       repository.UserRepository       = repository.NewUserRepository(db)
	bookRepository       repository.BookRepository       = repository.NewBookRepository(db)
	psySessionRepository repository.PsySessionRepository = repository.NewPsySessionRepository(db)
	receiptRepository    repository.ReceiptRepository    = repository.NewReceiptRepository(db)

	rd           *redis.Client      = config.SetupRedisConnection()
	receiptCache cache.ReceiptCache = cache.NewReceiptCache(rd, 1000) // random expiration lmao

	jwtService        service.JWTService        = service.NewJWTService()
	authService       service.AuthService       = service.NewAuthService(userRepository)
	userService       service.UserService       = service.NewUserService(userRepository)
	psySessionService service.PsySessionService = service.NewPsySessionService(psySessionRepository)
	bookService       service.BookService       = service.NewBookService(bookRepository)
	receiptService    service.ReceiptService    = service.NewReceiptService(receiptRepository)

	authController       controller.AuthController       = controller.NewAuthController(authService, jwtService)
	userController       controller.UserController       = controller.NewUserController(userService, jwtService)
	bookController       controller.BookController       = controller.NewBookController(bookService, jwtService)
	psySessionController controller.PsySessionController = controller.NewPsySessionController(psySessionService, jwtService)
	receiptController    controller.ReceiptController    = controller.NewReceiptController(receiptService, receiptCache)
)

func main() {
	defer config.CloseDatabaseConnection(db)
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowMethods:     []string{"POST", "PUT", "PATCH", "GET", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowAllOrigins:  true,
		AllowCredentials: true,
	}))

	authRoutes := r.Group("/api/auth")
	{
		authRoutes.POST("login", authController.Login)
		authRoutes.POST("register", authController.Register)
	}

	userRoutes := r.Group("/api/users", middleware.AuthorizeJWT(jwtService))
	{
		userRoutes.GET("profile", userController.Get)
		userRoutes.PUT("", userController.Update)
	}

	bookRoutes := r.Group("/api/books", middleware.AuthorizeJWT(jwtService))
	{
		bookRoutes.GET("", bookController.All)
		bookRoutes.POST("", bookController.Insert)
		bookRoutes.GET("/:id", bookController.Get)
		bookRoutes.PUT("", bookController.Update)
		bookRoutes.DELETE("/:id", bookController.Delete)
	}

	psySessionRoutes := r.Group("/api/sessions", middleware.AuthorizeJWT(jwtService))
	{
		psySessionRoutes.GET("", psySessionController.All)
		psySessionRoutes.POST("", psySessionController.Insert)
		psySessionRoutes.GET("/:id", psySessionController.Get)
		psySessionRoutes.PUT("", psySessionController.Update)
		psySessionRoutes.DELETE("/:id", psySessionController.Delete)
	}

	receiptRoutes := r.Group("/api/receipts")
	{
		receiptRoutes.GET("/all", receiptController.All)
		receiptRoutes.POST("", receiptController.Insert)
		receiptRoutes.GET("/:id", receiptController.Show)
		receiptRoutes.PUT("", receiptController.Update)
		receiptRoutes.DELETE("/:id", receiptController.Delete)
	}

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World",
		})
	})

	r.Run()
}

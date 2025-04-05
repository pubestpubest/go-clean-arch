package main

import (
	"context"
	"fmt"
	"net/http"

	"order-management/entity"
	shopDelivery "order-management/features/shop/delivery"
	shopRepository "order-management/features/shop/repository"
	shopUsecase "order-management/features/shop/usecase"
	userDelivery "order-management/features/user/delivery"
	userRepository "order-management/features/user/repository"
	userUsecase "order-management/features/user/usecase"

	"order-management/utils"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	// joonix "github.com/joonix/log"
	log "github.com/sirupsen/logrus"
)

var runEnv string
var DB *gorm.DB

func serveGracefulShutdown(e *echo.Echo) {
	go func() {
		var port string
		port = os.Getenv("HTTP_PORT")
		if port == "" {
			port = utils.ViperGetString("http.port")
		}

		if err := e.Start(port); err != nil {
			log.Println("shutting down the server", err)

		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	gracefulShutdownTimeout := 30 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), gracefulShutdownTimeout)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		log.Fatal(err.Error())
	}
}

func migrateDB() {
	DB.AutoMigrate(
		&entity.User{},
		&entity.Order{},
		&entity.Product{},
		&entity.Shop{},
	)
}

func init() {
	var err error
	runEnv = os.Getenv("RUN_ENV")
	if runEnv == "" {
		runEnv = "local"
	}

	utils.InitViper(runEnv)

	// secret, err := utils.GetSecret(os.Getenv("PROJECT_ID"), os.Getenv("SECRET_ID"), os.Getenv("SECRET_VERSION"))
	// if err != nil {
	// 	log.Fatalf("GetSecret error: cannot get secret:%s", err)
	// }

	// fmt.Println("secret: ", secret)

	// os.WriteFile("configs/.env", []byte(secret), 0644)
	// if err := godotenv.Load("configs/.env"); err != nil {
	// 	fmt.Println("Error loading .env file")
	// }

	if err = connectDB(); err != nil {
		log.Fatal(err)
	}

}

func main() {
	e := echo.New()

	// e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())

	// log.SetFormatter(joonix.NewFormatter())
	log.SetLevel(log.TraceLevel)
	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: false,
	})

	log.Info("Starting server")

	// Unauthenticated route
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{"success": true})
	})

	// Restricted group
	// v1 := e.Group("/v1")

	// adminGroup := v1.Group("")
	// adminGroup.Use(middleware.AdminAuth())

	// customerGroup := v1.Group("")
	// customerGroup.Use(middleware.CustomerAuth())

	shopGroup := e.Group("/shops")

	userGroup := e.Group("/users")

	userDelivery.NewHandler(userGroup,
		userUsecase.NewUserUsecase(
			userRepository.NewUserRepository(DB),
		),
	)

	shopDelivery.NewHandler(shopGroup,
		shopUsecase.NewShopUsecase(
			shopRepository.NewShopRepository(DB),
		),
	)

	serveGracefulShutdown(e)
}

func connectDB() error {
	connectionString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
		utils.ViperGetString("postgres.host"),
		utils.ViperGetString("postgres.user"),
		utils.ViperGetString("postgres.password"),
		utils.ViperGetString("postgres.dbname"),
		utils.ViperGetString("postgres.port"))

	var err error
	DB, err = gorm.Open(postgres.Open(connectionString), &gorm.Config{})

	migrateDB()

	return err
}

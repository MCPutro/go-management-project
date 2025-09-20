package main

import (
	"log"

	"github.com/MCPutro/go-management-project/internal/config"
	"github.com/MCPutro/go-management-project/internal/config/database"
	"github.com/MCPutro/go-management-project/internal/delivery/handler"
	"github.com/MCPutro/go-management-project/internal/delivery/router"
	"github.com/MCPutro/go-management-project/internal/repository"
	"github.com/MCPutro/go-management-project/internal/usecase"
	"github.com/gofiber/fiber/v2"
)

func main1() {

	//loadConfig := config.LoadConfig()
	//
	//postgresDB, err := database.NewPostgresDB(loadConfig)
	//if err != nil {
	//	fmt.Println("Error connecting to the database:", err)
	//	return
	//}
	//defer postgresDB.Close()
	//
	//userRepository := repository.NewUserRepository()
	//userUsecase := usecase.NewUserUsecase(postgresDB, userRepository)
	//userHandler := handler.NewUserHandler(userUsecase)
	//
	//app := fiber.New()
	//
	//router.RegisterUserRoutes(app, userHandler)
	//
	//err = app.Listen(":" + loadConfig.ServerPort)
	//if err != nil {
	//	fmt.Println("Error starting the server:", err)
	//	return
	//}
}

func main() {
	loadConfig, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("failed to load config:", err)
	}

	postgresDb, err := database.NewPostgresDB(loadConfig.GetDatabaseConfig())
	if err != nil {
		log.Fatalln("failed to connect to database:", err)
	}
	defer postgresDb.Close()

	userRepository := repository.NewUserRepository()
	userUsecase := usecase.NewUserUsecase(postgresDb, userRepository)
	userHandler := handler.NewUserHandler(userUsecase)

	app := fiber.New()

	router.RegisterUserRoutes(app, userHandler)

	err = app.Listen(":" + loadConfig.GetApplicationConfig().Port)
	if err != nil {
		log.Fatalln("failed to listen:", err)
	}
}

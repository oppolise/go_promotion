package main

import (
	"goPromotion/cmd/database"
	serverconfig "goPromotion/config/server_config"
	_ "goPromotion/docs"
	"goPromotion/handler"
	"goPromotion/pkg/repository"
	"goPromotion/pkg/service"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"
)

// @title Example API
// @version 1.0
// @description test swag
// @BasePath /

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := database.InitDatabase()

	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}

	if err := database.Migration(db); err != nil {
		log.Fatalf("Error migrating database: %v", err)
	}

	app := fiber.New()
	app.Get("/swagger/*", swagger.HandlerDefault)

	orderPepo := repository.NewOrderRepository(db)
	orderService := service.NewOrderImpService(orderPepo)
	orderHandler := handler.NewHttpOrderHandler(orderService)

	app.Get("/order/:id", orderHandler.GetOrder)

	app.Listen(":" + serverconfig.ServerConfig().PORT)
}

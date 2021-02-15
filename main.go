package main

import (
	"MartellX/avito-tech-task/controllers"
	"MartellX/avito-tech-task/models"
	"MartellX/avito-tech-task/services"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"os"
)

func main() {

	r, err := models.NewRepositoryFromEnvironments()
	if err != nil {
		panic("Failed to connect to database")
	}
	if r == nil {
		panic("one of env variables not set")
	}
	s := services.NewService(r)
	handler := controllers.NewHandler(s, r)
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/task", handler.NewTask)
	e.GET("/task", handler.GetTask)
	e.GET("/offers", handler.GetOffers)
	port, ok := os.LookupEnv("port")
	if !ok {
		port = "1323"
	}
	e.Logger.Fatal(e.Start(":" + port))
}

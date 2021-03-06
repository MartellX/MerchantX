package main

import (
	"MartellX/avito-tech-task/controllers"
	"MartellX/avito-tech-task/repositories"
	"MartellX/avito-tech-task/services"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"os"
)

func main() {

	r, err := repositories.NewRepositoryFromEnvironments()
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

	e.POST("/tasks", handler.NewTask)
	e.GET("/tasks", handler.GetTask)
	e.GET("/offers", handler.GetOffers)
	port, ok := os.LookupEnv("port")
	if !ok {
		port = "1323"
	}
	e.Logger.Fatal(e.Start(":" + port))
}

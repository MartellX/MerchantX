package main

import (
	"MartellX/avito-tech-task/controllers"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"os"
)

func main() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/task", controllers.NewTask)
	e.GET("/task", controllers.GetTask)
	e.GET("/offers", controllers.GetOffers)
	port, ok := os.LookupEnv("port")
	if !ok {
		port = "1323"
	}
	e.Logger.Fatal(e.Start(":" + port))
}

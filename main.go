package main

import (
	"strconv"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	"github.com/ravenhurst/golang-playground/config"
	"github.com/ravenhurst/golang-playground/handlers"
	"github.com/ravenhurst/golang-playground/middlewares"
)

var appConfig config.Config

func init() {
	appConfig = config.GetConfig()
}


func main() {
	e := echo.New()
	e.Use(middleware.BodyDump(middlewares.RequestResponseLogger))
	e.POST("/auth", handlers.AuthHandler)
	e.GET("/protected-resource", handlers.ProtectedResourceHandler, middlewares.TokenAuthMiddleware)
	e.GET("/logout", handlers.LogoutHandler)

	e.Logger.Fatal(e.Start(":" + strconv.Itoa(appConfig.Port)))
}

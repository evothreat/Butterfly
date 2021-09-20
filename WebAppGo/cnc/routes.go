package cnc

import "github.com/labstack/echo/v4"

func SetupRoutes(g *echo.Group) {
	g.GET("/login", Login)
	g.POST("/logout", Logout)

	g.GET("/workers", GetAllWorkers)
}

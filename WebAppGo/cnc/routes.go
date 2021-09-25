package cnc

import "github.com/labstack/echo/v4"

func SetupRoutes(g *echo.Group) {
	g.GET("/login", Login)
	g.POST("/login", Login)

	g.POST("/logout", Logout)

	g.Use(AuthCheck)

	g.GET("/workers", GetAllWorkers)
	g.GET("/workers/:wid", Interact)
}

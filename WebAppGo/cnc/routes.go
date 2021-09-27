package cnc

import "github.com/labstack/echo/v4"

func SetupRoutes(e *echo.Echo) {
	g := e.Group("/cnc")

	g.GET("/login", LoginPage)
	g.POST("/login", Login)

	g.POST("/logout", Logout)

	g.Use(AuthCheck)

	g.GET("/workers", GetAllWorkers)
	g.GET("/workers/:wid", Interact)
}

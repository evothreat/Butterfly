package api

import "github.com/labstack/echo/v4"

func Setup(e *echo.Echo) {
	setupDatabase()
	//addTestData()
	setupRoutes(e)
}

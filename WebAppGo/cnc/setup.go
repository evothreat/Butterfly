package cnc

import (
	"github.com/labstack/echo/v4"
	"path"
)

func Setup(e *echo.Echo) {
	e.Static("/static", path.Join(RESOURCE_DIR, "/static"))
	e.Renderer = parseTemplates()
	setupRoutes(e)
}

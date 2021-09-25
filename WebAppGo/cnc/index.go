package cnc

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetAllWorkers(c echo.Context) error {
	return c.Render(http.StatusOK, "workers", nil)
}

func Interact(c echo.Context) error {
	return c.Render(http.StatusOK, "interact", nil)
}

package cnc

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetAllWorkers(c echo.Context) error {
	fmt.Println("doing")
	return c.Render(http.StatusOK, "workers.html", "Workers")
}

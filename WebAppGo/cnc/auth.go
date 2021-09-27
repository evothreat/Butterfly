package cnc

import (
	"WebAppGo/cnc/types"
	"WebAppGo/utils"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

const (
	cookieName   = "CNCSESSID"
	cookiePath   = "/cnc"
	cookieMaxAge = 3600 // 1 hour
)

var cookieDataStore = types.NewCookieDataStore()

func hasValidCookie(c echo.Context) bool {
	cookie, err := c.Cookie(cookieName)
	if err == http.ErrNoCookie {
		return false
	}
	data := cookieDataStore.Get(cookie.Value)
	return data != nil && time.Now().Before(data.Expires)
}

func Login(c echo.Context) error {
	if hasValidCookie(c) {
		return c.Redirect(http.StatusSeeOther, "/cnc/workers")
	}
	if c.Request().Method == "GET" {
		return c.Render(http.StatusOK, "login", nil)
	}
	// else if c.Request().Method == "POST"
	login := c.FormValue("username")
	passwd := c.FormValue("password")

	if login == ADMIN_LOGIN && bcrypt.CompareHashAndPassword([]byte(ADMIN_PASSWD), []byte(passwd)) == nil {
		data := types.NewCookieData()
		data.Expires = time.Now().Add(time.Second * cookieMaxAge)
		key := utils.RandomHexString(16)
		cookieDataStore.Put(key, data)

		cookie := http.Cookie{
			Name:   cookieName,
			Value:  key,
			Path:   cookiePath,
			MaxAge: cookieMaxAge,
		}
		c.SetCookie(&cookie)
		return c.Redirect(http.StatusSeeOther, "/cnc/workers")
	}
	return c.Render(http.StatusUnauthorized, "login", nil)
}

func Logout(c echo.Context) error {
	cookie, err := c.Cookie(cookieName)
	if err == http.ErrNoCookie {
		return c.Redirect(http.StatusSeeOther, "/cnc/login")
	}
	cookieDataStore.Delete(cookie.Value)
	newCookie := http.Cookie{
		Name:  cookieName,
		Value: "",
		Path:  cookiePath,
		//Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
	}
	c.SetCookie(&newCookie)
	return c.Redirect(http.StatusSeeOther, "/cnc/login")
}

func AuthCheck(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if hasValidCookie(c) {
			return next(c)
		}
		return Logout(c)
	}
}

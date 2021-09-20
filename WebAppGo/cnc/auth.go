package cnc

import (
	"WebAppGo/utils"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

const cookieName = "CNCSESSID"

var cookieStore = sessions.NewCookieStore(utils.GetRandomBytes(16))

func Login(c echo.Context) error {
	login := c.FormValue("username")
	passwd := c.FormValue("password")

	if login == ADMIN_LOGIN && bcrypt.CompareHashAndPassword([]byte(ADMIN_PASSWD), []byte(passwd)) == nil {
		sess, _ := cookieStore.Get(c.Request(), cookieName)
		if !sess.IsNew {
			sess.Options.MaxAge = 86400 * 30 // extend to 1 month
		}
		sess.Values["authenticated"] = true
		if err := sess.Save(c.Request(), c.Response()); err != nil {
			return err
		}
		return c.Redirect(http.StatusFound, "/cnc/workers")
	}
	return c.Redirect(http.StatusUnauthorized, "/cnc/login")
}

func Logout(c echo.Context) error {
	sess, _ := cookieStore.Get(c.Request(), cookieName)
	sess.Options.MaxAge = -1
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return err
	}
	return c.Redirect(http.StatusUnauthorized, "/cnc/login")
}

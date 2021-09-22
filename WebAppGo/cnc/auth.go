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
	request := c.Request()
	if request.Method == "GET" {
		return c.Render(http.StatusOK, "login.html", "")
	}
	// else POST
	login := c.FormValue("username")
	passwd := c.FormValue("password")

	if login == ADMIN_LOGIN && bcrypt.CompareHashAndPassword([]byte(ADMIN_PASSWD), []byte(passwd)) == nil {
		sess, _ := cookieStore.Get(request, cookieName)
		if !sess.IsNew {
			sess.Options.MaxAge = 86400 * 30 // extend to 1 month
		}
		sess.Values["authenticated"] = true
		if err := sess.Save(request, c.Response()); err != nil {
			return err
		}
		return c.Redirect(http.StatusSeeOther, "/cnc/workers")
	}
	return c.Render(http.StatusUnauthorized, "login.html", "")
}

func Logout(c echo.Context) error {
	sess, _ := cookieStore.Get(c.Request(), cookieName)
	sess.Options.MaxAge = -1
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return err
	}
	return c.Redirect(http.StatusSeeOther, "/cnc/login")
}

package main

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

var Users = make(map[string]*User)

type (
	User struct {
		Email  string `json:email`
		Tokens []string
	}

	Device struct {
		Email string `json:email`
		Token string `json:token`
	}
)

func main() {
	// Echo instance
	e := echo.New()

	// Routes
	e.POST("user/register", RegisterUser)
	e.POST("device/register", RegisterDevice)
	e.POST("device/push", Push)
	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}

// Handler
func RegisterUser(c echo.Context) error {
	param := new(User)
	if err := c.Bind(param); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status": err.Error(),
		})
	}
	user := Users[param.Email]
	if user == nil {
		newUser := &User{Email: param.Email}
		Users[param.Email] = newUser
		return c.JSON(http.StatusOK, echo.Map{
			"status": "Dang ky user thanh cong",
			"email":  param.Email,
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"status": "User da ton tai trong he thong",
		"email":  param.Email,
	})

}

func RegisterDevice(c echo.Context) error {

	param := new(Device)
	if err := c.Bind(param); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status": err.Error(),
		})
	}

	if len(param.Email) == 0 || len(param.Token) == 0 {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status": "Email or token not empty",
		})
	}
	user := Users[param.Email]
	user.Tokens = append(user.Tokens, param.Token)

	return c.JSON(http.StatusOK, echo.Map{
		"status": "Thiet bi da dang ky thanh cong",
		"email":  param.Email,
	})
}
func Push(c echo.Context) error {
	return errors.New("abc")
}

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
)

var Users = make(map[string]*User)

type (
	User struct {
		Email  string `json:"email"`
		Tokens []string
	}

	Device struct {
		Email string `json:"email"`
		Token string `json:"token"`
	}
	Data struct {
		Feature string `json:"feature"`
		Body    string `json:"body"`
	}
	Payload struct {
		RegistrationIds []string `json:"registration_ids"`
		Data            Data     `json:"data"`
	}
	Student struct {
		Id   int    `json:id`
		Name string `json:name`
	}
)

func main() {
	// Echo instance
	e := echo.New()

	// Routes
	e.POST("user/register", RegisterUser)
	e.POST("device/register", RegisterDevice)
	e.POST("device/push", Push)
	e.GET("students", getStudents)
	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}

// Handler
func getStudents(c echo.Context) error {
	students := [8]Student{
		Student{Id: 1, Name: "a"},
		Student{Id: 1, Name: "a"},
		Student{Id: 1, Name: "a"},
		Student{Id: 1, Name: "a"},
		Student{Id: 1, Name: "a"},
		Student{Id: 1, Name: "a"},
		Student{Id: 1, Name: "a"},
		Student{Id: 1, Name: "a"},
	}
	return c.JSON(http.StatusOK, students)
}
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

	fmt.Println(param.Token)

	return c.JSON(http.StatusOK, echo.Map{
		"status": "Thiet bi da dang ky thanh cong",
		"email":  param.Email,
	})
}
func Push(c echo.Context) error {
	param := new(Device)
	if err := c.Bind(param); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"status": err.Error(),
		})
	}

	user := Users[param.Email]
	if user == nil {
		return c.JSON(http.StatusOK, echo.Map{
			"status": "User nay khong ton tai",
		})
	}
	if len(user.Tokens) == 0 {
		return c.JSON(http.StatusOK, echo.Map{
			"status": "Thiet bi nay chua dang ky",
		})
	}

	fmt.Println("token la ", user.Tokens[0])
	//payload
	payload := Payload{
		RegistrationIds: user.Tokens,
		Data: Data{
			Feature: "WELCOME",
			Body:    "Hello, " + user.Email,
		},
	}

	payloadByte, _ := json.Marshal(payload)

	go pushNoti(payloadByte)
	return c.JSON(http.StatusOK, echo.Map{
		"status": "Processing...",
	})
}

func pushNoti(payload []byte) {
	req, _ := http.NewRequest("POST", "https://fcm.googleapis.com/fcm/send", bytes.NewBuffer(payload))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "key=AAAAh3U5wVk:APA91bHlQLiXXR5iRjpSugLZxV36oCZpsDMhwYdq7Yn8GhYcWYRGYY11yKpfKOG-yjw-6sEMZMzfTppdQEJX0x-7_fKeOQr9XKNL-UF2qbRBiwHWIhvg8D6RMJqcOnaR0t7e5ffr1--a")

	client := &http.Client{}
	res, _ := client.Do(req)

	bytes, _ := ioutil.ReadAll(res.Body)
	resString := string(bytes)

	//fmt.Println("hghffgff")

	fmt.Println(resString, "    ghfgghf")
	defer res.Body.Close()

}

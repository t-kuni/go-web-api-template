package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
	"github.com/samber/do"
	"github.com/t-kuni/go-web-api-template/const/app"
	"github.com/t-kuni/go-web-api-template/domain/service"
	"net/http"
)

type HelloHandler struct {
	ExampleService service.IExampleService
}

type HelloResponse struct {
	AppName   string                 `json:"appName"`
	Status    string                 `json:"status"`
	Companies []HelloResponseCompany `json:"companies"`
}

type HelloResponseCompany struct {
	Name  string              `json:"name"`
	Users []HelloResponseUser `json:"users"`
}

type HelloResponseUser struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func NewHelloHandler(i *do.Injector) (*HelloHandler, error) {
	return &HelloHandler{
		do.MustInvoke[service.IExampleService](i),
	}, nil
}

func (h HelloHandler) Hello(c echo.Context) error {
	status, companies, err := h.ExampleService.Exec(c.Request().Context(), "BNB")
	if err != nil {
		return eris.Wrap(err, "")
	}

	var resp HelloResponse
	resp.AppName = app.AppName
	resp.Status = status

	var respCompanies []HelloResponseCompany
	for _, company := range companies {
		var users []HelloResponseUser
		for _, user := range company.Edges.Users {
			users = append(users, HelloResponseUser{
				Name: user.Name,
				Age:  user.Age,
			})
		}
		respCompanies = append(respCompanies, HelloResponseCompany{
			Name:  company.Name,
			Users: users,
		})
	}
	resp.Companies = respCompanies

	return c.JSON(http.StatusOK, resp)
}

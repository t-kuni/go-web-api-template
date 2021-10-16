package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/t-kuni/go-web-api-skeleton/domain/service"
	"net/http"
)

type HelloHandler struct {
	ExampleService service.ExampleServiceInterface
}

type HelloResponse struct {
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

func ProvideHello(exampleService service.ExampleServiceInterface) HelloHandler {
	return HelloHandler{exampleService}
}

func (h HelloHandler) Hello(c echo.Context) error {
	status, companies, err := h.ExampleService.Exec(c.Request().Context(), "BNB")
	if err != nil {
		return err
	}

	var resp HelloResponse
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

package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/t-kuni/go-web-api-skeleton/domain/infrastructure/db"
	"net/http"
)

type PostUserHandler struct {
	DbConnector db.ConnectorInterface
}

type PostUserRequest struct {
	Name      string `json:"name" validate:"required"`
	Age       int    `json:"age" validate:"required"`
	CompanyId int    `json:"companyId"`
}

type PostUserResponse struct {
	Status string `json:"status"`
}

func ProvidePostUserHandler(dbConnector db.ConnectorInterface) *PostUserHandler {
	return &PostUserHandler{dbConnector}
}

func (h PostUserHandler) PostUser(c echo.Context) error {
	var req PostUserRequest

	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	err := h.DbConnector.BeginTx()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	defer h.DbConnector.Rollback()

	_, err = h.DbConnector.GetDB().Exec("INSERT INTO users(name, age, created_at) VALUES(?, ?, '2010-12-31 23:59:59')", req.Name, req.Age)
	if err != nil {
		return err
	}
	//_, err = h.DbConnector.GetEnt().User.Create().
	//	SetAge(req.Age).
	//	SetName(req.Name).
	//	Save(c.Request().Context())

	err = h.DbConnector.Commit()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	var resp PostUserResponse
	resp.Status = "OK"

	return c.JSON(http.StatusOK, resp)
}

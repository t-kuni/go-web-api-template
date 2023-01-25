package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
	"github.com/t-kuni/go-web-api-template/domain/infrastructure/db"
	"github.com/t-kuni/go-web-api-template/ent"
	"net/http"
)

type PostUserHandler struct {
	DbConnector db.ConnectorInterface
}

type PostUserRequest struct {
	Name      *string `json:"name" validate:"required"`
	Age       *int    `json:"age" validate:"gte=18,lte=60,required"`
	CompanyId *int    `json:"companyId" validate:"required"`
}

type PostUserResponse struct {
	Status string `json:"status"`
}

func NewPostUserHandler(i *do.Injector) (*PostUserHandler, error) {
	return &PostUserHandler{
		do.MustInvoke[db.ConnectorInterface](i),
	}, nil
}

func (h PostUserHandler) PostUser(c echo.Context) error {
	var req PostUserRequest

	err := c.Bind(&req)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	err = c.Validate(req)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	err = h.DbConnector.Transaction(c.Request().Context(), func(tx *ent.Client) error {
		_, err := tx.User.Create().
			SetAge(*req.Age).
			SetName(*req.Name).
			Save(c.Request().Context())
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	var resp PostUserResponse
	resp.Status = "OK"

	return c.JSON(http.StatusOK, resp)
}

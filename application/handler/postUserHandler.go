package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
	"github.com/samber/do"
	"github.com/t-kuni/go-web-api-template/domain/infrastructure/db"
	"github.com/t-kuni/go-web-api-template/ent"
	"github.com/t-kuni/go-web-api-template/errors/types"
	"net/http"
)

type PostUserHandler struct {
	DbConnector db.Connector
}

type PostUserRequest struct {
	Name *string `json:"name" validate:"required"`
	Age  *int    `json:"age" validate:"gte=8,lte=60,required"`
}

type PostUserResponse struct {
	Status string `json:"status"`
}

func NewPostUserHandler(i *do.Injector) (*PostUserHandler, error) {
	return &PostUserHandler{
		do.MustInvoke[db.Connector](i),
	}, nil
}

func (h PostUserHandler) PostUser(c echo.Context) error {
	var req PostUserRequest

	err := c.Bind(&req)
	if err != nil {
		return eris.Wrap(err, "")
	}

	err = c.Validate(req)
	if err != nil {
		return eris.Wrap(err, "")
	}

	if *req.Name == "admin" {
		err := types.NewBasicBusinessError("Can't use name 'admin'", nil)
		return eris.Wrap(err, "")
	}

	err = h.DbConnector.Transaction(c.Request().Context(), func(tx *ent.Client) error {
		_, err := tx.User.Create().
			SetAge(*req.Age).
			SetName(*req.Name).
			Save(c.Request().Context())
		if err != nil {
			return eris.Wrap(err, "")
		}
		return nil
	})
	if err != nil {
		return eris.Wrap(err, "")
	}

	var resp PostUserResponse
	resp.Status = "OK"

	return c.JSON(http.StatusOK, resp)
}

package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
	"github.com/samber/do"
	"github.com/t-kuni/go-web-api-template/application/handler"
	"github.com/t-kuni/go-web-api-template/errors/types"
	"github.com/t-kuni/go-web-api-template/logger"
	"net/http"
)

type ErrorHandler struct {
}

func NewErrorHandler(i *do.Injector) (*ErrorHandler, error) {
	return &ErrorHandler{}, nil
}

func (h ErrorHandler) Handler(err error, c echo.Context) {
	shouldResponse := !c.Response().Committed

	var httpError *echo.HTTPError
	var bindingError *echo.BindingError
	var validationError validator.ValidationErrors
	var basicBusinessError *types.BasicBusinessError

	if eris.As(err, &httpError) {
		logger.WarnWithError(c, err, nil)
		if shouldResponse {
			err = responseHttpError(httpError, c)
		}
	} else if eris.As(err, &bindingError) {
		logger.WarnWithError(c, err, nil)
		if shouldResponse {
			err = c.NoContent(http.StatusUnsupportedMediaType)
		}
	} else if eris.As(err, &validationError) {
		logger.WarnWithError(c, err, nil)
		if shouldResponse {
			err = c.JSON(http.StatusBadRequest, validationError.Error())
		}
	} else if eris.As(err, &basicBusinessError) {
		logger.WarnWithError(c, err, basicBusinessError.Params)
		if shouldResponse {
			body := h.makeBasicBusinessErrorResponse(basicBusinessError)
			err = c.JSON(http.StatusUnprocessableEntity, body)
		}
	} else {
		logger.Error(c, err, nil)
		if shouldResponse {
			err = c.NoContent(http.StatusInternalServerError)
		}
	}

	if shouldResponse && err != nil {
		logger.Error(c, err, nil)
	}
}

func (h ErrorHandler) makeBasicBusinessErrorResponse(businessError *types.BasicBusinessError) *handler.BusinessErrorResponse {
	return &handler.BusinessErrorResponse{
		Message: businessError.Message,
	}
}

// responseHttpError Returns a response based on echo.HTTPError
// Referring to the DefaultHTTPErrorHandler implementation of echo.
//
// https://github.com/labstack/echo/blob/v4.5.0/echo.go#L360
func responseHttpError(he *echo.HTTPError, c echo.Context) (err error) {
	if he.Internal != nil {
		if herr, ok := he.Internal.(*echo.HTTPError); ok {
			he = herr
		}
	}

	code := he.Code
	message := he.Message
	if m, ok := he.Message.(string); ok {
		message = echo.Map{"message": m}
	}

	if c.Request().Method == http.MethodHead {
		err = c.NoContent(he.Code)
	} else {
		err = c.JSON(code, message)
	}
	return err
}

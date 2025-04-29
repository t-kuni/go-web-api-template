package handler

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/t-kuni/go-web-api-template/domain/service"
	customErrors "github.com/t-kuni/go-web-api-template/errors"
	"github.com/t-kuni/go-web-api-template/restapi/models"
	"github.com/t-kuni/go-web-api-template/restapi/operations/companies"
)

type GetCompanies struct {
	exampleService service.IExampleService
}

func NewGetCompanies(exampleService service.IExampleService) (*GetCompanies, error) {
	return &GetCompanies{
		exampleService,
	}, nil
}

func (u GetCompanies) Main(params companies.GetCompaniesParams) middleware.Responder {
	_, companies2, err := u.exampleService.Exec(params.HTTPRequest.Context(), "BNB")
	if err != nil {
		return customErrors.NewErrorResponder(err)
	}

	respCompanies := []*models.Company{}
	for _, company := range companies2 {
		id := strfmt.UUID(company.ID)

		// ユーザー情報を変換
		users := []*models.User{}
		for _, user := range company.Edges.Users {
			userId := strfmt.UUID(user.ID)
			age := int64(user.Age)
			gender := user.Gender

			users = append(users, &models.User{
				ID:     &userId,
				Name:   &user.Name,
				Age:    &age,
				Gender: &gender,
				// Company フィールドは循環参照を避けるため null に設定
				Company: nil,
			})
		}

		respCompanies = append(respCompanies, &models.Company{
			ID:    &id,
			Name:  &company.Name,
			Users: users,
		})
	}

	response := companies.GetCompaniesOKBody{
		Companies: respCompanies,
	}
	return companies.NewGetCompaniesOK().WithPayload(&response)
}

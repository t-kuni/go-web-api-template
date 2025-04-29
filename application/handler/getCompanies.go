package handler

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/t-kuni/go-web-api-template/domain/service"
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
		return middleware.Error(500, err)
	}

	respCompanies := []*models.Company{}
	for _, company := range companies2 {
		id := strfmt.UUID(company.ID)
		respCompanies = append(respCompanies, &models.Company{
			ID:   &id,
			Name: &company.Name,
		})
	}

	response := companies.GetCompaniesOKBody{
		Companies: respCompanies,
	}
	return companies.NewGetCompaniesOK().WithPayload(&response)
}

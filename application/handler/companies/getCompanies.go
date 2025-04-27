package companies

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/t-kuni/go-web-api-template/domain/service"
	"github.com/t-kuni/go-web-api-template/restapi/operations/companies"
	"github.com/t-kuni/go-web-api-template/util"
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

	respCompanies := []*companies.GetCompaniesOKBodyCompaniesItems0{}
	for _, company := range companies2 {
		respCompanies = append(respCompanies, &companies.GetCompaniesOKBodyCompaniesItems0{
			ID:   util.Ptr(int64(company.ID)),
			Name: &company.Name,
		})
	}

	response := companies.GetCompaniesOKBody{
		Companies: respCompanies,
	}
	return companies.NewGetCompaniesOK().WithPayload(&response)
}

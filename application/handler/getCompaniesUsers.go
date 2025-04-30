package handler

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/t-kuni/go-web-api-template/domain/infrastructure/db"
	customErrors "github.com/t-kuni/go-web-api-template/errors"
	"github.com/t-kuni/go-web-api-template/ent"
	"github.com/t-kuni/go-web-api-template/ent/company"
	"github.com/t-kuni/go-web-api-template/restapi/models"
	"github.com/t-kuni/go-web-api-template/restapi/operations/companies"
	"github.com/t-kuni/go-web-api-template/util"
)

type GetCompaniesUsers struct {
	dbConnector db.IConnector
}

func NewGetCompaniesUsers(dbConnector db.IConnector) (*GetCompaniesUsers, error) {
	return &GetCompaniesUsers{
		dbConnector,
	}, nil
}

func (u GetCompaniesUsers) Main(params companies.GetCompaniesUsersParams) middleware.Responder {
	ctx := params.HTTPRequest.Context()
	client := u.dbConnector.GetEnt()

	// 指定されたIDの会社を取得
	companyID := params.ID
	companyEntity, err := client.Company.Query().
		Where(company.ID(companyID)).
		WithUsers().
		First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			// 会社が見つからない場合は空のユーザーリストを返す
			response := companies.GetCompaniesUsersOKBody{
				Users: []*models.User{},
			}
			return companies.NewGetCompaniesUsersOK().WithPayload(&response)
		}
		return customErrors.NewErrorResponder(err)
	}

	// レスポンス用のユーザーリストを作成
	respUsers := make([]*models.User, 0, len(companyEntity.Edges.Users))
	for _, usr := range companyEntity.Edges.Users {
		// 会社情報を設定
		companyID := strfmt.UUID(companyEntity.ID)
		companyName := companyEntity.Name
		company := &models.Company{
			ID:    &companyID,
			Name:  &companyName,
		}

		// 性別を設定
		gender := usr.Gender
		if gender == "" {
			gender = "man" // デフォルト値
		}

		// ユーザー情報を設定
		userID := strfmt.UUID(usr.ID)
		respUsers = append(respUsers, &models.User{
			ID:      &userID,
			Name:    &usr.Name,
			Age:     util.Ptr(int64(usr.Age)),
			Gender:  &gender,
			Company: company,
		})
	}

	// レスポンスを作成
	response := companies.GetCompaniesUsersOKBody{
		Users: respUsers,
	}

	return companies.NewGetCompaniesUsersOK().WithPayload(&response)
}

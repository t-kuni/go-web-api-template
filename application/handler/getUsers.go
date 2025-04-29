package handler

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/t-kuni/go-web-api-template/domain/infrastructure/db"
	customErrors "github.com/t-kuni/go-web-api-template/errors"
	"github.com/t-kuni/go-web-api-template/ent"
	"github.com/t-kuni/go-web-api-template/ent/user"
	"github.com/t-kuni/go-web-api-template/restapi/models"
	userOps "github.com/t-kuni/go-web-api-template/restapi/operations/user"
	"github.com/t-kuni/go-web-api-template/util"
)

type GetUsers struct {
	dbConnector db.IConnector
}

func NewGetUsers(dbConnector db.IConnector) (*GetUsers, error) {
	return &GetUsers{
		dbConnector,
	}, nil
}

func (u GetUsers) Main(params userOps.GetUsersParams) middleware.Responder {
	ctx := params.HTTPRequest.Context()
	client := u.dbConnector.GetEnt()

	// デフォルトのページ番号を設定
	page := int64(1)
	if params.Page != nil {
		page = *params.Page
	}

	// ページサイズを設定（固定値）
	pageSize := int64(10)

	// 総ユーザー数を取得
	totalCount, err := client.User.Query().Count(ctx)
	if err != nil {
		return customErrors.NewErrorResponder(err)
	}

	// 最大ページ数を計算
	maxPage := (int64(totalCount) + pageSize - 1) / pageSize

	// ユーザー一覧を取得（ページネーション付き）
	offset := (page - 1) * pageSize
	users, err := client.User.Query().
		Limit(int(pageSize)).
		Offset(int(offset)).
		WithCompany().
		Order(ent.Asc(user.FieldID)).
		All(ctx)
	if err != nil {
		return customErrors.NewErrorResponder(err)
	}

	// レスポンス用のユーザーリストを作成
	respUsers := make([]*models.User, 0, len(users))
	for _, usr := range users {
		// 会社情報を取得
		var company *models.Company
		if usr.Edges.Company != nil {
			companyID := strfmt.UUID(usr.Edges.Company.ID)
			companyName := usr.Edges.Company.Name
			company = &models.Company{
				ID:    &companyID,
				Name:  &companyName,
				Users: []*models.User{},
			}
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
	response := userOps.GetUsersOKBody{
		Users:   respUsers,
		Page:    &page,
		MaxPage: &maxPage,
	}

	return userOps.NewGetUsersOK().WithPayload(&response)
}

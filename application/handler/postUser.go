package handler

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/rotisserie/eris"
	"github.com/t-kuni/go-web-api-template/domain/infrastructure/db"
	"github.com/t-kuni/go-web-api-template/domain/infrastructure/system"
	customErrors "github.com/t-kuni/go-web-api-template/errors"
	"github.com/t-kuni/go-web-api-template/ent"
	"github.com/t-kuni/go-web-api-template/ent/company"
	"github.com/t-kuni/go-web-api-template/ent/user"
	userOps "github.com/t-kuni/go-web-api-template/restapi/operations/user"
)

type PostUser struct {
	dbConnector   db.IConnector
	uuidGenerator system.IUuidGenerator
}

func NewPostUser(dbConnector db.IConnector, uuidGenerator system.IUuidGenerator) (*PostUser, error) {
	return &PostUser{
		dbConnector:   dbConnector,
		uuidGenerator: uuidGenerator,
	}, nil
}

func (u PostUser) Main(params userOps.PostUsersParams) middleware.Responder {
	ctx := params.HTTPRequest.Context()

	// トランザクション開始
	err := u.dbConnector.Transaction(ctx, func(tx *ent.Client) error {
		// UUIDの生成または既存IDの使用
		var userID string
		if params.Body.ID == nil || *params.Body.ID == "" {
			// 新規IDを生成
			generatedID, err := u.uuidGenerator.Generate()
			if err != nil {
				return eris.Wrap(err, "")
			}
			userID = generatedID
		} else {
			// 指定されたIDを使用
			userID = params.Body.ID.String()
		}

		// 会社の存在確認
		var companyID string
		if params.Body.Company != nil && params.Body.Company.ID != nil {
			companyID = params.Body.Company.ID.String()
			exists, err := tx.Company.Query().Where(
				company.ID(companyID),
			).Exist(ctx)
			if err != nil {
				return eris.Wrap(err, "")
			}
			if !exists {
				return eris.New("指定された会社が存在しません")
			}
		} else {
			return eris.New("会社情報が必要です")
		}

		// ユーザーの作成または更新
		gender := ""
		if params.Body.Gender != nil {
			gender = *params.Body.Gender
		}

		// ユーザーが既に存在するか確認
		exists, err := tx.User.Query().Where(
			user.ID(userID),
		).Exist(ctx)
		if err != nil {
			return eris.Wrap(err, "")
		}

		if exists {
			// 既存ユーザーの更新
			_, err = tx.User.UpdateOneID(userID).
				SetName(*params.Body.Name).
				SetAge(int(*params.Body.Age)).
				SetGender(gender).
				SetCompanyID(companyID).
				Save(ctx)
			if err != nil {
				return eris.Wrap(err, "")
			}
		} else {
			// 新規ユーザーの作成
			_, err = tx.User.Create().
				SetID(userID).
				SetName(*params.Body.Name).
				SetAge(int(*params.Body.Age)).
				SetGender(gender).
				SetCompanyID(companyID).
				Save(ctx)
			if err != nil {
				return eris.Wrap(err, "")
			}
		}

		return nil
	})

	if err != nil {
		return customErrors.NewErrorResponder(err)
	}

	// 成功レスポンスを返す
	return userOps.NewPostUsersOK()
}

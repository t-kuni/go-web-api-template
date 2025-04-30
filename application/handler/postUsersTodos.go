package handler

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/rotisserie/eris"
	"github.com/t-kuni/go-web-api-template/domain/infrastructure/db"
	"github.com/t-kuni/go-web-api-template/domain/infrastructure/system"
	"github.com/t-kuni/go-web-api-template/ent"
	"github.com/t-kuni/go-web-api-template/ent/schema"
	"github.com/t-kuni/go-web-api-template/ent/user"
	customErrors "github.com/t-kuni/go-web-api-template/errors"
	"github.com/t-kuni/go-web-api-template/restapi/operations"
)

type PostUsersTodos struct {
	dbConnector   db.IConnector
	uuidGenerator system.IUuidGenerator
}

func NewPostUsersTodos(dbConnector db.IConnector, uuidGenerator system.IUuidGenerator) (*PostUsersTodos, error) {
	return &PostUsersTodos{
		dbConnector:   dbConnector,
		uuidGenerator: uuidGenerator,
	}, nil
}

func (u PostUsersTodos) Main(params operations.PostUsersIDTodosParams) middleware.Responder {
	ctx := params.HTTPRequest.Context()

	// トランザクション開始
	err := u.dbConnector.Transaction(ctx, func(tx *ent.Client) error {
		// ユーザーの存在確認
		userID := params.ID
		exists, err := tx.User.Query().Where(
			user.ID(userID),
		).Exist(ctx)
		if err != nil {
			return eris.Wrap(err, "")
		}
		if !exists {
			return eris.New("指定されたユーザーが存在しません")
		}

		// Todoのステータスを確認
		status := schema.TodoStatus(*params.Body.Status)
		if status != schema.StatusNone && 
		   status != schema.StatusProgress && 
		   status != schema.StatusPending && 
		   status != schema.StatusCompleted {
			return eris.New("無効なステータスです")
		}

		// UUIDの生成または既存IDの使用
		var todoID string
		if params.Body.ID == nil || *params.Body.ID == "" {
			// 新規IDを生成
			generatedID, err := u.uuidGenerator.Generate()
			if err != nil {
				return eris.Wrap(err, "")
			}
			todoID = generatedID
		} else {
			// 指定されたIDを使用
			todoID = params.Body.ID.String()
		}

		// Todoの作成
		_, err = tx.Todo.Create().
			SetID(todoID).
			SetTitle(*params.Body.Title).
			SetDescription(*params.Body.Description).
			SetStatus(status).
			SetCompleted(false). // 新規作成時は未完了
			SetOwnerID(userID).
			Save(ctx)
		if err != nil {
			return eris.Wrap(err, "")
		}

		return nil
	})

	if err != nil {
		return customErrors.NewErrorResponder(err)
	}

	// 成功レスポンスを返す
	return operations.NewPostUsersIDTodosOK()
}

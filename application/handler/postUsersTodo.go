package handler

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/rotisserie/eris"
	"github.com/t-kuni/go-web-api-template/domain/infrastructure/db"
	"github.com/t-kuni/go-web-api-template/domain/infrastructure/system"
	"github.com/t-kuni/go-web-api-template/ent"
	"github.com/t-kuni/go-web-api-template/ent/schema"
	"github.com/t-kuni/go-web-api-template/ent/todo"
	"github.com/t-kuni/go-web-api-template/ent/user"
	customErrors "github.com/t-kuni/go-web-api-template/errors"
	"github.com/t-kuni/go-web-api-template/restapi/operations"
	"time"
)

type PostUsersTodo struct {
	dbConnector   db.IConnector
	uuidGenerator system.IUuidGenerator
	timer         system.ITimer
}

func NewPostUsersTodo(dbConnector db.IConnector, uuidGenerator system.IUuidGenerator, timer system.ITimer) (*PostUsersTodo, error) {
	return &PostUsersTodo{
		dbConnector:   dbConnector,
		uuidGenerator: uuidGenerator,
		timer:         timer,
	}, nil
}

func (u PostUsersTodo) Main(params operations.PostUsersIDTodosParams) middleware.Responder {
	ctx := params.HTTPRequest.Context()

	// トランザクション開始
	err := u.dbConnector.Transaction(ctx, func(tx *ent.Client) error {
		// ユーザーの存在確認
		userExists, err := tx.User.Query().Where(
			user.ID(params.ID),
		).Exist(ctx)
		if err != nil {
			return eris.Wrap(err, "")
		}
		if !userExists {
			return eris.New("指定されたユーザーが存在しません")
		}

		// TodoのUUIDの生成または既存IDの使用
		var todoID string
		if params.Body.ID == nil || params.Body.ID.String() == "" {
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

		// Todoのステータスを取得
		status := schema.StatusNone
		if params.Body.Status != nil {
			switch *params.Body.Status {
			case "none":
				status = schema.StatusNone
			case "progress":
				status = schema.StatusProgress
			case "pending":
				status = schema.StatusPending
			case "complated":
				status = schema.StatusCompleted
			default:
				return eris.New("無効なステータスが指定されました")
			}
		}

		// Todoが既に存在するか確認
		todoExists, err := tx.Todo.Query().Where(
			todo.ID(todoID),
		).Exist(ctx)
		if err != nil {
			return eris.Wrap(err, "")
		}

		// completedフラグの設定
		completed := false
		var completedAt *time.Time
		if status == schema.StatusCompleted {
			completed = true
			now := u.timer.Now()
			completedAt = &now
		}

		if todoExists {
			// 既存Todoの更新
			_, err = tx.Todo.UpdateOneID(todoID).
				SetTitle(*params.Body.Title).
				SetDescription(*params.Body.Description).
				SetStatus(status).
				SetCompleted(completed).
				SetNillableCompletedAt(completedAt).
				SetOwnerID(params.ID).
				Save(ctx)
			if err != nil {
				return eris.Wrap(err, "")
			}
		} else {
			// 新規Todoの作成
			_, err = tx.Todo.Create().
				SetID(todoID).
				SetTitle(*params.Body.Title).
				SetDescription(*params.Body.Description).
				SetStatus(status).
				SetCompleted(completed).
				SetNillableCompletedAt(completedAt).
				SetOwnerID(params.ID).
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
	return operations.NewPostUsersIDTodosOK()
}

package handler

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/rotisserie/eris"
	"github.com/t-kuni/go-web-api-template/domain/infrastructure/db"
	"github.com/t-kuni/go-web-api-template/domain/infrastructure/system"
	customErrors "github.com/t-kuni/go-web-api-template/errors"
	"github.com/t-kuni/go-web-api-template/ent"
	"github.com/t-kuni/go-web-api-template/ent/schema"
	"github.com/t-kuni/go-web-api-template/ent/todo"
	"github.com/t-kuni/go-web-api-template/ent/user"
	"github.com/t-kuni/go-web-api-template/restapi/operations"
)

type PostUsersIdTodos struct {
	dbConnector   db.IConnector
	uuidGenerator system.IUuidGenerator
}

func NewPostUsersIdTodos(dbConnector db.IConnector, uuidGenerator system.IUuidGenerator) (*PostUsersIdTodos, error) {
	return &PostUsersIdTodos{
		dbConnector:   dbConnector,
		uuidGenerator: uuidGenerator,
	}, nil
}

func (u PostUsersIdTodos) Main(params operations.PostUsersIDTodosParams) middleware.Responder {
	ctx := params.HTTPRequest.Context()

	err := u.dbConnector.Transaction(ctx, func(tx *ent.Client) error {
		userID := params.ID

		exists, err := tx.User.Query().Where(user.ID(userID)).Exist(ctx)
		if err != nil {
			return eris.Wrap(err, "")
		}
		if !exists {
			return eris.New("指定されたユーザーが存在しません")
		}

		var todoID string
		if params.Body.ID == nil || params.Body.ID.String() == "" {
			generatedID, err := u.uuidGenerator.Generate()
			if err != nil {
				return eris.Wrap(err, "")
			}
			todoID = generatedID
		} else {
			todoID = params.Body.ID.String()
		}

		status := schema.StatusNone
		if params.Body.Status != nil {
			status = schema.TodoStatus(*params.Body.Status)
		}

		todoExists, err := tx.Todo.Query().Where(todo.ID(todoID)).Exist(ctx)
		if err != nil {
			return eris.Wrap(err, "")
		}

		if todoExists {
			_, err = tx.Todo.UpdateOneID(todoID).
				SetTitle(*params.Body.Title).
				SetDescription(*params.Body.Description).
				SetStatus(status).
				SetOwnerID(userID).
				Save(ctx)
			if err != nil {
				return eris.Wrap(err, "")
			}
		} else {
			_, err = tx.Todo.Create().
				SetID(todoID).
				SetTitle(*params.Body.Title).
				SetDescription(*params.Body.Description).
				SetStatus(status).
				SetOwnerID(userID).
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

	return operations.NewPostUsersIDTodosOK()
}
package handler

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/rotisserie/eris"
	"github.com/t-kuni/go-web-api-template/domain/infrastructure/db"
	"github.com/t-kuni/go-web-api-template/domain/infrastructure/system"
	customErrors "github.com/t-kuni/go-web-api-template/errors"
	"github.com/t-kuni/go-web-api-template/ent"
	"github.com/t-kuni/go-web-api-template/ent/user"
	"github.com/t-kuni/go-web-api-template/ent/schema"
	"github.com/t-kuni/go-web-api-template/restapi/models"
	"github.com/t-kuni/go-web-api-template/restapi/operations/user"
	"net/http"
	"time"
)

type PostUserTodos struct {
	dbConnector   db.IConnector
	uuidGenerator system.IUuidGenerator
	timer         system.ITimer
}

func NewPostUserTodos(dbConnector db.IConnector, uuidGenerator system.IUuidGenerator, timer system.ITimer) (*PostUserTodos, error) {
	return &PostUserTodos{
		dbConnector:   dbConnector,
		uuidGenerator: uuidGenerator,
		timer:         timer,
	}, nil
}

type PostUsersIDTodosParams struct {
	HTTPRequest *http.Request

	ID string

	Body *models.Todo
}

func NewPostUsersIDTodosOK() middleware.Responder {
	return middleware.ResponderFunc(func(rw http.ResponseWriter, p interface{}) {
		rw.WriteHeader(200)
	})
}

func (u PostUserTodos) Main(params user.PostUsersIDTodosParams) middleware.Responder {
	ctx := params.HTTPRequest.Context()

	err := u.dbConnector.Transaction(ctx, func(tx *ent.Client) error {
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

		var todoID string
		if params.Body.ID == nil || *params.Body.ID == "" {
			generatedID, err := u.uuidGenerator.Generate()
			if err != nil {
				return eris.Wrap(err, "")
			}
			todoID = generatedID
		} else {
			todoID = params.Body.ID.String()
		}

		status := string(*params.Body.Status)
		completed := false
		var completedAt *time.Time

		if status == string(schema.StatusCompleted) {
			completed = true
			now := u.timer.Now()
			completedAt = &now
		}

		todoCreator := tx.Todo.Create().
			SetID(todoID).
			SetTitle(*params.Body.Title).
			SetDescription(*params.Body.Description).
			SetStatus(status).
			SetCompleted(completed).
			SetNillableCompletedAt(completedAt).
			SetOwnerID(userID)

		_, err = todoCreator.Save(ctx)
		if err != nil {
			return eris.Wrap(err, "")
		}

		return nil
	})

	if err != nil {
		return customErrors.NewErrorResponder(err)
	}

	return NewPostUsersIDTodosOK()
}

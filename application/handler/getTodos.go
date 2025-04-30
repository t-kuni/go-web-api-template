package handler

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/t-kuni/go-web-api-template/domain/infrastructure/db"
	"github.com/t-kuni/go-web-api-template/ent"
	"github.com/t-kuni/go-web-api-template/ent/schema"
	"github.com/t-kuni/go-web-api-template/ent/todo"
	customErrors "github.com/t-kuni/go-web-api-template/errors"
	"github.com/t-kuni/go-web-api-template/restapi/models"
	"github.com/t-kuni/go-web-api-template/restapi/operations/todos"
	"github.com/t-kuni/go-web-api-template/util"
)

type ListTodos struct {
	DBConnector db.IConnector
}

func NewListTodos(conn db.IConnector) (*ListTodos, error) {
	return &ListTodos{
		conn,
	}, nil
}

func (u ListTodos) Main(params todos.GetTodosParams) middleware.Responder {
	// Todoクエリを作成し、所有者情報を取得
	todoQuery := u.DBConnector.GetEnt().Todo.Query().WithOwner(func(query *ent.UserQuery) {
		query.WithCompany() // 所有者の会社情報も取得
	})

	// Handle status filter
	if params.Status != nil {
		todoQuery = todoQuery.Where(todo.StatusEQ(schema.TodoStatus(*params.Status)))
	}

	// Handle pagination
	perPage := int64(10)
	//if params.PerPage != nil {
	//	perPage = *params.PerPage
	//}

	page := int64(1)
	if params.Page != nil {
		page = *params.Page
	}

	offset := (page - 1) * perPage
	todoList, err := todoQuery.Limit(int(perPage)).Offset(int(offset)).All(params.HTTPRequest.Context())
	if err != nil {
		return customErrors.NewErrorResponder(err)
	}

	total, err := todoQuery.Count(params.HTTPRequest.Context())
	if err != nil {
		return customErrors.NewErrorResponder(err)
	}

	maxPage := (int64(total) + perPage - 1) / perPage

	// Convert ent.Todo to models.Todo
	var apiTodoList []*models.Todo
	for _, entTodo := range todoList {
		id, err := util.StringToStrfmtUUID(entTodo.ID)
		if err != nil {
			return customErrors.NewErrorResponder(err)
		}

		// 所有者情報を設定
		var owner *models.User
		if entTodo.Edges.Owner != nil {
			ownerID := strfmt.UUID(entTodo.Edges.Owner.ID)
			ownerName := entTodo.Edges.Owner.Name
			gender := entTodo.Edges.Owner.Gender
			if gender == "" {
				gender = "man" // デフォルト値
			}

			owner = &models.User{
				ID:     &ownerID,
				Name:   &ownerName,
				Age:    util.Ptr(int64(entTodo.Edges.Owner.Age)),
				Gender: &gender,
			}

			// 所有者の会社情報も設定（存在する場合）
			if entTodo.Edges.Owner.Edges.Company != nil {
				companyID := strfmt.UUID(entTodo.Edges.Owner.Edges.Company.ID)
				companyName := entTodo.Edges.Owner.Edges.Company.Name
				owner.Company = &models.Company{
					ID:   &companyID,
					Name: &companyName,
				}
			}
		}

		apiTodoList = append(apiTodoList, &models.Todo{
			ID:          util.Ptr(id),
			Title:       &entTodo.Title,
			Description: &entTodo.Description,
			Status:      util.Ptr(string(entTodo.Status)),
			Completed:   entTodo.Completed,
			CompletedAt: func() strfmt.DateTime {
				if entTodo.CompletedAt != nil {
					return strfmt.DateTime(*entTodo.CompletedAt)
				}
				return strfmt.DateTime{}
			}(),
			Owner: owner,
		})
	}

	return todos.NewGetTodosOK().WithPayload(&todos.GetTodosOKBody{
		Todos:   apiTodoList,
		Page:    &page,
		MaxPage: &maxPage,
	})
}

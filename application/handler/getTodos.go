package handler

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/t-kuni/go-web-api-template/domain/infrastructure/db"
	"github.com/t-kuni/go-web-api-template/ent/schema"
	"github.com/t-kuni/go-web-api-template/ent/todo"
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
	todoQuery := u.DBConnector.GetEnt().Todo.Query()

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
		return middleware.Error(500, err)
	}

	total, err := todoQuery.Count(params.HTTPRequest.Context())
	if err != nil {
		return middleware.Error(500, err)
	}

	maxPage := (int64(total) + perPage - 1) / perPage

	// Convert ent.Todo to models.Todo
	var apiTodoList []*models.Todo
	for _, entTodo := range todoList {
		id, err := util.StringToStrfmtUUID(entTodo.ID)
		if err != nil {
			return middleware.Error(500, err)
		}

		apiTodoList = append(apiTodoList, &models.Todo{
			ID:          util.Ptr(id),
			Title:       &entTodo.Title,
			Description: &entTodo.Description,
			Status:      util.Ptr(string(entTodo.Status)),
			Completed:   entTodo.Completed,
		})
	}

	return todos.NewGetTodosOK().WithPayload(&todos.GetTodosOKBody{
		Todos:   apiTodoList,
		Page:    &page,
		MaxPage: &maxPage,
	})
}

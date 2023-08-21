package todos

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/samber/do"
	"github.com/t-kuni/go-web-api-template/domain/service"
	"github.com/t-kuni/go-web-api-template/restapi/models"
	"github.com/t-kuni/go-web-api-template/restapi/operations/todos"
	"github.com/t-kuni/go-web-api-template/util"
)

type Find struct {
	exampleService service.ExampleServiceInterface
}

func NewFind(i *do.Injector) (*Find, error) {
	return &Find{
		do.MustInvoke[service.ExampleServiceInterface](i),
	}, nil
}

func (u Find) Main(params todos.FindTodosParams) middleware.Responder {
	response := []*models.Item{
		{
			Completed:   true,
			Description: util.Ptr("aaa"),
			ID:          0,
		},
	}
	return todos.NewFindTodosOK().WithPayload(response)
}

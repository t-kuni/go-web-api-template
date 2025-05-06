//go:build feature

package handler_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/t-kuni/go-web-api-template/application/handler"
	"github.com/t-kuni/go-web-api-template/domain/infrastructure/db"
	"github.com/t-kuni/go-web-api-template/ent"
	"github.com/t-kuni/go-web-api-template/ent/schema"
	"github.com/t-kuni/go-web-api-template/restapi/models"
	"github.com/t-kuni/go-web-api-template/restapi/operations"
	"github.com/t-kuni/go-web-api-template/testUtil"
	"github.com/t-kuni/go-web-api-template/util"
)

func Test_PostUsersTodos(t *testing.T) {
	t.Run("正常にTodoを作成できること", func(t *testing.T) {
		//
		// Arrange
		//
		cont := testUtil.Prepare(t)
		defer cont.Finish()

		cont.SetTime("2020-04-10T00:00:00+09:00")

		// ユーザーデータを作成
		cont.PrepareTestData(func(db *ent.Client) {
			// 会社データを作成
			db.Company.Create().SetID("UUID-1").SetName("テスト株式会社").SaveX(t.Context())
			// ユーザーデータを作成
			db.User.Create().SetID("UUID-2").SetName("テストユーザー").SetAge(30).SetGender("man").SetCompanyID("UUID-1").SaveX(t.Context())
		})

		cont.Exec(func(testee *handler.PostUsersTodos, conn db.IConnector) {
			//
			// Act
			//
			title := "テストTodo"
			description := "これはテスト用のTodoです"
			status := string(schema.StatusProgress)

			params := operations.PostUsersIDTodosParams{
				HTTPRequest: util.Ptr(http.Request{}),
				ID:          "UUID-2", // ユーザーID
				Body: &models.Todo{
					Title:       &title,
					Description: &description,
					Status:      &status,
				},
			}
			resp := testee.Main(params)

			//
			// Assert
			//
			_, ok := resp.(*operations.PostUsersIDTodosOK)
			assert.True(t, ok)

			db := conn.GetEnt()
			todos, err := db.Todo.Query().WithOwner().All(t.Context())
			assert.NoError(t, err)
			assert.Equal(t, 1, len(todos))
			assert.Equal(t, "テストTodo", todos[0].Title)
			assert.Equal(t, "これはテスト用のTodoです", todos[0].Description)
			assert.Equal(t, schema.StatusProgress, todos[0].Status)
			assert.False(t, todos[0].Completed)
			assert.Nil(t, todos[0].CompletedAt)
			assert.Equal(t, "UUID-2", todos[0].Edges.Owner.ID)
		})
	})
}

package handler_test

import (
	"net/http"
	"testing"

	"github.com/go-openapi/strfmt"
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

func Test_PostUsersIdTodos(t *testing.T) {
	t.Run("正常にtodoを作成できること", func(t *testing.T) {
		//
		// Arrange
		//
		cont := testUtil.Prepare(t)
		defer cont.Finish()

		cont.SetTime("2020-04-10T00:00:00+09:00")

		// 会社データとユーザーデータを作成
		cont.PrepareTestData(func(db *ent.Client) {
			db.Company.Create().SetID("UUID-1").SetName("テスト株式会社").SaveX(t.Context())
			db.User.Create().SetID("UUID-2").SetName("テストユーザー").SetAge(30).SetGender("man").SetCompanyID("UUID-1").SaveX(t.Context())
		})

		cont.Exec(func(testee *handler.PostUsersIdTodos, conn db.IConnector) {
			//
			// Act
			//
			params := operations.PostUsersIDTodosParams{
				HTTPRequest: util.Ptr(http.Request{}),
				ID:          "UUID-2",
				Body: &models.Todo{
					Title:       util.Ptr("新しいタスク"),
					Description: util.Ptr("テスト用のタスクです"),
					Status:      util.Ptr("pending"),
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
			assert.Equal(t, "新しいタスク", todos[0].Title)
			assert.Equal(t, "テスト用のタスクです", todos[0].Description)
			assert.Equal(t, schema.StatusPending, todos[0].Status)
			assert.Equal(t, "UUID-2", todos[0].Edges.Owner.ID)
		})
	})

	t.Run("既存のtodoを更新できること", func(t *testing.T) {
		//
		// Arrange
		//
		cont := testUtil.Prepare(t)
		defer cont.Finish()

		cont.SetTime("2020-04-10T00:00:00+09:00")

		// 会社データ、ユーザーデータ、todoデータを作成
		cont.PrepareTestData(func(db *ent.Client) {
			db.Company.Create().SetID("UUID-1").SetName("テスト株式会社").SaveX(t.Context())
			db.User.Create().SetID("UUID-2").SetName("テストユーザー").SetAge(30).SetGender("man").SetCompanyID("UUID-1").SaveX(t.Context())
			db.Todo.Create().SetID("UUID-3").SetTitle("更新前タスク").SetDescription("更新前の説明").SetStatus(schema.StatusNone).SetOwnerID("UUID-2").SaveX(t.Context())
		})

		cont.Exec(func(testee *handler.PostUsersIdTodos, conn db.IConnector) {
			//
			// Act
			//
			params := operations.PostUsersIDTodosParams{
				HTTPRequest: util.Ptr(http.Request{}),
				ID:          "UUID-2",
				Body: &models.Todo{
					ID:          util.Ptr(strfmt.UUID("UUID-3")),
					Title:       util.Ptr("更新後タスク"),
					Description: util.Ptr("更新後の説明"),
					Status:      util.Ptr("progress"),
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
			assert.Equal(t, "UUID-3", todos[0].ID)
			assert.Equal(t, "更新後タスク", todos[0].Title)
			assert.Equal(t, "更新後の説明", todos[0].Description)
			assert.Equal(t, schema.StatusProgress, todos[0].Status)
			assert.Equal(t, "UUID-2", todos[0].Edges.Owner.ID)
		})
	})

	t.Run("存在しないユーザーでエラーが返ること", func(t *testing.T) {
		//
		// Arrange
		//
		cont := testUtil.Prepare(t)
		defer cont.Finish()

		cont.SetTime("2020-04-10T00:00:00+09:00")

		cont.Exec(func(testee *handler.PostUsersIdTodos, conn db.IConnector) {
			//
			// Act
			//
			params := operations.PostUsersIDTodosParams{
				HTTPRequest: util.Ptr(http.Request{}),
				ID:          "non-existent-user",
				Body: &models.Todo{
					Title:       util.Ptr("新しいタスク"),
					Description: util.Ptr("テスト用のタスクです"),
					Status:      util.Ptr("pending"),
				},
			}
			resp := testee.Main(params)

			//
			// Assert
			//
			_, ok := resp.(*operations.PostUsersIDTodosOK)
			assert.False(t, ok)
		})
	})
}
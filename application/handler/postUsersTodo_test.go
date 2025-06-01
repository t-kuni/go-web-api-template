package handler_test

import (
	"net/http"
	"testing"
	"time"

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

func Test_PostUsersTodo(t *testing.T) {
	t.Run("正常にTodoを作成できること", func(t *testing.T) {
		//
		// Arrange
		//
		cont := testUtil.Prepare(t)
		defer cont.Finish()

		cont.SetTime("2020-04-10T00:00:00+09:00")

		// 会社とユーザーデータを作成
		cont.PrepareTestData(func(db *ent.Client) {
			db.Company.Create().SetID("UUID-1").SetName("テスト株式会社").SaveX(t.Context())
			db.User.Create().SetID("UUID-2").SetName("山田太郎").SetAge(30).SetGender("man").SetCompanyID("UUID-1").SaveX(t.Context())
		})

		cont.Exec(func(testee *handler.PostUsersTodo, conn db.IConnector) {
			//
			// Act
			//
			params := operations.PostUsersIDTodosParams{
				HTTPRequest: util.Ptr(http.Request{}),
				ID:          "UUID-2",
				Body: &models.Todo{
					Title:       util.Ptr("新規タスク"),
					Description: util.Ptr("新規タスクの説明"),
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
			assert.Equal(t, "新規タスク", todos[0].Title)
			assert.Equal(t, "新規タスクの説明", todos[0].Description)
			assert.Equal(t, schema.StatusProgress, todos[0].Status)
			assert.Equal(t, false, todos[0].Completed)
			assert.Nil(t, todos[0].CompletedAt)
			assert.Equal(t, "UUID-2", todos[0].Edges.Owner.ID)
		})
	})

	t.Run("完了状態のTodoを作成できること", func(t *testing.T) {
		//
		// Arrange
		//
		cont := testUtil.Prepare(t)
		defer cont.Finish()

		expectedTime := time.Date(2020, 4, 10, 0, 0, 0, 0, time.UTC)
		cont.SetTime("2020-04-10T00:00:00+00:00")

		// 会社とユーザーデータを作成
		cont.PrepareTestData(func(db *ent.Client) {
			db.Company.Create().SetID("UUID-1").SetName("テスト株式会社").SaveX(t.Context())
			db.User.Create().SetID("UUID-2").SetName("山田太郎").SetAge(30).SetGender("man").SetCompanyID("UUID-1").SaveX(t.Context())
		})

		cont.Exec(func(testee *handler.PostUsersTodo, conn db.IConnector) {
			//
			// Act
			//
			params := operations.PostUsersIDTodosParams{
				HTTPRequest: util.Ptr(http.Request{}),
				ID:          "UUID-2",
				Body: &models.Todo{
					Title:       util.Ptr("完了タスク"),
					Description: util.Ptr("完了タスクの説明"),
					Status:      util.Ptr("complated"),
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
			assert.Equal(t, "完了タスク", todos[0].Title)
			assert.Equal(t, "完了タスクの説明", todos[0].Description)
			assert.Equal(t, schema.StatusCompleted, todos[0].Status)
			assert.Equal(t, true, todos[0].Completed)
			assert.NotNil(t, todos[0].CompletedAt)
			assert.Equal(t, expectedTime, *todos[0].CompletedAt)
			assert.Equal(t, "UUID-2", todos[0].Edges.Owner.ID)
		})
	})

	t.Run("既存のTodoを更新できること", func(t *testing.T) {
		//
		// Arrange
		//
		cont := testUtil.Prepare(t)
		defer cont.Finish()

		expectedTime := time.Date(2020, 4, 10, 0, 0, 0, 0, time.UTC)
		cont.SetTime("2020-04-10T00:00:00+00:00")

		// 会社、ユーザー、Todoデータを作成
		cont.PrepareTestData(func(db *ent.Client) {
			db.Company.Create().SetID("UUID-1").SetName("テスト株式会社").SaveX(t.Context())
			db.User.Create().SetID("UUID-2").SetName("山田太郎").SetAge(30).SetGender("man").SetCompanyID("UUID-1").SaveX(t.Context())
			db.Todo.Create().
				SetID("UUID-3").
				SetTitle("更新前タスク").
				SetDescription("更新前タスクの説明").
				SetStatus(schema.StatusNone).
				SetCompleted(false).
				SetOwnerID("UUID-2").
				SaveX(t.Context())
		})

		cont.Exec(func(testee *handler.PostUsersTodo, conn db.IConnector) {
			//
			// Act
			//
			params := operations.PostUsersIDTodosParams{
				HTTPRequest: util.Ptr(http.Request{}),
				ID:          "UUID-2",
				Body: &models.Todo{
					ID:          util.Ptr(strfmt.UUID("UUID-3")),
					Title:       util.Ptr("更新後タスク"),
					Description: util.Ptr("更新後タスクの説明"),
					Status:      util.Ptr("complated"),
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
			assert.Equal(t, "更新後タスクの説明", todos[0].Description)
			assert.Equal(t, schema.StatusCompleted, todos[0].Status)
			assert.Equal(t, true, todos[0].Completed)
			assert.NotNil(t, todos[0].CompletedAt)
			assert.Equal(t, expectedTime, *todos[0].CompletedAt)
			assert.Equal(t, "UUID-2", todos[0].Edges.Owner.ID)
		})
	})
}

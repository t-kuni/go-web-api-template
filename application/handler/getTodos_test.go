//go:build feature

package handler_test

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/t-kuni/go-web-api-template/application/handler"
	"github.com/t-kuni/go-web-api-template/ent"
	"github.com/t-kuni/go-web-api-template/ent/schema"
	"github.com/t-kuni/go-web-api-template/restapi/operations/todos"
	"github.com/t-kuni/go-web-api-template/testUtil"
	"github.com/t-kuni/go-web-api-template/util"
	"net/http"
	"testing"
	"time"
)

func Test_GetTodos(t *testing.T) {
	t.Run("正常にレスポンスを返すこと", func(t *testing.T) {
		//
		// Arrange
		//
		cont := testUtil.Prepare(t)
		defer cont.Finish()

		cont.SetTime("2020-04-10T00:00:00+09:00")

		completedAt := time.Date(2020, 4, 1, 12, 0, 0, 0, time.UTC)

		cont.PrepareTestData(func(db *ent.Client) {
			db.Company.Create().SetID("UUID-1").SetName("テスト株式会社").SaveX(t.Context())
			db.User.Create().SetID("UUID-2").SetName("山田太郎").SetAge(30).SetGender("man").SetCompanyID("UUID-1").SaveX(t.Context())
			db.User.Create().SetID("UUID-3").SetName("佐藤花子").SetAge(25).SetGender("woman").SetCompanyID("UUID-1").SaveX(t.Context())
			db.Todo.Create().
				SetID("UUID-4").
				SetTitle("タスク1").
				SetDescription("タスク1の説明").
				SetStatus(schema.StatusNone).
				SetCompleted(false).
				SetOwnerID("UUID-2").
				SaveX(t.Context())
			db.Todo.Create().
				SetID("UUID-5").
				SetTitle("タスク2").
				SetDescription("タスク2の説明").
				SetStatus(schema.StatusProgress).
				SetCompleted(false).
				SetOwnerID("UUID-3").
				SaveX(t.Context())
			db.Todo.Create().
				SetID("UUID-6").
				SetTitle("タスク3").
				SetDescription("タスク3の説明").
				SetStatus(schema.StatusCompleted).
				SetCompleted(true).
				SetCompletedAt(completedAt).
				SetOwnerID("UUID-2").
				SaveX(t.Context())
		})

		cont.Exec(func(testee *handler.ListTodos) {
			//
			// Act
			//
			params := todos.GetTodosParams{
				HTTPRequest: util.Ptr(http.Request{}),
				Page:        util.Ptr(int64(1)),
			}
			resp := testee.Main(params)

			//
			// Assert
			//
			okResp, ok := resp.(*todos.GetTodosOK)
			assert.True(t, ok)
			actualBody, err := json.Marshal(okResp.Payload)
			assert.NoError(t, err)

			expectBody := `
{
  "todos": [
    {
      "id": "UUID-4",
      "title": "タスク1",
      "description": "タスク1の説明",
      "status": "none",
      "completed": false,
      "completed_at": "0001-01-01T00:00:00.000Z",
      "owner": {
        "id": "UUID-2",
        "name": "山田太郎",
        "age": 30,
        "gender": "man",
        "company": {
          "id": "UUID-1",
          "name": "テスト株式会社"
        }
      }
    },
    {
      "id": "UUID-5",
      "title": "タスク2",
      "description": "タスク2の説明",
      "status": "progress",
      "completed": false,
      "completed_at": "0001-01-01T00:00:00.000Z",
      "owner": {
        "id": "UUID-3",
        "name": "佐藤花子",
        "age": 25,
        "gender": "woman",
        "company": {
          "id": "UUID-1",
          "name": "テスト株式会社"
        }
      }
    },
    {
      "id": "UUID-6",
      "title": "タスク3",
      "description": "タスク3の説明",
      "status": "completed",
      "completed": true,
      "completed_at": "2020-04-01T12:00:00.000Z",
      "owner": {
        "id": "UUID-2",
        "name": "山田太郎",
        "age": 30,
        "gender": "man",
        "company": {
          "id": "UUID-1",
          "name": "テスト株式会社"
        }
      }
    }
  ],
  "page": 1,
  "maxPage": 1
}`
			assert.JSONEq(t, expectBody, string(actualBody))
		})
	})
}

//go:build feature

package handler_test

import (
	"github.com/go-openapi/runtime"
	"github.com/stretchr/testify/assert"
	"github.com/t-kuni/go-web-api-template/application/handler"
	"github.com/t-kuni/go-web-api-template/ent"
	"github.com/t-kuni/go-web-api-template/restapi/operations/user"
	"github.com/t-kuni/go-web-api-template/testUtil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_GetUsers(t *testing.T) {
	t.Run("正常にレスポンスを返すこと", func(t *testing.T) {
		//
		// Arrange
		//
		cont := testUtil.Prepare(t)
		defer cont.Finish()

		cont.SetTime("2020-04-10T00:00:00+09:00")

		cont.PrepareTestData(func(db *ent.Client) {
			// 会社データを作成
			company := db.Company.Create().SetID("UUID-1").SetName("テスト株式会社").SaveX(t.Context())
			
			// ユーザーデータを作成
			db.User.Create().SetID("UUID-2").SetName("山田太郎").SetAge(30).SetGender("man").SetCompany(company).SaveX(t.Context())
			db.User.Create().SetID("UUID-3").SetName("佐藤花子").SetAge(25).SetGender("woman").SetCompany(company).SaveX(t.Context())
		})

		cont.Exec(func(testee *handler.GetUsers) {
			//
			// Act
			//
			req, err := http.NewRequest(http.MethodGet, "http://example.com/users", nil)
			assert.NoError(t, err)
			
			// ページパラメータを設定
			page := int64(1)
			params := user.GetUsersParams{
				HTTPRequest: req,
				Page:        &page,
			}
			
			resp := testee.Main(params)

			recorder := httptest.NewRecorder()
			producer := runtime.JSONProducer()
			resp.WriteResponse(recorder, producer)
			actualBody := recorder.Body.String()

			//
			// Assert
			//
			expectBody := `
{
  "users": [
    {
      "id": "UUID-2",
      "name": "山田太郎",
      "age": 30,
      "gender": "man",
      "company": {
        "id": "UUID-1",
        "name": "テスト株式会社"
      }
    },
    {
      "id": "UUID-3",
      "name": "佐藤花子",
      "age": 25,
      "gender": "woman",
      "company": {
        "id": "UUID-1",
        "name": "テスト株式会社"
      }
    }
  ],
  "page": 1,
  "maxPage": 1
}`
			assert.Equal(t, http.StatusOK, recorder.Code)
			assert.JSONEq(t, expectBody, actualBody)
		})
	})
}

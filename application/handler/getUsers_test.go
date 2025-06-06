package handler_test

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/t-kuni/go-web-api-template/application/handler"
	"github.com/t-kuni/go-web-api-template/ent"
	"github.com/t-kuni/go-web-api-template/restapi/operations/user"
	"github.com/t-kuni/go-web-api-template/testUtil"
	"github.com/t-kuni/go-web-api-template/util"
	"net/http"
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
			db.Company.Create().SetID("UUID-1").SetName("テスト株式会社").SaveX(t.Context())
			db.User.Create().SetID("UUID-2").SetName("山田太郎").SetAge(30).SetGender("man").SetCompanyID("UUID-1").SaveX(t.Context())
			db.User.Create().SetID("UUID-3").SetName("佐藤花子").SetAge(25).SetGender("woman").SetCompanyID("UUID-1").SaveX(t.Context())
		})

		cont.Exec(func(testee *handler.GetUsers) {
			//
			// Act
			//
			params := user.GetUsersParams{
				HTTPRequest: util.Ptr(http.Request{}),
				Page:        util.Ptr(int64(1)),
			}
			resp := testee.Main(params)

			//
			// Assert
			//
			okResp, ok := resp.(*user.GetUsersOK)
			assert.True(t, ok)
			actualBody, err := json.Marshal(okResp.Payload)
			assert.NoError(t, err)

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
			assert.JSONEq(t, expectBody, string(actualBody))
		})
	})
}

//go:build feature

package handler_test

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/t-kuni/go-web-api-template/application/handler"
	"github.com/t-kuni/go-web-api-template/ent"
	"github.com/t-kuni/go-web-api-template/restapi/operations/companies"
	"github.com/t-kuni/go-web-api-template/testUtil"
	"github.com/t-kuni/go-web-api-template/util"
	"net/http"
	"testing"
)

func Test_GetCompaniesUsers(t *testing.T) {
	t.Run("正常にレスポンスを返すこと", func(t *testing.T) {
		//
		// Arrange
		//
		cont := testUtil.Prepare(t)
		defer cont.Finish()

		cont.SetTime("2020-04-10T00:00:00+09:00")

		cont.PrepareTestData(func(db *ent.Client) {
			company := db.Company.Create().SetID("UUID-1").SetName("テスト株式会社").SaveX(t.Context())
			db.User.Create().SetID("UUID-2").SetName("山田太郎").SetAge(30).SetGender("man").SetCompany(company).SaveX(t.Context())
			db.User.Create().SetID("UUID-3").SetName("佐藤花子").SetAge(25).SetGender("woman").SetCompany(company).SaveX(t.Context())
		})

		cont.Exec(func(testee *handler.GetCompaniesUsers) {
			//
			// Act
			//
			params := companies.GetCompaniesUsersParams{
				HTTPRequest: util.Ptr(http.Request{}),
				ID:          "UUID-1",
			}
			resp := testee.Main(params)

			//
			// Assert
			//
			okResp, ok := resp.(*companies.GetCompaniesUsersOK)
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
  ]
}`
			assert.JSONEq(t, expectBody, string(actualBody))
		})
	})

	t.Run("存在しない会社IDの場合は空のユーザーリストを返すこと", func(t *testing.T) {
		//
		// Arrange
		//
		cont := testUtil.Prepare(t)
		defer cont.Finish()

		cont.SetTime("2020-04-10T00:00:00+09:00")

		cont.PrepareTestData(func(db *ent.Client) {
			// テストデータは作成しない
		})

		cont.Exec(func(testee *handler.GetCompaniesUsers) {
			//
			// Act
			//
			params := companies.GetCompaniesUsersParams{
				HTTPRequest: util.Ptr(http.Request{}),
				ID:          "non-existent-id",
			}
			resp := testee.Main(params)

			//
			// Assert
			//
			okResp, ok := resp.(*companies.GetCompaniesUsersOK)
			assert.True(t, ok)
			actualBody, err := json.Marshal(okResp.Payload)
			assert.NoError(t, err)

			expectBody := `
{
  "users": []
}`
			assert.JSONEq(t, expectBody, string(actualBody))
		})
	})
}

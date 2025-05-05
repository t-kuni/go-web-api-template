package handler_test

import (
	"encoding/json"
	"github.com/t-kuni/go-web-api-template/util"
	"net/http"
	"testing"

	"github.com/t-kuni/go-web-api-template/application/handler"

	"github.com/stretchr/testify/assert"
	"github.com/t-kuni/go-web-api-template/ent"
	"github.com/t-kuni/go-web-api-template/restapi/operations/companies"
	"github.com/t-kuni/go-web-api-template/testUtil"
)

func Test_GetCompanies(t *testing.T) {
	t.Run("正常にレスポンスを返すこと", func(t *testing.T) {
		//
		// Arrange
		//
		cont := testUtil.Prepare(t)
		defer cont.Finish()

		cont.SetTime("2020-04-10T00:00:00+09:00")

		cont.PrepareTestData(func(db *ent.Client) {
			db.Company.Create().SetID("UUID-1").SetName("NAME1").SaveX(t.Context())
			db.User.Create().SetID("UUID-10").SetCompanyID("UUID-1").SetName("NAME1").SetAge(20).SetGender("GENDER_1").SaveX(t.Context())
		})

		cont.Exec(func(testee *handler.GetCompanies) {
			//
			// Act
			//
			resp := testee.Main(companies.GetCompaniesParams{
				HTTPRequest: util.Ptr(http.Request{}),
			})

			//
			// Assert
			//
			okResp, ok := resp.(*companies.GetCompaniesOK)
			assert.True(t, ok)
			actualBody, err := json.Marshal(okResp.Payload)
			assert.NoError(t, err)

			expectBody := `
{
  "companies": [
    {
      "id": "UUID-1",
      "name": "NAME1"
    }
  ]
}`
			assert.JSONEq(t, expectBody, string(actualBody))
		})
	})
}

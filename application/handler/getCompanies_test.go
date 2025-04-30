//go:build feature

package handler_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/t-kuni/go-web-api-template/application/handler"

	"github.com/go-openapi/runtime"
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
			body := `{
		"key": "value"
	}`
			req, err := http.NewRequest(http.MethodPost, "http://example.com", bytes.NewBuffer([]byte(body)))
			assert.NoError(t, err)
			resp := testee.Main(companies.GetCompaniesParams{
				HTTPRequest: req,
			})

			recorder := httptest.NewRecorder()
			producer := runtime.JSONProducer()
			resp.WriteResponse(recorder, producer)
			actualBody := recorder.Body.String()

			//
			// Assert
			//
			expectBody := `
{
  "companies": [
    {
      "id": "UUID-1",
      "name": "NAME1"
    }
  ]
}`
			assert.Equal(t, http.StatusOK, recorder.Code)
			assert.JSONEq(t, expectBody, actualBody)
		})
	})
}

//go:build feature

package handler_test

import (
	"bytes"
	"github.com/t-kuni/go-web-api-template/application/handler"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-openapi/runtime"
	"github.com/stretchr/testify/assert"
	"github.com/t-kuni/go-web-api-template/domain/service"
	"github.com/t-kuni/go-web-api-template/ent"
	"github.com/t-kuni/go-web-api-template/restapi/operations/companies"
	"github.com/t-kuni/go-web-api-template/testUtil"
	"go.uber.org/mock/gomock"
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
			db.Company.Create().SetID(1).SetName("NAME1").SaveX(t.Context())
		})

		{
			mock := service.NewMockIExampleService(cont.MockCtrl)
			mock.
				EXPECT().
				Exec(gomock.Any(), gomock.Eq("BNB")).
				Return("DUMMY", []*ent.Company{
					{
						ID:        1,
						Name:      "TEST",
						CreatedAt: testUtil.MustNewDateTime("2006-01-02T15:04:05+09:00"),
						Edges:     ent.CompanyEdges{},
					},
				}, nil)
			testUtil.Override[service.IExampleService](cont, mock)
		}

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
      "id": 1,
      "name": "TEST"
    }
  ]
}`
			assert.Equal(t, http.StatusOK, recorder.Code)
			assert.JSONEq(t, expectBody, actualBody)
		})
	})
}

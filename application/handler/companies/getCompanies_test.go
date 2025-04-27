//go:build feature

package companies_test

import (
	"bytes"
	"github.com/go-openapi/runtime"
	"github.com/stretchr/testify/assert"
	"github.com/t-kuni/go-web-api-template/application/handler/companies"
	"github.com/t-kuni/go-web-api-template/domain/service"
	"github.com/t-kuni/go-web-api-template/ent"
	companies2 "github.com/t-kuni/go-web-api-template/restapi/operations/companies"
	"github.com/t-kuni/go-web-api-template/testUtil"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func Test_a(t *testing.T) {
	// Arrange
	cont := testUtil.Prepare(t)
	defer cont.Finish()

	cont.SetTime("2020-04-10T00:00:00+09:00")

	cont.PrepareTestData(func(db *ent.Client) {
		db.Company.Create().SetID(1).SetName("NAME1").SaveX(t.Context())
	})

	exampleServiceMock := service.NewMockIExampleService(cont.MockCtrl)
	createdAt, err := time.Parse("2006-01-02 15:04:05 MST", "2014-12-31 12:31:24 JST")
	assert.NoError(t, err)
	exampleServiceMock.
		EXPECT().
		Exec(gomock.Any(), gomock.Eq("BNB")).
		Return("DUMMY", []*ent.Company{
			{
				ID:        1,
				Name:      "TEST",
				CreatedAt: createdAt,
				Edges:     ent.CompanyEdges{},
			},
		}, nil)

	// Act
	body := `{
		"key": "value"
	}`
	testee, err := companies.NewGetCompanies(exampleServiceMock)
	assert.NoError(t, err)
	req, err := http.NewRequest(http.MethodPost, "http://example.com", bytes.NewBuffer([]byte(body)))
	assert.NoError(t, err)
	resp := testee.Main(companies2.GetCompaniesParams{
		HTTPRequest: req,
	})

	recorder := httptest.NewRecorder()
	producer := runtime.JSONProducer()
	resp.WriteResponse(recorder, producer)
	actualBody := recorder.Body.String()

	// Assert
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
}

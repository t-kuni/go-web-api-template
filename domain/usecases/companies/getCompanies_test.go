//go:build feature

package companies_test

import (
	"bytes"
	"database/sql"
	"github.com/go-openapi/runtime"
	"github.com/golang/mock/gomock"
	"github.com/joho/godotenv"
	"github.com/samber/do"
	"github.com/stretchr/testify/assert"
	"github.com/t-kuni/go-web-api-template/di"
	"github.com/t-kuni/go-web-api-template/domain/infrastructure/db"
	"github.com/t-kuni/go-web-api-template/domain/usecases/companies"
	dbImpl "github.com/t-kuni/go-web-api-template/infrastructure/db"
	companies2 "github.com/t-kuni/go-web-api-template/restapi/operations/companies"
	"github.com/t-kuni/go-web-api-template/testutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
)

func Test_a(t *testing.T) {
	// Arrange
	godotenv.Load(filepath.Join("..", "..", "..", ".env.feature"))

	app := di.NewApp()
	defer app.Shutdown()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	do.Override[db.Connector](app, dbImpl.NewTestConnector)

	d := do.MustInvoke[db.Connector](app).GetDB()
	testutil.PrepareTestData(d, func(db *sql.DB) {
		testutil.MustInsert(d, "companies", []map[string]interface{}{
			{"id": 1, "name": "NAME1", "created_at": "2020-05-10 10:00:00"},
		})
	})

	// Act
	body := `{
		"key": "value"
	}`
	testee, err := companies.NewGetCompanies(app)
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
      "name": "NAME1"
    }
  ]
}`
	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.JSONEq(t, expectBody, actualBody)
}

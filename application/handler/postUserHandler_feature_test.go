//go:build feature

package handler_test

import (
	"bytes"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
	"github.com/stretchr/testify/assert"
	"github.com/t-kuni/go-web-api-template/application/handler"
	"github.com/t-kuni/go-web-api-template/di"
	db2 "github.com/t-kuni/go-web-api-template/domain/infrastructure/db"
	"github.com/t-kuni/go-web-api-template/ent/user"
	"github.com/t-kuni/go-web-api-template/infrastructure/db"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

// テスト対象にトランザクション＋コミットが含まれるパターン
func TestPostUserHandler_PostUser(t *testing.T) {
	//
	// Prepare
	//
	app := di.NewApp()
	defer app.Shutdown()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	do.Override[db2.ConnectorInterface](app, db.NewTestConnector)

	body, err := json.Marshal(handler.PostUserRequest{
		Name:      "TEST",
		Age:       10,
		CompanyId: 1,
	})
	if err != nil {
		t.Fatal(err)
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	dbConnector := do.MustInvoke[db2.ConnectorInterface](app)
	db := dbConnector.GetDB()

	_, err = db.Exec("SET FOREIGN_KEY_CHECKS = 0")
	if err != nil {
		t.Fatal(err)
	}

	tx, err := dbConnector.GetEnt().Tx(c.Request().Context())
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { tx.Rollback() })

	_, err = db.Exec("SET FOREIGN_KEY_CHECKS = 1")
	if err != nil {
		t.Fatal(err)
	}

	//
	// Execute
	//
	h := do.MustInvoke[*handler.PostUserHandler](app)
	err = h.PostUser(c)
	if err != nil {
		t.Fatal(err)
	}

	//
	// Assert
	//
	var res handler.PostUserResponse
	buf, err := ioutil.ReadAll(rec.Body)
	if err != nil {
		t.Fatal(err)
	}
	err = json.Unmarshal(buf, &res)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "OK", res.Status)

	users, err := tx.User.Query().
		Where(user.Name("TEST"), user.Age(10)).
		All(c.Request().Context())
	if err != nil {
		t.Fatal(err)
	}
	assert.Len(t, users, 1)
}

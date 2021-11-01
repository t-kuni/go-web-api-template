//go:build feature

package handler_test

import (
	"bytes"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/t-kuni/go-web-api-skeleton/ent/user"
	db2 "github.com/t-kuni/go-web-api-skeleton/infrastructure/db"
	"github.com/t-kuni/go-web-api-skeleton/interface/handler"
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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

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

	// トランザクションをモック化
	transactionMock := db2.NewMockTransactionInterface(ctrl)
	dbConnector.Tx = transactionMock
	transactionMock.EXPECT().Begin(gomock.Any()).Return(tx, nil)
	transactionMock.EXPECT().Commit(gomock.Any()).Return(nil)
	transactionMock.EXPECT().Rollback(gomock.Any()).Return(nil)

	//
	// Execute
	//
	h := handler.ProvidePostUserHandler(dbConnector)
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

// //go:build feature
package handler_test

//
//import (
//	"bytes"
//	"encoding/json"
//	"github.com/golang/mock/gomock"
//	"github.com/labstack/echo/v4"
//	"github.com/samber/do"
//	"github.com/stretchr/testify/assert"
//	"github.com/t-kuni/go-web-api-template/application/handler"
//	"github.com/t-kuni/go-web-api-template/di"
//	db2 "github.com/t-kuni/go-web-api-template/domain/infrastructure/db"
//	"github.com/t-kuni/go-web-api-template/ent/user"
//	"github.com/t-kuni/go-web-api-template/infrastructure/db"
//	"github.com/t-kuni/go-web-api-template/server"
//	"github.com/t-kuni/go-web-api-template/util"
//	"io/ioutil"
//	"net/http"
//	"net/http/httptest"
//	"testing"
//)
//
//// テスト対象にトランザクション＋コミットが含まれるパターン
//func TestPostUserHandler_PostUser(t *testing.Param) {
//	//
//	// Prepare
//	//
//	app := di.NewApp()
//	defer app.Shutdown()
//
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	do.Override[db2.Connector](app, db.NewTestConnector)
//
//	body, err := json.Marshal(handler.PostUserRequest{
//		Name: util.Ptr[string]("TEST"),
//		Age:  util.Ptr[int](10),
//	})
//	assert.NoError(t, err)
//
//	e := do.MustInvoke[*server.Server](app).Echo
//	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
//	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
//	rec := httptest.NewRecorder()
//	c := e.NewContext(req, rec)
//
//	dbConnector := do.MustInvoke[db2.Connector](app)
//	db := dbConnector.GetDB()
//
//	_, err = db.Exec("SET FOREIGN_KEY_CHECKS = 0")
//	assert.NoError(t, err)
//
//	tx, err := dbConnector.GetEnt().Tx(c.Request().Context())
//	assert.NoError(t, err)
//	t.Cleanup(func() { tx.Rollback() })
//
//	_, err = db.Exec("SET FOREIGN_KEY_CHECKS = 1")
//	assert.NoError(t, err)
//
//	//
//	// Execute
//	//
//	h := do.MustInvoke[*handler.PostUserHandler](app)
//	err = h.PostUser(c)
//	assert.NoError(t, err)
//
//	//
//	// Assert
//	//
//	var res handler.PostUserResponse
//	buf, err := ioutil.ReadAll(rec.Body)
//	assert.NoError(t, err)
//	err = json.Unmarshal(buf, &res)
//	assert.NoError(t, err)
//
//	assert.Equal(t, "OK", res.Status)
//
//	users, err := tx.User.Query().
//		Where(user.Name("TEST"), user.Age(10)).
//		All(c.Request().Context())
//	assert.NoError(t, err)
//	assert.Len(t, users, 1)
//}

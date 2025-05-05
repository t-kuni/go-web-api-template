package handler_test

import (
	"net/http"
	"testing"

	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/assert"
	"github.com/t-kuni/go-web-api-template/application/handler"
	"github.com/t-kuni/go-web-api-template/domain/infrastructure/db"
	"github.com/t-kuni/go-web-api-template/ent"
	"github.com/t-kuni/go-web-api-template/restapi/models"
	"github.com/t-kuni/go-web-api-template/restapi/operations/user"
	"github.com/t-kuni/go-web-api-template/testUtil"
	"github.com/t-kuni/go-web-api-template/util"
)

func Test_PostUser(t *testing.T) {
	t.Run("正常にユーザーを作成できること", func(t *testing.T) {
		//
		// Arrange
		//
		cont := testUtil.Prepare(t)
		defer cont.Finish()

		cont.SetTime("2020-04-10T00:00:00+09:00")

		// 会社データを作成
		cont.PrepareTestData(func(db *ent.Client) {
			db.Company.Create().SetID("UUID-1").SetName("テスト株式会社").SaveX(t.Context())
		})

		cont.Exec(func(testee *handler.PostUser, conn db.IConnector) {
			//
			// Act
			//
			params := user.PostUsersParams{
				HTTPRequest: util.Ptr(http.Request{}),
				Body: &models.User{
					Company: &models.Company{
						ID:   util.Ptr(strfmt.UUID("UUID-1")),
						Name: util.Ptr("テスト株式会社"),
					},
					Name:   util.Ptr("新規ユーザー"),
					Age:    util.Ptr(int64(35)),
					Gender: util.Ptr("man"),
				},
			}
			resp := testee.Main(params)

			//
			// Assert
			//
			_, ok := resp.(*user.PostUsersOK)
			assert.True(t, ok)

			db := conn.GetEnt()
			users, err := db.User.Query().WithCompany().All(t.Context())
			assert.NoError(t, err)
			assert.Equal(t, 1, len(users))
			assert.Equal(t, "新規ユーザー", users[0].Name)
			assert.Equal(t, 35, users[0].Age)
			assert.Equal(t, "man", users[0].Gender)
			assert.Equal(t, "UUID-1", users[0].Edges.Company.ID)
		})
	})

	t.Run("既存のユーザーを更新できること", func(t *testing.T) {
		//
		// Arrange
		//
		cont := testUtil.Prepare(t)
		defer cont.Finish()

		cont.SetTime("2020-04-10T00:00:00+09:00")

		// 会社データとユーザーデータを作成
		cont.PrepareTestData(func(db *ent.Client) {
			db.Company.Create().SetID("UUID-1").SetName("テスト株式会社").SaveX(t.Context())
			db.Company.Create().SetID("UUID-2").SetName("テスト株式会社2").SaveX(t.Context())
			db.User.Create().SetID("UUID-2").SetName("更新前ユーザー").SetAge(30).SetGender("man").SetCompanyID("UUID-1").SaveX(t.Context())
		})

		cont.Exec(func(testee *handler.PostUser, conn db.IConnector) {
			//
			// Act
			//
			params := user.PostUsersParams{
				HTTPRequest: util.Ptr(http.Request{}),
				Body: &models.User{
					ID: util.Ptr(strfmt.UUID("UUID-2")),
					Company: &models.Company{
						ID: util.Ptr(strfmt.UUID("UUID-2")),
					},
					Name:   util.Ptr("更新後ユーザー"),
					Age:    util.Ptr(int64(35)),
					Gender: util.Ptr("woman"),
				},
			}
			resp := testee.Main(params)

			//
			// Assert
			//
			_, ok := resp.(*user.PostUsersOK)
			assert.True(t, ok)

			db := conn.GetEnt()
			users, err := db.User.Query().WithCompany().All(t.Context())
			assert.NoError(t, err)
			assert.Equal(t, 1, len(users))
			assert.Equal(t, "UUID-2", users[0].ID)
			assert.Equal(t, "更新後ユーザー", users[0].Name)
			assert.Equal(t, 35, users[0].Age)
			assert.Equal(t, "woman", users[0].Gender)
			assert.Equal(t, "UUID-2", users[0].Edges.Company.ID)
		})
	})
}

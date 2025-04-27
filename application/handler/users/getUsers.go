package users

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/rotisserie/eris"
	"github.com/samber/do"
	"github.com/t-kuni/go-web-api-template/domain/infrastructure/db"
	"github.com/t-kuni/go-web-api-template/ent"
	"github.com/t-kuni/go-web-api-template/ent/user"
)

type GetUsersHandler struct {
	DbConnector db.Connector
}

type EmailResponse struct {
	Address string `json:"address"`
	Type    string `json:"type"`
}

type UserResponse struct {
	ID     int             `json:"id"`
	Name   string          `json:"name"`
	Gender string          `json:"gender"`
	Age    int             `json:"age"`
	Emails []EmailResponse `json:"emails"`
}

type GetUsersResponse struct {
	Users   []UserResponse `json:"users"`
	Page    int            `json:"page"`
	MaxPage int            `json:"maxPage"`
}

func NewGetUsersHandler(i *do.Injector) (*GetUsersHandler, error) {
	return &GetUsersHandler{
		DbConnector: do.MustInvoke[db.Connector](i),
	}, nil
}

func (h GetUsersHandler) GetUsers(c echo.Context) error {
	// ページネーションの処理
	pageStr := c.QueryParam("page")
	page := 1 // デフォルトは1ページ目
	if pageStr != "" {
		var err error
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			return eris.Wrap(err, "Invalid page parameter")
		}
	}

	// 1ページあたりの件数
	const perPage = 10
	offset := (page - 1) * perPage

	var users []*ent.User
	var totalCount int
	var maxPage int

	err := h.DbConnector.Transaction(c.Request().Context(), func(tx *ent.Client) error {
		// 総件数を取得
		var err error
		totalCount, err = tx.User.Query().Count(c.Request().Context())
		if err != nil {
			return eris.Wrap(err, "Failed to count users")
		}

		// 最大ページ数を計算
		maxPage = (totalCount + perPage - 1) / perPage
		if maxPage == 0 {
			maxPage = 1
		}

		// ユーザー一覧を取得
		users, err = tx.User.Query().
			Order(ent.Asc(user.FieldID)).
			Limit(perPage).
			Offset(offset).
			WithEmails().
			All(c.Request().Context())
		if err != nil {
			return eris.Wrap(err, "Failed to query users")
		}

		return nil
	})
	if err != nil {
		return eris.Wrap(err, "")
	}

	// レスポンスの構築
	userResponses := make([]UserResponse, 0, len(users))
	for _, u := range users {
		emails := make([]EmailResponse, 0, len(u.Edges.Emails))
		for _, e := range u.Edges.Emails {
			emailType := "sub"
			if e.IsPrimary {
				emailType = "main"
			}
			emails = append(emails, EmailResponse{
				Address: e.Email,
				Type:    emailType,
			})
		}

		userResponses = append(userResponses, UserResponse{
			ID:     u.ID,
			Name:   u.Name,
			Gender: u.Gender,
			Age:    u.Age,
			Emails: emails,
		})
	}

	response := GetUsersResponse{
		Users:   userResponses,
		Page:    page,
		MaxPage: maxPage,
	}

	return c.JSON(http.StatusOK, response)
}

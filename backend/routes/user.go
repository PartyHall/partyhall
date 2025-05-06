package routes

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/partyhall/partyhall/api_errors"
	"github.com/partyhall/partyhall/config"
	"github.com/partyhall/partyhall/dal"
	"github.com/partyhall/partyhall/middlewares"
	"github.com/partyhall/partyhall/models"
	"github.com/partyhall/partyhall/services"
	"github.com/partyhall/partyhall/utils"
)

type RoutesUser struct{}

func (h RoutesUser) Register(router *gin.RouterGroup) {
	router.GET(
		"",
		middlewares.Authorized(models.ROLE_ADMIN),
		h.getCollection,
	)

	router.GET(
		":userId",
		middlewares.Authorized(models.ROLE_ADMIN),
		h.get,
	)

	router.POST(
		"",
		middlewares.Authorized(models.ROLE_ADMIN),
		h.create,
	)

	router.PUT(
		":userId",
		middlewares.Authorized(models.ROLE_ADMIN),
		h.update,
	)

	router.DELETE(
		":userId",
		middlewares.Authorized(models.ROLE_ADMIN),
		h.delete,
	)
}

func (h RoutesUser) getCollection(c *gin.Context) {
	page, offset, err := utils.ParsePageOffset(c)
	if err != nil {
		return
	}

	resp, err := dal.USERS.GetCollection(config.AMT_RESULTS_PER_PAGE, offset)
	if err != nil {
		c.Render(http.StatusInternalServerError, api_errors.DATABASE_ERROR.WithExtra(map[string]any{"err": err.Error()}))
		return
	}

	resp.Page = page

	c.JSON(http.StatusOK, resp)
}

func (h RoutesUser) get(c *gin.Context) {
	id, parseFailed := utils.ParamAsIntOrError(c, "userId")
	if parseFailed {
		api_errors.BAD_REQUEST.WithExtra(map[string]any{"err": "Invalid user id"}).Render(c.Writer)

		return
	}

	user, err := dal.USERS.Get(int(id))
	if err != nil {
		c.Render(http.StatusInternalServerError, api_errors.DATABASE_ERROR.WithExtra(map[string]any{"err": err.Error()}))
		return
	}

	if user == nil {
		c.Status(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h RoutesUser) create(c *gin.Context) {
	var req struct {
		Username string   `json:"username" binding:"required"`
		Name     string   `json:"name"`
		Password string   `json:"password" binding:"required"`
		Roles    []string `json:"roles"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		api_errors.RenderValidationErr(c, err)
		return
	}

	hashed, err := services.GetArgon().Hash(req.Password)
	if err != nil {
		c.Render(http.StatusInternalServerError, api_errors.DATABASE_ERROR.WithExtra(map[string]any{"err": err.Error()}))
		return
	}

	userDTO := models.User{
		Username: req.Username,
		Name:     req.Name,
		Password: hashed,
		Roles:    models.Roles(req.Roles),
	}

	created, err := dal.USERS.Create(userDTO)
	if err != nil {
		c.Render(http.StatusInternalServerError, api_errors.DATABASE_ERROR.WithExtra(map[string]any{"err": err.Error()}))
		return
	}

	c.JSON(http.StatusOK, created)
}

func (h RoutesUser) update(c *gin.Context) {
	id, parseFailed := utils.ParamAsIntOrError(c, "userId")
	if parseFailed {
		api_errors.BAD_REQUEST.WithExtra(map[string]any{"err": "Invalid user id"}).Render(c.Writer)

		return
	}

	var req struct {
		Username string   `json:"username" binding:"required"`
		Name     string   `json:"name"`
		Password *string  `json:"password"`
		Roles    []string `json:"roles"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		api_errors.RenderValidationErr(c, err)

		return
	}

	existing, err := dal.USERS.Get(int(id))
	if err != nil {
		c.Render(
			http.StatusInternalServerError,
			api_errors.DATABASE_ERROR.WithExtra(
				map[string]any{
					"err": err.Error(),
				},
			),
		)

		return
	}
	if existing == nil {
		c.Status(http.StatusNotFound)
		return
	}

	var passHash string
	if req.Password != nil {
		passHash, err = services.GetArgon().Hash(*req.Password)
		if err != nil {
			c.Render(
				http.StatusInternalServerError,
				api_errors.DATABASE_ERROR.WithExtra(
					map[string]any{
						"err": err.Error(),
					},
				),
			)

			return
		}
	} else {
		passHash = existing.Password
	}

	userDTO := models.User{
		Id:       int(id),
		Username: req.Username,
		Name:     req.Name,
		Password: passHash,
		Roles:    models.Roles(req.Roles),
	}

	updated, err := dal.USERS.Update(userDTO)
	if err != nil {
		c.Render(
			http.StatusInternalServerError,
			api_errors.DATABASE_ERROR.WithExtra(
				map[string]any{
					"err": err.Error(),
				},
			))
		return
	}

	c.JSON(http.StatusOK, updated)
}

func (h RoutesUser) delete(c *gin.Context) {
	id, parseFailed := utils.ParamAsIntOrError(c, "userId")
	if parseFailed {
		api_errors.BAD_REQUEST.WithExtra(map[string]any{"err": "Invalid user id"}).Render(c.Writer)

		return
	}

	err := dal.USERS.Delete(int(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.Status(http.StatusNotFound)
			return
		}

		c.Render(http.StatusInternalServerError, api_errors.DATABASE_ERROR.WithExtra(map[string]any{"err": err.Error()}))
		return
	}

	c.Status(http.StatusNoContent)
}

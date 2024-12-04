package api_errors

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

const DEFAULT_ERR_TYPE = "generic-error"
const NEXUS_ERR_TYPE = "nexus-error"

type JsonProblem struct {
	Type      string         `json:"type"`
	Title     string         `json:"title"`
	Detail    string         `json:"detail"`
	ExtraData map[string]any `json:"extra_data"`
}

func (r JsonProblem) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)

	data, _ := json.Marshal(r)
	w.Write(data)

	return nil
}

func (r JsonProblem) WriteContentType(w http.ResponseWriter) {
	w.Header()["Content-Type"] = []string{"application/problem+json"}
}

func (r JsonProblem) WithExtra(data map[string]any) JsonProblem {
	r.ExtraData = data

	return r
}

func ApiErrorWithData(c *gin.Context, code int, errType string, title string, detail string, extraData map[string]any) error {
	if len(errType) == 0 {
		errType = DEFAULT_ERR_TYPE
	}

	c.Render(code, JsonProblem{
		Type:      "https://github.com/partyhall/partyhall/" + errType,
		Title:     title,
		Detail:    detail,
		ExtraData: extraData,
	})

	return nil
}

func ApiError(c *gin.Context, code int, errType string, title string, detail string) error {
	return ApiErrorWithData(c, code, errType, title, detail, map[string]any{})
}

func RenderValidationErr(c *gin.Context, err error) {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		c.Render(http.StatusBadRequest, BAD_REQUEST.WithExtra(map[string]any{
			"err": validationErrors.Error(),
		}))
	} else {
		c.Render(http.StatusBadRequest, BAD_REQUEST.WithExtra(map[string]any{
			"err": err.Error(),
		}))
	}
}

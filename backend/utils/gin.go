package utils

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/partyhall/partyhall/api_errors"
)

func ParamAsInt(c *gin.Context, paramName string) (int64, error) {
	str := strings.TrimSpace(c.Params.ByName(paramName))
	if len(str) == 0 {
		return 0, errors.New("should be filled")
	}

	valInt, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, errors.New("should be a valid integer")
	}

	return valInt, nil
}

func ParamAsIntOrError(c *gin.Context, paramName string) (int64, bool) {
	val, err := ParamAsInt(c, paramName)
	if err != nil {
		c.Render(http.StatusBadRequest, api_errors.INVALID_PARAMETERS.WithExtra(map[string]any{
			paramName: err.Error(),
		}))

		return 0, true
	}

	return val, false
}

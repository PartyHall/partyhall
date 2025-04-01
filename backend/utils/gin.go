package utils

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/partyhall/partyhall/api_errors"
	"github.com/partyhall/partyhall/config"
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

/**
 * @returns page, offset, error
 */
func ParsePageOffset(c *gin.Context) (int, int, error) {
	offset := 0
	page := c.Query("page")

	var pageInt int = 1

	var err error
	if len(page) > 0 {
		pageInt, err = strconv.Atoi(page)
		if err != nil {
			c.Render(http.StatusBadRequest, api_errors.INVALID_PARAMETERS.WithExtra(map[string]any{
				"page": "The page should be an integer",
			}))

			return pageInt, 0, err
		}

		offset = (pageInt - 1) * config.AMT_RESULTS_PER_PAGE
	}

	return pageInt, offset, nil
}

func ParseTrilean(c *gin.Context, queryParamName string) (Thrilean, error) {
	valStr := c.Query(queryParamName)

	val := Thrilean{IsNull: true}

	if len(valStr) > 0 {
		valBool, err := strconv.ParseBool(valStr)
		if err != nil {
			c.Render(http.StatusBadRequest, api_errors.INVALID_PARAMETERS.WithExtra(map[string]any{
				"has_vocals": "Should be either left blank, true or false",
			}))

			return val, err
		}

		val = Thrilean{
			Value:  valBool,
			IsNull: false,
		}
	}

	return val, nil
}

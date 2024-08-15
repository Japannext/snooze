package apiserver

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"time"
	"strconv"

	"github.com/gin-gonic/gin"

	api "github.com/japannext/snooze/pkg/common/api/v2"
)

func extractBase64(c *gin.Context, field string) ([]byte, error) {
	str := c.Param(field)
	if str == "" {
		return []byte{}, nil
	}
	b, err := base64.URLEncoding.DecodeString(str)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("Value %s=%s is not a valid Base64 encoded data: %s", field, str, err))
		return []byte{}, err
	}

	return b, nil
}

func parseIntegerOrDefault(c *gin.Context, field string, defaultValue int) int {
	str := c.Param(field)
	if str == "" {
		return defaultValue
	}
	i, err := strconv.Atoi(str)
	if err != nil {
		return defaultValue
	}
	return i
}

func parsePerPage(c *gin.Context) int {
	return parseIntegerOrDefault(c, "per_page", 10)
}

// Parse the page index/number
func parsePageNumber(c *gin.Context) int {
	return parseIntegerOrDefault(c, "page_nb", 0)
}

func parseTimestamp(c *gin.Context, field string) time.Time {
	text := c.Param(field)
	if text == "" {
		return time.UnixMilli(0)
	}
	i, err := strconv.Atoi(text)
	if err != nil {
		return time.UnixMilli(0)
	}
	return time.UnixMilli(int64(i))
}

// Parse the sorting order (ascending/descending)
func parseAscending(c *gin.Context) bool {
	text := c.Param("asc")
	if text == "true" {
		return true
	}
	return false
}

// Parse all the pagination-related parameters
func parsePagination(c *gin.Context) api.Pagination {
	perPage := parsePerPage(c)
	pageNumber := parsePageNumber(c)
	orderBy := c.Param("order_by")
	asc := parseAscending(c)

	return api.Pagination{
		PageNumber: pageNumber,
		PerPage: perPage,
		OrderBy: orderBy,
		Ascending: asc,
	}
}

func parseTimeRange(c *gin.Context) api.TimeRange {
	start := parseTimestamp(c, "start_millis")
	end := parseTimestamp(c, "end_millis")
	return api.TimeRange{Start: start, End: end}
}

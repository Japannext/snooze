package apiserver

import (
  "fmt"
  "encoding/base64"
  "net/http"
  "strconv"

  "github.com/gin-gonic/gin"
)

const perPageDefault = 10

func extractInteger(c *gin.Context, field string, defaultValue int) (int, error) {
  str := c.Param(field)
  if str == "" {
    return defaultValue, nil
  }
  i, err := strconv.Atoi(str)
  if err != nil {
    c.String(http.StatusBadRequest, fmt.Sprintf("Value %s=%s is not a valid integer: %s", field, str, err))
    return 0, nil
  }
  return i, nil
}

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

func extractPerPage(c *gin.Context) (int, error) {
  return extractInteger(c, "per_page", 10)
}

func extractPage(c *gin.Context) (int, error) {
  return extractInteger(c, "page", 0)
}


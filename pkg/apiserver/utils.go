package apiserver

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type searchParams struct {
	Search string `form:"search"`
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

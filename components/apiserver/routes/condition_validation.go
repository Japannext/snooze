package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/japannext/snooze/internal/lang"
)

// The body of a getConditionValidation request
type getConditionValidationBody struct {
	If string `json:"if"`
}

// Validate a given condition
func getConditionValidation(c *gin.Context) {
	_, span := tracer.Start(c.Request.Context(), "getConditionValidation")
	defer span.End()

	data := getConditionValidationBody{}
	if err := c.ShouldBindJSON(&data); err != nil {
		c.String(http.StatusBadRequest, "no condition given")
		return
	}

	_, err := lang.NewCondition(data.If)
	if err != nil {
		c.String(http.StatusBadRequest, "invalid condition: %s", err)
		return
	}

	c.String(http.StatusOK, "OK")
}

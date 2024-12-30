package alertmanager

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PostableAlerts struct {
	Alerts []PostableAlert `json:",inline"`
}

type PostableAlert struct {
	StartsAt     string            `json:"startsAt"`
	EndsAt       string            `json:"endsAt"`
	GeneratorURL string            `json:"generatorURL"`
	Labels       map[string]string `json:"labels"`
	Annotations  map[string]string `json:"annotations"`
}

func postAlert(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusInternalServerError, "failed to read body: %w", err)
		return
	}

	var alerts []PostableAlert
	if err := json.Unmarshal(body, &alerts); err != nil {
		c.String(http.StatusBadRequest, "Failed to unmarshal body into Alerts: %w", err)
		return
	}

	for _, alert := range alerts {
		ctx := context.TODO()
		alertHandler(ctx, alert)
	}
}

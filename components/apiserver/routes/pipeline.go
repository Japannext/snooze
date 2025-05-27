package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/japannext/snooze/internal/pipeline"
)

// Return the list of pipeline groupings from the static configuration
func getPipelineGroupings(c *gin.Context) {
	pl, err := pipeline.Reload()
	if err != nil {
		c.String(http.StatusInternalServerError, "configuration error: %s", err)
	}
	c.JSON(http.StatusOK, pl.Groupings)
}

// Return the list of mappings from the static configuration
func getPipelineMappings(c *gin.Context) {
	pl, err := pipeline.Reload()
	if err != nil {
		c.String(http.StatusInternalServerError, "configuration error: %s", err)
	}
	c.JSON(http.StatusOK, pl.Mappings)
}

// Return the list of transforms from the static configuration
func getPipelineTransforms(c *gin.Context) {
	pl, err := pipeline.Reload()
	if err != nil {
		c.String(http.StatusInternalServerError, "configuration error: %s", err)
	}
	c.JSON(http.StatusOK, pl.Transforms)
}

// Return the list of profiles from the static configuration
func getPipelineProfiles(c *gin.Context) {
	pl, err := pipeline.Reload()
	if err != nil {
		c.String(http.StatusInternalServerError, "configuration error: %s", err)
	}
	c.JSON(http.StatusOK, pl.Profiles)
}

// Return the list of silences from the static configuration
func getPipelineSilences(c *gin.Context) {
	pl, err := pipeline.Reload()
	if err != nil {
		c.String(http.StatusInternalServerError, "configuration error: %s", err)
	}
	c.JSON(http.StatusOK, pl.Silences)
}

// Return the list of ratelimits from the static configuration
func getPipelineRatelimits(c *gin.Context) {
	pl, err := pipeline.Reload()
	if err != nil {
		c.String(http.StatusInternalServerError, "configuration error: %s", err)
	}
	c.JSON(http.StatusOK, pl.Ratelimits)
}

// Return the list of notifications from the static configuration
func getPipelineNotifications(c *gin.Context) {
	pl, err := pipeline.Reload()
	if err != nil {
		c.String(http.StatusInternalServerError, "configuration error: %s", err)
	}
	c.JSON(http.StatusOK, pl.Notifications)
}

// Return the list of default destinations from the static configuration
func getPipelineDefaultDestinations(c *gin.Context) {
	pl, err := pipeline.Reload()
	if err != nil {
		c.String(http.StatusInternalServerError, "configuration error: %s", err)
	}
	c.JSON(http.StatusOK, pl.DefaultDestinations)
}

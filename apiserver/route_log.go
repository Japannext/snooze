package api

import (
  "net/http"

  "github.com/gin-gonic/gin"
)

func searchLogV2(c *gin.Context) {

  // Extract parameters
  search := c.Param("search")
  sortBy := c.Param("sort_by")
  pp, err := extractPerPage(c); if err != nil { return }
  page, err := extractPage(c); if err != nil { return }

  ll, err := database.LogV2Search(c, search, sortBy, page, pp)
  if err != nil {
    c.String(http.StatusInternalServerError, "Error fetching logs from database: %w", err)
  }

  c.JSON(http.StatusOK, ll)
}

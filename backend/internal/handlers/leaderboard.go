package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/wcpredictions/backend/internal/repository"
)

type LeaderboardHandler struct{ lb *repository.LeaderboardRepo }

func NewLeaderboardHandler(l *repository.LeaderboardRepo) *LeaderboardHandler {
	return &LeaderboardHandler{lb: l}
}

func (h *LeaderboardHandler) Top(c *gin.Context) {
	limit := 100
	if v := c.Query("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 500 {
			limit = n
		}
	}
	out, err := h.lb.Top(c, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"leaderboard": out})
}

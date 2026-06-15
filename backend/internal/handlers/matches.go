package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/wcpredictions/backend/internal/repository"
)

type MatchHandler struct{ matches *repository.MatchRepo }

func NewMatchHandler(m *repository.MatchRepo) *MatchHandler { return &MatchHandler{matches: m} }

func (h *MatchHandler) List(c *gin.Context) {
	out, err := h.matches.List(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"matches": out})
}

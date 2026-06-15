package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/wcpredictions/backend/internal/auth"
	"github.com/wcpredictions/backend/internal/repository"
)

type PredictionHandler struct {
	preds   *repository.PredictionRepo
	matches *repository.MatchRepo
}

func NewPredictionHandler(p *repository.PredictionRepo, m *repository.MatchRepo) *PredictionHandler {
	return &PredictionHandler{preds: p, matches: m}
}

type predictionReq struct {
	PredHome int `json:"pred_home" binding:"min=0,max=20"`
	PredAway int `json:"pred_away" binding:"min=0,max=20"`
}

func (h *PredictionHandler) Upsert(c *gin.Context) {
	uid := auth.UserID(c)
	matchID, err := strconv.ParseInt(c.Param("match_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad match_id"})
		return
	}
	var req predictionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.preds.AssertNotLocked(c, matchID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "match not found"})
			return
		}
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	p, err := h.preds.Upsert(c, uid, matchID, req.PredHome, req.PredAway)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, p)
}

func (h *PredictionHandler) Mine(c *gin.Context) {
	uid := auth.UserID(c)
	out, err := h.preds.ByUser(c, uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"predictions": out})
}

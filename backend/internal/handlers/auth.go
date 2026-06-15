package handlers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/wcpredictions/backend/internal/auth"
	"github.com/wcpredictions/backend/internal/repository"
)

type AuthHandler struct {
	users *repository.UserRepo
	jwt   *auth.JWTService
}

func NewAuthHandler(users *repository.UserRepo, j *auth.JWTService) *AuthHandler {
	return &AuthHandler{users: users, jwt: j}
}

type registerReq struct {
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=8"`
	DisplayName string `json:"display_name" binding:"required,min=2,max=40"`
}

type loginReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type tokenResp struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req registerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "hash failed"})
		return
	}
	u, err := h.users.Create(c, req.Email, req.DisplayName, hash)
	if err != nil {
		// likely unique violation
		c.JSON(http.StatusConflict, gin.H{"error": "email already registered"})
		return
	}
	access, _ := h.jwt.IssueAccess(u.ID)
	refresh, _ := h.jwt.IssueRefresh(u.ID)
	c.JSON(http.StatusCreated, gin.H{"user": u, "tokens": tokenResp{access, refresh}})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	u, err := h.users.ByEmail(c, req.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	if !auth.CheckPassword(u.PasswordHash, req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	access, _ := h.jwt.IssueAccess(u.ID)
	refresh, _ := h.jwt.IssueRefresh(u.ID)
	c.JSON(http.StatusOK, gin.H{"user": u, "tokens": tokenResp{access, refresh}})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var body struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	claims, err := h.jwt.Parse(body.RefreshToken)
	if err != nil || claims.Type != "refresh" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}
	access, _ := h.jwt.IssueAccess(claims.UserID)
	c.JSON(http.StatusOK, gin.H{"access_token": access})
}

func (h *AuthHandler) Me(c *gin.Context) {
	uid := auth.UserID(c)
	u, err := h.users.ByID(c, uid)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, u)
}

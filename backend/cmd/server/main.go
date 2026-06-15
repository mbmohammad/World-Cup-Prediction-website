package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/wcpredictions/backend/internal/auth"
	"github.com/wcpredictions/backend/internal/config"
	"github.com/wcpredictions/backend/internal/db"
	"github.com/wcpredictions/backend/internal/handlers"
	"github.com/wcpredictions/backend/internal/repository"
	"github.com/wcpredictions/backend/internal/scoring"
	"github.com/wcpredictions/backend/internal/sync"
)

func main() {
	_ = godotenv.Load()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	cfg, err := config.Load()
	if err != nil {
		slog.Error("config load failed", "err", err)
		os.Exit(1)
	}

	pool, err := db.Connect(context.Background(), cfg.DatabaseURL)
	if err != nil {
		slog.Error("db connect failed", "err", err)
		os.Exit(1)
	}
	defer pool.Close()

	userRepo := repository.NewUserRepo(pool)
	matchRepo := repository.NewMatchRepo(pool)
	predRepo := repository.NewPredictionRepo(pool)
	lbRepo := repository.NewLeaderboardRepo(pool)

	jwtSvc := auth.NewJWTService(cfg.JWTSecret, cfg.JWTAccessTTL, cfg.JWTRefreshTTL)
	scorer := scoring.New(cfg.PointsExact, cfg.PointsGD, cfg.PointsOutcome)

	authH := handlers.NewAuthHandler(userRepo, jwtSvc)
	matchH := handlers.NewMatchHandler(matchRepo)
	predH := handlers.NewPredictionHandler(predRepo, matchRepo)
	lbH := handlers.NewLeaderboardHandler(lbRepo)

	if cfg.FootballDataAPIKey != "" {
		syncer := sync.NewFootballDataSyncer(cfg.FootballDataAPIKey, matchRepo, predRepo, scorer)
		go syncer.RunPeriodic(context.Background(), 10*time.Minute)
	} else {
		slog.Warn("FOOTBALL_DATA_API_KEY not set — match sync disabled (admin endpoints still work)")
	}

	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(requestLogger())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.AllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/healthz", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })
	r.GET("/readyz", func(c *gin.Context) {
		if err := pool.Ping(c); err != nil {
			c.JSON(503, gin.H{"status": "db unreachable"})
			return
		}
		c.JSON(200, gin.H{"status": "ready"})
	})

	v1 := r.Group("/api/v1")
	{
		v1.POST("/auth/register", authH.Register)
		v1.POST("/auth/login", authH.Login)
		v1.POST("/auth/refresh", authH.Refresh)

		v1.GET("/matches", matchH.List)
		v1.GET("/leaderboard", lbH.Top)

		authed := v1.Group("/")
		authed.Use(auth.RequireAuth(jwtSvc))
		{
			authed.GET("/me", authH.Me)
			authed.GET("/predictions", predH.Mine)
			authed.PUT("/predictions/:match_id", predH.Upsert)
		}
	}

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		slog.Info("server starting", "port", cfg.Port, "env", cfg.Env)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("listen failed", "err", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("shutdown failed", "err", err)
	}
}

func requestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		slog.Info("http",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"dur_ms", time.Since(start).Milliseconds(),
		)
	}
}

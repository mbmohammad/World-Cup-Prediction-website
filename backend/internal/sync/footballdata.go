package sync

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/wcpredictions/backend/internal/repository"
	"github.com/wcpredictions/backend/internal/scoring"
)

// FootballDataSyncer pulls match data from football-data.org and scores predictions.
// API docs: https://www.football-data.org/documentation/quickstart
// Free tier: 10 req/min, no historical data older than 12 months.
//
// World Cup 2026 competition code is not yet known at the time of writing — set
// FOOTBALL_DATA_COMP env (e.g. "WC") once the org publishes it.
type FootballDataSyncer struct {
	apiKey  string
	client  *http.Client
	matches *repository.MatchRepo
	preds   *repository.PredictionRepo
	scorer  *scoring.Scorer
	comp    string
}

func NewFootballDataSyncer(apiKey string, m *repository.MatchRepo, p *repository.PredictionRepo, s *scoring.Scorer) *FootballDataSyncer {
	return &FootballDataSyncer{
		apiKey:  apiKey,
		client:  &http.Client{Timeout: 15 * time.Second},
		matches: m,
		preds:   p,
		scorer:  s,
		comp:    "WC",
	}
}

type fdMatchesResp struct {
	Matches []struct {
		ID       int64  `json:"id"`
		UtcDate  string `json:"utcDate"`
		Status   string `json:"status"`
		Stage    string `json:"stage"`
		Group    string `json:"group"`
		HomeTeam struct {
			ID   int64  `json:"id"`
			Name string `json:"name"`
			Tla  string `json:"tla"`
		} `json:"homeTeam"`
		AwayTeam struct {
			ID   int64  `json:"id"`
			Name string `json:"name"`
			Tla  string `json:"tla"`
		} `json:"awayTeam"`
		Score struct {
			FullTime struct {
				Home *int `json:"home"`
				Away *int `json:"away"`
			} `json:"fullTime"`
		} `json:"score"`
	} `json:"matches"`
}

func (s *FootballDataSyncer) RunPeriodic(ctx context.Context, interval time.Duration) {
	t := time.NewTicker(interval)
	defer t.Stop()
	if err := s.SyncOnce(ctx); err != nil {
		slog.Error("sync failed", "err", err)
	}
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			if err := s.SyncOnce(ctx); err != nil {
				slog.Error("sync failed", "err", err)
			}
		}
	}
}

func (s *FootballDataSyncer) SyncOnce(ctx context.Context) error {
	url := fmt.Sprintf("https://api.football-data.org/v4/competitions/%s/matches", s.comp)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	req.Header.Set("X-Auth-Token", s.apiKey)

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("football-data status %d", resp.StatusCode)
	}

	var body fdMatchesResp
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return err
	}

	for _, m := range body.Matches {
		_ = m // TODO: map football-data team IDs to local team IDs, then UpsertExternal
	}

	return s.ScoreFinishedMatches(ctx)
}

// ScoreFinishedMatches finds matches with results, scores any predictions whose
// points_awarded is still NULL, and writes them back. Idempotent.
func (s *FootballDataSyncer) ScoreFinishedMatches(ctx context.Context) error {
	finished, err := s.matches.RecentlyFinished(ctx)
	if err != nil {
		return err
	}
	for _, m := range finished {
		if m.HomeScore == nil || m.AwayScore == nil {
			continue
		}
		preds, err := s.preds.ByMatch(ctx, m.ID)
		if err != nil {
			return err
		}
		for _, p := range preds {
			if p.PointsAwarded != nil {
				continue
			}
			pts := s.scorer.Points(p.PredHome, p.PredAway, *m.HomeScore, *m.AwayScore)
			if err := s.preds.SetPoints(ctx, p.ID, pts); err != nil {
				slog.Error("set points failed", "prediction_id", p.ID, "err", err)
			}
		}
	}
	return nil
}

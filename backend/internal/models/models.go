package models

import "time"

type User struct {
	ID           int64     `json:"id"`
	Email        string    `json:"email"`
	DisplayName  string    `json:"display_name"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

type Team struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Code    string `json:"code"`
	FlagURL string `json:"flag_url"`
	GroupID string `json:"group"`
}

type MatchStatus string

const (
	MatchScheduled MatchStatus = "scheduled"
	MatchLive      MatchStatus = "live"
	MatchFinished  MatchStatus = "finished"
)

type Match struct {
	ID         int64       `json:"id"`
	ExternalID *int64      `json:"external_id,omitempty"`
	HomeTeam   Team        `json:"home_team"`
	AwayTeam   Team        `json:"away_team"`
	KickoffUTC time.Time   `json:"kickoff_utc"`
	Stage      string      `json:"stage"`
	GroupID    string      `json:"group,omitempty"`
	HomeScore  *int        `json:"home_score,omitempty"`
	AwayScore  *int        `json:"away_score,omitempty"`
	Status     MatchStatus `json:"status"`
}

type Prediction struct {
	ID             int64     `json:"id"`
	UserID         int64     `json:"user_id"`
	MatchID        int64     `json:"match_id"`
	PredHome       int       `json:"pred_home"`
	PredAway       int       `json:"pred_away"`
	PointsAwarded  *int      `json:"points_awarded,omitempty"`
	SubmittedAt    time.Time `json:"submitted_at"`
}

type LeaderboardEntry struct {
	Rank        int    `json:"rank"`
	UserID      int64  `json:"user_id"`
	DisplayName string `json:"display_name"`
	Points      int    `json:"points"`
	Predictions int    `json:"predictions"`
}

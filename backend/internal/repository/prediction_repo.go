package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/wcpredictions/backend/internal/models"
)

type PredictionRepo struct{ pool *pgxpool.Pool }

func NewPredictionRepo(p *pgxpool.Pool) *PredictionRepo { return &PredictionRepo{pool: p} }

func (r *PredictionRepo) Upsert(ctx context.Context, userID, matchID int64, home, away int) (*models.Prediction, error) {
	p := &models.Prediction{}
	err := r.pool.QueryRow(ctx, `
		INSERT INTO predictions (user_id, match_id, pred_home, pred_away)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id, match_id) DO UPDATE SET
			pred_home = EXCLUDED.pred_home,
			pred_away = EXCLUDED.pred_away,
			submitted_at = NOW()
		RETURNING id, user_id, match_id, pred_home, pred_away, points_awarded, submitted_at
	`, userID, matchID, home, away).Scan(&p.ID, &p.UserID, &p.MatchID, &p.PredHome, &p.PredAway, &p.PointsAwarded, &p.SubmittedAt)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *PredictionRepo) ByUser(ctx context.Context, userID int64) ([]models.Prediction, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, user_id, match_id, pred_home, pred_away, points_awarded, submitted_at
		FROM predictions WHERE user_id = $1 ORDER BY submitted_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.Prediction
	for rows.Next() {
		var p models.Prediction
		if err := rows.Scan(&p.ID, &p.UserID, &p.MatchID, &p.PredHome, &p.PredAway, &p.PointsAwarded, &p.SubmittedAt); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func (r *PredictionRepo) ByMatch(ctx context.Context, matchID int64) ([]models.Prediction, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, user_id, match_id, pred_home, pred_away, points_awarded, submitted_at
		FROM predictions WHERE match_id = $1
	`, matchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.Prediction
	for rows.Next() {
		var p models.Prediction
		if err := rows.Scan(&p.ID, &p.UserID, &p.MatchID, &p.PredHome, &p.PredAway, &p.PointsAwarded, &p.SubmittedAt); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func (r *PredictionRepo) SetPoints(ctx context.Context, predictionID int64, points int) error {
	tag, err := r.pool.Exec(ctx, `UPDATE predictions SET points_awarded = $1 WHERE id = $2`, points, predictionID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// AssertNotLocked returns an error if the match has already started.
func (r *PredictionRepo) AssertNotLocked(ctx context.Context, matchID int64) error {
	var locked bool
	err := r.pool.QueryRow(ctx, `
		SELECT kickoff_utc <= NOW() FROM matches WHERE id = $1
	`, matchID).Scan(&locked)
	if errors.Is(err, pgx.ErrNoRows) {
		return ErrNotFound
	}
	if err != nil {
		return err
	}
	if locked {
		return errors.New("predictions are locked: match has started")
	}
	return nil
}

package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/wcpredictions/backend/internal/models"
)

type LeaderboardRepo struct{ pool *pgxpool.Pool }

func NewLeaderboardRepo(p *pgxpool.Pool) *LeaderboardRepo { return &LeaderboardRepo{pool: p} }

func (r *LeaderboardRepo) Top(ctx context.Context, limit int) ([]models.LeaderboardEntry, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT
			RANK() OVER (ORDER BY COALESCE(SUM(p.points_awarded), 0) DESC) AS rank,
			u.id, u.display_name,
			COALESCE(SUM(p.points_awarded), 0)::int AS points,
			COUNT(p.id)::int AS predictions
		FROM users u
		LEFT JOIN predictions p ON p.user_id = u.id AND p.points_awarded IS NOT NULL
		GROUP BY u.id, u.display_name
		ORDER BY points DESC, u.id ASC
		LIMIT $1
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.LeaderboardEntry
	for rows.Next() {
		var e models.LeaderboardEntry
		if err := rows.Scan(&e.Rank, &e.UserID, &e.DisplayName, &e.Points, &e.Predictions); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}

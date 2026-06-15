package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/wcpredictions/backend/internal/models"
)

type MatchRepo struct{ pool *pgxpool.Pool }

func NewMatchRepo(p *pgxpool.Pool) *MatchRepo { return &MatchRepo{pool: p} }

func (r *MatchRepo) List(ctx context.Context) ([]models.Match, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT
			m.id, m.external_id, m.kickoff_utc, m.stage, m.group_id,
			m.home_score, m.away_score, m.status,
			ht.id, ht.name, ht.code, ht.flag_url, ht.group_id,
			at.id, at.name, at.code, at.flag_url, at.group_id
		FROM matches m
		JOIN teams ht ON ht.id = m.home_team_id
		JOIN teams at ON at.id = m.away_team_id
		ORDER BY m.kickoff_utc ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.Match
	for rows.Next() {
		var m models.Match
		var status string
		if err := rows.Scan(
			&m.ID, &m.ExternalID, &m.KickoffUTC, &m.Stage, &m.GroupID,
			&m.HomeScore, &m.AwayScore, &status,
			&m.HomeTeam.ID, &m.HomeTeam.Name, &m.HomeTeam.Code, &m.HomeTeam.FlagURL, &m.HomeTeam.GroupID,
			&m.AwayTeam.ID, &m.AwayTeam.Name, &m.AwayTeam.Code, &m.AwayTeam.FlagURL, &m.AwayTeam.GroupID,
		); err != nil {
			return nil, err
		}
		m.Status = models.MatchStatus(status)
		out = append(out, m)
	}
	return out, rows.Err()
}

func (r *MatchRepo) ByID(ctx context.Context, id int64) (*models.Match, error) {
	m := &models.Match{}
	var status string
	err := r.pool.QueryRow(ctx, `
		SELECT
			m.id, m.external_id, m.kickoff_utc, m.stage, m.group_id,
			m.home_score, m.away_score, m.status,
			ht.id, ht.name, ht.code, ht.flag_url, ht.group_id,
			at.id, at.name, at.code, at.flag_url, at.group_id
		FROM matches m
		JOIN teams ht ON ht.id = m.home_team_id
		JOIN teams at ON at.id = m.away_team_id
		WHERE m.id = $1
	`, id).Scan(
		&m.ID, &m.ExternalID, &m.KickoffUTC, &m.Stage, &m.GroupID,
		&m.HomeScore, &m.AwayScore, &status,
		&m.HomeTeam.ID, &m.HomeTeam.Name, &m.HomeTeam.Code, &m.HomeTeam.FlagURL, &m.HomeTeam.GroupID,
		&m.AwayTeam.ID, &m.AwayTeam.Name, &m.AwayTeam.Code, &m.AwayTeam.FlagURL, &m.AwayTeam.GroupID,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	m.Status = models.MatchStatus(status)
	return m, nil
}

func (r *MatchRepo) UpsertExternal(ctx context.Context, externalID int64, homeTeamID, awayTeamID int64, kickoff_utc, stage, groupID string, homeScore, awayScore *int, status string) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO matches (external_id, home_team_id, away_team_id, kickoff_utc, stage, group_id, home_score, away_score, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (external_id) DO UPDATE SET
			kickoff_utc = EXCLUDED.kickoff_utc,
			home_score = EXCLUDED.home_score,
			away_score = EXCLUDED.away_score,
			status = EXCLUDED.status
	`, externalID, homeTeamID, awayTeamID, kickoff_utc, stage, groupID, homeScore, awayScore, status)
	return err
}

func (r *MatchRepo) RecentlyFinished(ctx context.Context) ([]models.Match, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, home_score, away_score
		FROM matches
		WHERE status = 'finished' AND home_score IS NOT NULL AND away_score IS NOT NULL
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []models.Match
	for rows.Next() {
		var m models.Match
		if err := rows.Scan(&m.ID, &m.HomeScore, &m.AwayScore); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	return out, rows.Err()
}

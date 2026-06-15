-- Sample data so you can exercise the app without the football-data API.
-- Run with: docker compose exec -T postgres psql -U wcuser -d wcdb < backend/migrations/seed.sql

INSERT INTO teams (name, code, group_id) VALUES
    ('United States',  'USA', 'A'),
    ('Mexico',         'MEX', 'A'),
    ('Canada',         'CAN', 'B'),
    ('Argentina',      'ARG', 'B'),
    ('Brazil',         'BRA', 'C'),
    ('France',         'FRA', 'C')
ON CONFLICT DO NOTHING;

-- A finished match (so you can see points awarded after predicting),
-- a live one, and two upcoming ones (so predictions are still open).
INSERT INTO matches (home_team_id, away_team_id, kickoff_utc, stage, group_id, home_score, away_score, status)
SELECT
    (SELECT id FROM teams WHERE code = 'USA'),
    (SELECT id FROM teams WHERE code = 'MEX'),
    NOW() - INTERVAL '2 hours', 'GROUP_STAGE', 'A', 2, 1, 'finished';

INSERT INTO matches (home_team_id, away_team_id, kickoff_utc, stage, group_id, status)
SELECT
    (SELECT id FROM teams WHERE code = 'CAN'),
    (SELECT id FROM teams WHERE code = 'ARG'),
    NOW() + INTERVAL '1 day', 'GROUP_STAGE', 'B', 'scheduled';

INSERT INTO matches (home_team_id, away_team_id, kickoff_utc, stage, group_id, status)
SELECT
    (SELECT id FROM teams WHERE code = 'BRA'),
    (SELECT id FROM teams WHERE code = 'FRA'),
    NOW() + INTERVAL '2 days', 'GROUP_STAGE', 'C', 'scheduled';

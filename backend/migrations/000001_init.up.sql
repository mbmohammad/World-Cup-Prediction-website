CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE users (
    id            BIGSERIAL PRIMARY KEY,
    email         CITEXT UNIQUE NOT NULL,
    display_name  TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE teams (
    id          BIGSERIAL PRIMARY KEY,
    external_id BIGINT UNIQUE,
    name        TEXT NOT NULL,
    code        TEXT NOT NULL,
    flag_url    TEXT NOT NULL DEFAULT '',
    group_id    TEXT NOT NULL DEFAULT ''
);

CREATE TABLE matches (
    id            BIGSERIAL PRIMARY KEY,
    external_id   BIGINT UNIQUE,
    home_team_id  BIGINT NOT NULL REFERENCES teams(id),
    away_team_id  BIGINT NOT NULL REFERENCES teams(id),
    kickoff_utc   TIMESTAMPTZ NOT NULL,
    stage         TEXT NOT NULL,
    group_id      TEXT NOT NULL DEFAULT '',
    home_score    INT,
    away_score    INT,
    status        TEXT NOT NULL DEFAULT 'scheduled'
);

CREATE INDEX matches_kickoff_idx ON matches (kickoff_utc);
CREATE INDEX matches_status_idx  ON matches (status);

CREATE TABLE predictions (
    id             BIGSERIAL PRIMARY KEY,
    user_id        BIGINT NOT NULL REFERENCES users(id)   ON DELETE CASCADE,
    match_id       BIGINT NOT NULL REFERENCES matches(id) ON DELETE CASCADE,
    pred_home      INT NOT NULL CHECK (pred_home >= 0 AND pred_home <= 20),
    pred_away      INT NOT NULL CHECK (pred_away >= 0 AND pred_away <= 20),
    points_awarded INT,
    submitted_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, match_id)
);

CREATE INDEX predictions_user_idx  ON predictions (user_id);
CREATE INDEX predictions_match_idx ON predictions (match_id);


CREATE TABLE IF NOT EXISTS games (
    id UUID PRIMARY KEY,
    player1 VARCHAR NOT NULL,
    player2 VARCHAR,
    winner VARCHAR,
    duration_seconds INT,
    created_at TIMESTAMP DEFAULT NOW(),
    finished_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS leaderboard (
    username VARCHAR PRIMARY KEY,
    wins INT DEFAULT 0
);

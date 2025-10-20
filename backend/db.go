
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/jackc/pgx/v5"
)

var dbConn *pgx.Conn

func InitDB() error {
    dsn := os.Getenv("DATABASE_URL")
    if dsn == "" {
        dsn = "postgres://fourinarow:fourinarow@localhost:5432/fourinarow"
    }
    conn, err := pgx.Connect(context.Background(), dsn)
    if err != nil {
        return err
    }
    dbConn = conn
    // create schema if needed
    _, err = conn.Exec(context.Background(), string(getSchema()))
    if err != nil {
        log.Println("schema exec err:", err)
    }
    return nil
}

func getSchema() []byte {
    return []byte(`
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
`)
}

func SaveCompletedGame(g *Game) {
    if dbConn == nil {
        fmt.Println("DB not configured; skipping SaveCompletedGame")
        return
    }
    // insert game and update leaderboard
    _, err := dbConn.Exec(context.Background(),
        `INSERT INTO games (id, player1, player2, winner, duration_seconds, finished_at) VALUES ($1,$2,$3,$4,$5,NOW())`,
        g.ID, g.Player1, g.Player2, g.Winner, int(0),
    )
    if err != nil {
        log.Println("Insert game err:", err)
    }
    if g.Winner != "" {
        _, err := dbConn.Exec(context.Background(),
            `INSERT INTO leaderboard (username, wins) VALUES ($1,1) ON CONFLICT (username) DO UPDATE SET wins = leaderboard.wins + 1`,
            g.Winner,
        )
        if err != nil {
            log.Println("leaderboard update err:", err)
        }
    }
}

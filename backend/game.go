
package main

import (
    "time"

    "github.com/google/uuid"
)

type GameStatus string

const (
    StatusOngoing GameStatus = "ongoing"
    StatusFinished GameStatus = "finished"
)

type Game struct {
    ID        string
    Player1   string
    Player2   string
    Board     [6][7]int // 0 empty, 1 player1, 2 player2
    NextTurn  string
    StartedAt time.Time
    Status    GameStatus
    Winner    string
}

func NewGame(p1, p2 string) *Game {
    return &Game{
        ID:       uuid.New().String(),
        Player1:  p1,
        Player2:  p2,
        NextTurn: p1,
        StartedAt: time.Now(),
        Status:   StatusOngoing,
    }
}

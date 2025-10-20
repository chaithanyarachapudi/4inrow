package main

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

const (
	StatusWaiting  = "waiting"
	StatusOngoing  = "ongoing"
	StatusFinished = "finished"
)

type Game struct {
	ID         string
	Player1    string
	Player2    string
	Board      [][]int
	NextTurn   string
	Status     string
	Winner     string
	RematchReq map[string]bool
}

func NewGame(p1, p2 string) *Game {
	board := make([][]int, 6)
	for i := range board {
		board[i] = make([]int, 7)
	}
	return &Game{
		ID:         generateGameID(),
		Player1:    p1,
		Player2:    p2,
		Board:      board,
		NextTurn:   p1,
		Status:     StatusOngoing,
		RematchReq: make(map[string]bool),
	}
}

type GameManager struct {
	mu      sync.Mutex
	queue   []string
	clients map[string]*Client
	games   map[string]*Game
}

func NewGameManager() *GameManager {
	return &GameManager{
		clients: make(map[string]*Client),
		games:   make(map[string]*Game),
	}
}

// Matchmaking
func (gm *GameManager) JoinQueue(username string, c *Client) {
	gm.mu.Lock()
	gm.clients[username] = c
	gm.queue = append(gm.queue, username)
	gm.mu.Unlock()
	gm.tryMatch()
	// Auto match with BOT after 10s if not matched
	go func() {
		time.Sleep(10 * time.Second)
		gm.mu.Lock()
		defer gm.mu.Unlock()
		for i, u := range gm.queue {
			if u == username {
				gm.queue = append(gm.queue[:i], gm.queue[i+1:]...)
				game := NewGame(username, "BOT_AI")
				gm.games[game.ID] = game
				c.gameId = game.ID
				c.send <- map[string]interface{}{"type": "matched", "gameId": game.ID, "opponent": "BOT_AI", "you": 1}
				gm.broadcastGameState(game)
				return
			}
		}
	}()
}

func (gm *GameManager) tryMatch() {
	gm.mu.Lock()
	defer gm.mu.Unlock()
	if len(gm.queue) >= 2 {
		p1 := gm.queue[0]
		p2 := gm.queue[1]
		gm.queue = gm.queue[2:]
		game := NewGame(p1, p2)
		gm.games[game.ID] = game

		if c1, ok := gm.clients[p1]; ok {
			c1.gameId = game.ID
			c1.send <- map[string]interface{}{"type": "matched", "gameId": game.ID, "opponent": p2, "you": 1}
		}
		if c2, ok := gm.clients[p2]; ok {
			c2.gameId = game.ID
			c2.send <- map[string]interface{}{"type": "matched", "gameId": game.ID, "opponent": p1, "you": 2}
		}

		gm.broadcastGameState(game)
	}
}

// Handle player move
func (gm *GameManager) HandleMove(gameId, username string, col int) {
	gm.mu.Lock()
	game, ok := gm.games[gameId]
	gm.mu.Unlock()
	if !ok || game.Status != StatusOngoing {
		return
	}

	var pid int
	if username == game.Player1 {
		pid = 1
	} else {
		pid = 2
	}

	if game.NextTurn != username {
		return
	}

	if !dropPiece(game.Board, col, pid) {
		if c, ok := gm.clients[username]; ok {
			c.send <- map[string]interface{}{"type": "error", "message": "invalid move"}
		}
		return
	}

	// Switch turn
	if game.NextTurn == game.Player1 {
		game.NextTurn = game.Player2
	} else {
		game.NextTurn = game.Player1
	}

	gm.broadcastGameState(game)

	// Check winner
	if winnerPid, _ := checkWin(game.Board); winnerPid != 0 {
		game.Status = StatusFinished
		if winnerPid == 1 {
			game.Winner = game.Player1
		} else {
			game.Winner = game.Player2
		}
		for _, uname := range []string{game.Player1, game.Player2} {
			if c, ok := gm.clients[uname]; ok {
				c.send <- map[string]interface{}{"type": "result", "result": "win", "winner": game.Winner}
			}
		}
		return
	}

	if boardFull(game.Board) {
		game.Status = StatusFinished
		for _, uname := range []string{game.Player1, game.Player2} {
			if c, ok := gm.clients[uname]; ok {
				c.send <- map[string]interface{}{"type": "result", "result": "draw"}
			}
		}
		return
	}

	// Bot move
	if game.Player2 == "BOT_AI" && username != "BOT_AI" {
		go func(game *Game) {
			time.Sleep(400 * time.Millisecond)
			col := BotChooseMove(game.Board, 2, 1)
			gm.HandleMove(game.ID, "BOT_AI", col)
		}(game)
	}
}

// Rematch request
func (gm *GameManager) HandleRematchRequest(gameId, username string) {
	gm.mu.Lock()
	defer gm.mu.Unlock()

	game, ok := gm.games[gameId]
	if !ok {
		if c, ok := gm.clients[username]; ok {
			c.send <- map[string]interface{}{"type": "error", "message": "game not found"}
		}
		return
	}

	// If BOT, start immediately
	if game.Player2 == "BOT_AI" {
		resetBoard(game)
		gm.sendRematchStart(game)
		return
	}

	game.RematchReq[username] = true

	// Notify opponent
	opponent := game.Player1
	if username == game.Player1 {
		opponent = game.Player2
	}
	if c, ok := gm.clients[opponent]; ok {
		c.send <- map[string]interface{}{"type": "info", "message": username + " requested a rematch. Click Rematch to accept."}
	}

	// If both requested, start rematch
	if game.RematchReq[game.Player1] && game.RematchReq[game.Player2] {
		resetBoard(game)
		gm.sendRematchStart(game)
	}
}

// Reset game board
func resetBoard(game *Game) {
	game.Board = make([][]int, 6)
	for i := range game.Board {
		game.Board[i] = make([]int, 7)
	}
	game.Status = StatusOngoing
	game.NextTurn = game.Player1
	game.RematchReq = make(map[string]bool)
}

// Send rematch_start to both clients
func (gm *GameManager) sendRematchStart(game *Game) {
	for _, u := range []string{game.Player1, game.Player2} {
		if c, ok := gm.clients[u]; ok {
			c.send <- map[string]interface{}{"type": "rematch_start", "message": "Rematch started!"}
		}
	}
	gm.broadcastGameState(game)
}

// Broadcast game state
func (gm *GameManager) broadcastGameState(game *Game) {
	for _, u := range []string{game.Player1, game.Player2} {
		if c, ok := gm.clients[u]; ok {
			c.send <- map[string]interface{}{"type": "state", "board": game.Board, "nextTurn": game.NextTurn}
		}
	}
}

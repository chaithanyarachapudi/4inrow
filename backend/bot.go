
package main

// A simple strategic bot:
// 1. If bot can win in one move, play that move.
// 2. If opponent can win in one move, block it.
// 3. Otherwise, prefer center column, then near-center.

func BotChooseMove(g Game, botPid, oppPid int) int {
    // check winning move for bot
    for col:=0; col<7; col++ {
        copyBoard := g.Board
        if dropPiece(&copyBoard, col, botPid) {
            if win, _ := checkWin(copyBoard); win == botPid {
                return col
            }
        }
    }
    // block opponent
    for col:=0; col<7; col++ {
        copyBoard := g.Board
        if dropPiece(&copyBoard, col, oppPid) {
            if win, _ := checkWin(copyBoard); win == oppPid {
                return col
            }
        }
    }
    // prefer center columns order
    order := []int{3,2,4,1,5,0,6}
    for _, c := range order {
        copyBoard := g.Board
        if dropPiece(&copyBoard, c, botPid) {
            return c
        }
    }
    return 0
}

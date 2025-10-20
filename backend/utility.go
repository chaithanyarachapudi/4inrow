
package main

// dropPiece: drop in column col (0-based), returns true if success
func dropPiece(board *[6][7]int, col, pid int) bool {
    if col < 0 || col > 6 { return false }
    for r := 5; r >= 0; r-- {
        if board[r][col] == 0 {
            board[r][col] = pid
            return true
        }
    }
    return false
}

// simple checkWin scanning for 4 in row. returns winning pid and list of cells (r,c)
func checkWin(board [6][7]int) (int, [][2]int) {
    // directions: right, down, diag-down-right, diag-down-left
    dirs := [][2]int{{0,1},{1,0},{1,1},{1,-1}}
    for r:=0; r<6; r++ {
        for c:=0; c<7; c++ {
            pid := board[r][c]
            if pid == 0 { continue }
            for _, d := range dirs {
                cells := [][2]int{}
                rr, cc := r, c
                for k:=0; k<4; k++ {
                    if rr<0 || rr>=6 || cc<0 || cc>=7 { break }
                    if board[rr][cc] != pid { break }
                    cells = append(cells, [2]int{rr,cc})
                    rr += d[0]
                    cc += d[1]
                }
                if len(cells) == 4 {
                    return pid, cells
                }
            }
        }
    }
    return 0, nil
}

func boardFull(board [6][7]int) bool {
    for c:=0; c<7; c++ {
        if board[0][c] == 0 { return false }
    }
    return true
}

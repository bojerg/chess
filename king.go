package main

import "github.com/hajimehoshi/ebiten/v2"

type King struct {
	Piece
}

// Moves returns a slice of all valid moves for given Piece. Each valid move in the slice is stored
// in an array with a length of two-- Row and Col.
func (p *King) Moves(g Game) [][2]int {
	moves := make([][2]int, 0)

	possibleMoves := [8][2]int{
		{p.row, p.col + 1},
		{p.row, p.col - 1},
		{p.row + 1, p.col},
		{p.row - 1, p.col},
		{p.row - 1, p.col + 1},
		{p.row - 1, p.col - 1},
		{p.row + 1, p.col + 1},
		{p.row + 1, p.col - 1},
	}

	for _, move := range possibleMoves {
		if IsInBounds(move[0], move[1]) {
			otherPiece := GetPieceOnSquare(move[0], move[1], g.pieces)
			if otherPiece == nil || otherPiece.White() != p.white {
				moves = append(moves, [2]int{move[0], move[1]})
			}
		}
	}

	//evaluate castle moves
	//here, we only worry if we are in check, if previous moves invalidated castling (rook/king cant have moved),
	//and if there are any pieces blocking the move. Additional constraints evaluated by MakeMoveIfLegal
	if !g.inCheck {
		if p.white {
			if g.whiteCastles[0] {
				//rook at 7,0
				//check all positions between for pieces
				legal := true
				for checkCol := 1; checkCol < 4; checkCol++ {
					if GetPieceOnSquare(p.row, checkCol, g.pieces) != nil {
						legal = false
						break
					}
				}
				//add appropriate king move to slice
				if legal {
					moves = append(moves, [2]int{p.row, p.col - 2})
				}
			}
			if g.whiteCastles[1] {
				//rook at 7,7
				//check all positions between for pieces
				legal := true
				for checkCol := 6; checkCol > 4; checkCol-- {
					if GetPieceOnSquare(p.row, checkCol, g.pieces) != nil {
						legal = false
						break
					}
				}
				//add appropriate king move to slice
				if legal {
					moves = append(moves, [2]int{p.row, p.col + 2})
				}
			}
		} else {
			if g.blackCastles[0] {
				//rook at 0,0
				//check all positions between for pieces
				legal := true
				for checkCol := 1; checkCol < 4; checkCol++ {
					if GetPieceOnSquare(p.row, checkCol, g.pieces) != nil {
						legal = false
						break
					}
				}
				//add appropriate king move to slice
				if legal {
					moves = append(moves, [2]int{p.row, p.col - 2})
				}
			}
			if g.blackCastles[1] {
				//rook at 0,7
				//check all positions between for pieces
				legal := true
				for checkCol := 6; checkCol > 4; checkCol-- {
					if GetPieceOnSquare(p.row, checkCol, g.pieces) != nil {
						legal = false
						break
					}
				}
				//add appropriate king move to slice
				if legal {
					moves = append(moves, [2]int{p.row, p.col + 2})
				}
			}
		}
	}

	return moves
}

func (p *King) Name() string {
	if p.white {
		return "White king"
	} else {
		return "Black king"
	}
}

func (p *King) Image() *ebiten.Image {
	filepathStr := "images/"
	if p.white {
		filepathStr += "whiteKing.png"
	} else {
		filepathStr += "blackKing.png"
	}
	// Reusing the Piece GetImage for filesystem functionality
	return GetImage(filepathStr)
}

func (p *King) Col() int {
	return p.col
}

func (p *King) Row() int {
	return p.row
}

func (p *King) SetCol(c int) {
	p.col = c
}

func (p *King) SetRow(r int) {
	p.row = r
}

func (p *King) White() bool {
	return p.white
}

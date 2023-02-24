package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Pawn struct {
	Piece
}

// GetMoves returns a slice of all possible moves (may include invalid moves) for given Piece. Each valid move in the slice is stored
// in an array with a length of two-- Row and Col.
func (p *Pawn) GetMoves(pieces [32]ChessPiece) [][2]int {
	moves := make([][2]int, 0)

	// TODO add en passant, or however you spell that move
	if p.white {
		// white pawn on starting position, so could move forward one or two
		if p.row == 6 {
			if GetPieceOnSquare(4, p.col, pieces) == nil {
				moves = append(moves, [2]int{4, p.col})
			}
			if GetPieceOnSquare(5, p.col, pieces) == nil {
				moves = append(moves, [2]int{5, p.col})
			}

		} else {
			// white pawn not on starting position
			if GetPieceOnSquare(p.row-1, p.col, pieces) == nil {
				moves = append(moves, [2]int{p.row - 1, p.col})
			}
		}

		//now checking for takes
		otherPiece := GetPieceOnSquare(p.row-1, p.col+1, pieces)
		if otherPiece != nil && otherPiece.White() != p.white {
			moves = append(moves, [2]int{p.row - 1, p.col + 1})
		}

		otherPiece = GetPieceOnSquare(p.row-1, p.col-1, pieces)
		if otherPiece != nil && otherPiece.White() != p.white {
			moves = append(moves, [2]int{p.row - 1, p.col - 1})
		}

	} else {
		// black pawn on starting position, so could move forward one or two
		if p.row == 1 {
			if GetPieceOnSquare(3, p.col, pieces) == nil {
				moves = append(moves, [2]int{3, p.col})
			}
			if GetPieceOnSquare(2, p.col, pieces) == nil {
				moves = append(moves, [2]int{2, p.col})
			}
		} else {
			// black pawn not on starting position
			if GetPieceOnSquare(p.row+1, p.col, pieces) == nil {
				moves = append(moves, [2]int{p.row + 1, p.col})
			}
		}

		//now checking for takes
		otherPiece := GetPieceOnSquare(p.row+1, p.col+1, pieces)
		if otherPiece != nil && otherPiece.White() != p.white {
			moves = append(moves, [2]int{p.row + 1, p.col + 1})
		}

		otherPiece = GetPieceOnSquare(p.row+1, p.col-1, pieces)
		if otherPiece != nil && otherPiece.White() != p.white {
			moves = append(moves, [2]int{p.row + 1, p.col - 1})
		}
	}

	//should loop through and remove invalid indices (outside bounds of the board)
	for i, move := range moves {
		if move[0] > 8 || move[0] < 0 || move[1] > 8 || move[1] < 0 {
			//https://www.geeksforgeeks.org/delete-elements-in-a-slice-in-golang/
			//removing the move from the slice
			moves = append(moves[:i], moves[i+1:]...)
		}
	}

	return moves
}

// GetName primarily intended for debugging
func (p *Pawn) GetName() string {
	if p.white {
		return "White pawn"
	} else {
		return "Black pawn"
	}
}

func (p *Pawn) GetImage() *ebiten.Image {
	filepathStr := "images/"
	if p.white {
		filepathStr += "whitePawn.png"
	} else {
		filepathStr += "blackPawn.png"
	}
	// Reusing the piece.go GetImage for filesystem functionality
	return GetImage(filepathStr)
}

func (p *Pawn) GetCol() int {
	return p.col
}

func (p *Pawn) GetRow() int {
	return p.row
}

func (p *Pawn) SetCol(c int) {
	p.col = c
}

func (p *Pawn) SetRow(r int) {
	p.row = r
}

func (p *Pawn) White() bool {
	return p.white
}

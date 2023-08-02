package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Pawn struct {
	Piece
}

// Moves returns a slice of all possible moves (may include invalid moves) for given Piece. Each valid move in the slice is stored
// in an array with a length of two-- Row and Col.
func (p *Pawn) Moves(g Game) [][2]int {
	moves := make([][2]int, 0)

	if p.white {
		// white pawn on starting position, so could move forward one or two
		// not checking bounds because there's no way to move out of bounds with hardcoded moves (I hope!)
		if p.row == 6 {
			if GetPieceOnSquare(5, p.col, g.pieces) == nil {
				moves = append(moves, [2]int{5, p.col})

				//nested check to ensure we don't jump over a piece
				if GetPieceOnSquare(4, p.col, g.pieces) == nil {
					moves = append(moves, [2]int{4, p.col})
				}
			}
		} else {
			// white pawn not on starting position
			// now we check bounds
			if GetPieceOnSquare(p.row-1, p.col, g.pieces) == nil && IsInBounds(p.row-1, p.col) {
				moves = append(moves, [2]int{p.row - 1, p.col})
			}
		}

		//now checking for takes
		if IsInBounds(p.row-1, p.col+1) {
			otherPiece := GetPieceOnSquare(p.row-1, p.col+1, g.pieces)
			if otherPiece != nil && otherPiece.White() != p.white {
				moves = append(moves, [2]int{p.row - 1, p.col + 1})
			}
		}

		if IsInBounds(p.row-1, p.col-1) {
			otherPiece := GetPieceOnSquare(p.row-1, p.col-1, g.pieces)
			if otherPiece != nil && otherPiece.White() != p.white {
				moves = append(moves, [2]int{p.row - 1, p.col - 1})
			}
		}

		//don't forget to check the for an en passant too
		if p.row == 3 {

			if IsInBounds(p.row, p.col+1) {
				otherPiece := GetPieceOnSquare(p.row, p.col+1, g.pieces)
				if otherPiece != nil && otherPiece.White() != p.white {
					if otherPiece.Row() == g.enPassantLocation[0] && otherPiece.Col() == g.enPassantLocation[1] {
						moves = append(moves, [2]int{p.row - 1, p.col + 1})
					}
				}
			}

			if IsInBounds(p.row, p.col-1) {
				otherPiece := GetPieceOnSquare(p.row, p.col-1, g.pieces)
				if otherPiece != nil && otherPiece.White() != p.white {
					if otherPiece.Row() == g.enPassantLocation[0] && otherPiece.Col() == g.enPassantLocation[1] {
						moves = append(moves, [2]int{p.row - 1, p.col - 1})
					}
				}
			}
		}

	} else {
		// black pawn on starting position, so could move forward one or two
		if p.row == 1 {
			if GetPieceOnSquare(2, p.col, g.pieces) == nil {
				moves = append(moves, [2]int{2, p.col})

				if GetPieceOnSquare(3, p.col, g.pieces) == nil {
					moves = append(moves, [2]int{3, p.col})
				}
			}
		} else {
			// black pawn not on starting position
			if GetPieceOnSquare(p.row+1, p.col, g.pieces) == nil && IsInBounds(p.row+1, p.col) {
				moves = append(moves, [2]int{p.row + 1, p.col})
			}
		}

		//now checking for takes
		if IsInBounds(p.row+1, p.col+1) {
			otherPiece1 := GetPieceOnSquare(p.row+1, p.col+1, g.pieces)
			if otherPiece1 != nil && otherPiece1.White() != p.white {
				moves = append(moves, [2]int{p.row + 1, p.col + 1})
			}
		}

		if IsInBounds(p.row+1, p.col-1) {
			otherPiece2 := GetPieceOnSquare(p.row+1, p.col-1, g.pieces)
			if otherPiece2 != nil && otherPiece2.White() != p.white {
				moves = append(moves, [2]int{p.row + 1, p.col - 1})
			}
		}

		//don't forget to check the for an en passant too
		if p.row == 4 {

			if IsInBounds(p.row, p.col+1) {
				otherPiece := GetPieceOnSquare(p.row, p.col+1, g.pieces)
				if otherPiece != nil && otherPiece.White() != p.white {
					if otherPiece.Row() == g.enPassantLocation[0] && otherPiece.Col() == g.enPassantLocation[1] {
						moves = append(moves, [2]int{p.row + 1, p.col + 1})
					}
				}
			}

			if IsInBounds(p.row, p.col-1) {
				otherPiece := GetPieceOnSquare(p.row, p.col-1, g.pieces)
				if otherPiece != nil && otherPiece.White() != p.white {
					if otherPiece.Row() == g.enPassantLocation[0] && otherPiece.Col() == g.enPassantLocation[1] {
						moves = append(moves, [2]int{p.row + 1, p.col - 1})
					}
				}
			}
		}

	}

	return moves
}

func (p *Pawn) Name() string {
	if p.white {
		return "White pawn"
	} else {
		return "Black pawn"
	}
}

func (p *Pawn) Image() *ebiten.Image {
	filepathStr := "images/"
	if p.white {
		filepathStr += "whitePawn.png"
	} else {
		filepathStr += "blackPawn.png"
	}
	// Reusing the piece.go GetImage for filesystem functionality
	return GetImage(filepathStr)
}

func (p *Pawn) Col() int {
	return p.col
}

func (p *Pawn) Row() int {
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

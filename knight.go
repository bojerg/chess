package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Knight struct {
	Piece
}

// GetMoves returns a slice of all valid moves for given Piece. Each valid move in the slice is stored
// in an array with a length of two-- Row and Col.
func (p *Knight) GetMoves(pieces [32]ChessPiece) [][2]int {
	moves := make([][2]int, 0)

	possibleMoves := [8][2]int{
		{p.row + 2, p.col + 1},
		{p.row + 2, p.col - 1},
		{p.row + 1, p.col + 2},
		{p.row + 1, p.col - 2},
		{p.row - 1, p.col + 2},
		{p.row - 1, p.col - 2},
		{p.row - 2, p.col + 1},
		{p.row - 2, p.col - 1},
	}

	for _, move := range possibleMoves {
		if IsInBounds(move[0], move[1]) {
			otherPiece := GetPieceOnSquare(move[0], move[1], pieces)
			if otherPiece == nil || otherPiece.White() != p.white {
				moves = append(moves, [2]int{move[0], move[1]})
			}
		}
	}

	return moves
}

func (p *Knight) IsKing() bool {
	return false
}

// GetName primarily intended for debugging
func (p *Knight) GetName() string {
	if p.white {
		return "White knight"
	} else {
		return "Black knight"
	}
}

func (p *Knight) GetImage() *ebiten.Image {
	filepathStr := "images/"
	if p.white {
		filepathStr += "whiteKnight.png"
	} else {
		filepathStr += "blackKnight.png"
	}
	// Reusing the Piece GetImage for filesystem functionality
	return GetImage(filepathStr)
}

func (p *Knight) GetCol() int {
	return p.col
}

func (p *Knight) GetRow() int {
	return p.row
}

func (p *Knight) SetCol(c int) {
	p.col = c
}

func (p *Knight) SetRow(r int) {
	p.row = r
}

func (p *Knight) White() bool {
	return p.white
}

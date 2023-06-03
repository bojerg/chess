package main

import "github.com/hajimehoshi/ebiten/v2"

type King struct {
	Piece
}

// GetMoves returns a slice of all valid moves for given Piece. Each valid move in the slice is stored
// in an array with a length of two-- Row and Col.
// TODO castling
func (p *King) GetMoves(pieces [32]ChessPiece) [][2]int {
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
			otherPiece := GetPieceOnSquare(move[0], move[1], pieces)
			if otherPiece == nil || otherPiece.White() != p.white {
				moves = append(moves, [2]int{move[0], move[1]})
			}
		}
	}

	return moves
}

func (p *King) IsKing() bool {
	return true
}

// GetName primarily intended for debugging
func (p *King) GetName() string {
	if p.white {
		return "White pawn"
	} else {
		return "Black pawn"
	}
}

func (p *King) GetImage() *ebiten.Image {
	filepathStr := "images/"
	if p.white {
		filepathStr += "whiteKing.png"
	} else {
		filepathStr += "blackKing.png"
	}
	// Reusing the Piece GetImage for filesystem functionality
	return GetImage(filepathStr)
}

func (p *King) GetCol() int {
	return p.col
}

func (p *King) GetRow() int {
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

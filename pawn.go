package main

import "github.com/hajimehoshi/ebiten/v2"

type Pawn struct {
	Piece
}

// GetMoves returns a slice of all valid moves for given Piece. Each valid move in the slice is stored
// in an array with a length of two-- Row and Col.
// TODO Add checkmate checks
func (p *Pawn) GetMoves(pieces [32]ChessPiece) [][2]int {
	moves := make([][2]int, 2)

	// TODO add en passant, or however you spell that move
	if p.white {
		// white pawn on starting position, so could move forward one or two
		if p.row == 6 {
			moves = append(moves, [2]int{4, p.col}, [2]int{5, p.col})
			//TODO check for available pawn takes
			//GetPiecesOnSquare + validate position is in bounds of the board
		}
	} else {
		// black pawn on starting position, so could move forward one or two
		if p.row == 1 {
			moves = append(moves, [2]int{3, p.col}, [2]int{2, p.col})
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
	// Reusing the Piece GetImage for filesystem functionality
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

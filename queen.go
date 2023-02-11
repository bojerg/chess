package main

import "github.com/hajimehoshi/ebiten/v2"

type Queen struct {
	Piece
}

// GetMoves returns a slice of all valid moves for given Piece. Each valid move in the slice is stored
// in an array with a length of two-- Row and Col.
// TODO everything
func (p *Queen) GetMoves(pieces [32]ChessPiece) [][2]int {
	moves := make([][2]int, 2)
	return moves
}

// GetName primarily intended for debugging
func (p *Queen) GetName() string {
	if p.white {
		return "White queen"
	} else {
		return "Black queen"
	}
}

func (p *Queen) GetImage() *ebiten.Image {
	filepathStr := "images/"
	if p.white {
		filepathStr += "whiteQueen.png"
	} else {
		filepathStr += "blackQueen.png"
	}
	// Reusing the Piece GetImage for filesystem functionality
	return GetImage(filepathStr)
}

func (p *Queen) GetCol() int {
	return p.col
}

func (p *Queen) GetRow() int {
	return p.row
}

func (p *Queen) SetCol(c int) {
	p.col = c
}

func (p *Queen) SetRow(r int) {
	p.row = r
}

func (p *Queen) White() bool {
	return p.white
}

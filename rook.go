package main

import "github.com/hajimehoshi/ebiten/v2"

type Rook struct {
	Piece
}

// GetMoves returns a slice of all valid moves for given Piece. Each valid move in the slice is stored
// in an array with a length of two-- Row and Col.
// TODO everything
func (p *Rook) GetMoves(pieces [32]ChessPiece) [][2]int {
	//moves := make([][2]int, 2)
	return nil
}

// GetName primarily intended for debugging
func (p *Rook) GetName() string {
	if p.white {
		return "White rook"
	} else {
		return "Black rook"
	}
}

func (p *Rook) GetImage() *ebiten.Image {
	filepathStr := "images/"
	if p.white {
		filepathStr += "whiteRook.png"
	} else {
		filepathStr += "blackRook.png"
	}
	// Reusing the Piece GetImage for filesystem functionality
	return GetImage(filepathStr)
}

func (p *Rook) GetCol() int {
	return p.col
}

func (p *Rook) GetRow() int {
	return p.row
}

func (p *Rook) SetCol(c int) {
	p.col = c
}

func (p *Rook) SetRow(r int) {
	p.row = r
}

func (p *Rook) White() bool {
	return p.white
}

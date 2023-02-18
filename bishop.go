package main

import "github.com/hajimehoshi/ebiten/v2"

type Bishop struct {
	Piece
}

// GetMoves returns a slice of all valid moves for given Piece. Each valid move in the slice is stored
// in an array with a length of two-- Row and Col.
// TODO everything
func (p *Bishop) GetMoves(pieces [32]ChessPiece) [][2]int {
	//moves := make([][2]int, 2)
	return nil
}

// GetName primarily intended for debugging
func (p *Bishop) GetName() string {
	if p.white {
		return "White bishop"
	} else {
		return "Black bishop"
	}
}

func (p *Bishop) GetImage() *ebiten.Image {
	filepathStr := "images/"
	if p.white {
		filepathStr += "whiteBishop.png"
	} else {
		filepathStr += "blackBishop.png"
	}
	// Reusing the Piece GetImage for filesystem functionality
	return GetImage(filepathStr)
}

func (p *Bishop) GetCol() int {
	return p.col
}

func (p *Bishop) GetRow() int {
	return p.row
}

func (p *Bishop) SetCol(c int) {
	p.col = c
}

func (p *Bishop) SetRow(r int) {
	p.row = r
}

func (p *Bishop) White() bool {
	return p.white
}

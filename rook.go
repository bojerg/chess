package main

import "github.com/hajimehoshi/ebiten/v2"

type Rook struct {
	Piece
}

// Moves returns a slice of all valid moves for given Piece. Each valid move in the slice is stored
// in an array with a length of two-- Row and Col.
func (p *Rook) Moves(g Game) [][2]int {
	moves := make([][2]int, 0)

	for col := p.col - 1; col >= 0; col-- {
		otherPiece := GetPieceOnSquare(p.row, col, g.pieces)
		if otherPiece == nil {
			//No piece encountered, valid move
			moves = append(moves, [2]int{p.row, col})
		} else {
			//Piece encountered... is it on the other team?
			if otherPiece.White() != p.white {
				moves = append(moves, [2]int{p.row, col})
			}
			// Can't go further in this loop / can't go further this direction on the board
			break
		}
	}

	for col := p.col + 1; col <= 7; col++ {
		otherPiece := GetPieceOnSquare(p.row, col, g.pieces)
		if otherPiece == nil {
			//No piece encountered, valid move
			moves = append(moves, [2]int{p.row, col})
		} else {
			//Piece encountered... is it on the other team?
			if otherPiece.White() != p.white {
				moves = append(moves, [2]int{p.row, col})
			}
			// Can't go further in this loop / can't go further this direction on the board
			break
		}
	}

	for row := p.row - 1; row >= 0; row-- {
		otherPiece := GetPieceOnSquare(row, p.col, g.pieces)
		if otherPiece == nil {
			//No piece encountered, valid move
			moves = append(moves, [2]int{row, p.col})
		} else {
			//Piece encountered... is it on the other team?
			if otherPiece.White() != p.white {
				moves = append(moves, [2]int{row, p.col})
			}
			// Can't go further in this loop / can't go further this direction on the board
			break
		}
	}

	for row := p.row + 1; row <= 7; row++ {
		otherPiece := GetPieceOnSquare(row, p.col, g.pieces)
		if otherPiece == nil {
			//No piece encountered, valid move
			moves = append(moves, [2]int{row, p.col})
		} else {
			//Piece encountered... is it on the other team?
			if otherPiece.White() != p.white {
				moves = append(moves, [2]int{row, p.col})
			}
			// Can't go further in this loop / can't go further this direction on the board
			break
		}
	}

	return moves
}

func (p *Rook) Name() string {
	if p.white {
		return "White rook"
	} else {
		return "Black rook"
	}
}

func (p *Rook) Image() *ebiten.Image {
	filepathStr := "images/"
	if p.white {
		filepathStr += "whiteRook.png"
	} else {
		filepathStr += "blackRook.png"
	}
	// Reusing the Piece GetImage for filesystem functionality
	return GetImage(filepathStr)
}

func (p *Rook) Col() int {
	return p.col
}

func (p *Rook) Row() int {
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

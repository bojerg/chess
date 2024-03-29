package main

import "github.com/hajimehoshi/ebiten/v2"

type Bishop struct {
	Piece
}

// Moves returns a slice of all valid moves for given Piece. Each valid move in the slice is stored
// in an array with a length of two-- Row and Col.
func (p *Bishop) Moves(g Game) [][2]int {
	moves := make([][2]int, 0)

	//The following are for loops in the diagonals
	col := p.col
	row := p.row
	for {

		col++
		row++
		if col > 7 || row > 7 {
			break
		}

		otherPiece := GetPieceOnSquare(row, col, g.pieces)
		if otherPiece == nil {
			//No piece encountered, valid move
			moves = append(moves, [2]int{row, col})
		} else {
			//Piece encountered... is it on the other team?
			if otherPiece.White() != p.white {
				moves = append(moves, [2]int{row, col})
			}
			// Can't go further in this loop / can't go further this direction on the board
			break
		}
	}

	col = p.col
	row = p.row
	for {

		col++
		row--
		if col > 7 || row < 0 {
			break
		}

		otherPiece := GetPieceOnSquare(row, col, g.pieces)
		if otherPiece == nil {
			//No piece encountered, valid move
			moves = append(moves, [2]int{row, col})
		} else {
			//Piece encountered... is it on the other team?
			if otherPiece.White() != p.white {
				moves = append(moves, [2]int{row, col})
			}
			// Can't go further in this loop / can't go further this direction on the board
			break
		}

	}

	col = p.col
	row = p.row
	for {

		col--
		row--
		if col < 0 || row < 0 {
			break
		}

		otherPiece := GetPieceOnSquare(row, col, g.pieces)
		if otherPiece == nil {
			//No piece encountered, valid move
			moves = append(moves, [2]int{row, col})
		} else {
			//Piece encountered... is it on the other team?
			if otherPiece.White() != p.white {
				moves = append(moves, [2]int{row, col})
			}
			// Can't go further in this loop / can't go further this direction on the board
			break
		}

	}

	col = p.col
	row = p.row
	for {

		col--
		row++
		if col < 0 || row > 7 {
			break
		}

		otherPiece := GetPieceOnSquare(row, col, g.pieces)
		if otherPiece == nil {
			//No piece encountered, valid move
			moves = append(moves, [2]int{row, col})
		} else {
			//Piece encountered... is it on the other team?
			if otherPiece.White() != p.white {
				moves = append(moves, [2]int{row, col})
			}
			// Can't go further in this loop / can't go further this direction on the board
			break
		}

	}

	return moves
}

func (p *Bishop) Name() string {
	if p.white {
		return "White bishop"
	} else {
		return "Black bishop"
	}
}

func (p *Bishop) Image() *ebiten.Image {
	filepathStr := "images/"
	if p.white {
		filepathStr += "whiteBishop.png"
	} else {
		filepathStr += "blackBishop.png"
	}
	// Reusing the Piece GetImage for filesystem functionality
	return GetImage(filepathStr)
}

func (p *Bishop) Col() int {
	return p.col
}

func (p *Bishop) Row() int {
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

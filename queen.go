package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Queen struct {
	Piece
}

// GetMoves returns a slice of all valid moves for given Piece. Each valid move in the slice is stored
// in an array with a length of two-- Row and Col.
// TODO everything
func (p *Queen) GetMoves(pieces [32]ChessPiece) [][2]int {
	moves := make([][2]int, 0)

	//Queen can go any direction until it encounters a piece. If not on it's team, can take.
	for col := p.col - 1; col >= 0; col-- {
		otherPiece := GetPieceOnSquare(p.row, col, pieces)
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
		otherPiece := GetPieceOnSquare(p.row, col, pieces)
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
		otherPiece := GetPieceOnSquare(row, p.col, pieces)
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
		otherPiece := GetPieceOnSquare(row, p.col, pieces)
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

	//The following are for loops in the diagonals
	col := p.col
	row := p.row
	for {

		col++
		row++
		if col > 7 || row > 7 {
			break
		}

		otherPiece := GetPieceOnSquare(row, col, pieces)
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

		otherPiece := GetPieceOnSquare(row, col, pieces)
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

		otherPiece := GetPieceOnSquare(row, col, pieces)
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

		otherPiece := GetPieceOnSquare(row, col, pieces)
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

func (p *Queen) IsKing() bool {
	return false
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

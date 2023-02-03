package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	_ "image/png" // required for ebitenutil/NewImageFromFile
	"log"
	"path/filepath"
	"strconv"
)

// Piece
// id: 0 pawn, 1 knight, 2 bishop, 3 rook, 4 queen, 5 king, 6 none/taken
// col, row: identifies location on chess board
// white: If on team white, is true
type Piece struct {
	id    int
	col   int
	row   int
	white bool
}

// GetMoves returns a slice of all valid moves for given Piece. Each valid move in the slice is stored
// in an array with a length of two-- Row and Col.
// TODO implement this where needed, and add checkmate checks
func (p *Piece) GetMoves(pieces [32]*Piece) [][2]int {
	moves := make([][2]int, 2)
	done := false // loop sentinel value

	switch p.id {
	case 0: // pawn
		// TODO add en passant, or however you spell that move
		if p.white {
			// white pawn on starting position, so could move forward one or two
			if p.row == 6 {
				moves = append(moves, [2]int{4, p.col}, [2]int{5, p.col})
			}
			// TODO check if there's a piece to take
		} else {
			// black pawn on starting position, so could move forward one or two
			if p.row == 1 {
				moves = append(moves, [2]int{3, p.col}, [2]int{2, p.col})
			}
			// TODO check if there's a piece to take
		}
	case 1: // knight

		for !done {
			done = true
		}
	case 2: // bishop

		for !done {
			done = true
		}
	case 3: // rook

		for !done {
			done = true
		}
	case 4: // queen

		for !done {
			done = true
		}
	case 5: // king

		for !done {
			done = true
		}
	}

	return moves
}

// GetName primarily intended for debugging
func (p *Piece) GetName() string {
	var ret string

	switch p.id {
	case 0:
		ret = "Pawn"
	case 1:
		ret = "Knight"
	case 2:
		ret = "Bishop"
	case 3:
		ret = "Rook"
	case 4:
		ret = "Queen"
	case 5:
		ret = "King"
	}

	if p.white {
		ret += ", white"
	} else {
		ret += ", black"
	}

	ret += "\trow:" + strconv.Itoa(int(p.row)) + " col:" + strconv.Itoa(int(p.col))

	return ret
}

// GetImage returns the corresponding ebiten image
func (p *Piece) GetImage() *ebiten.Image {
	// https://commons.wikimedia.org/wiki/Category:PNG_chess_pieces/Standard_transparent
	filepathStr := "images/"
	switch p.id {
	case 0:
		if p.white {
			filepathStr += "whitePawn.png"
		} else {
			filepathStr += "blackPawn.png"
		}
	case 1:
		if p.white {
			filepathStr += "whiteKnight.png"
		} else {
			filepathStr += "blackKnight.png"
		}

	case 2:
		if p.white {
			filepathStr += "whiteBishop.png"
		} else {
			filepathStr += "blackBishop.png"
		}

	case 3:
		if p.white {
			filepathStr += "whiteRook.png"
		} else {
			filepathStr += "blackRook.png"
		}

	case 4:
		if p.white {
			filepathStr += "whiteQueen.png"
		} else {
			filepathStr += "blackQueen.png"
		}

	case 5:
		if p.white {
			filepathStr += "whiteKing.png"
		} else {
			filepathStr += "blackKing.png"
		}
	default:
		fmt.Println("piece.GetImage(), switch input: ", p.id)
	}

	var err error = nil
	// path/filepath creates filepath for any OS, supposedly
	fileLoc, _ := filepath.Abs(filepathStr)
	pieceImage, _, err := ebitenutil.NewImageFromFile(fileLoc)
	if err != nil {
		log.Fatal(err)
		return nil
	} else {
		return pieceImage
	}

}

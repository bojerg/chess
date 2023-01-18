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
	id    int8
	col   int8
	row   int8
	white bool
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
		fmt.Println("...wtf")
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

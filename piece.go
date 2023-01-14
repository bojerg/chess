package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"log"
)

// Piece
// id: 0 pawn, 1 knight, 2 bishop, 3 rook, 4 queen, 5 king
// col, row: identifies location on chess board
// white: If on team white, is true
type Piece struct {
	id    int8
	col   int8
	row   int8
	white bool
}

func (p *Piece) DrawImage(pieceImage *ebiten.Image) {
	// https://commons.wikimedia.org/wiki/Category:PNG_chess_pieces/Standard_transparent
	filepath := "/images/"
	switch p.id {
	case 0:
		if p.white {
			filepath += "whitePawn.png"
		} else {
			filepath += "blackPawn.png"
		}
	case 1:
		if p.white {
			filepath += "whiteKnight.png"
		} else {
			filepath += "blackKnight.png"
		}

	case 2:
		if p.white {
			filepath += "whiteBishop.png"
		} else {
			filepath += "blackBishop.png"
		}

	case 3:
		if p.white {
			filepath += "whiteRook.png"
		} else {
			filepath += "blackRook.png"
		}

	case 4:
		if p.white {
			filepath += "whiteQueen.png"
		} else {
			filepath += "blackQueen.png"
		}

	case 5:
		if p.white {
			filepath += "whiteKing.png"
		} else {
			filepath += "blackQueen.png"
		}
	}

	var err error = nil
	pieceImage, _, err = ebitenutil.NewImageFromFile(filepath)
	if err != nil {
		log.Fatal(err)
	}

}

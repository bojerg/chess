package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
)

type Board struct {
}

func (b *Board) Draw(boardImage *ebiten.Image, pieces [32]*Piece, selected [2]float64, selectedPiece int, boardCol int, boardRow int) {

	darkColor := color.RGBA{R: 0xbb, G: 0x99, B: 0x55, A: 0xff}
	lightColor := color.RGBA{R: 0xcb, G: 0xbe, B: 0xb5, A: 0xff}
	selectedColor := color.RGBA{R: 0xea, G: 0xdd, B: 0x23, A: 0xff}

	tileSize := 128
	tileImage := ebiten.NewImage(tileSize, tileSize)
	backImage := ebiten.NewImage(tileSize*8, tileSize*8)

	// Drawing one big light square to cut down on draw ops
	opBackground := &ebiten.DrawImageOptions{}
	opBackground.GeoM.Translate(448, 28)
	backImage.Fill(lightColor)
	boardImage.DrawImage(backImage, opBackground)
	// Row, Column
	// Draw tiles for pieces
	for r := 0; r < 8; r++ {
		for c := 0; c < 8; c++ {
			opDark := &ebiten.DrawImageOptions{}
			opDark.GeoM.Translate(float64(c*tileSize+448), float64(r*tileSize+28))
			//TODO color in available moves for a selected piece
			if r == boardRow && c == boardCol {
				tileImage.Fill(selectedColor)
				boardImage.DrawImage(tileImage, opDark)
			} else if (r%2 == 0 && c%2 != 0) || (r%2 != 0 && c%2 == 0) {
				tileImage.Fill(darkColor)
				boardImage.DrawImage(tileImage, opDark)
			}

		}
	}

	//Draw pieces

	for i := 0; i < len(pieces); i++ {
		var tx float64
		var ty float64

		//if a piece has been selected we want to follow that piece the mouse instead
		if i == selectedPiece {
			tx = selected[0] * 1.5
			ty = selected[1] * 1.5
		} else {
			tx = float64(pieces[i].col*tileSize) + 465
			ty = float64(pieces[i].row*tileSize) + 42
		}

		opPiece := &ebiten.DrawImageOptions{}
		opPiece.GeoM.Scale(1.5, 1.5) //essentially W x H = 90 x 90
		opPiece.GeoM.Translate(tx, ty)

		// id of 6 means piece has been taken
		if pieces[i].id != 6 {
			boardImage.DrawImage(pieces[i].GetImage(), opPiece)
		}
	}

}

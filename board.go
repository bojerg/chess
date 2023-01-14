package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
)

type Board struct {
	pieces [32]Piece
}

var lightColor = color.RGBA{R: 0xbb, G: 0x99, B: 0x55, A: 0xff}
var darkColor = color.RGBA{R: 0xcb, G: 0xbe, B: 0xb5, A: 0xff}

func (b *Board) Draw(boardImage *ebiten.Image) {
	tileSize := 128
	pieceSize := 120
	tileImage := ebiten.NewImage(tileSize, tileSize)
	pieceImage := ebiten.NewImage(pieceSize, pieceSize)

	// Row, Column
	// Draw tiles for pieces
	for r := 0; r < 8; r++ {
		for c := 0; c < 8; c++ {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(c*tileSize+448), float64(r*tileSize+28))
			//determine if tile is light or dark
			if r%2 == 0 {
				if c%2 == 0 {
					tileImage.Fill(lightColor)
				} else {
					tileImage.Fill(darkColor)
				}
			} else {
				if c%2 == 0 {
					tileImage.Fill(darkColor)
				} else {
					tileImage.Fill(lightColor)
				}
			}

			boardImage.DrawImage(tileImage, op)
		}
	}

	fmt.Println(len(b.pieces))

	//Draw pieces
	for i := 0; i < len(b.pieces); i++ {
		b.pieces[i].DrawImage(pieceImage)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(int(b.pieces[i].col)*tileSize+448), float64(int(b.pieces[i].row)*tileSize+28))
		//boardImage.DrawImage(pieceImage, op)
	}

}

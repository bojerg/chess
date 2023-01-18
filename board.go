package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
)

type Board struct {
}

var lightColor = color.RGBA{R: 0xbb, G: 0x99, B: 0x55, A: 0xff}
var darkColor = color.RGBA{R: 0xcb, G: 0xbe, B: 0xb5, A: 0xff}

func (b *Board) Draw(boardImage *ebiten.Image, pieces [32]*Piece) {
	tileSize := 128

	tileImage := ebiten.NewImage(tileSize, tileSize)

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

	//Draw pieces
	//TODO reformat to allow for missing pieces
	for i := 0; i < len(pieces); i++ {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(pieces[i].col)*85.3+310, float64(pieces[i].row)*85.3+28)
		//essentially W x H = 90 x 90
		op.GeoM.Scale(1.5, 1.5)
		boardImage.DrawImage(pieces[i].GetImage(), op)
	}

}

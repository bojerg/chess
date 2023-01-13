package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
)

type Board struct {
}

func (b *Board) Draw(boardImage *ebiten.Image) {
	tileSize := 128
	tileImage := ebiten.NewImage(tileSize, tileSize)
	// Row, Column
	for r := 0; r < 8; r++ {
		for c := 0; c < 8; c++ {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(c*tileSize+448), float64(r*tileSize+28))

			//determine if tile is light or dark
			if r%2 == 0 {
				if c%2 == 0 {
					tileImage.Fill(color.White)
				} else {
					tileImage.Fill(color.Black)
				}
			} else {
				if c%2 == 0 {
					tileImage.Fill(color.Black)
				} else {
					tileImage.Fill(color.White)
				}
			}

			boardImage.DrawImage(tileImage, op)
		}
	}
}

package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
)

type Board struct {
}

var lightColor = color.RGBA{R: 0xbb, G: 0x99, B: 0x55, A: 0xff}
var darkColor = color.RGBA{R: 0xcb, G: 0xbe, B: 0xb5, A: 0xff}

func (b *Board) Draw(boardImage *ebiten.Image, pieces [32]*Piece, selected [2]float64, selectedPiece int) {
	tileSize := 128
	tileImage := ebiten.NewImage(tileSize, tileSize)
	// Row, Column
	// Draw tiles for pieces
	for r := 0; r < 8; r++ {
		for c := 0; c < 8; c++ {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(c*tileSize+448), float64(r*tileSize+28))
			//determine if tile is light or dark, or selected
			//TODO color in available moves for a selected piece
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
	for i := 0; i < len(pieces); i++ {
		var tx float64
		var ty float64

		//if a piece has been selected we want to follow that piece the mouse instead
		if i == selectedPiece {
			tx = selected[0]
			ty = selected[1]
		} else {
			tx = float64(pieces[i].col)*85.33 + 310
			ty = float64(pieces[i].row)*85.33 + 28
		}

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(tx, ty)
		//essentially W x H = 90 x 90
		op.GeoM.Scale(1.5, 1.5)
		// id of 6 means piece has been taken
		if pieces[i].id != 6 {
			boardImage.DrawImage(pieces[i].GetImage(), op)
		}
	}

}

package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
)

type Board struct {
}

func (b *Board) Draw(gameImage *ebiten.Image, pieces [32]*Piece, selected [2]float64, selectedPiece int, boardCol int, boardRow int) {
	tileImage := ebiten.NewImage(TILE_SIZE, TILE_SIZE)

	//gameImage.Clear()

	// Row, Column
	// Draw selected tile
	for r := 0; r < 8; r++ {
		for c := 0; c < 8; c++ {
			//TODO color in available moves for a selected piece as well?
			if r == boardRow && c == boardCol {
				opTile := &ebiten.DrawImageOptions{}
				opTile.GeoM.Translate(float64(c*TILE_SIZE+448), float64(r*TILE_SIZE+28))
				tileImage.Fill(color.RGBA{R: 0xea, G: 0xdd, B: 0x23, A: 0xff})
				gameImage.DrawImage(tileImage, opTile)
				break
			}
		}
	}

	//Draw selected (moving) piece

	for i := 0; i < len(pieces); i++ {
		//if a piece has been selected we want to follow that piece the mouse instead
		if i == selectedPiece {
			tx := selected[0] * 1.5
			ty := selected[1] * 1.5
			opPiece := &ebiten.DrawImageOptions{}
			opPiece.GeoM.Scale(1.5, 1.5) //essentially W x H = 90 x 90
			opPiece.GeoM.Translate(tx, ty)
			gameImage.DrawImage(pieces[i].GetImage(), opPiece)
			break
		}
	}

}

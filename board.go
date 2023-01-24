package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
)

type Board struct {
}

func (b *Board) DrawStaticPieces(pieceImage *ebiten.Image, pieces [32]*Piece, selectedPiece int) {
	pieceImage.Clear()

	for i := 0; i < len(pieces); i++ {
		// Don't draw selected (moving) piece, or any pieces with id of 6 (taken)
		if i != selectedPiece && pieces[i].id != 6 {
			tx := float64(pieces[i].col*TileSize) + 465
			ty := float64(pieces[i].row*TileSize) + 42
			opPiece := &ebiten.DrawImageOptions{}
			opPiece.GeoM.Scale(1.5, 1.5) //essentially W x H = 90 x 90
			opPiece.GeoM.Translate(tx, ty)
			pieceImage.DrawImage(pieces[i].GetImage(), opPiece)
		}
	}
}

func (b *Board) DrawMovingPiece(gameImage *ebiten.Image, pieces [32]*Piece, selected [2]float64, selectedPiece int) {
	for i := 0; i < len(pieces); i++ {
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

func (b *Board) DrawHighlightedTiles(gameImage *ebiten.Image, selectedCol int, selectedRow int) {
	tileImage := ebiten.NewImage(TileSize, TileSize)
	gameImage.Clear()

	// Row, Column
	// Draw selected tile
	for r := 0; r < 8; r++ {
		for c := 0; c < 8; c++ {
			//TODO color in available moves for a selected piece as well?
			if r == selectedRow && c == selectedCol {
				opTile := &ebiten.DrawImageOptions{}
				opTile.GeoM.Translate(float64(c*TileSize+448), float64(r*TileSize+28))
				tileImage.Fill(color.RGBA{R: 0xea, G: 0xdd, B: 0x23, A: 0xff})
				gameImage.DrawImage(tileImage, opTile)
				break
			}
		}
	}
}

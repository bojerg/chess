package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
	"math"
)

// Board
// A collection of helper functions to render the images displayed every frame
// whitesTurn... true if it's white's turn, duh!
// inCheck signifies the player who is up to move is in check. I don't think there's any way both players could be
// in check simultaneously
// scheduleDraw bool is a sentinel value to indicate that the piece locations have changed
// and the pieceImage should be redrawn
type Board struct {
	whitesTurn   bool
	inCheck      bool
	scheduleDraw bool
}

func (b *Board) DrawStaticPieces(pieceImage *ebiten.Image, pieces [32]ChessPiece, selectedPiece int) {
	pieceImage.Clear()

	xOffset := 465.0
	yOffset := 42.0
	rotate := 0.0

	if !b.whitesTurn {
		rotate = math.Pi
		xOffset += TileSize - 34
		yOffset += TileSize - 28
	}

	for i, piece := range pieces {
		// Don't draw selected (moving) piece, or any pieces with id of 6 (taken)
		if i != selectedPiece && piece.GetCol() != -1 {
			tx := float64(pieces[i].GetCol()*TileSize) + xOffset
			ty := float64(pieces[i].GetRow()*TileSize) + yOffset
			opPiece := &ebiten.DrawImageOptions{}
			opPiece.GeoM.Rotate(rotate)
			opPiece.GeoM.Scale(1.5, 1.5) //essentially W x H = 90 x 90
			opPiece.GeoM.Translate(tx, ty)
			pieceImage.DrawImage(pieces[i].GetImage(), opPiece)
		}
	}
}

func (b *Board) DrawMovingPiece(movingImage *ebiten.Image, pieces [32]ChessPiece, selected [2]float64, selectedPiece int) {

	for i, _ := range pieces {
		if i == selectedPiece {
			tx := selected[0] * 1.5
			ty := selected[1] * 1.5
			opPiece := &ebiten.DrawImageOptions{}
			opPiece.GeoM.Scale(1.5, 1.5) //essentially W x H = 90 x 90
			opPiece.GeoM.Translate(tx, ty)
			movingImage.DrawImage(pieces[i].GetImage(), opPiece)
			break
		}
	}
}

func (b *Board) DrawHighlightedTiles(gameImage *ebiten.Image, selectRow int, selectCol int, selectIndex int, pieces [32]ChessPiece) {
	tileImage := ebiten.NewImage(TileSize, TileSize)
	gameImage.Clear()

	// drawing highlighted tiles (available moves in red)
	if selectIndex >= 0 {
		availableMoves := pieces[selectIndex].GetMoves(pieces)
		if availableMoves != nil {
			for _, move := range availableMoves {
				opTile := &ebiten.DrawImageOptions{}
				opTile.GeoM.Translate(float64(move[1]*TileSize+448), float64(move[0]*TileSize+28))
				tileImage.Fill(color.RGBA{R: 0xff, G: 0x06, B: 0x03, A: 0xba})
				gameImage.DrawImage(tileImage, opTile)
			}
		}

	}

	// Draw hovered tile (in highlighter yellow)
	for r := 0; r < 8; r++ {
		for c := 0; c < 8; c++ {
			if r == selectRow && c == selectCol {
				opTile := &ebiten.DrawImageOptions{}
				opTile.GeoM.Translate(float64(c*TileSize+448), float64(r*TileSize+28))
				tileImage.Fill(color.RGBA{R: 0xea, G: 0xdd, B: 0x23, A: 0xff})
				gameImage.DrawImage(tileImage, opTile)
				break
			}
		}
	}

	//highlight a king in check (purple)
	if b.inCheck {
		for _, piece := range pieces {
			if piece.IsKing() && piece.White() == b.whitesTurn {
				opTile := &ebiten.DrawImageOptions{}
				opTile.GeoM.Translate(float64(piece.GetCol()*TileSize+448), float64(piece.GetRow()*TileSize+28))
				tileImage.Fill(color.RGBA{R: 0xbf, G: 0x00, B: 0xe6, A: 0xff})
				gameImage.DrawImage(tileImage, opTile)
				break
			}
		}
	}

}

func (b *Board) DrawBoard(boardImage *ebiten.Image) {
	darkColor := color.RGBA{R: 0xbb, G: 0x99, B: 0x55, A: 0xff}
	lightColor := color.RGBA{R: 0xcb, G: 0xbe, B: 0xb5, A: 0xff}

	lightImage := ebiten.NewImage(TileSize*8, TileSize*8)
	darkImage := ebiten.NewImage(TileSize, TileSize)

	// Drawing one big light square to (slightly) cut down on draw ops
	opLight := &ebiten.DrawImageOptions{}
	opLight.GeoM.Translate(448, 28)
	lightImage.Fill(lightColor)
	boardImage.DrawImage(lightImage, opLight)
	for row := 0; row < 8; row++ {
		for col := 0; col < 8; col++ {
			if (row%2 == 0 && col%2 != 0) || (row%2 != 0 && col%2 == 0) {
				opDark := &ebiten.DrawImageOptions{}
				opDark.GeoM.Translate(float64(col*TileSize+448), float64(row*TileSize+28))
				darkImage.Fill(darkColor)
				boardImage.DrawImage(darkImage, opDark)
			}

		}
	}
}

func (b *Board) DrawUI(uiImage *ebiten.Image, gameOver bool) {

	if gameOver {
		// draw centered block for game over message and controls
	}
}

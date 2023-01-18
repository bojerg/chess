package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
	"log"
)

type Game struct {
	board      *Board
	boardImage *ebiten.Image
	pieces     [32]*Piece
}

const (
	WIDTH  = 1920
	HEIGHT = 1080
)

func (g *Game) Update() error {

	// No way to exit fullscreen without this for now
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		ebiten.SetFullscreen(false)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	if g.boardImage == nil {
		g.boardImage = ebiten.NewImage(WIDTH, HEIGHT)
		g.InitPieces()
	}

	screen.Fill(color.RGBA{R: 0x13, G: 0x33, B: 0x31, A: 0xff})
	g.board.Draw(g.boardImage, g.pieces)

	op := &ebiten.DrawImageOptions{}
	sw, sh := screen.Size()
	bw, bh := g.boardImage.Size()
	x := (sw - bw) / 2
	y := (sh - bh) / 2
	op.GeoM.Translate(float64(x), float64(y))

	screen.DrawImage(g.boardImage, op)
}

func (g *Game) InitPieces() {
	g.pieces[0] = &Piece{3, 0, 0, false}
	g.pieces[1] = &Piece{1, 1, 0, false}
	g.pieces[2] = &Piece{2, 2, 0, false}
	g.pieces[3] = &Piece{4, 3, 0, false}
	g.pieces[4] = &Piece{5, 4, 0, false}
	g.pieces[5] = &Piece{2, 5, 0, false}
	g.pieces[6] = &Piece{1, 6, 0, false}
	g.pieces[7] = &Piece{3, 7, 0, false}
	g.pieces[8] = &Piece{0, 0, 1, false}
	g.pieces[9] = &Piece{0, 1, 1, false}
	g.pieces[10] = &Piece{0, 2, 1, false}
	g.pieces[11] = &Piece{0, 3, 1, false}
	g.pieces[12] = &Piece{0, 4, 1, false}
	g.pieces[13] = &Piece{0, 5, 1, false}
	g.pieces[14] = &Piece{0, 6, 1, false}
	g.pieces[15] = &Piece{0, 7, 1, false}
	g.pieces[16] = &Piece{0, 0, 6, true}
	g.pieces[17] = &Piece{0, 1, 6, true}
	g.pieces[18] = &Piece{0, 2, 6, true}
	g.pieces[19] = &Piece{0, 3, 6, true}
	g.pieces[20] = &Piece{0, 4, 6, true}
	g.pieces[21] = &Piece{0, 5, 6, true}
	g.pieces[22] = &Piece{0, 6, 6, true}
	g.pieces[23] = &Piece{0, 7, 6, true}
	g.pieces[24] = &Piece{3, 0, 7, true}
	g.pieces[25] = &Piece{1, 1, 7, true}
	g.pieces[26] = &Piece{2, 2, 7, true}
	g.pieces[26] = &Piece{2, 2, 7, true}
	g.pieces[27] = &Piece{4, 3, 7, true}
	g.pieces[28] = &Piece{5, 4, 7, true}
	g.pieces[29] = &Piece{2, 5, 7, true}
	g.pieces[30] = &Piece{1, 6, 7, true}
	g.pieces[31] = &Piece{3, 7, 7, true}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return WIDTH, HEIGHT
}

func main() {
	game := &Game{}
	ebiten.SetWindowSize(WIDTH/2, HEIGHT/2)
	ebiten.SetWindowTitle("chess")
	ebiten.SetFullscreen(true)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

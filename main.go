package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
	"log"
)

type Game struct {
	board      *Board
	boardImage *ebiten.Image
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
	g.board.Draw(g.boardImage)

	op := &ebiten.DrawImageOptions{}
	sw, sh := screen.Size()
	bw, bh := g.boardImage.Size()
	x := (sw - bw) / 2
	y := (sh - bh) / 2
	op.GeoM.Translate(float64(x), float64(y))

	screen.DrawImage(g.boardImage, op)
}

func (g *Game) InitPieces() {
	var pieces [32]Piece
	pieces[0] = Piece{3, 0, 0, true}
	pieces[1] = Piece{1, 1, 0, true}
	pieces[2] = Piece{2, 2, 0, true}
	pieces[3] = Piece{4, 3, 0, true}
	pieces[4] = Piece{5, 4, 0, true}
	pieces[5] = Piece{2, 5, 0, true}
	pieces[6] = Piece{1, 6, 0, true}
	pieces[7] = Piece{3, 7, 0, true}
	pieces[8] = Piece{0, 0, 1, true}
	pieces[9] = Piece{0, 1, 1, true}
	pieces[10] = Piece{0, 2, 1, true}
	pieces[11] = Piece{0, 3, 1, true}
	pieces[12] = Piece{0, 4, 1, true}
	pieces[13] = Piece{0, 5, 1, true}
	pieces[14] = Piece{0, 6, 1, true}
	pieces[15] = Piece{0, 7, 1, true}
	pieces[16] = Piece{0, 0, 6, false}
	pieces[17] = Piece{0, 1, 6, false}
	pieces[18] = Piece{0, 2, 6, false}
	pieces[19] = Piece{0, 3, 6, false}
	pieces[20] = Piece{0, 4, 6, false}
	pieces[21] = Piece{0, 5, 6, false}
	pieces[22] = Piece{0, 6, 6, false}
	pieces[23] = Piece{0, 7, 6, false}
	pieces[24] = Piece{3, 0, 7, false}
	pieces[25] = Piece{1, 1, 7, false}
	pieces[26] = Piece{2, 2, 7, false}
	pieces[27] = Piece{4, 3, 7, false}
	pieces[28] = Piece{5, 4, 7, false}
	pieces[29] = Piece{2, 5, 7, false}
	pieces[30] = Piece{1, 6, 7, false}
	pieces[31] = Piece{3, 7, 7, false}

	g.board.pieces = pieces
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

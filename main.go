//go:build ebitenginesinglethread

package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	_ "github.com/silbinarywolf/preferdiscretegpu" // Fix for discrete GPUs in windows
	"image/color"
	"log"
	"math"
)

// Game
// board is the chess board object/logic
// boardImage is, well, the board image
// pieces is an array of all the pieces...
// selected is for the x, y values of a piece in motion
// selectedPiece is the index of the selected piece... -1 means none
// boardCol, boardRow is the hovered over/selected board square
type Game struct {
	board         *Board
	boardImage    *ebiten.Image
	pieces        [32]*Piece
	selected      [2]float64
	selectedPiece int
	boardCol      int
	boardRow      int
}

const (
	WIDTH  = 1920
	HEIGHT = 1080
)

func (g *Game) Update() error {

	// mouse position and relative board position
	x, y := ebiten.CursorPosition()

	// fancy min max floor math to determine the closest board square to the cursor
	g.boardCol = int(math.Floor(math.Min(math.Max(float64((x-448)/128), 0), 7)))
	g.boardRow = int(math.Floor(math.Min(math.Max(float64((y-28)/128), 0), 7)))

	// No way to exit fullscreen without this for now
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		ebiten.SetFullscreen(false)
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {

		if g.selectedPiece == -1 {
			// No piece selected but left mouse is held down
			for i := 0; i < len(g.pieces); i++ {
				if g.pieces[i].col == g.boardCol {
					if g.pieces[i].row == g.boardRow {
						g.selectedPiece = i
						g.selected[0] = float64(x)/1.5 - 30
						g.selected[1] = float64(y)/1.5 - 30
						break
					}
				}
			}
		} else {
			// Piece is being held by the mouse
			// the selected array is passed to board's draw method for x y coordinates of moving piece
			g.selected[0] = float64(x)/1.5 - 30
			g.selected[1] = float64(y)/1.5 - 30
		}
	} else {
		if g.selectedPiece != -1 {
			// Piece is selected but left mouse is now released
			g.CheckPieces(g.boardRow, g.boardCol, true)
			g.pieces[g.selectedPiece].row = g.boardRow
			g.pieces[g.selectedPiece].col = g.boardCol
			//TODO make this work and delete the following
			g.selectedPiece = -1
		}
	}

	return nil
}

// CheckPieces checks if there is a piece on the square and will set any piece there to id = 6 (taken) if
// the "takeIt" bool is true. It returns the index of the piece in the game's pieces array, or -1 if none found.
// The game's logic should prevent this from running if no piece index is stored in selectedPiece int...
func (g *Game) CheckPieces(row int, col int, takeIt bool) int {
	for i, piece := range g.pieces {
		if piece.row == row && piece.col == col {
			if takeIt && i != g.selectedPiece {
				piece.id = 6 // 6 = piece taken
			}
			return i
		}
	}

	return -1
}

func (g *Game) Draw(screen *ebiten.Image) {

	if g.boardImage == nil {
		g.boardImage = ebiten.NewImage(WIDTH, HEIGHT)
	}

	g.boardImage.Fill(color.RGBA{R: 0x13, G: 0x33, B: 0x31, A: 0xff})
	g.board.Draw(g.boardImage, g.pieces, g.selected, g.selectedPiece, g.boardCol, g.boardRow)

	op := &ebiten.DrawImageOptions{}
	sw, sh := screen.Size()
	bw, bh := g.boardImage.Size()
	x := (sw - bw) / 2
	y := (sh - bh) / 2
	op.GeoM.Translate(float64(x), float64(y))

	screen.DrawImage(g.boardImage, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return WIDTH, HEIGHT
}

func main() {
	game := &Game{}
	game.InitPieces()
	ebiten.SetWindowSize(WIDTH/2, HEIGHT/2)
	ebiten.SetWindowTitle("chess")
	ebiten.SetFullscreen(true)
	//ebiten.SetScreenClearedEveryFrame(false)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

func (g *Game) InitPieces() {
	g.selectedPiece = -1
	g.selected[0] = 0.0
	g.selected[1] = 0.0
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

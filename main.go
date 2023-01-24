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
// gameImage is the most foreground-- like moving pieces, selected tiles, UI, etc
// board is becoming just draw functions for pieces on the board, needs a name refactor
// boardImage is, well, the board image
// pieceImage is a static image of where pieces lay
// pieces is an array of all the pieces...
// selected is for the x, y values of a piece in motion
// selectedPiece is the index of the selected piece... -1 means none
// selectedCol, selectedRow is the hovered over/selected board square
// scheduleDraw bool is a sentinel value to indicate that the piece locations have changed
// and the pieceImage should be redrawn
type Game struct {
	gameImage     *ebiten.Image
	board         *Board
	boardImage    *ebiten.Image
	pieceImage    *ebiten.Image
	pieces        [32]*Piece
	selected      [2]float64
	selectedPiece int
	selectedCol   int
	selectedRow   int
	scheduleDraw  bool
}

const (
	Width    = 1920
	Height   = 1080
	TileSize = 128
)

// Update
// Required function by ebitengine. Contains the logic ran every tick of the game.
func (g *Game) Update() error {
	x, y := ebiten.CursorPosition()

	// fancy min max floor math to determine the closest board square to the cursor, even
	// when the mouse is not over the board
	g.selectedCol = int(math.Floor(math.Min(math.Max(float64((x-448)/128), 0), 7)))
	g.selectedRow = int(math.Floor(math.Min(math.Max(float64((y-28)/128), 0), 7)))

	// No way to exit fullscreen without this for now
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		ebiten.SetFullscreen(false)
	}

	// left click hold and drag
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		// No piece selected but left mouse is held down
		if g.selectedPiece == -1 {
			// figure out which piece is selected
			for i := 0; i < len(g.pieces); i++ {
				if g.pieces[i].col == g.selectedCol {
					if g.pieces[i].row == g.selectedRow {
						g.selectedPiece = i
						// store the xy coordinates of the cursor
						g.selected[0] = float64(x)/1.5 - 30
						g.selected[1] = float64(y)/1.5 - 30
						g.scheduleDraw = true
						break
					}
				}
			}
		} else {
			// update current mouse position because piece is still selected
			g.selected[0] = float64(x)/1.5 - 30
			g.selected[1] = float64(y)/1.5 - 30
		}
	} else { // MouseButtonLeft is not pressed
		if g.selectedPiece != -1 {
			// piece is asking to be let go of at it the current mouse position
			// TODO rules check goes here
			g.CheckPieces(g.selectedRow, g.selectedCol, true)
			g.pieces[g.selectedPiece].row = g.selectedRow
			g.pieces[g.selectedPiece].col = g.selectedCol
			g.selectedPiece = -1
			g.scheduleDraw = true
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
	g.board.DrawHighlightedTiles(g.gameImage, g.selectedCol, g.selectedRow)

	// no moving pieces
	if g.selectedPiece != -1 {
		g.board.DrawMovingPiece(g.gameImage, g.pieces, g.selected, g.selectedPiece)
	}

	// if game logic signals that piece locations have changed, then DrawStaticPieces, and...
	// stop scheduling draw because we did the update to the static piece image
	if g.scheduleDraw {
		g.board.DrawStaticPieces(g.pieceImage, g.pieces, g.selectedPiece)
		g.scheduleDraw = false
	}

	op := &ebiten.DrawImageOptions{}
	sw, sh := screen.Size()
	bw, bh := g.boardImage.Size()
	x := (sw - bw) / 2
	y := (sh - bh) / 2
	op.GeoM.Translate(float64(x), float64(y))

	screen.DrawImage(g.boardImage, op)
	screen.DrawImage(g.gameImage, op)
	screen.DrawImage(g.pieceImage, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return Width, Height
}

func main() {
	game := &Game{}
	game.InitBoard()
	game.InitPieces()
	ebiten.SetWindowSize(Width/2, Height/2)
	ebiten.SetWindowTitle("chess")
	ebiten.SetFullscreen(true)
	//ebiten.SetScreenClearedEveryFrame(false)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

// InitPieces
// Assign variables with starting positions. Render the first image of the pieces in their starting positions.
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

	g.pieceImage = ebiten.NewImage(Width, Height)
	g.board.DrawStaticPieces(g.pieceImage, g.pieces, g.selectedPiece)
	g.scheduleDraw = false
}

// InitBoard
// The "Board" struct currently provides only helper functions which are drawn over the board's spaces. The board
// image which is seen on screen is rendered in this function, after assigning a few variables.
func (g *Game) InitBoard() {
	lightImage := ebiten.NewImage(TileSize*8, TileSize*8)
	darkImage := ebiten.NewImage(TileSize, TileSize)
	g.gameImage = ebiten.NewImage(Width, Height)
	g.boardImage = ebiten.NewImage(Width, Height)
	g.boardImage.Fill(color.RGBA{R: 0x13, G: 0x33, B: 0x31, A: 0xff})
	darkColor := color.RGBA{R: 0xbb, G: 0x99, B: 0x55, A: 0xff}
	lightColor := color.RGBA{R: 0xcb, G: 0xbe, B: 0xb5, A: 0xff}

	// Drawing one big light square to (slightly) cut down on draw ops
	opLight := &ebiten.DrawImageOptions{}
	opLight.GeoM.Translate(448, 28)
	lightImage.Fill(lightColor)
	g.boardImage.DrawImage(lightImage, opLight)
	for row := 0; row < 8; row++ {
		for col := 0; col < 8; col++ {
			if (row%2 == 0 && col%2 != 0) || (row%2 != 0 && col%2 == 0) {
				opDark := &ebiten.DrawImageOptions{}
				opDark.GeoM.Translate(float64(col*TileSize+448), float64(row*TileSize+28))
				darkImage.Fill(darkColor)
				g.boardImage.DrawImage(darkImage, opDark)
			}

		}
	}
}

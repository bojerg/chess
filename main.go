//go:build ebitenginesinglethread

package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	_ "github.com/silbinarywolf/preferdiscretegpu" // Fix for discrete GPUs in windows
	"log"
	"math"
)

// Game
// gameImage is the most foreground-- like moving pieces, selected tiles, UI, etc
// board is becoming just draw functions for pieces on the board, needs a name refactor
// boardImage is, well, the board image
// movingImage is the moving piece
// pieceImage is a static image of where pieces lay
// pieces is an array of all the pieces...
// selected is for the x, y values of a piece in motion
// selectedPiece is the index of the selected piece... -1 means none
// selectedCol, selectedRow is the hovered over/selected board square

type Game struct {
	gameImage     *ebiten.Image
	board         *Board
	boardImage    *ebiten.Image
	movingImage   *ebiten.Image
	pieceImage    *ebiten.Image
	pieces        [32]ChessPiece
	selected      [2]float64
	selectedPiece int
	selectedCol   int
	selectedRow   int
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
			// Match the selected tile to a piece location. Then, ensure the piece belongs to the
			// team whose turn it currently is, and that it is still in play.
			for i := 0; i < len(g.pieces); i++ {
				if g.pieces[i].GetCol() == g.selectedCol && g.pieces[i].GetRow() == g.selectedRow {
					if g.pieces[i].GetCol() != -1 && g.board.whitesTurn == g.pieces[i].White() {
						g.selectedPiece = i
						// store the xy coordinates of the cursor
						g.selected[0] = float64(x)/1.5 - 30
						g.selected[1] = float64(y)/1.5 - 30
						g.board.scheduleDraw = true
						break
					}
				}
			}
		} else {
			// update current mouse position because piece is still selected
			// and the mouse may be moving!
			g.selected[0] = float64(x)/1.5 - 30
			g.selected[1] = float64(y)/1.5 - 30
		}
	} else { // MouseButtonLeft is not pressed
		if g.selectedPiece != -1 {
			// piece is asking to be let go of at it the current mouse position
			// TODO rules check goes here
			g.CheckPieces(g.selectedRow, g.selectedCol, true)
			g.pieces[g.selectedPiece].SetRow(g.selectedRow)
			g.pieces[g.selectedPiece].SetCol(g.selectedCol)
			g.selectedPiece = -1
			g.board.scheduleDraw = true
		}
	}

	return nil
}

// CheckPieces checks if there is a piece on the square and will set any piece there to id = 6 (taken) if
// the "takeIt" bool is true. It returns the index of the piece in the game's pieces array, or -1 if none found.
// The game's logic should prevent this from running if no piece index is stored in selectedPiece int...
func (g *Game) CheckPieces(row int, col int, takeIt bool) int {
	for i, piece := range g.pieces {
		if piece.GetRow() == row && piece.GetCol() == col {
			if takeIt && i != g.selectedPiece {
				piece.SetCol(-1)                         // Col of -1 is de facto notation for piece taken
				g.board.whitesTurn = !g.board.whitesTurn //switch turns
				return i
			} else if i == g.selectedPiece {
				return i
			}
		}
	}

	//We get here if there was no piece on the square
	g.board.whitesTurn = !g.board.whitesTurn //switch turns
	return -1
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.board.DrawHighlightedTiles(g.gameImage, g.selectedRow, g.selectedCol, g.selectedPiece, g.pieces)

	// no moving pieces
	g.movingImage.Clear()
	if g.selectedPiece != -1 {
		g.board.DrawMovingPiece(g.movingImage, g.pieces, g.selected, g.selectedPiece)
	}

	// if game logic signals that piece locations have changed, then DrawStaticPieces, and...
	// stop scheduling draw because we did the update to the static piece image
	if g.board.scheduleDraw {
		g.board.DrawStaticPieces(g.pieceImage, g.pieces, g.selectedPiece)
		g.board.scheduleDraw = false
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
	screen.DrawImage(g.movingImage, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return Width, Height
}

func main() {
	game := &Game{}
	game.InitBoard()
	ebiten.SetWindowSize(Width/2, Height/2)
	ebiten.SetWindowTitle("chess")
	ebiten.SetFullscreen(true)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

// InitBoard
// The "Board" struct currently provides only helper functions which are drawn over the board's spaces. The board
// image which is seen on screen is rendered in this function, after assigning a few variables.
func (g *Game) InitBoard() {
	g.board = &Board{}
	g.boardImage = ebiten.NewImage(Width, Height)
	g.gameImage = ebiten.NewImage(Width, Height)
	g.pieceImage = ebiten.NewImage(Width, Height)
	g.movingImage = ebiten.NewImage(Width, Height)
	g.InitPieces()
	g.board.DrawBoard(g.boardImage)
	g.board.DrawStaticPieces(g.pieceImage, g.pieces, g.selectedPiece)
	g.board.whitesTurn = true
	g.board.scheduleDraw = false
}

// InitPieces
// Assign variables with starting positions.
func (g *Game) InitPieces() {
	g.selectedPiece = -1
	g.selected[0] = 0.0
	g.selected[1] = 0.0
	g.pieces[0] = &Rook{Piece{0, 0, false}}
	g.pieces[1] = &Knight{Piece{1, 0, false}}
	g.pieces[2] = &Bishop{Piece{2, 0, false}}
	g.pieces[3] = &Queen{Piece{3, 0, false}}
	g.pieces[4] = &King{Piece{4, 0, false}}
	g.pieces[5] = &Bishop{Piece{5, 0, false}}
	g.pieces[6] = &Knight{Piece{6, 0, false}}
	g.pieces[7] = &Rook{Piece{7, 0, false}}
	g.pieces[8] = &Pawn{Piece{0, 1, false}}
	g.pieces[9] = &Pawn{Piece{1, 1, false}}
	g.pieces[10] = &Pawn{Piece{2, 1, false}}
	g.pieces[11] = &Pawn{Piece{3, 1, false}}
	g.pieces[12] = &Pawn{Piece{4, 1, false}}
	g.pieces[13] = &Pawn{Piece{5, 1, false}}
	g.pieces[14] = &Pawn{Piece{6, 1, false}}
	g.pieces[15] = &Pawn{Piece{7, 1, false}}
	g.pieces[16] = &Pawn{Piece{0, 6, true}}
	g.pieces[17] = &Pawn{Piece{1, 6, true}}
	g.pieces[18] = &Pawn{Piece{2, 6, true}}
	g.pieces[19] = &Pawn{Piece{3, 6, true}}
	g.pieces[20] = &Pawn{Piece{4, 6, true}}
	g.pieces[21] = &Pawn{Piece{5, 6, true}}
	g.pieces[22] = &Pawn{Piece{6, 6, true}}
	g.pieces[23] = &Pawn{Piece{7, 6, true}}
	g.pieces[24] = &Rook{Piece{0, 7, true}}
	g.pieces[25] = &Knight{Piece{1, 7, true}}
	g.pieces[26] = &Bishop{Piece{2, 7, true}}
	g.pieces[27] = &Queen{Piece{3, 7, true}}
	g.pieces[28] = &King{Piece{4, 7, true}}
	g.pieces[29] = &Bishop{Piece{5, 7, true}}
	g.pieces[30] = &Knight{Piece{6, 7, true}}
	g.pieces[31] = &Rook{Piece{0, 7, true}}
}

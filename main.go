package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text"
	_ "github.com/silbinarywolf/preferdiscretegpu" // Fix for discrete GPUs in windows
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"image/color"
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
// checkmateNotChecked is a sentinel boolean to skip checkmate checks on update cycles if
// the current turn has already been checked. It is set to true every new turn.

type Game struct {
	gameImage           *ebiten.Image
	board               *Board
	uiFont              font.Face
	boardImage          *ebiten.Image
	movingImage         *ebiten.Image
	pieceImage          *ebiten.Image
	uiImage             *ebiten.Image
	pieces              [32]ChessPiece
	selected            [2]float64
	selectedPiece       int
	selectedCol         int
	selectedRow         int
	checkmateNotChecked bool
	gameOver            bool
	gameOverMsg         string
}

const (
	Width    = 1920
	Height   = 1080
	TileSize = 128
	FontDPI  = 72
)

// Update
// Required function by ebitengine. Contains the logic ran every tick of the game.
func (g *Game) Update() error {
	x, y := ebiten.CursorPosition()

	// fancy min max floor math to determine the closest board square to the cursor, even
	// when the mouse is not over the board
	g.selectedCol = int(math.Floor(math.Min(math.Max(float64((x-448)/128), 0), 7)))
	g.selectedRow = int(math.Floor(math.Min(math.Max(float64((y-28)/128), 0), 7)))

	// invert selected row and col when the board is rotated
	if !g.board.whitesTurn {
		g.selectedCol = (g.selectedCol - 7) * -1
		g.selectedRow = (g.selectedRow - 7) * -1
	}

	// No way to exit fullscreen without this for now
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		ebiten.SetFullscreen(false)
	}

	// Checks for a checkmate if in check (only once per turn)
	if g.board.inCheck && g.checkmateNotChecked {
		g.checkmateNotChecked = false
		//function evaluates if checkmate and flags for the game to end if it is
		g.IsCheckmate()
	}

	//TODO: Stalemate

	// left click hold and drag
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		// No piece selected but left mouse is held down
		if g.selectedPiece == -1 {
			// Match the selected tile to a piece location. Then, ensure the piece belongs to the
			// team whose turn it currently is, and that it is still in play.
			for i, piece := range g.pieces {
				if piece.GetCol() == g.selectedCol && piece.GetRow() == g.selectedRow {
					if piece.GetCol() != -1 && g.board.whitesTurn == piece.White() {
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

		//If we do have a piece selected
		if g.selectedPiece != -1 {

			// piece is asking to be let go of at it the current mouse position
			// Verify the move if the piece is being set down on a different square than it started on
			if g.pieces[g.selectedPiece].GetCol() != g.selectedCol || g.pieces[g.selectedPiece].GetRow() != g.selectedRow {
				g.MakeMoveIfLegal(g.selectedRow, g.selectedCol)
			}

			//Either way, we need to update the board image and clear selectedPiece index
			g.board.scheduleDraw = true
			g.selectedPiece = -1
		}
	}

	return nil
}

// MakeMoveIfLegal handles three things: Checking if a move is legal, removing a taken piece from the
// game if the move was legal, and handling the switching of turns.
func (g *Game) MakeMoveIfLegal(row int, col int) {
	//check if move is legal
	//first, make sure the tile it's being set on is possible by comparing it to the Piece's GetMoves function
	//second, don't allow the player to put themselves into check, and see if they are putting their opponent in check
	possibleMoves := g.pieces[g.selectedPiece].GetMoves(g.pieces)
	legal := false

	for _, move := range possibleMoves {
		if move[0] == row && move[1] == col {
			//we found the move in list of possible moves
			legal = true
			break
		}
	}

	if legal {
		//we should save the old piece position then set the new position and make sure the move is still legal
		startingPos := [2]int{g.pieces[g.selectedPiece].GetRow(), g.pieces[g.selectedPiece].GetCol()}
		g.pieces[g.selectedPiece].SetRow(g.selectedRow)
		g.pieces[g.selectedPiece].SetCol(g.selectedCol)

		// If there's a piece on the square we moved to, we need to take it away!
		var capturedPiece *ChessPiece
		capturedOldCol := -1
		for i, piece := range g.pieces {
			if piece.GetRow() == row && piece.GetCol() == col {
				if i != g.selectedPiece {
					capturedOldCol = piece.GetCol()
					piece.SetCol(-1) // Col of -1 is de facto notation for piece taken
					capturedPiece = &piece
					break
				}
			}
		}

		// for each piece on opposing team, does it have possible move to check this player after the move?
		// reminder, a piece with col of -1 has been taken
		for _, piece := range g.pieces {
			if piece.White() != g.board.whitesTurn && piece.GetCol() != -1 {

				//check possible moves for each valid piece and see if any would check the king
				for _, move := range piece.GetMoves(g.pieces) {
					otherPiece := GetPieceOnSquare(move[0], move[1], g.pieces)
					if otherPiece != nil && otherPiece.White() == g.board.whitesTurn && otherPiece.IsKing() == true {
						legal = false
						break
					}
				}
			}
			// idk but it's probably slightly faster to break the loop as soon as move is shown to be illegal
			// honestly it just bothers me to break the inner loop but continue on the outer loop for no reason
			if !legal {
				break
			}
		}

		//using same variable for the condition, yuck! Not going to change it tho :P
		//look a few lines up if this is confusing
		if !legal {
			//put the piece back
			g.pieces[g.selectedPiece].SetRow(startingPos[0])
			g.pieces[g.selectedPiece].SetCol(startingPos[1])

			//put the captured piece back too, if it was taken
			if capturedPiece != nil {
				(*capturedPiece).SetCol(capturedOldCol)
			}

		} else {
			g.board.inCheck = false
			g.checkmateNotChecked = true
			g.board.whitesTurn = !g.board.whitesTurn //switch turns

			//now checking if this move puts the opponent in check
			//note we switched turns in the logic just before this loop
			for _, piece := range g.pieces {
				if piece.White() != g.board.whitesTurn && piece.GetCol() != -1 {

					//check possible moves for each valid piece and see if any would check the king
					for _, move := range piece.GetMoves(g.pieces) {
						otherPiece := GetPieceOnSquare(move[0], move[1], g.pieces)
						if otherPiece != nil && otherPiece.White() == g.board.whitesTurn && otherPiece.IsKing() == true {
							g.board.inCheck = true
							break
						}
					}
				}
			}

		}
	}
}

func (g *Game) IsCheckmate() {
	checkmate := true
	// Try every possible move and see if still in check
	for _, piece := range g.pieces {
		if piece.White() == g.board.whitesTurn && piece.GetCol() != -1 {

			//save the original position so we can put the piece back after checking moves
			startingPos := [2]int{piece.GetRow(), piece.GetCol()}

			//capturedPiece is a placeholder to save pieces that are taken by potential moves
			//we should put it back after running our check
			var capturedPiece *ChessPiece
			capturedOldCol := -1

			for _, move := range piece.GetMoves(g.pieces) {

				//here we simulate each move and see if it gets them out of check
				//if it does, we put the pieces back and exit the loop, indicating it's not checkmate
				//otherwise, we'll put the pieces back and try the next move

				//simulate move
				otherPiece := GetPieceOnSquare(move[0], move[1], g.pieces)
				if otherPiece != nil && otherPiece.White() != g.board.whitesTurn {
					capturedOldCol = otherPiece.GetCol()
					otherPiece.SetCol(-1)
					capturedPiece = &otherPiece
				}
				piece.SetRow(move[0])
				piece.SetCol(move[1])

				//check if this move still leaves them in check
				thisMoveInCheck := false
				for _, nestedPiece := range g.pieces {
					if nestedPiece.White() != g.board.whitesTurn && nestedPiece.GetCol() != -1 {
						//check possible moves for each valid piece and see if any would check the king
						for _, nestedMove := range nestedPiece.GetMoves(g.pieces) {
							otherOtherPiece := GetPieceOnSquare(nestedMove[0], nestedMove[1], g.pieces)
							if otherOtherPiece != nil && otherOtherPiece.White() == g.board.whitesTurn && otherOtherPiece.IsKing() == true {
								thisMoveInCheck = true
								break
							}
						}
						//break out of the upper loop too if check was found for this move
						if thisMoveInCheck {
							break
						}
					}
				}

				//put our pieces back
				piece.SetRow(startingPos[0])
				piece.SetCol(startingPos[1])

				if capturedPiece != nil {
					(*capturedPiece).SetCol(capturedOldCol)
				}

				//if true, we found a move that gets the player out of check
				if !thisMoveInCheck {
					checkmate = false
					break
				}
			}
		}
		if !checkmate {
			break
		}
	}
	if checkmate {
		g.gameOver = true
		g.gameOverMsg = "Checkmate, "
		if g.board.whitesTurn {
			g.gameOverMsg += "Black wins!"
		} else {
			g.gameOverMsg += "White wins!"
		}
	}
}

func (g *Game) Draw(screen *ebiten.Image) {

	g.board.DrawHighlightedTiles(g.gameImage, g.selectedRow, g.selectedCol, g.selectedPiece, g.pieces)
	g.board.DrawUI(g.uiImage, g.gameOver, g.pieces)

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

	//Draw operation settings & execution
	sw, sh := screen.Size()
	bw, bh := g.boardImage.Size()

	uiOp := &ebiten.DrawImageOptions{}
	uiOpX := (sw - bw) / 2
	uiOpY := (sh - bh) / 2

	boardOp := &ebiten.DrawImageOptions{}
	boardOpX := uiOpX
	boardOpY := uiOpY
	boardOpRotate := 0.0

	//flipping the board
	if !g.board.whitesTurn {
		boardOpRotate = math.Pi
		//bring the board back into view after rotating
		boardOpX += bw
		boardOpY += bh
	}

	boardOp.GeoM.Rotate(boardOpRotate)
	boardOp.GeoM.Translate(float64(boardOpX), float64(boardOpY))

	screen.Fill(color.RGBA{R: 0x13, G: 0x33, B: 0x31, A: 0xff})
	screen.DrawImage(g.boardImage, boardOp)
	screen.DrawImage(g.gameImage, boardOp)
	screen.DrawImage(g.pieceImage, boardOp)
	screen.DrawImage(g.uiImage, uiOp)
	screen.DrawImage(g.movingImage, uiOp)

	if g.gameOver {
		text.Draw(screen, g.gameOverMsg, g.uiFont, (sw/2)-len(g.gameOverMsg)*12, sh/2+12, colornames.Darkred)
	}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return Width, Height
}

func main() {
	game := &Game{}
	game.InitBoard()
	ebiten.SetWindowSize(Width, Height)
	ebiten.SetWindowTitle("chess")
	ebiten.SetFullscreen(true)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

func (g *Game) InitBoard() {

	g.board = &Board{}
	g.boardImage = ebiten.NewImage(Width, Height)
	g.gameImage = ebiten.NewImage(Width, Height)
	g.pieceImage = ebiten.NewImage(Width, Height)
	g.movingImage = ebiten.NewImage(Width, Height)
	g.uiImage = ebiten.NewImage(Width, Height)

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
	g.pieces[31] = &Rook{Piece{7, 7, true}}

	//Attempt to load font
	tt, err := opentype.Parse(fonts.PressStart2P_ttf)
	if err != nil {
		log.Fatal(err)
	}

	g.uiFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     FontDPI,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		log.Fatal(err)
	}

	g.checkmateNotChecked = true
	g.gameOver = false
	g.gameOverMsg = ""
	g.board.inCheck = false
	g.board.whitesTurn = true

	g.board.DrawBoard(g.boardImage)
	g.board.DrawUI(g.uiImage, g.gameOver, g.pieces)
	g.board.DrawStaticPieces(g.pieceImage, g.pieces, g.selectedPiece)
	g.board.scheduleDraw = false
}

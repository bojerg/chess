package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text"
	_ "github.com/silbinarywolf/preferdiscretegpu" // Fix for discrete GPUs in windows
	"golang.org/x/image/colornames"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"image/color"
	"log"
	"math"
	"path/filepath"
	"sort"
)

// Game
// gameType indicates the selected game mode. -1 = main menu, 0 = local multiplayer, 1 = versus a bot.
// gameImage, among the other image variables, are for rendering various "layers" of the game.
// scheduleDraw is a sentinel value to indicate when static images need to be refreshed.
// checkmateNotChecked is false until we evaluate if the previous move ends the game.
// selectedLocations is for the x, y values of a piece in motion.
// selectedPiece is the index of the selected piece (-1 indicates none selected).
// selectedCol, selectedRow is the hovered over/selected board square.
// The unmentioned variables seem straightforward enough.
type Game struct {
	gameType            int
	gameImage           *ebiten.Image
	boardImage          *ebiten.Image
	movingImage         *ebiten.Image
	pieceImage          *ebiten.Image
	uiImage             *ebiten.Image
	menuBgImage         *ebiten.Image
	pieces              [32]ChessPiece
	scheduleDraw        bool
	whitesTurn          bool
	inCheck             bool
	checkmateNotChecked bool
	selectedLocation    [2]float64
	selectedPiece       int
	selectedCol         int
	selectedRow         int
	gameOver            bool
	gameOverMsg         string
	uiFontBig           font.Face
	uiFont              font.Face
	uiFontSmall         font.Face
	btnHoverIndex       int
	btnPrimary          *ebiten.Image
	btnPrimaryHover     *ebiten.Image
	btnInfo             *ebiten.Image
	btnInfoHover        *ebiten.Image
}

const (
	Width    = 1920
	Height   = 1080
	TileSize = 128
	FontDPI  = 72
)

// Draw
// Draws stuff. Required by ebitengine.
func (g *Game) Draw(screen *ebiten.Image) {

	screen.Fill(color.RGBA{R: 0x13, G: 0x33, B: 0x31, A: 0xff})

	switch g.gameType {
	case -1:
		g.DrawMainMenu(false) //Prints to g.uiImage and g.menuBgImage

		//Using selectedCol as a counter for infinite scroll of the background
		g.selectedCol = (g.selectedCol + 1) % 50
		opMenuBg := &ebiten.DrawImageOptions{}
		opMenuBg.GeoM.Translate(float64(-g.selectedCol*2), float64(-150+g.selectedCol*2))
		opMenuBg.GeoM.Scale(1.3, 1.3)
		screen.DrawImage(g.menuBgImage, opMenuBg)

		screen.DrawImage(g.uiImage, &ebiten.DrawImageOptions{})
		text.Draw(screen, "Chess", g.uiFontBig, 770, 432, colornames.White)
		text.Draw(screen, "by bojerg", g.uiFont, 960, 432, colornames.Whitesmoke)

	default:
		//play the game
		g.movingImage.Clear()
		g.DrawHighlightedTiles()
		g.DrawUI()

		if g.selectedPiece != -1 {
			g.DrawMovingPiece()
		}

		// if game logic signals that piece locations have changed, then DrawStaticPieces, and...
		// stop scheduling draw because we updated the static piece image
		if g.scheduleDraw {
			g.DrawStaticPieces()
			g.scheduleDraw = false
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
		if !g.whitesTurn && g.gameType == 1 {
			boardOpRotate = math.Pi
			//bring the board back into view after rotating
			boardOpX += bw
			boardOpY += bh
		}

		boardOp.GeoM.Rotate(boardOpRotate)
		boardOp.GeoM.Translate(float64(boardOpX), float64(boardOpY))

		screen.DrawImage(g.boardImage, boardOp)
		screen.DrawImage(g.gameImage, boardOp)
		screen.DrawImage(g.pieceImage, boardOp)
		screen.DrawImage(g.uiImage, uiOp)
		screen.DrawImage(g.movingImage, uiOp)

		if g.gameOver {
			text.Draw(screen, g.gameOverMsg, g.uiFont, (sw/2)-len(g.gameOverMsg)*12, sh/2+12, colornames.Darkred)
		}
	}
}

// Update
// Required function by ebitengine. Contains the logic ran every tick of the game.
func (g *Game) Update() error {
	x, y := ebiten.CursorPosition()

	// No way to exit fullscreen without this for now
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		ebiten.SetFullscreen(false)
	}

	switch g.gameType {
	case -1:
		//at main menu

		//Indicating which controls are hovered over
		if x > Width/2-90 && x < Width/2+90 && y > Height/2-40 && y < Height/2+40 {
			g.btnHoverIndex = 1
		} else {
			g.btnHoverIndex = -1
		}

		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			if g.btnHoverIndex != -1 {
				g.gameType = g.btnHoverIndex
				g.InitPiecesAndImages()
				g.btnHoverIndex = -1
				break
			}
		}

	default:
		//playing the game

		// Checks for a checkmate if in check (only once per turn)
		if g.inCheck && g.checkmateNotChecked {
			g.checkmateNotChecked = false
			//function evaluates if checkmate and flags for the game to end if it is
			g.IsCheckmate()
		}

		//TODO: Stalemate

		// XY locations reflect the two buttons drawn on screen
		// This code block determines what the mouse is interacting with and updates the appropriate parameter
		if x > TileSize+36 && x < TileSize+216 && y > TileSize*2.5 && y < TileSize*2.5+80 {
			g.btnHoverIndex = 1
		} else if x > TileSize+36 && x < TileSize+216 && y > TileSize*4.5 && y < TileSize*4.5+80 {
			g.btnHoverIndex = 2
		} else {
			g.btnHoverIndex = -1
			// fancy min max floor math to determine the closest board square to the cursor, even
			// when the mouse is not over the board
			g.selectedCol = int(math.Floor(math.Min(math.Max(float64((x-448)/128), 0), 7)))
			g.selectedRow = int(math.Floor(math.Min(math.Max(float64((y-28)/128), 0), 7)))

			// invert selected row and col when the board is rotated
			if !g.whitesTurn && g.gameType == 1 {
				g.selectedCol = (g.selectedCol - 7) * -1
				g.selectedRow = (g.selectedRow - 7) * -1
			}
		}

		// left click hold and drag
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {

			if g.btnHoverIndex != -1 {

				if g.btnHoverIndex == 1 {
					//TODO: Return to menu
				} else if g.btnHoverIndex == 2 {
					//TODO: Start new game of same type
				}

			} else {
				//Not clicking on a button...

				if g.selectedPiece == -1 {
					// No piece selected but left mouse is held down
					// Match the selected tile to a piece location. Then, ensure the piece belongs to the
					// team whose turn it currently is, and that it is still in play.
					for i, piece := range g.pieces {
						if piece.GetCol() == g.selectedCol && piece.GetRow() == g.selectedRow {
							if piece.GetCol() != -1 && g.whitesTurn == piece.White() {
								g.selectedPiece = i
								// store the xy coordinates of the cursor
								g.selectedLocation[0] = float64(x)/1.5 - 30
								g.selectedLocation[1] = float64(y)/1.5 - 30
								g.scheduleDraw = true
								break
							}
						}
					}
				} else {
					// update current mouse position because piece is still selected
					// and the mouse may be moving!
					g.selectedLocation[0] = float64(x)/1.5 - 30
					g.selectedLocation[1] = float64(y)/1.5 - 30
				}
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
				g.scheduleDraw = true
				g.selectedPiece = -1
			}
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
			if piece.White() != g.whitesTurn && piece.GetCol() != -1 {

				//check possible moves for each valid piece and see if any would check the king
				for _, move := range piece.GetMoves(g.pieces) {
					otherPiece := GetPieceOnSquare(move[0], move[1], g.pieces)
					if otherPiece != nil && otherPiece.White() == g.whitesTurn && otherPiece.IsKing() == true {
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
			g.inCheck = false
			g.checkmateNotChecked = true
			g.whitesTurn = !g.whitesTurn //switch turns

			//now checking if this move puts the opponent in check
			//note we switched turns in the logic just before this loop
			for _, piece := range g.pieces {
				if piece.White() != g.whitesTurn && piece.GetCol() != -1 {

					//check possible moves for each valid piece and see if any would check the king
					for _, move := range piece.GetMoves(g.pieces) {
						otherPiece := GetPieceOnSquare(move[0], move[1], g.pieces)
						if otherPiece != nil && otherPiece.White() == g.whitesTurn && otherPiece.IsKing() == true {
							g.inCheck = true
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
		if piece.White() == g.whitesTurn && piece.GetCol() != -1 {

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
				if otherPiece != nil && otherPiece.White() != g.whitesTurn {
					capturedOldCol = otherPiece.GetCol()
					otherPiece.SetCol(-1)
					capturedPiece = &otherPiece
				}
				piece.SetRow(move[0])
				piece.SetCol(move[1])

				//check if this move still leaves them in check
				thisMoveInCheck := false
				for _, nestedPiece := range g.pieces {
					if nestedPiece.White() != g.whitesTurn && nestedPiece.GetCol() != -1 {
						//check possible moves for each valid piece and see if any would check the king
						for _, nestedMove := range nestedPiece.GetMoves(g.pieces) {
							otherOtherPiece := GetPieceOnSquare(nestedMove[0], nestedMove[1], g.pieces)
							if otherOtherPiece != nil && otherOtherPiece.White() == g.whitesTurn && otherOtherPiece.IsKing() == true {
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
		if g.whitesTurn {
			g.gameOverMsg += "Black wins!"
		} else {
			g.gameOverMsg += "White wins!"
		}
	}
}

func (g *Game) DrawStaticPieces() {
	g.pieceImage.Clear()

	xOffset := 465.0
	yOffset := 42.0
	rotate := 0.0

	// Rotate the board for local multiplayer (gameType 1)
	if !g.whitesTurn && g.gameType == 1 {
		rotate = math.Pi
		xOffset += TileSize - 34
		yOffset += TileSize - 28
	}

	for i, piece := range g.pieces {
		// Don't draw selected (moving) piece, or any pieces with id of 6 (taken)
		if i != g.selectedPiece && piece.GetCol() != -1 {
			tx := float64(g.pieces[i].GetCol()*TileSize) + xOffset
			ty := float64(g.pieces[i].GetRow()*TileSize) + yOffset
			opPiece := &ebiten.DrawImageOptions{}
			opPiece.GeoM.Rotate(rotate)
			opPiece.GeoM.Scale(1.5, 1.5) //essentially W x H = 90 x 90
			opPiece.GeoM.Translate(tx, ty)
			g.pieceImage.DrawImage(g.pieces[i].GetImage(), opPiece)
		}
	}
}

func (g *Game) DrawMovingPiece() {
	for i, _ := range g.pieces {
		if i == g.selectedPiece {
			tx := g.selectedLocation[0] * 1.5
			ty := g.selectedLocation[1] * 1.5
			opPiece := &ebiten.DrawImageOptions{}
			opPiece.GeoM.Scale(1.5, 1.5) //essentially W x H = 90 x 90
			opPiece.GeoM.Translate(tx, ty)
			g.movingImage.DrawImage(g.pieces[i].GetImage(), opPiece)
			break
		}
	}
}

func (g *Game) DrawHighlightedTiles() {
	tileImage := ebiten.NewImage(TileSize, TileSize)
	g.gameImage.Clear()

	// drawing highlighted tiles (available moves in red)
	if g.selectedPiece >= 0 {
		availableMoves := g.pieces[g.selectedPiece].GetMoves(g.pieces)
		if availableMoves != nil {
			for _, move := range availableMoves {
				opTile := &ebiten.DrawImageOptions{}
				opTile.GeoM.Translate(float64(move[1]*TileSize+448), float64(move[0]*TileSize+28))
				tileImage.Fill(color.RGBA{R: 0xff, G: 0x06, B: 0x03, A: 0xba})
				g.gameImage.DrawImage(tileImage, opTile)
			}
		}

	}

	// Draw hovered tile (in highlighter yellow) if not hovering a button
	if g.btnHoverIndex == -1 {
		for r := 0; r < 8; r++ {
			for c := 0; c < 8; c++ {
				if r == g.selectedRow && c == g.selectedCol {
					opTile := &ebiten.DrawImageOptions{}
					opTile.GeoM.Translate(float64(c*TileSize+448), float64(r*TileSize+28))
					tileImage.Fill(color.RGBA{R: 0xea, G: 0xdd, B: 0x23, A: 0xff})
					g.gameImage.DrawImage(tileImage, opTile)
					break
				}
			}
		}
	}

	//highlight a king in check (purple)
	if g.inCheck {
		for _, piece := range g.pieces {
			if piece.IsKing() && piece.White() == g.whitesTurn {
				opTile := &ebiten.DrawImageOptions{}
				opTile.GeoM.Translate(float64(piece.GetCol()*TileSize+448), float64(piece.GetRow()*TileSize+28))
				tileImage.Fill(color.RGBA{R: 0xbf, G: 0x00, B: 0xe6, A: 0xff})
				g.gameImage.DrawImage(tileImage, opTile)
				break
			}
		}
	}

}

func (g *Game) DrawBoard() {
	darkColor := color.RGBA{R: 0xbb, G: 0x99, B: 0x55, A: 0xff}
	lightColor := color.RGBA{R: 0xcb, G: 0xbe, B: 0xb5, A: 0xff}

	lightImage := ebiten.NewImage(TileSize*8, TileSize*8)
	darkImage := ebiten.NewImage(TileSize, TileSize)

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

func (g *Game) DrawUI() {
	g.uiImage.Clear()

	//Arranging taken pieces into two structs to sort by value and team for display
	var whitePieces []ChessPiece
	var blackPieces []ChessPiece
	for _, piece := range g.pieces {
		if piece.GetCol() == -1 {
			if piece.White() {
				whitePieces = append(whitePieces, piece)
			} else {
				blackPieces = append(blackPieces, piece)
			}
		}
	}

	sort.Slice(whitePieces, func(p, q int) bool {
		return GetWeighting(whitePieces[p]) < GetWeighting(whitePieces[q])
	})

	sort.Slice(blackPieces, func(p, q int) bool {
		return GetWeighting(blackPieces[p]) < GetWeighting(blackPieces[q])
	})

	//The following offsets and modifiers help to dynamically grow the column of taken pieces and flip them
	//as the board is flipped
	whiteXOffset := 384
	blackXOffset := 448 + TileSize*8
	var whiteYOffset float64 = 28
	var blackYOffset float64 = TileSize * 8
	whiteGrowth := -16
	blackGrowth := 16

	// Rotate the board for local multiplayer (gameType 1)
	if !g.whitesTurn && g.gameType == 1 {
		whiteXOffset = blackXOffset
		whiteYOffset = blackYOffset
		blackXOffset = 384
		blackYOffset = 28
		whiteGrowth *= -1
		blackGrowth *= -1
	}

	//Draw the lists of taken piece images
	for i, p := range whitePieces {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(0.7, 0.7)
		op.GeoM.Translate(float64((len(whitePieces)-i)*whiteGrowth+whiteXOffset), whiteYOffset)
		g.uiImage.DrawImage(p.GetImage(), op)
	}

	for i, p := range blackPieces {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(0.7, 0.7)
		op.GeoM.Translate(float64((len(blackPieces)-i)*blackGrowth+blackXOffset), blackYOffset)
		g.uiImage.DrawImage(p.GetImage(), op)
	}

	opMenuBtn := &ebiten.DrawImageOptions{}
	opMenuBtn.GeoM.Translate(TileSize+36, TileSize*2.5)
	if g.btnHoverIndex == 1 {
		g.uiImage.DrawImage(g.btnPrimaryHover, opMenuBtn)
		text.Draw(g.uiImage, "Main Menu", g.uiFontSmall, TileSize+70, TileSize*2.5+50, colornames.Gray)
	} else {
		g.uiImage.DrawImage(g.btnPrimary, opMenuBtn)
		text.Draw(g.uiImage, "Main Menu", g.uiFontSmall, TileSize+70, TileSize*2.5+50, colornames.Whitesmoke)
	}

	opMenuBtn.GeoM.Translate(0, TileSize*2)
	if g.btnHoverIndex == 2 {
		g.uiImage.DrawImage(g.btnInfoHover, opMenuBtn)
		text.Draw(g.uiImage, "New Game", g.uiFontSmall, TileSize+74, TileSize*4.5+50, colornames.Gray)
	} else {
		g.uiImage.DrawImage(g.btnInfo, opMenuBtn)
		text.Draw(g.uiImage, "New Game", g.uiFontSmall, TileSize+74, TileSize*4.5+50, colornames.Whitesmoke)
	}

}

func (g *Game) DrawMainMenu(generate bool) {
	g.uiImage.Clear()
	//We will draw a scrolling background of chess pieces and place button images on top
	//Generate the image once to significantly improve performance and thus appearance
	if generate {
		g.selectedCol = 0 //reset this because we use it as a counter for scrolling effect
		g.pieces[0] = &Pawn{Piece{0, 0, false}}
		g.pieces[1] = &Pawn{Piece{0, 0, true}}
		g.pieces[2] = &Rook{Piece{0, 0, false}}
		g.pieces[3] = &Knight{Piece{0, 0, true}}
		g.pieces[4] = &Bishop{Piece{0, 0, false}}
		g.pieces[5] = &Queen{Piece{0, 0, true}}
		g.pieces[6] = &King{Piece{0, 0, false}}
		g.pieces[7] = &Bishop{Piece{0, 0, true}}
		g.pieces[8] = &Knight{Piece{0, 0, false}}
		g.pieces[9] = &Rook{Piece{0, 0, true}}

		for y := 0; y < 24; y++ {
			for x := 0; x < 24; x++ {
				opPiece := &ebiten.DrawImageOptions{}
				opPiece.GeoM.Scale(1.8, 1.8)
				opPiece.GeoM.Translate(float64(x*100), float64(y*100))
				opPiece.ColorM.Translate(0, 0, 0, -.7)
				g.menuBgImage.DrawImage(g.pieces[(x+y)%10].GetImage(), opPiece)
			}
		}

	}

	opButton := &ebiten.DrawImageOptions{}
	opButton.GeoM.Translate(Width/2-90, Height/2-40)
	if g.btnHoverIndex == 1 {
		g.uiImage.DrawImage(g.btnPrimaryHover, opButton)
		text.Draw(g.uiImage, "Local Match", g.uiFontSmall, Width/2-76, Height/2+8, colornames.Whitesmoke)
	} else {
		g.uiImage.DrawImage(g.btnPrimary, opButton)
		text.Draw(g.uiImage, "Local Match", g.uiFontSmall, Width/2-76, Height/2+8, colornames.Gray)
	}

	opButton.GeoM.Translate(0, 110)
	opButton.ColorM.Translate(-.1, -.1, -.1, -.5)
	g.uiImage.DrawImage(g.btnPrimary, opButton)
	text.Draw(g.uiImage, "Versus Bot", g.uiFontSmall, Width/2-69, Height/2+118, colornames.Gray)

}

func (g *Game) InitPiecesAndImages() {
	g.selectedPiece = -1
	g.selectedLocation[0] = 0.0
	g.selectedLocation[1] = 0.0
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

	g.checkmateNotChecked = true
	g.gameOver = false
	g.gameOverMsg = ""
	g.inCheck = false
	g.whitesTurn = true

	g.DrawBoard()
	g.DrawStaticPieces()
	g.scheduleDraw = false
}

func (g *Game) InitGame() {
	g.gameType = -1

	g.boardImage = ebiten.NewImage(Width, Height)
	g.gameImage = ebiten.NewImage(Width, Height)
	g.pieceImage = ebiten.NewImage(Width, Height)
	g.movingImage = ebiten.NewImage(Width, Height)
	g.uiImage = ebiten.NewImage(Width, Height)
	g.menuBgImage = ebiten.NewImage(Width, Height)

	//Attempt to load font
	tt, err := opentype.Parse(fonts.PressStart2P_ttf)
	if err != nil {
		log.Fatal(err)
	}

	g.uiFontBig, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    34,
		DPI:     FontDPI,
		Hinting: font.HintingVertical,
	})
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

	g.uiFontSmall, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    14,
		DPI:     FontDPI,
		Hinting: font.HintingVertical,
	})
	if err != nil {
		log.Fatal(err)
	}

	fileLoc, _ := filepath.Abs("images/btnPrimary.png")
	g.btnPrimary, _, err = ebitenutil.NewImageFromFile(fileLoc)
	if err != nil {
		return
	}
	if err != nil {
		log.Fatal(err)
	}

	fileLoc, _ = filepath.Abs("images/btnPrimaryHover.png")
	g.btnPrimaryHover, _, err = ebitenutil.NewImageFromFile(fileLoc)
	if err != nil {
		log.Fatal(err)
	}

	fileLoc, _ = filepath.Abs("images/btnInfo.png")
	g.btnInfo, _, err = ebitenutil.NewImageFromFile(fileLoc)
	if err != nil {
		log.Fatal(err)
	}

	fileLoc, _ = filepath.Abs("images/btnInfoHover.png")
	g.btnInfoHover, _, err = ebitenutil.NewImageFromFile(fileLoc)
	if err != nil {
		log.Fatal(err)
	}

	g.DrawMainMenu(true)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return Width, Height
}

func main() {
	game := &Game{}
	game.InitGame()
	ebiten.SetWindowSize(Width, Height)
	ebiten.SetWindowTitle("chess")
	ebiten.SetFullscreen(true)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}

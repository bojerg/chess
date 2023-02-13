package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	_ "image/png" // required for ebitenutil/NewImageFromFile
	"log"
	"path/filepath"
)

// Piece
// Used for each piece class to inherit from.
// col, row: identifies location on chess board.
// white: If on team white, is true.
type Piece struct {
	col   int
	row   int
	white bool
}

// ChessPiece is the collection of all pieces and their shared functionality
type ChessPiece interface {
	GetCol() int
	SetCol(int)
	GetRow() int
	SetRow(int)
	White() bool
	GetImage() *ebiten.Image
	GetMoves([32]ChessPiece) [][2]int
	GetName() string
}

// GetImage returns the corresponding ebiten image from filepath argument
func GetImage(filepathStr string) *ebiten.Image {
	// https://commons.wikimedia.org/wiki/Category:PNG_chess_pieces/Standard_transparent
	var err error = nil
	// path/filepath creates filepath for any OS, supposedly
	fileLoc, _ := filepath.Abs(filepathStr)
	pieceImage, _, err := ebitenutil.NewImageFromFile(fileLoc)
	if err != nil {
		log.Fatal(err)
		return nil
	} else {
		return pieceImage
	}

}

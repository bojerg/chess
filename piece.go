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
	Col() int
	SetCol(int)
	Row() int
	SetRow(int)
	White() bool
	Image() *ebiten.Image
	Moves(Game) [][2]int
	Name() string
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

func IsPawn(piece ChessPiece) bool {
	return piece.Name()[6:] == "pawn"
}

func IsKing(piece ChessPiece) bool {
	return piece.Name()[6:] == "king"
}

func IsInBounds(row int, col int) bool {
	return row <= 7 && row >= 0 && col <= 7 && col >= 0
}

// GetPieceOnSquare checks list of ChessPiece against provided row and col positions. If a match
// is found, a copy of that piece is returned. Otherwise, returns nil.
func GetPieceOnSquare(row int, col int, pieces [32]ChessPiece) ChessPiece {
	for _, piece := range pieces {
		if piece.Col() == col && piece.Row() == row {
			return piece
		}
	}
	return nil
}

// GetWeighting is implemented so that we can sort pieces by value for UI purposes
func GetWeighting(piece ChessPiece) int {
	//	ex. piece.Name() >> "White pawn"
	switch piece.Name()[6:] {
	case "pawn":
		return 1
	case "knight":
		return 2
	case "bishop":
		return 3
	case "rook":
		return 4
	case "queen":
		return 5
	default:
		return 0
	}
}

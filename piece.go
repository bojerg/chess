package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	_ "image/png" // required for ebitenutil/NewImageFromFile
	"log"
	"path/filepath"
	"strings"
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
	IsKing() bool
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

func IsInBounds(row int, col int) bool {
	return row <= 7 && row >= 0 && col <= 7 && col >= 0
}

// GetPieceOnSquare checks list of ChessPiece against provided row and col positions. If a match
// is found, a copy of that piece is returned. Otherwise, returns nil.
func GetPieceOnSquare(row int, col int, pieces [32]ChessPiece) ChessPiece {
	for _, piece := range pieces {
		if piece.GetCol() == col && piece.GetRow() == row {
			return piece
		}
	}
	return nil
}

// GetWeighting is implemented so that we can sort pieces by value for UI purposes
//
//	Splitting the name after the space to get the piece type  ex. GetName() = "White pawn"
func GetWeighting(piece ChessPiece) int {
	switch strings.Split(piece.GetName(), " ")[1] {
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

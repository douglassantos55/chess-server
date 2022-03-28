package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Board struct {
	matrix [8]map[string]Piece
}

type Piece struct {
	Type  string
	Color Color
}

func CreatePiece(name string, color Color) Piece {
	return Piece{Type: name, Color: color}
}

func (p *Piece) String() string {
	if p.Type == "" {
		return ""
	}
	return fmt.Sprintf("%s %s", p.Color, p.Type)
}

func Empty() Piece {
	return CreatePiece("", White)
}

func Rook(color Color) Piece {
	return CreatePiece("R", color)
}
func Knight(color Color) Piece {
	return CreatePiece("N", color)
}
func Bishop(color Color) Piece {
	return CreatePiece("B", color)
}
func Queen(color Color) Piece {
	return CreatePiece("Q", color)
}
func King(color Color) Piece {
	return CreatePiece("K", color)
}
func Pawn(color Color) Piece {
	return CreatePiece("p", color)
}

func NewBoard() *Board {
	return &Board{
		matrix: [8]map[string]Piece{
			{
				"a": Rook(White),
				"b": Knight(White),
				"c": Bishop(White),
				"d": Queen(White),
				"e": King(White),
				"f": Bishop(White),
				"g": Knight(White),
				"h": Rook(White),
			},
			{
				"a": Pawn(White),
				"b": Pawn(White),
				"c": Pawn(White),
				"d": Pawn(White),
				"e": Pawn(White),
				"f": Pawn(White),
				"g": Pawn(White),
				"h": Pawn(White),
			},
			{"a": Empty(), "b": Empty(), "c": Empty(), "d": Empty(), "e": Empty(), "f": Empty(), "g": Empty(), "h": Empty()},
			{"a": Empty(), "b": Empty(), "c": Empty(), "d": Empty(), "e": Empty(), "f": Empty(), "g": Empty(), "h": Empty()},
			{"a": Empty(), "b": Empty(), "c": Empty(), "d": Empty(), "e": Empty(), "f": Empty(), "g": Empty(), "h": Empty()},
			{"a": Empty(), "b": Empty(), "c": Empty(), "d": Empty(), "e": Empty(), "f": Empty(), "g": Empty(), "h": Empty()},
			{
				"a": Pawn(Black),
				"b": Pawn(Black),
				"c": Pawn(Black),
				"d": Pawn(Black),
				"e": Pawn(Black),
				"f": Pawn(Black),
				"g": Pawn(Black),
				"h": Pawn(Black),
			},
			{
				"a": Rook(Black),
				"b": Knight(Black),
				"c": Bishop(Black),
				"d": Queen(Black),
				"e": King(Black),
				"f": Bishop(Black),
				"g": Knight(Black),
				"h": Rook(Black),
			},
		},
	}
}

func (b *Board) Square(square string) Piece {
	row, col := b.parseSquare(square)
	return b.matrix[row][col]
}

func (b *Board) parseSquare(square string) (int64, string) {
	coord := strings.Split(square, "")
	x, _ := strconv.ParseInt(coord[1], 10, 64)

	return x - 1, coord[0]
}

func (b *Board) Move(from, to string) {
	// TODO: check if it's a valid move
	piece := b.Square(from)

	destRow, destCol := b.parseSquare(to)
	sourceRow, sourceCol := b.parseSquare(from)

	b.matrix[destRow][destCol] = piece
	b.matrix[sourceRow][sourceCol] = Empty()
}

package main

import (
	"fmt"
	"strconv"
	"strings"
)

func parseSquare(square string) (int, rune) {
	coord := strings.Split(square, "")
	x, _ := strconv.ParseInt(coord[1], 10, 64)

	return int(x - 1), []rune(coord[0])[0]
}

type Board struct {
	matrix [8]map[rune]Piece
}

func NewBoard() *Board {
	return &Board{
		matrix: [8]map[rune]Piece{
			{
				'a': Rook(White),
				'b': Knight(White),
				'c': Bishop(White),
				'd': Queen(White),
				'e': King(White),
				'f': Bishop(White),
				'g': Knight(White),
				'h': Rook(White),
			},
			{
				'a': Pawn(White),
				'b': Pawn(White),
				'c': Pawn(White),
				'd': Pawn(White),
				'e': Pawn(White),
				'f': Pawn(White),
				'g': Pawn(White),
				'h': Pawn(White),
			},
			{'a': Empty(), 'b': Empty(), 'c': Empty(), 'd': Empty(), 'e': Empty(), 'f': Empty(), 'g': Empty(), 'h': Empty()},
			{'a': Empty(), 'b': Empty(), 'c': Empty(), 'd': Empty(), 'e': Empty(), 'f': Empty(), 'g': Empty(), 'h': Empty()},
			{'a': Empty(), 'b': Empty(), 'c': Empty(), 'd': Empty(), 'e': Empty(), 'f': Empty(), 'g': Empty(), 'h': Empty()},
			{'a': Empty(), 'b': Empty(), 'c': Empty(), 'd': Empty(), 'e': Empty(), 'f': Empty(), 'g': Empty(), 'h': Empty()},
			{
				'a': Pawn(Black),
				'b': Pawn(Black),
				'c': Pawn(Black),
				'd': Pawn(Black),
				'e': Pawn(Black),
				'f': Pawn(Black),
				'g': Pawn(Black),
				'h': Pawn(Black),
			},
			{
				'a': Rook(Black),
				'b': Knight(Black),
				'c': Bishop(Black),
				'd': Queen(Black),
				'e': King(Black),
				'f': Bishop(Black),
				'g': Knight(Black),
				'h': Rook(Black),
			},
		},
	}
}

func (b *Board) Square(square string) Piece {
	row, col := parseSquare(square)
	return b.matrix[row][col]
}

func (b *Board) IsThreatned(square string, color Color) bool {
	for i := int('a'); i <= int('h'); i++ {
		for j := 1; j <= 8; j++ {
			piece := b.matrix[j-1][rune(i)]
			from := fmt.Sprintf("%s%d", string(rune(i)), j)

			if piece != Empty() && piece.Color != color {
				if piece.Sees(from, square, b) {
					return true
				}
			}
		}
	}
	return false
}

func (b *Board) Move(from, to string) {
	piece := b.Square(from)

	if piece.Move(from, to, b) {
		destRow, destCol := parseSquare(to)
		sourceRow, sourceCol := parseSquare(from)

		b.matrix[destRow][destCol] = piece
		b.matrix[sourceRow][sourceCol] = Empty()
	}
}

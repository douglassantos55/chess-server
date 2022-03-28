package main

import (
	"strconv"
	"strings"
)

type Board struct {
	matrix [8]map[string]string
}

func NewBoard() *Board {
	return &Board{
		matrix: [8]map[string]string{
			{"a": "R", "b": "N", "c": "B", "d": "Q", "e": "K", "f": "B", "g": "N", "h": "R"},
			{"a": "p", "b": "p", "c": "p", "d": "p", "e": "p", "f": "p", "g": "p", "h": "p"},
			{"a": "", "b": "", "c": "", "d": "", "e": "", "f": "", "g": "", "h": ""},
			{"a": "", "b": "", "c": "", "d": "", "e": "", "f": "", "g": "", "h": ""},
			{"a": "", "b": "", "c": "", "d": "", "e": "", "f": "", "g": "", "h": ""},
			{"a": "", "b": "", "c": "", "d": "", "e": "", "f": "", "g": "", "h": ""},
			{"a": "p", "b": "p", "c": "p", "d": "p", "e": "p", "f": "p", "g": "p", "h": "p"},
			{"a": "R", "b": "N", "c": "B", "d": "Q", "e": "K", "f": "B", "g": "N", "h": "R"},
		},
	}
}

func (b *Board) Square(square string) string {
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
	b.matrix[sourceRow][sourceCol] = ""
}

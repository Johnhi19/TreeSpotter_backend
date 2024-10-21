package models

type Position struct {
	X int
	Y int
}

type Tree struct {
	ID       int
	Type     string
	Age      int
	Position Position
}

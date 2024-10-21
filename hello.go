package main

import (
	"fmt"

	"github.com/Johnhi19/TreeSpotter_backend/models"
)

func main() {
	newTree := models.Tree{
		Type:     "Oak",
		Age:      50,
		Position: models.Position{X: 1, Y: 2},
	}

	fmt.Printf("%+v\n", newTree)
}

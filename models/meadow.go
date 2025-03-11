package models

import (
	"go.mongodb.org/mongo-driver/bson"
)

type Meadow struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	TreeIds  []int  `json:"treeIds"`
	Size     []int  `json:"size"`
	Location string `json:"location"`
}

func TransformMeadowToBson(meadow Meadow) bson.D {
	return bson.D{
		{Key: "ID", Value: meadow.ID},
		{Key: "Name", Value: meadow.Name},
		{Key: "TreeIds", Value: meadow.TreeIds},
		{Key: "Size", Value: meadow.Size},
		{Key: "Location", Value: meadow.Location},
	}

}

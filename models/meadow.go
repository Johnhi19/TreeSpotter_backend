package models

import "go.mongodb.org/mongo-driver/bson"

type Meadow struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	NumberTrees int    `json:"numberTrees"`
}

func TransformToBson(meadow Meadow) bson.D {
	return bson.D{
		{"ID", meadow.ID},
		{"Name", meadow.Name},
		{"NumberTrees", meadow.NumberTrees},
	}

}

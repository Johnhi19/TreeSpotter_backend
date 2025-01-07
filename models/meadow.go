package models

import "go.mongodb.org/mongo-driver/bson"

type Meadow struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	NumberTrees int    `json:"numberTrees"`
}

func TransformMeadowToBson(meadow Meadow) bson.D {
	return bson.D{
		{Key: "ID", Value: meadow.ID},
		{Key: "Name", Value: meadow.Name},
		{Key: "NumberTrees", Value: meadow.NumberTrees},
	}

}

package models

import "go.mongodb.org/mongo-driver/bson"

type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Tree struct {
	ID       int      `json:"id"`
	Type     string   `json:"type"`
	Age      int      `json:"age"`
	MeadowId int      `json:"meadowId"`
	Position Position `json:"position"`
}

func TransformTreeToBson(tree Tree) bson.D {
	return bson.D{
		{Key: "ID", Value: tree.ID},
		{Key: "Type", Value: tree.Type},
		{Key: "Age", Value: tree.Age},
		{Key: "MeadowId", Value: tree.MeadowId},
		{Key: "Position", Value: bson.D{
			{Key: "X", Value: tree.Position.X},
			{Key: "Y", Value: tree.Position.Y}}},
	}

}

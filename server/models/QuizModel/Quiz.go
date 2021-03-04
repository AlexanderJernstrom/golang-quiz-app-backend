package Quiz

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)


type Alternative struct {
	ID   primitive.ObjectID `json:"_id,omitempty" bson:"_id, omitempty"`
	Text string             `json:"text,omitempty" bson:"string, omitempty"`
}

type Question struct {
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Text         string             `json:"text,omitempty" bson:"text,omitempty"`
	Answer       string             `json:"answer,omitempty" bson:"answer,omitempty"`
	Alternatives []*Alternative
}

type PlayedBy struct {
	Name string `json:"name,omitempty" bson:"name,omitempty"`
	RightAmount int `json:"rightAmount,omitEmpty" bson:"rightAmount,omitempty"`
	Date time.Time `json:"date,omitempty" bson:"date,omitempty"`
}


type Quiz struct {
	Name      string `json:"name,omitempty" bson:"name,omitempty"`
	Questions []*Question
	ID           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	CreatedBy 		string             `json:"createdBy,omitempty" bson:"createdBy,omitempty"`
}
package UserModel

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Name           *string `json:"name" bson:"name,omitempty"`
	Email          *string `json:"email" bson:"email,omitempty"`
	Password       *string `json:"password" bson:"password,omitempty"`
	CreatedQuizzes *[]primitive.ObjectID `json:"createdQuizzes" bson:"createdQuizzes,omitempty"`
	PlayedQuizzes *[]primitive.ObjectID `json:"playedQuizzes" bson:"playedQuizzes,omitempty"`
	ID *primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
}
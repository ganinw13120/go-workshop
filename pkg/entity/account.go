package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type Account struct {
	Id              primitive.ObjectID
	DisplayName     string
	Username        string
	ProfileImageUrl string
	Description     string
	Follower        int
	Following       int
}

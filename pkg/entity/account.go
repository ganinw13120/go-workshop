package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type Account struct {
	Id              primitive.ObjectID `bson:"_id"`
	DisplayName     string             `bson:"display_name"`
	Username        string             `bson:"username"`
	ProfileImageUrl string             `bson:"profile_image_url"`
	Description     string             `bson:"description"`
	Follower        int                `bson:"follower"`
	Following       int                `bson:"following"`
}

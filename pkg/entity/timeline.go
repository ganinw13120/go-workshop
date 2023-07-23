package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type Thread struct {
	Text        string
	UserId      primitive.ObjectID
	Likes       int
	Replies     []*Thread
	RepostCount int
}

package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type Thread struct {
	Id           primitive.ObjectID  `bson:"_id" json:"_id"`
	Text         string              `bson:"text" json:"text"`
	UserId       primitive.ObjectID  `bson:"user_id" json:"user_id"`
	Likes        int                 `bson:"likes" json:"likes"`
	ParentThread *primitive.ObjectID `bson:"parent_thread" json:"parent_thread"`
	RepostCount  int                 `bson:"repost_count" json:"repost_count"`
}

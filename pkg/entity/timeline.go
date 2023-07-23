package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type Thread struct {
	Id           primitive.ObjectID  `bson:"_id"`
	Text         string              `bson:"text"`
	UserId       primitive.ObjectID  `bson:"user_id"`
	Likes        int                 `bson:"likes"`
	ParentThread *primitive.ObjectID `bson:"parent_thread"`
	RepostCount  int                 `bson:"repost_count"`
}

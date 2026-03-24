package models

import "go.mongodb.org/mongo-driver/bson/primitive"

var (
	UseMemoryStore bool

	AdminUserID       = primitive.NewObjectID()
	DemoUserID        = primitive.NewObjectID()
	DemoCompetitionID = primitive.NewObjectID()
	DemoRecordID      = primitive.NewObjectID()
)

func InitMemoryStore() {
	UseMemoryStore = true
}

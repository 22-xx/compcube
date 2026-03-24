package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Model struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Name         string             `bson:"modelname"`
	Introduction string             `bson:"introduction"`
	IsDelete     bool               `bson:"is_delete"`
	CreateTime   time.Time          `bson:"create_time"`
	LatestTime   time.Time          `bson:"latest_time"`
}

func ModelInit() *Model {
	return &Model{
		IsDelete: false,
	}
}

func (model *Model) GetStringID() string {
	return model.ID.Hex()
}

func (model *Model) GetInfo() map[string]any {
	return map[string]any{
		"id":           model.GetStringID(),
		"dataname":     model.Name,
		"introduction": model.Introduction,
		"is_delete":    model.IsDelete,
		"create_time":  model.CreateTime.Format("2006-01-02 15:04:05"),
		"latest_time":  model.LatestTime.Format("2006-01-02 15:04:05"),
	}
}

func (model *Model) GetShortInfo() map[string]any {
	return map[string]any{
		"id":           model.GetStringID(),
		"dataname":     model.Name,
		"introduction": model.Introduction,
	}
}

func (model *Model) SetValue(modelname string, introduction string) {
	if modelname != "" {
		model.Name = modelname
	}
	if introduction != "" {
		model.Introduction = introduction
	}
}

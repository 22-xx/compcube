package data

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Data struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Name         string             `bson:"dataname"`
	Path         string             `bson:"datapath"`
	Introduction string             `bson:"introduction"`
	IsDelete     bool               `bson:"is_delete"`
	CreateTime   time.Time          `bson:"create_time"`
	LatestTime   time.Time          `bson:"latest_time"`
}

func DataInit() *Data {
	return &Data{
		IsDelete: false,
	}
}

func (data *Data) GetStringID() string {
	return data.ID.Hex()
}

func (data *Data) GetInfo() map[string]any {
	return map[string]any{
		"id":           data.GetStringID(),
		"dataname":     data.Name,
		"path":         data.Path,
		"introduction": data.Introduction,
		"is_delete":    data.IsDelete,
		"create_time":  data.CreateTime.Format("2006-01-02 15:04:05"),
		"latest_time":  data.LatestTime.Format("2006-01-02 15:04:05"),
	}
}

func (data *Data) GetShortInfo() map[string]any {
	return map[string]any{
		"id":           data.GetStringID(),
		"dataname":     data.Name,
		"path":         data.Path,
		"introduction": data.Introduction,
	}
}

func (data *Data) SetValue(dataname string, path string, introduction string) {
	if dataname != "" {
		data.Name = dataname
	}
	if path != "" {
		data.Path = path
	}
	if introduction != "" {
		data.Introduction = introduction
	}
}

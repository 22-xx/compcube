package data

import (
	"LadderCompetitionPlatform/config"
	"LadderCompetitionPlatform/models"
	"LadderCompetitionPlatform/utils"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func SelectAll(filter bson.M) []Data {
	collection := models.Client.Database(config.Config.DB.Name).Collection("data")
	find, err := collection.Find(context.Background(), filter)
	if err != nil {
		utils.Logger.Errorf("Filter: %s\nError: %s", filter, err)
	}
	defer func(find *mongo.Cursor, ctx context.Context) {
		err := find.Close(ctx)
		if err != nil {
			utils.Logger.Errorf("Error: %s", err)
		}
	}(find, context.TODO())

	// 遍历查询结果
	var resList []Data
	var res Data
	for find.Next(context.Background()) {
		_ = find.Decode(&res)
		resList = append(resList, res)
	}
	return resList
}

func SelectOne(filter bson.M) (Data, bool) {
	collection := models.Client.Database(config.Config.DB.Name).Collection("data")

	var res Data
	err := collection.FindOne(context.Background(), filter).Decode(&res)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return res, false
		} else {
			utils.Logger.Errorf("Filter: %s\nError: %s", filter, err)
			return res, false
		}
	}

	return res, true
}

func InsertOne(data Data) bool {
	collection := models.Client.Database(config.Config.DB.Name).Collection("data")

	data.CreateTime = time.Now()
	data.LatestTime = time.Now()
	_, err := collection.InsertOne(context.TODO(), data)
	if err != nil {
		utils.Logger.Errorf("Data: %s\nError: %s", data.GetInfo(), err)
		return false
	}
	return true
}

func UpdateOne(data Data) bool {
	collection := models.Client.Database(config.Config.DB.Name).Collection("data")
	data.LatestTime = time.Now()
	opts := options.Update().SetUpsert(false)
	_, err := collection.UpdateOne(context.TODO(), bson.M{"_id": data.ID}, bson.M{"$set": data}, opts)
	if err != nil {
		utils.Logger.Errorf("Data: %s\nError: %s", data.GetInfo(), err)
		return false
	}
	return true
}

func DeleteOne(filter bson.M) bool {
	collection := models.Client.Database(config.Config.DB.Name).Collection("data")
	opts := options.Update().SetUpsert(false)
	_, err := collection.UpdateOne(context.TODO(), filter, bson.M{"$set": bson.M{"is_delete": true}}, opts)
	if err != nil {
		utils.Logger.Errorf("Filter: %s\nError: %s", filter, err)
		return false
	}
	return true
}

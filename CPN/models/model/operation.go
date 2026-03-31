package model

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

func SelectAll(filter bson.M) []Model {
	collection := models.Client.Database(config.Config.DB.Name).Collection("model")
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
	var resList []Model
	var res Model
	for find.Next(context.Background()) {
		_ = find.Decode(&res)
		resList = append(resList, res)
	}
	return resList
}

func SelectOne(filter bson.M) (Model, bool) {
	collection := models.Client.Database(config.Config.DB.Name).Collection("model")

	var res Model
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

func InsertOne(model Model) bool {
	collection := models.Client.Database(config.Config.DB.Name).Collection("model")

	model.CreateTime = time.Now()
	model.LatestTime = time.Now()
	_, err := collection.InsertOne(context.TODO(), model)
	if err != nil {
		utils.Logger.Errorf("Model: %s\nError: %s", model.GetInfo(), err)
		return false
	}
	return true
}

func UpdateOne(model Model) bool {
	collection := models.Client.Database(config.Config.DB.Name).Collection("model")
	model.LatestTime = time.Now()
	opts := options.Update().SetUpsert(false)
	_, err := collection.UpdateOne(context.TODO(), bson.M{"_id": model.ID}, bson.M{"$set": model}, opts)
	if err != nil {
		utils.Logger.Errorf("Model: %s\nError: %s", model.GetInfo(), err)
		return false
	}
	return true
}

func DeleteOne(filter bson.M) bool {
	collection := models.Client.Database(config.Config.DB.Name).Collection("model")
	opts := options.Update().SetUpsert(false)
	_, err := collection.UpdateOne(context.TODO(), filter, bson.M{"$set": bson.M{"is_delete": true}}, opts)
	if err != nil {
		utils.Logger.Errorf("Filter: %s\nError: %s", filter, err)
		return false
	}
	return true
}

//go:build ignore

package competition

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

func SelectAll(filter bson.M) []Competition {
	collection := models.Client.Database(config.Config.DB.Name).Collection("competition")
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"create_time", -1}})

	find, err := collection.Find(context.Background(), filter, findOptions)
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
	var resList []Competition
	var res Competition
	for find.Next(context.Background()) {
		_ = find.Decode(&res)
		resList = append(resList, res)
	}
	return resList
}

func SelectOne(filter bson.M) (Competition, bool) {
	collection := models.Client.Database(config.Config.DB.Name).Collection("competition")

	var res Competition
	err := collection.FindOne(context.Background(), filter).Decode(&res)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return res, false
		} else {
			// 其他查询错误
			utils.Logger.Errorf("Filter: %s\nError: %s", filter, err)
			return res, false
		}
	}

	return res, true
}

func InsertOne(competition Competition) bool {
	collection := models.Client.Database(config.Config.DB.Name).Collection("competition")
	competition.CreateTime = time.Now()
	competition.LatestTime = time.Now()
	_, err := collection.InsertOne(context.TODO(), competition)
	if err != nil {
		utils.Logger.Errorf("Competition: %s\nError: %s", competition.GetInfo(), err)
		return false
	}
	return true
}

func UpdateOne(competition Competition) bool {
	collection := models.Client.Database(config.Config.DB.Name).Collection("competition")
	competition.LatestTime = time.Now()
	opts := options.Update().SetUpsert(false)
	_, err := collection.UpdateOne(context.TODO(), bson.M{"_id": competition.ID}, bson.M{"$set": competition}, opts)
	if err != nil {
		utils.Logger.Errorf("Competition: %s\nError: %s", competition.GetInfo(), err)
		return false
	}
	return true
}

func DeleteOne(filter bson.M) bool {
	collection := models.Client.Database(config.Config.DB.Name).Collection("competition")
	opts := options.Update().SetUpsert(false)
	_, err := collection.UpdateOne(context.TODO(), filter, bson.M{"$set": bson.M{"is_delete": true}}, opts)
	if err != nil {
		utils.Logger.Errorf("Filter: %s\nError: %s", filter, err)
		return false
	}
	return true
}

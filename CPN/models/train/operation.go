package train

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

func SelectAll(filter bson.M) []Train {
	collection := models.Client.Database(config.Config.DB.Name).Collection("train")
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
	var resList []Train
	var res Train
	for find.Next(context.Background()) {
		_ = find.Decode(&res)
		resList = append(resList, res)
	}
	return resList
}

func SelectOne(filter bson.M) (Train, bool) {
	collection := models.Client.Database(config.Config.DB.Name).Collection("train")

	var res Train
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

func InsertOne(train Train) bool {
	collection := models.Client.Database(config.Config.DB.Name).Collection("train")

	train.CreateTime = time.Now()
	train.LatestTime = time.Now()
	_, err := collection.InsertOne(context.TODO(), train)
	if err != nil {
		utils.Logger.Errorf("Train: %s\nError: %s", train.GetInfo(), err)
		return false
	}
	return true
}

func InsertOneReturnCreateTime(train Train) (time.Time,bool) {
	collection := models.Client.Database(config.Config.DB.Name).Collection("train")

	train.CreateTime = time.Now()
	train.LatestTime = time.Now()
	_, err := collection.InsertOne(context.TODO(), train)
	if err != nil {
		utils.Logger.Errorf("Train: %s\nError: %s", train.GetInfo(), err)
		return train.CreateTime,false
	}
	return train.CreateTime,true
}

func UpdateOne(train Train) bool {
	collection := models.Client.Database(config.Config.DB.Name).Collection("train")
	train.LatestTime = time.Now()
	opts := options.Update().SetUpsert(false)
	_, err := collection.UpdateOne(context.TODO(), bson.M{"_id": train.ID}, bson.M{"$set": train}, opts)
	if err != nil {
		utils.Logger.Errorf("Train: %s\nError: %s", train.GetInfo(), err)
		return false
	}
	return true
}

func DeleteOne(filter bson.M) bool {
	collection := models.Client.Database(config.Config.DB.Name).Collection("train")
	opts := options.Update().SetUpsert(false)
	_, err := collection.UpdateOne(context.TODO(), filter, bson.M{"$set": bson.M{"is_delete": true}}, opts)
	if err != nil {
		utils.Logger.Errorf("Filter: %s\nError: %s", filter, err)
		return false
	}
	return true
}

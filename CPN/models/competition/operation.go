package competition

import (
	"LadderCompetitionPlatform/config"
	"LadderCompetitionPlatform/models"
	"LadderCompetitionPlatform/utils"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sort"
	"time"
)

var memoryCompetitions = []Competition{
	{
		ID:          models.DemoCompetitionID,
		Author:      models.AdminUserID,
		Title:       "Demo Competition",
		Abstract:    "Local development fallback competition.",
		SortOrder:   "降序",
		TimeLimit:   20,
		DockerImage: "python:3.10",
		Status:      "进行中",
		IsDelete:    false,
		CreateTime:  time.Now(),
		LatestTime:  time.Now(),
	},
}

func SelectAll(filter bson.M) []Competition {
	if models.UseMemoryStore {
		models.InitMemoryStore()
		filtered := filterCompetitions(memoryCompetitions, filter)
		sort.Slice(filtered, func(i, j int) bool {
			return filtered[i].CreateTime.After(filtered[j].CreateTime)
		})
		return filtered
	}
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
	if models.UseMemoryStore {
		models.InitMemoryStore()
		for _, item := range memoryCompetitions {
			if matchCompetition(item, filter) {
				return item, true
			}
		}
		return Competition{}, false
	}
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
	if models.UseMemoryStore {
		models.InitMemoryStore()
		competition.ID = primitive.NewObjectID()
		competition.CreateTime = time.Now()
		competition.LatestTime = time.Now()
		memoryCompetitions = append(memoryCompetitions, competition)
		return true
	}
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
	if models.UseMemoryStore {
		models.InitMemoryStore()
		for i, item := range memoryCompetitions {
			if item.ID == competition.ID {
				competition.LatestTime = time.Now()
				memoryCompetitions[i] = competition
				return true
			}
		}
		return false
	}
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
	if models.UseMemoryStore {
		models.InitMemoryStore()
		for i, item := range memoryCompetitions {
			if matchCompetition(item, filter) {
				item.IsDelete = true
				item.LatestTime = time.Now()
				memoryCompetitions[i] = item
				return true
			}
		}
		return false
	}
	collection := models.Client.Database(config.Config.DB.Name).Collection("competition")
	opts := options.Update().SetUpsert(false)
	_, err := collection.UpdateOne(context.TODO(), filter, bson.M{"$set": bson.M{"is_delete": true}}, opts)
	if err != nil {
		utils.Logger.Errorf("Filter: %s\nError: %s", filter, err)
		return false
	}
	return true
}

func filterCompetitions(source []Competition, filter bson.M) []Competition {
	var result []Competition
	for _, item := range source {
		if matchCompetition(item, filter) {
			result = append(result, item)
		}
	}
	return result
}

func matchCompetition(item Competition, filter bson.M) bool {
	for key, value := range filter {
		switch key {
		case "_id":
			objID, ok := value.(primitive.ObjectID)
			if !ok || item.ID != objID {
				return false
			}
		case "title":
			title, ok := value.(string)
			if !ok || item.Title != title {
				return false
			}
		case "is_delete":
			isDelete, ok := value.(bool)
			if !ok || item.IsDelete != isDelete {
				return false
			}
		}
	}
	return true
}

package record

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

var memoryRecords = []Record{
	{
		ID:            models.DemoRecordID,
		UserID:        models.DemoUserID,
		CompetitionID: models.DemoCompetitionID,
		Submission:    "demo.zip",
		Score:         95,
		Error:         "",
		RunTime:       2,
		Status:        "运行完成",
		IsDelete:      false,
		CreateTime:    time.Now(),
		LatestTime:    time.Now(),
		FinishTime:    time.Now(),
	},
}

func SelectAll(filter bson.M) []Record {
	if models.UseMemoryStore {
		models.InitMemoryStore()
		filtered := filterRecords(memoryRecords, filter)
		sort.Slice(filtered, func(i, j int) bool {
			return filtered[i].CreateTime.After(filtered[j].CreateTime)
		})
		return filtered
	}
	collection := models.Client.Database(config.Config.DB.Name).Collection("record")
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
	var resList []Record
	var res Record
	for find.Next(context.Background()) {
		_ = find.Decode(&res)
		resList = append(resList, res)
	}
	return resList
}

func SelectAllOrder(competitionID string) []map[string]any {
	if models.UseMemoryStore {
		models.InitMemoryStore()
		compID, err := primitive.ObjectIDFromHex(competitionID)
		if err != nil {
			return nil
		}
		filtered := filterRecords(memoryRecords, bson.M{
			"competition": compID,
			"is_delete":   false,
		})
		bestByUser := map[primitive.ObjectID]Record{}
		for _, item := range filtered {
			if item.Score < 0 {
				continue
			}
			existing, ok := bestByUser[item.UserID]
			if !ok || item.Score > existing.Score || (item.Score == existing.Score && item.CreateTime.Before(existing.CreateTime)) {
				bestByUser[item.UserID] = item
			}
		}
		var records []Record
		for _, item := range bestByUser {
			records = append(records, item)
		}
		sort.Slice(records, func(i, j int) bool {
			if records[i].Score == records[j].Score {
				return records[i].CreateTime.Before(records[j].CreateTime)
			}
			return records[i].Score > records[j].Score
		})
		var result []map[string]any
		for i, item := range records {
			info := item.GetInfo()
			info["rank"] = i + 1
			result = append(result, info)
		}
		return result
	}
	collection := models.Client.Database(config.Config.DB.Name).Collection("record")
	objID, _ := primitive.ObjectIDFromHex(competitionID)

	pipeLine := mongo.Pipeline{
		bson.D{{
			"$match", bson.D{
				{"competition", objID},
				{"is_delete", false},
				{"score", bson.D{{"$gte", 0}}},
			},
		}},
		bson.D{{
			"$group", bson.D{
				{"_id", "$user"},
				{"maxScore", bson.D{{"$min", "$score"}}},
			},
		}},
		bson.D{{
			"$lookup", bson.D{
				{"from", "record"},
				{"let", bson.D{
					{"user", "$_id"},
					{"score", "$maxScore"},
				}},
				{"pipeline", bson.A{
					bson.D{{
						"$match", bson.D{{
							"$expr", bson.D{{
								"$and", bson.A{
									bson.D{{"$eq", bson.A{"$score", "$$score"}}},
									bson.D{{"$eq", bson.A{"$user", "$$user"}}},
								},
							}},
						}},
					}},
					bson.D{{
						"$project", bson.D{
							{"_id", 1},
							{"user", 1},
							{"competition", 1},
							{"submission", 1},
							{"score", 1},
							{"run_time", 1},
							{"status", 1},
							{"errors", 1},
							{"is_delete", 1},
							{"create_time", 1},
							{"latest_time", 1},
							{"finish_time", 1},
						},
					}},
					bson.D{{
						"$sort", bson.D{{"create_time", 1}},
					}},
				}},
				{"as", "recordDetail"},
			},
		}},
		bson.D{{
			"$sort", bson.D{{
				"maxScore", 1,
			}},
		}},
	}

	find, err := collection.Aggregate(context.TODO(), pipeLine)
	if err != nil {
		utils.Logger.Errorf("CompetitionID: %s\nError: %s", competitionID, err)
	}
	defer func(find *mongo.Cursor, ctx context.Context) {
		err := find.Close(ctx)
		if err != nil {
			utils.Logger.Errorf("Error: %s", err)
		}
	}(find, context.TODO())

	// 遍历查询结果
	var tempList []bson.M
	find.All(context.Background(), &tempList)

	type TempStruct struct {
		ID           primitive.ObjectID `bson:"_id"`
		MaxScore     int                `bson:"maxScore"`
		RecordDetail []Record
	}

	var resList []map[string]any
	var tempStruct Record
	rank := 1
	for _, tempRes := range tempList {
		tempBytes, _ := bson.Marshal(tempRes["recordDetail"].(primitive.A)[0])
		bson.Unmarshal(tempBytes, &tempStruct)
		recordInfo := tempStruct.GetInfo()
		recordInfo["rank"] = rank
		resList = append(resList, recordInfo)
		rank += 1
	}

	return resList
}

func SelectOne(filter bson.M) (Record, bool) {
	if models.UseMemoryStore {
		models.InitMemoryStore()
		for _, item := range memoryRecords {
			if matchRecord(item, filter) {
				return item, true
			}
		}
		return Record{}, false
	}
	collection := models.Client.Database(config.Config.DB.Name).Collection("record")

	var res Record
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

func InsertOne(record Record) bool {
	if models.UseMemoryStore {
		models.InitMemoryStore()
		record.ID = primitive.NewObjectID()
		record.CreateTime = time.Now()
		record.LatestTime = time.Now()
		memoryRecords = append(memoryRecords, record)
		return true
	}
	collection := models.Client.Database(config.Config.DB.Name).Collection("record")
	record.CreateTime = time.Now()
	record.LatestTime = time.Now()
	_, err := collection.InsertOne(context.TODO(), record)
	if err != nil {
		utils.Logger.Errorf("Record: %s\nError: %s", record.GetInfo(), err)
		return false
	}
	return true
}

func UpdateOne(record Record) bool {
	if models.UseMemoryStore {
		models.InitMemoryStore()
		for i, item := range memoryRecords {
			if item.ID == record.ID {
				record.LatestTime = time.Now()
				memoryRecords[i] = record
				return true
			}
		}
		return false
	}
	collection := models.Client.Database(config.Config.DB.Name).Collection("record")
	record.LatestTime = time.Now()
	opts := options.Update().SetUpsert(false)
	_, err := collection.UpdateOne(context.TODO(), bson.M{"_id": record.ID}, bson.M{"$set": record}, opts)
	if err != nil {
		utils.Logger.Errorf("Record: %s\nError: %s", record.GetInfo(), err)
		return false
	}
	return true
}

func DeleteOne(filter bson.M) bool {
	if models.UseMemoryStore {
		models.InitMemoryStore()
		for i, item := range memoryRecords {
			if matchRecord(item, filter) {
				item.IsDelete = true
				item.LatestTime = time.Now()
				memoryRecords[i] = item
				return true
			}
		}
		return false
	}
	collection := models.Client.Database(config.Config.DB.Name).Collection("record")
	opts := options.Update().SetUpsert(false)
	_, err := collection.UpdateOne(context.TODO(), filter, bson.M{"$set": bson.M{"is_delete": true}}, opts)
	if err != nil {
		utils.Logger.Errorf("Filter: %s\nError: %s", filter, err)
		return false
	}
	return true
}

func filterRecords(source []Record, filter bson.M) []Record {
	var result []Record
	for _, item := range source {
		if matchRecord(item, filter) {
			result = append(result, item)
		}
	}
	return result
}

func matchRecord(item Record, filter bson.M) bool {
	for key, value := range filter {
		switch key {
		case "_id":
			objID, ok := value.(primitive.ObjectID)
			if !ok || item.ID != objID {
				return false
			}
		case "user":
			objID, ok := value.(primitive.ObjectID)
			if !ok || item.UserID != objID {
				return false
			}
		case "competition":
			objID, ok := value.(primitive.ObjectID)
			if !ok || item.CompetitionID != objID {
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

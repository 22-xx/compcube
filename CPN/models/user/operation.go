package user

import (
	"LadderCompetitionPlatform/config"
	"LadderCompetitionPlatform/models"
	"LadderCompetitionPlatform/utils"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var memoryUsers = []User{
	{
		ID:         models.AdminUserID,
		Name:       "administrator",
		Password:   "iThings666",
		School:     "TJU",
		Email:      "administrator@example.com",
		Role:       "admin",
		Source:     "Competition_Platform",
		IsDelete:   false,
		CreateTime: time.Now(),
		LatestTime: time.Now(),
	},
	{
		ID:         models.DemoUserID,
		Name:       "demo_user",
		Password:   "demo123456",
		School:     "TJU",
		Email:      "demo_user@example.com",
		Role:       "user",
		Source:     "Competition_Platform",
		IsDelete:   false,
		CreateTime: time.Now(),
		LatestTime: time.Now(),
	},
}

func SelectAll(filter bson.M) []User {
	if models.UseMemoryStore {
		models.InitMemoryStore()
		filtered := filterUsers(memoryUsers, filter)
		res := make([]User, len(filtered))
		copy(res, filtered)
		return res
	}
	collection := models.Client.Database(config.Config.DB.Name).Collection("user")
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
	var resList []User
	var res User
	for find.Next(context.Background()) {
		_ = find.Decode(&res)
		resList = append(resList, res)
	}
	return resList
}

func SelectOne(filter bson.M) (User, bool) {
	if models.UseMemoryStore {
		models.InitMemoryStore()
		for _, item := range memoryUsers {
			if matchUser(item, filter) {
				return item, true
			}
		}
		return User{}, false
	}
	collection := models.Client.Database(config.Config.DB.Name).Collection("user")

	var res User
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

func InsertOne(user User) bool {
	if models.UseMemoryStore {
		models.InitMemoryStore()
		user.ID = primitive.NewObjectID()
		user.CreateTime = time.Now()
		user.LatestTime = time.Now()
		memoryUsers = append(memoryUsers, user)
		return true
	}
	collection := models.Client.Database(config.Config.DB.Name).Collection("user")

	user.CreateTime = time.Now()
	user.LatestTime = time.Now()
	_, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		utils.Logger.Errorf("User: %s\nError: %s", user.GetInfo(), err)
		return false
	}
	return true
}

func UpdateOne(user User) bool {
	if models.UseMemoryStore {
		models.InitMemoryStore()
		for i, item := range memoryUsers {
			if item.ID == user.ID {
				user.LatestTime = time.Now()
				memoryUsers[i] = user
				return true
			}
		}
		return false
	}
	collection := models.Client.Database(config.Config.DB.Name).Collection("user")
	user.LatestTime = time.Now()
	opts := options.Update().SetUpsert(false)
	_, err := collection.UpdateOne(context.TODO(), bson.M{"_id": user.ID}, bson.M{"$set": user}, opts)
	if err != nil {
		utils.Logger.Errorf("User: %s\nError: %s", user.GetInfo(), err)
		return false
	}
	return true
}

func DeleteOne(filter bson.M) bool {
	if models.UseMemoryStore {
		models.InitMemoryStore()
		for i, item := range memoryUsers {
			if matchUser(item, filter) {
				item.IsDelete = true
				item.LatestTime = time.Now()
				memoryUsers[i] = item
				return true
			}
		}
		return false
	}
	collection := models.Client.Database(config.Config.DB.Name).Collection("user")
	opts := options.Update().SetUpsert(false)
	_, err := collection.UpdateOne(context.TODO(), filter, bson.M{"$set": bson.M{"is_delete": true}}, opts)
	if err != nil {
		utils.Logger.Errorf("Filter: %s\nError: %s", filter, err)
		return false
	}
	return true
}

func filterUsers(source []User, filter bson.M) []User {
	var result []User
	for _, item := range source {
		if matchUser(item, filter) {
			result = append(result, item)
		}
	}
	return result
}

func matchUser(item User, filter bson.M) bool {
	for key, value := range filter {
		switch key {
		case "_id":
			objID, ok := value.(primitive.ObjectID)
			if !ok || item.ID != objID {
				return false
			}
		case "username":
			name, ok := value.(string)
			if !ok || item.Name != name {
				return false
			}
		case "password":
			password, ok := value.(string)
			if !ok || item.Password != password {
				return false
			}
		case "email":
			email, ok := value.(string)
			if !ok || item.Email != email {
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

package models

import (
	"LadderCompetitionPlatform/config"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var Client *mongo.Client

type Demo struct {
	Id   string `json:"id" bson:"_id"`
	Name string `json:"username" bson:"username"`
}

func DBConnectionInit() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(config.Config.DB.Host).SetAuth(
		options.Credential{
			Username: config.Config.DB.Username,
			Password: config.Config.DB.Password,
		}))
	if err != nil {
		fmt.Println(err)
	}
	//defer client.Disconnect(context.TODO())

	// 测试连接
	err = client.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		fmt.Println(err)
		UseMemoryStore = true
		InitMemoryStore()
		fmt.Println("mongo unavailable, fallback to in-memory store")
	} else {
		fmt.Println("connect success!!!")
	}

	Client = client
}

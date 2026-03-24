package train

import (
	"LadderCompetitionPlatform/models/data"
	"LadderCompetitionPlatform/models/model"
	"LadderCompetitionPlatform/utils"
	"os/exec"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Train struct {
	ID                primitive.ObjectID `bson:"_id,omitempty"`
	UserID            primitive.ObjectID `bson:"userid"`
	Name              string             `bson:"train_name"`
	Model             primitive.ObjectID `bson:"model"`
	Data              primitive.ObjectID `bson:"data"`
	Introduction      string             `bson:"introduction"`
	DockerImage       string             `bson:"dockerimage"`
	Decision          string             `bson:"decision"`
	DataQueue         string             `bson:"data_queue"`
	TimeQueue         string             `bson:"time_queue"`
	EnergyQueue       string             `bson:"energy_queue"`
	TimeDelay         string             `bson:"time_delay"`
	EnergyConsumption string             `bson:"energy_consumption"`
	Objective         string             `bson:"objective"`
	Cost              string             `bson:"cost"`
	Status            string             `bson:"status"`
	IsDelete          bool               `bson:"is_delete"`
	Error             string             `bson:"errors"`
	RunTime           int                `bson:"run_time"`
	CreateTime        time.Time          `bson:"create_time"`
	LatestTime        time.Time          `bson:"latest_time"`
	FinishTime        time.Time          `bson:"finish_time"`
}

func TrainInit() *Train {
	return &Train{
		Error:    "",
		RunTime:  -1,
		Status:   "训练中",
		IsDelete: false,
	}
}

func (train *Train) GetStringID() string {
	return train.ID.Hex()
}

func (train *Train) GetInfo() map[string]any {
	model, _ := model.SelectOne(bson.M{"_id": train.Model})
	data, _ := data.SelectOne(bson.M{"_id": train.Data})
	return map[string]any{
		"id":           train.GetStringID(),
		"user":         train.UserID,
		"trainname":    train.Name,
		"model":        model.GetShortInfo(),
		"data":         data.GetShortInfo(),
		"introduction": train.Introduction,
		"status":       train.Status,
		"is_delete":    train.IsDelete,
		"run_time":     train.RunTime,
		"create_time":  train.CreateTime.Format("2006-01-02 15:04:05"),
		"latest_time":  train.LatestTime.Format("2006-01-02 15:04:05"),
		"finish_time":  train.FinishTime.Format("2006-01-02 15:04:05"),
	}
}

func (train *Train) SetValue(userID string, trainName string, modelID string, dataID string, introduction string, dockerImage string, status string, decision string, dataQueue string, timeQueue string, energyQueue string, timeDelay string, energyConsumption string, objective string, cost string, finishTime string, runTime string) {
	if userID != "" {
		user, _ := primitive.ObjectIDFromHex(userID)
		train.UserID = user
	}
	if trainName != "" {
		train.Name = trainName
	}
	if modelID != "" {
		model, _ := primitive.ObjectIDFromHex(modelID)
		train.Model = model
	}
	if dataID != "" {
		data, _ := primitive.ObjectIDFromHex(dataID)
		train.Data = data
	}
	if introduction != "" {
		train.Introduction = introduction
	}
	if dockerImage != "" {
		train.DockerImage = dockerImage
	}
	if status != "" {
		train.Status = status
	}
	if decision != "" {
		train.Decision = decision
	}
	if dataQueue != "" {
		train.DataQueue = dataQueue
	}
	if timeQueue != "" {
		train.TimeQueue = timeQueue
	}
	if energyQueue != "" {
		train.EnergyQueue = energyQueue
	}
	if timeDelay != "" {
		train.TimeDelay = timeDelay
	}
	if energyConsumption != "" {
		train.EnergyConsumption = energyConsumption
	}
	if objective != "" {
		train.Objective = objective
	}
	if cost != "" {
		train.Cost = cost
	}
	if finishTime != "" {
		train.FinishTime = time.Now()
	}
	if runTime != "" {
		runTimeInt, _ := strconv.Atoi(runTime)
		train.RunTime = runTimeInt
	}
}

// 接口1： 启动docker
func (train *Train) GetDockerCmd() string {
	dockerImage := train.DockerImage
	trianID := train.GetStringID()
	return "docker run " + "-e " + "trainID" + "=" + trianID + " " + dockerImage + " &"
}

func (train *Train) RunDocker() bool {
	dockerImage := train.DockerImage
	trianID := train.GetStringID()
	cmd := exec.Command("docker", "run", "-e", "trainID="+trianID, dockerImage)
	if err := cmd.Start(); err != nil { // 运行命令
		utils.Logger.Errorf("Error: %s", err)
	}
	return false
}

//接口2：接受docker的结果，似乎不需要这个后端来做？

//接口3：将结果返回前端

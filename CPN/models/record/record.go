package record

import (
	"LadderCompetitionPlatform/config"
	"LadderCompetitionPlatform/models/competition"
	"LadderCompetitionPlatform/models/user"
	"path"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Record struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	UserID        primitive.ObjectID `bson:"user"`
	CompetitionID primitive.ObjectID `bson:"competition"`

	Submission string `bson:"submission"`
	Score      int    `bson:"score"`
	Error      string `bson:"errors"`
	RunTime    int    `bson:"run_time"`
	Status     string `bson:"status"`

	IsDelete   bool      `bson:"is_delete"`
	CreateTime time.Time `bson:"create_time"`
	LatestTime time.Time `bson:"latest_time"`
	FinishTime time.Time `bson:"finish_time"`
}

func RecordInit() *Record {
	return &Record{
		Score:    -1,
		Error:    "",
		RunTime:  -1,
		Status:   "上传完成",
		IsDelete: false,
	}
}

func (record *Record) GetStringID() string {
	return record.ID.Hex()
}

func (record *Record) GetInfo() map[string]any {
	user_, _ := user.SelectOne(bson.M{"_id": record.UserID})
	competition_, _ := competition.SelectOne(bson.M{"_id": record.CompetitionID})
	return map[string]any{
		"id":          record.GetStringID(),
		"user":        user_.GetShortInfo(),
		"competition": competition_.GetInfo(),
		"status":      record.Status,
		"run_time":    record.RunTime,
		"score":       record.Score,
		"errors":      record.Error,
		"create_time": record.CreateTime.Format("2006-01-02 15:04:05"),
		"latest_time": record.LatestTime.Format("2006-01-02 15:04:05"),
		"finish_time": record.LatestTime.Format("2006-01-02 15:04:05"),
	}
}

func (record *Record) GetPath(folderName string) string {
	return path.Join(config.Config.Path.RootPath, "user", record.UserID.Hex(), record.CompetitionID.Hex(), record.GetStringID(), folderName)
}

func (record *Record) GetDockerCmd() string {
	competition_, _ := competition.SelectOne(bson.M{"_id": record.CompetitionID})
	return "docker run " +
		"-v " + competition_.GetPath() + ":" + config.Config.Docker.InputPath + ":ro " +
		"-v " + record.GetPath("userData") + ":" + config.Config.Docker.UserPath + " " +
		"-v " + record.GetPath("output") + ":" + config.Config.Docker.OutputPath + " " +
		"-v " + record.GetPath("logs") + ":" + config.Config.Docker.LogPath + " " +
		competition_.DockerImage + " " + record.GetStringID() + " " + strconv.Itoa(competition_.TimeLimit) + " &"
}

func (record *Record) SetValue(
	userID string,
	competitionID string,
	submission string,
	score string,
	error string,
	runTime string,
	status string,
	finishTime string,
) {
	if userID != "" {
		objID, _ := primitive.ObjectIDFromHex(userID)
		record.UserID = objID
	}
	if competitionID != "" {
		objID, _ := primitive.ObjectIDFromHex(competitionID)
		record.CompetitionID = objID
	}
	if submission != "" {
		record.Submission = submission
	}
	if score != "" {
		scoreInt, _ := strconv.Atoi(score)
		record.Score = scoreInt
	}
	if error != "" {
		record.Error = error
	}
	if runTime != "" {
		runTimeInt, _ := strconv.Atoi(runTime)
		record.RunTime = runTimeInt
	}
	if status != "" {
		record.Status = status
	}
	if finishTime != "" {
		record.FinishTime = time.Now()
	}
}

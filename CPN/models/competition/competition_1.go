//go:build ignore

package competition

import (
	"LadderCompetitionPlatform/config"
	"LadderCompetitionPlatform/models/user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"path"
	"strconv"
	"time"
)

type Competition struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`

	Author primitive.ObjectID `bson:"author"`

	Title       string `bson:"title"`
	Abstract    string `bson:"abstract"`
	SortOrder   string `bson:"sort_order"`
	TimeLimit   int    `bson:"time_limit"`
	DockerImage string `bson:"docker_image"`
	Status      string `bson:"status"`

	IsDelete   bool      `bson:"is_delete"`
	CreateTime time.Time `bson:"create_time"`
	LatestTime time.Time `bson:"latest_time"`
}

func CompetitionInit() *Competition {
	return &Competition{
		SortOrder: "降序",
		TimeLimit: 20,
		Status:    "准备中",
		IsDelete:  false,
	}
}

func (competition *Competition) GetStringID() string {
	return competition.ID.Hex()
}

func (competition *Competition) GetInfo() map[string]any {
	author, _ := user.SelectOne(bson.M{"_id": competition.Author})
	return map[string]any{
		"id":          competition.GetStringID(),
		"author":      author.GetShortInfo(),
		"title":       competition.Title,
		"abstract":    competition.Abstract,
		"sort_order":  competition.SortOrder,
		"time_limit":  competition.TimeLimit,
		"status":      competition.Status,
		"create_time": competition.CreateTime.Format("2006-01-02 15:04:05"),
		"latest_time": competition.LatestTime.Format("2006-01-02 15:04:05"),
	}
}

func (competition *Competition) SetValue(
	title string,
	abstract string,
	authorID string,
	sortOrder string,
	timeLimit string,
	dockerImage string,
	statue string) {
	if title != "" {
		competition.Title = title
	}
	if abstract != "" {
		competition.Abstract = abstract
	}
	if authorID != "" {
		objID, _ := primitive.ObjectIDFromHex(authorID)
		competition.Author = objID
	}
	if sortOrder != "" {
		if sortOrder != "升序" {
			sortOrder = "降序"
		}
		competition.SortOrder = sortOrder
	}
	if timeLimit != "" {
		timeLimitInt, _ := strconv.Atoi(timeLimit)
		competition.TimeLimit = timeLimitInt
	}
	if dockerImage != "" {
		competition.DockerImage = dockerImage
	}
	if statue != "" {
		if statue != "进行中" && statue != "已结束" {
			sortOrder = "准备中"
		}
		competition.Status = statue
	}
}

func (competition *Competition) GetPath() string {
	return path.Join(config.Config.Path.RootPath, "competition", competition.GetStringID(), "input")
}

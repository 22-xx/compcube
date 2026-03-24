package app

import (
	"LadderCompetitionPlatform/models/competition"
	"LadderCompetitionPlatform/models/user"
	"LadderCompetitionPlatform/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os"
	"strconv"
)

func CompetitionAppInit(router *gin.Engine) {
	router.GET("/competition", competitionList)
	router.GET("/competition/:competitionID", competitionRetrieve)
	router.POST("/competition", competitionCreate)
	router.PUT("/competition/:competitionID", competitionUpdate)
	router.DELETE("/competition/:competitionID", competitionDelete)
}

// @Summary 查询所有比赛信息
// @Tags Competition
// @Produce json
// @Param pageNum query int false "页数" default(1)
// @Param pageSize query int false "每页数据条数" default(10)
// @Success 200 {object} map[string]any "success"
// @Router /competition [get]
func competitionList(context *gin.Context) {
	pageNumString, _ := context.GetQuery("pageNum")
	pageSizeString, _ := context.GetQuery("pageSize")

	if pageNumString == "" {
		pageNumString = "1"
	}
	if pageSizeString == "" {
		pageSizeString = "10"
	}
	pageNum, _ := strconv.Atoi(pageNumString)
	pageSize, _ := strconv.Atoi(pageSizeString)

	selectResultList := competition.SelectAll(bson.M{"is_delete": false})
	resNum := len(selectResultList)

	if ((pageNum-1)*pageSize) < 0 || ((pageNum-1)*pageSize) >= len(selectResultList) {
		selectResultList = selectResultList[0:0]
	} else if (pageNum * pageSize) < len(selectResultList) {
		selectResultList = selectResultList[((pageNum - 1) * pageSize):(pageNum * pageSize)]
	} else {
		selectResultList = selectResultList[((pageNum - 1) * pageSize):]
	}

	var resList []map[string]any
	for _, res := range selectResultList {
		resList = append(resList, res.GetInfo())
	}

	context.JSON(200, utils.SuccessResponse(map[string]any{
		"total":           resNum,
		"competitionList": resList,
	}))
}

// @Summary 查询单个比赛详细信息
// @Tags Competition
// @Produce json
// @Param competitionID path string true "赛题ID"
// @Success 200 {object} competition.Competition "success"
// @Failure 500 {object} map[string]any "该赛题不存在"
// @Router /competition/{competitionID} [get]
func competitionRetrieve(context *gin.Context) {
	competitionID := context.Param("competitionID")
	objID, _ := primitive.ObjectIDFromHex(competitionID)

	if selectResult, isSuccess := competition.SelectOne(bson.M{"_id": objID}); isSuccess {
		context.JSON(200, utils.SuccessResponse(selectResult.GetInfo()))
	} else {
		context.JSON(500, utils.ErrorResponse(nil, "该赛题不存在"))
	}
}

// @Summary 管理员创建新比赛
// @Tags Competition
// @Accept multipart/form-data
// @Produce json
// @Param title formData string true "比赛题目"
// @Param abstract formData string true "比赛简介（文档链接）"
// @Param sortOrder formData string false "成绩排序方式" Enums(升序, 降序)	default(降序)
// @Param timeLimit formData string false "时间要求" default(20)
// @Param dockerImage formData string true "赛题docker镜像名"
// @Success 200 {object} map[string]any "赛题添加成功"
// @Failure 500 {object} map[string]any "题目已被占用"
// @Failure 500 {object} map[string]any "赛题添加失败"
// @Router /competition [post]
func competitionCreate(context *gin.Context) {
	currentUser := context.MustGet("currentUser").(user.User)
	title, _ := context.GetPostForm("title")
	abstract, _ := context.GetPostForm("abstract")
	sortOrder, _ := context.GetPostForm("sortOrder")
	timeLimit, _ := context.GetPostForm("timeLimit")
	dockerImage, _ := context.GetPostForm("dockerImage")

	newCompetition := competition.CompetitionInit()
	newCompetition.SetValue(title, abstract, currentUser.GetStringID(), sortOrder, timeLimit, dockerImage, "")

	if _, isSuccess := competition.SelectOne(bson.M{"title": title}); isSuccess {
		context.JSON(500, utils.ErrorResponse(nil, "题目已被占用"))
		return
	}

	if competition.InsertOne(*newCompetition) {
		filePath := newCompetition.GetPath()
		_, err := os.Stat(filePath)
		if err != nil {
			if os.IsNotExist(err) {
				errDir := os.MkdirAll(filePath, 0755)
				if errDir != nil {
					utils.Logger.Errorf("创建文件夹: %s\nError: %s", filePath, errDir)
				}
			}
		}

		context.JSON(200, utils.SuccessResponse("赛题添加成功"))
	} else {
		context.JSON(500, utils.ErrorResponse(nil, "赛题添加失败"))
	}
}

// @Summary 管理员修改比赛赛题
// @Tags Competition
// @Accept multipart/form-data
// @Produce json
// @Param competitionID path string true "赛题ID"
// @Param title formData string false "比赛题目"
// @Param abstract formData string false "比赛简介（文档链接）"
// @Param sortOrder formData string false "成绩排序方式" Enums(升序, 降序)
// @Param timeLimit formData string false "时间要求"
// @Param dockerImage formData string false "赛题docker镜像名"
// @Param status formData string false "赛题状态" Enums(准备中, 进行中, 已结束)
// @Success 200 {object} map[string]any "赛题更新成功"
// @Failure 500 {object} map[string]any "该赛题不存在"
// @Failure 500 {object} map[string]any "题目已被占用"
// @Failure 500 {object} map[string]any "赛题更新失败"
// @Router /competition/{competitionID} [put]
func competitionUpdate(context *gin.Context) {
	competitionID := context.Param("competitionID")
	objID, _ := primitive.ObjectIDFromHex(competitionID)
	title, _ := context.GetPostForm("title")
	abstract, _ := context.GetPostForm("abstract")
	sortOrder, _ := context.GetPostForm("sortOrder")
	timeLimit, _ := context.GetPostForm("timeLimit")
	dockerImage, _ := context.GetPostForm("dockerImage")
	status, _ := context.GetPostForm("status")

	updateCompetition, isSuccess := competition.SelectOne(bson.M{"_id": objID})
	if !isSuccess {
		context.JSON(500, utils.ErrorResponse(nil, "该赛题不存在"))
		return
	}
	updateCompetition.SetValue(title, abstract, "", sortOrder, timeLimit, dockerImage, status)

	if _, isSuccess := competition.SelectOne(bson.M{"title": title}); isSuccess {
		context.JSON(500, utils.ErrorResponse(nil, "题目已被占用"))
		return
	}

	if competition.UpdateOne(updateCompetition) {
		context.JSON(200, utils.SuccessResponse("赛题更新成功"))
	} else {
		context.JSON(500, utils.ErrorResponse(nil, "赛题更新失败"))
	}
}

// @Summary 管理员删除赛题
// @Tags Competition
// @Produce json
// @Param competitionID path string true "赛题ID"
// @Success 200 {object} map[string]any "赛题删除成功"
// @Failure 500 {object} map[string]any "该赛题不存在"
// @Failure 500 {object} map[string]any "赛题删除失败"
// @Router /competition/{competitionID} [delete]
func competitionDelete(context *gin.Context) {
	competitionID := context.Param("competitionID")
	objID, _ := primitive.ObjectIDFromHex(competitionID)

	if _, isSuccess := competition.SelectOne(bson.M{"_id": objID}); !isSuccess {
		context.JSON(500, utils.ErrorResponse(nil, "该赛题不存在"))
		return
	}

	if competition.DeleteOne(bson.M{"_id": objID}) {
		context.JSON(200, utils.SuccessResponse("赛题删除成功"))
	} else {
		context.JSON(500, utils.ErrorResponse(nil, "赛题删除失败"))
	}
}

package app

import (
	"LadderCompetitionPlatform/models/competition"
	"LadderCompetitionPlatform/models/record"
	"LadderCompetitionPlatform/models/user"
	"LadderCompetitionPlatform/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"os"
	"strconv"
	"strings"
)

func RecordAppInit(router *gin.Engine) {
	router.GET("/record", recordList)
	router.GET("/competition/:competitionID/record", recordOrder)
	router.GET("/competition/:competitionID/record/:recordID", recordRetrieve)
	router.POST("/competition/:competitionID/record", recordCreate)
	router.PUT("/competition/:competitionID/record/:recordID", recordUpdate)
	router.DELETE("/competition/:competitionID/record/:recordID", recordDelete)
}

// @Summary 管理员查询所有提交记录，用户查询个人提交记录
// @Tags Record
// @Produce json
// @Param pageNum query int false "页数" default(1)
// @Param pageSize query int false "每页数据条数" default(10)
// @Success 200 {object} map[string]any "success"
// @Router /record [get]
func recordList(context *gin.Context) {
	currentUser := context.MustGet("currentUser").(user.User)
	role := currentUser.Role
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

	var selectResultList []record.Record
	if role == "admin" {
		selectResultList = record.SelectAll(bson.M{"is_delete": false})
	} else if role == "user" {
		selectResultList = record.SelectAll(bson.M{"user": currentUser.ID, "is_delete": false})
	} else {
		context.JSON(500, utils.ErrorResponse(nil, "未知身份"))
		return
	}

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
		"total":      resNum,
		"recordList": resList,
	}))
}

// @Summary 查询某个赛题的提交排名
// @Tags Record
// @Produce json
// @Param pageNum query int false "页数" default(1)
// @Param pageSize query int false "每页数据条数" default(10)
// @Param competitionID path string true "赛题ID"
// @Success 200 {object} map[string]any "success"
// @Router /competition/{competitionID}/record [get]
func recordOrder(context *gin.Context) {
	competitionID := context.Param("competitionID")
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

	selectResultList := record.SelectAllOrder(competitionID)

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
		resList = append(resList, res)
	}

	context.JSON(200, utils.SuccessResponse(map[string]any{
		"total":      resNum,
		"recordList": resList,
	}))
}

// @Summary 查询单个提交
// @Tags Record
// @Produce json
// @Param competitionID path string true "赛题ID"
// @Param recordID path string true "提交记录ID"
// @Success 200 {object} record.Record "success"
// @Failure 500 {object} map[string]any "该次提交不存在"
// @Router /competition/{competitionID}/record/{recordID} [get]
func recordRetrieve(context *gin.Context) {
	recordID := context.Param("recordID")


	
	objID, _ := primitive.ObjectIDFromHex(recordID)
	if selectResult, isSuccess := record.SelectOne(bson.M{"_id": objID, "is_delete": false}); isSuccess {
		context.JSON(200, utils.SuccessResponse(selectResult.GetInfo()))
	} else {
		context.JSON(500, utils.ErrorResponse(nil, "该次提交不存在"))
	}
}

// @Summary 进行一次提交
// @Tags Record
// @Accept multipart/form-data
// @Produce json
// @Param competitionID path string true "赛题ID"
// @Param submission formData file true "提交文件"
// @Success 200 {object} map[string]any "提交成功"
// @Failure 500 {object} map[string]any "提交失败，该比赛还在准备中"
// @Failure 500 {object} map[string]any "提交失败，该比赛已经结束"
// @Failure 500 {object} map[string]any "该比赛不存在或已被删除"
// @Failure 500 {object} map[string]any "上传文件非zip文件"
// @Failure 500 {object} map[string]any "提交失败"
// @Router /competition/{competitionID}/record [post]
func recordCreate(context *gin.Context) {
	currentUser := context.MustGet("currentUser").(user.User)
	competitionID := context.Param("competitionID")
	submission, _ := context.FormFile("submission")

	objID, _ := primitive.ObjectIDFromHex(competitionID)
	selectResult, isSuccess := competition.SelectOne(bson.M{"_id": objID})
	if isSuccess {
		if selectResult.Status == "准备中" {
			context.JSON(500, "提交失败，该比赛还在准备中")
			return
		} else if selectResult.Status == "已结束" {
			context.JSON(500, "提交失败，该比赛已经结束")
			return
		}
	} else {
		context.JSON(500, "该比赛不存在或已被删除")
		return
	}

	if fileNameSlice := strings.Split(submission.Filename, "."); fileNameSlice[len(fileNameSlice)-1] != "zip" {
		context.JSON(500, "上传文件非zip文件")
		return
	}

	newRecord := record.RecordInit()
	newRecord.SetValue(currentUser.GetStringID(), competitionID, submission.Filename, "", "", "", "", "")

	if record.InsertOne(*newRecord) {
		filePath := newRecord.GetPath("userData")
		_, err := os.Stat(filePath)
		if err != nil {
			if os.IsNotExist(err) {
				errDir := os.MkdirAll(filePath, 0755)
				if errDir != nil {
					utils.Logger.Errorf("创建文件夹: %s\nError: %s", filePath, errDir)
				}
			}
		}

		context.JSON(200, utils.SuccessResponse("提交成功"))
		utils.DockerManager(newRecord.GetDockerCmd())
	} else {
		context.JSON(500, utils.ErrorResponse(nil, "提交失败"))
	}
}

// @Summary 赛题docker更新提交跑分结果
// @Tags Record
// @Accept multipart/form-data
// @Produce json
// @Param competitionID path string true "赛题ID"
// @Param recordID path string true "提交记录ID"
// @Param score formData string true "得分"
// @Param error formData string true "错误信息"
// @Param runTime formData string true "运行时间"
// @Param status formData string true "状态"
// @Success 200 {object} map[string]any "更新成功"
// @Failure 500 {object} map[string]any "更新失败"
// @Router /competition/{competitionID}/record/{recordID} [put]
func recordUpdate(context *gin.Context) {
	recordID := context.Param("recordID")
	score, _ := context.GetPostForm("score")
	error_, _ := context.GetPostForm("error")
	runTime, _ := context.GetPostForm("runTime")
	status, _ := context.GetPostForm("status")

	objID, _ := primitive.ObjectIDFromHex(recordID)
	updateRecord, _ := record.SelectOne(bson.M{"_id": objID})
	updateRecord.SetValue("", "", "", score, error_, runTime, status, "now")

	if record.UpdateOne(updateRecord) {
		context.JSON(200, utils.SuccessResponse("更新成功"))
	} else {
		context.JSON(500, utils.ErrorResponse(nil, "更新失败"))
	}
}

// @Summary 管理员删除提交记录
// @Tags Record
// @Accept multipart/form-data
// @Produce json
// @Param competitionID path string true "赛题ID"
// @Param recordID path string true "提交记录ID"
// @Success 200 {object} map[string]any "该次提交删除成功"
// @Failure 500 {object} map[string]any "该次提交不存在"
// @Failure 500 {object} map[string]any "该次提交删除失败"
// @Router /competition/{competitionID}/record/{recordID} [delete]
func recordDelete(context *gin.Context) {
	recordID := context.Param("recordID")
	objID, _ := primitive.ObjectIDFromHex(recordID)

	if _, isSuccess := record.SelectOne(bson.M{"_id": objID}); !isSuccess {
		context.JSON(500, utils.ErrorResponse(nil, "该次提交不存在"))
		return
	}

	if record.DeleteOne(bson.M{"_id": objID}) {
		context.JSON(200, utils.SuccessResponse("该次提交删除成功"))
	} else {
		context.JSON(500, utils.ErrorResponse(nil, "该次提交删除失败"))
	}
}

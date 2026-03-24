package app

import (
	"LadderCompetitionPlatform/models/train"
	"LadderCompetitionPlatform/models/user"
	"LadderCompetitionPlatform/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TrainAppInit(router *gin.Engine) {
	router.GET("/train", trainList)
	router.GET("/train/:trainID", trainRetrieve)
	router.POST("/train/:modelID/:dataID", trainCreate)
	router.PUT("/train/:trainID", trainUpdate)
	router.DELETE("/train/:trainID", trainDelete)
}

// @Summary 管理员查询所有提交记录，用户查询个人提交记录
// @Tags Record
// @Produce json
// @Param pageNum query int false "页数" default(1)
// @Param pageSize query int false "每页数据条数" default(10)
// @Success 200 {object} map[string]any "success"
// @Router /record [get]
func trainList(context *gin.Context) {
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

	var selectResultList []train.Train
	if role == "admin" {
		selectResultList = train.SelectAll(bson.M{"is_delete": false})
	} else if role == "user" {
		selectResultList = train.SelectAll(bson.M{"userid": currentUser.ID, "is_delete": false})
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
		"total":     resNum,
		"trainList": resList,
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
func trainRetrieve(context *gin.Context) {
	trainID := context.Param("trainID")

	objID, _ := primitive.ObjectIDFromHex(trainID)
	if selectResult, isSuccess := train.SelectOne(bson.M{"_id": objID, "is_delete": false}); isSuccess {
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
func trainCreate(context *gin.Context) {
	currentUser := context.MustGet("currentUser").(user.User)
	trainName, _ := context.GetPostForm("trainname")
	modelID := context.Param("modelID")
	dataID := context.Param("dataID")
	introduction, _ := context.GetPostForm("introduction")
	dockerImage, _ := context.GetPostForm("dockerImage")

	newTrain := train.TrainInit()
	newTrain.SetValue(currentUser.GetStringID(), trainName, modelID, dataID, introduction, dockerImage, "", "", "", "","", "", "", "","", "", "")

	createTime,ok := train.InsertOneReturnCreateTime(*newTrain)
	if ok {
		context.JSON(200, utils.SuccessResponse("提交成功"))
		trainSuccessInsert,_ := train.SelectOne(bson.M{"userid": newTrain.UserID,"train_name":newTrain.Name,"model":newTrain.Model,"data":newTrain.Data,"introduction":newTrain.Introduction,"dockerimage":newTrain.DockerImage,"create_time":createTime})
		
		//执行docker语句
		trainSuccessInsert.RunDocker()
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

//decision string,dataQueue string,timeQueue string,energyQueue string,timeDelay string,energyConsumption string,objective string,cost string,

func trainUpdate(context *gin.Context) {
	trainID := context.Param("trainID")
	decision, _ := context.GetPostForm("decision")
	dataQueue, _ := context.GetPostForm("dataQueue")
	timeQueue, _ := context.GetPostForm("timeQueue")
	energyQueue, _ := context.GetPostForm("energyQueue")
	timeDelay, _ := context.GetPostForm("timeDelay")
	energyConsumption, _ := context.GetPostForm("energyConsumption")
	objective, _ := context.GetPostForm("objective")
	cost, _ := context.GetPostForm("cost")
	runTime, _ := context.GetPostForm("runTime")
	status, _ := context.GetPostForm("status")

	objID, _ := primitive.ObjectIDFromHex(trainID)
	updateTrain, _ := train.SelectOne(bson.M{"_id": objID})
	updateTrain.SetValue("", "", "", "", "", "", status, decision,  dataQueue, timeQueue, energyQueue, timeDelay, energyConsumption, objective, cost,"now", runTime)

	if train.UpdateOne(updateTrain) {
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
func trainDelete(context *gin.Context) {
	trainID := context.Param("trainID")
	objID, _ := primitive.ObjectIDFromHex(trainID)

	if _, isSuccess := train.SelectOne(bson.M{"_id": objID}); !isSuccess {
		context.JSON(500, utils.ErrorResponse(nil, "该次提交不存在"))
		return
	}

	if train.DeleteOne(bson.M{"_id": objID}) {
		context.JSON(200, utils.SuccessResponse("该次提交删除成功"))
	} else {
		context.JSON(500, utils.ErrorResponse(nil, "该次提交删除失败"))
	}
}

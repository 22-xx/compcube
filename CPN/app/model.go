package app

import (
	"LadderCompetitionPlatform/models/model"
	"LadderCompetitionPlatform/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ModelAppInit(router *gin.Engine) {
	router.GET("/model", modelList)
	router.GET("/model/:modelID", modelRetrieve)
	router.POST("/model", modelCreate)
	router.PUT("/model/:modelID", modelUpdate)
	router.DELETE("/model/:modelID", modelDelete)
}

// @Summary 管理员查询所有数据集信息
// @Tags Model
// @Produce json
// @Param pageNum query int false "页数" default(1)
// @Param pageSize query int false "每页数据条数" default(10)
// @Success 200 {object} map[string]any "success"
// @Router /model [get]
func modelList(context *gin.Context) {
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

	selectResultList := model.SelectAll(bson.M{"is_delete": false})
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
		"total":    resNum,
		"modelList": resList,
	}))
}

// @Summary 查询当前数据集的详细信息
// @Tags Data
// @Produce json
// @Param userID path string true "用户ID(任意值)"
// @Success 200 {object} user.User "success"
// @Router /user/{userID} [get]
func modelRetrieve(context *gin.Context) {
	// currentUser := context.MustGet("currentUser").(user.User)
	URLDataID := context.Param("modelID")

	objID, _ := primitive.ObjectIDFromHex(URLDataID)
	if selectResult, isSuccess := model.SelectOne(bson.M{"_id": objID, "is_delete": false}); isSuccess {
		context.JSON(200, utils.SuccessResponse(selectResult.GetShortInfo()))
	}else {
		context.JSON(500, utils.ErrorResponse(nil, "该模型不存在"))
	}

	//
	//if URLUserID == "0" {
	// context.JSON(200, utils.SuccessResponse(currentUser.GetShortInfo()))
	//} else if currentUser.Role == "admin" {
	//	objID, _ := primitive.ObjectIDFromHex(URLUserID)
	//	if selectResult, isSuccess := user.SelectOne(bson.M{"_id": objID, "is_delete": false}); isSuccess {
	//		context.JSON(200, utils.SuccessResponse(selectResult.GetShortInfo()))
	//	} else {
	//		context.JSON(500, utils.ErrorResponse(nil, "该用户不存在"))
	//	}
	//} else {
	//	context.JSON(500, utils.ErrorResponse(nil, "当前用户权限不足"))
	//}
}

// @Summary 管理员新建数据集
// @Tags Data
// @Accept multipart/form-data
// @Produce json
// @Param dataname formData string true "用户名"
// @Param datapath formData string true "密码"
// @Param introduction formData string true "邮箱"
// @Success 200 {object} map[string]any "数据集添加成功"
// @Failure 500 {object} map[string]any "数据集名已被占用"
// @Failure 500 {object} map[string]any "用户添加失败"
// @Router /user [post]
func modelCreate(context *gin.Context) {
	modelname, _ := context.GetPostForm("modelname")
	introduction, _ := context.GetPostForm("introduction")


	newModel := model.ModelInit()
	newModel.SetValue(modelname, introduction)

	if _, isSuccess := model.SelectOne(bson.M{"modelname": modelname}); isSuccess {
		context.JSON(500, utils.ErrorResponse(nil, "模型名已被占用"))
		return
	}


	if model.InsertOne(*newModel) {
		context.JSON(200, utils.SuccessResponse("模型添加成功"))
	} else {
		context.JSON(500, utils.ErrorResponse(nil, "模型添加失败"))
	}
}

// @Summary 修改数据集
// @Tags User
// @Accept multipart/form-data
// @Produce json
// @Param userID path string true "用户ID(任意值)"
// @Param username formData string false "用户名"
// @Param password formData string false "密码"
// @Param role formData string false "用户角色" Enums(admin, user)
// @Param email formData string false "邮箱"
// @Param school formData string false "学校"
// @Success 200 {object} map[string]any "用户更新成功"
// @Failure 500 {object} map[string]any "用户名已被占用"
// @Failure 500 {object} map[string]any "邮箱已被占用"
// @Failure 500 {object} map[string]any "用户更新失败"
// @Router /user/{userID} [put]
func modelUpdate(context *gin.Context) {

	modelname, _ := context.GetPostForm("modelname")
	introduction, _ := context.GetPostForm("introduction")

	URLModelID := context.Param("modelID")

	objID, _ := primitive.ObjectIDFromHex(URLModelID)
	currentModel,_ := model.SelectOne(bson.M{"_id": objID, "is_delete": false})

	if _, isSuccess := model.SelectOne(bson.M{"modelname": modelname}); isSuccess {
		context.JSON(500, utils.ErrorResponse(nil, "模型名已被占用"))
		return
	}

	currentModel.SetValue(modelname, introduction)

	if model.UpdateOne(currentModel) {
		context.JSON(200, utils.SuccessResponse("模型更新成功"))
	} else {
		context.JSON(500, utils.ErrorResponse(nil, "模型更新失败"))
	}
}

// @Summary 管理员删除数据集
// @Tags Data
// @Produce json
// @Param userID path string true "用户ID(任意值)"
// @Success 200 {object} map[string]any "用户删除成功"
// @Failure 500 {object} map[string]any "该用户不存在"
// @Failure 500 {object} map[string]any "用户删除失败"
// @Router /user/{userID} [delete]
func modelDelete(context *gin.Context) {
	modelID := context.Param("modelID")
	objID, _ := primitive.ObjectIDFromHex(modelID)

	if _, isSuccess := model.SelectOne(bson.M{"_id": objID}); !isSuccess {
		context.JSON(500, utils.ErrorResponse(nil, "该模型不存在"))
		return
	}

	if model.DeleteOne(bson.M{"_id": objID}) {
		context.JSON(200, utils.SuccessResponse("模型删除成功"))
	} else {
		context.JSON(500, utils.ErrorResponse(nil, "模型删除失败"))
	}
}

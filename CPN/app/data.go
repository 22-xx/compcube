package app

import (
	"LadderCompetitionPlatform/models/data"
	"LadderCompetitionPlatform/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func DataAppInit(router *gin.Engine) {
	router.GET("/data", dataList)
	router.GET("/data/:dataID", dataRetrieve)
	router.POST("/data", dataCreate)
	router.PUT("/data/:dataID", dataUpdate)
	router.DELETE("/data/:dataID", dataDelete)
}

// @Summary 管理员查询所有数据集信息
// @Tags Data
// @Produce json
// @Param pageNum query int false "页数" default(1)
// @Param pageSize query int false "每页数据条数" default(10)
// @Success 200 {object} map[string]any "success"
// @Router /data [get]
func dataList(context *gin.Context) {
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

	selectResultList := data.SelectAll(bson.M{"is_delete": false})
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
		"dataList": resList,
	}))
}

// @Summary 查询当前数据集的详细信息
// @Tags Data
// @Produce json
// @Param userID path string true "用户ID(任意值)"
// @Success 200 {object} user.User "success"
// @Router /user/{userID} [get]
func dataRetrieve(context *gin.Context) {
	// currentUser := context.MustGet("currentUser").(user.User)
	URLDataID := context.Param("dataID")

	objID, _ := primitive.ObjectIDFromHex(URLDataID)
	if selectResult, isSuccess := data.SelectOne(bson.M{"_id": objID, "is_delete": false}); isSuccess {
		context.JSON(200, utils.SuccessResponse(selectResult.GetShortInfo()))
	}else {
		context.JSON(500, utils.ErrorResponse(nil, "该数据集不存在"))
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
func dataCreate(context *gin.Context) {
	dataname, _ := context.GetPostForm("dataname")
	datapath, _ := context.GetPostForm("datapath")
	introduction, _ := context.GetPostForm("introduction")


	newData := data.DataInit()
	newData.SetValue(dataname, datapath, introduction)

	if _, isSuccess := data.SelectOne(bson.M{"dataname": dataname}); isSuccess {
		context.JSON(500, utils.ErrorResponse(nil, "数据集名已被占用"))
		return
	}


	if data.InsertOne(*newData) {
		context.JSON(200, utils.SuccessResponse("数据集添加成功"))
	} else {
		context.JSON(500, utils.ErrorResponse(nil, "数据集添加失败"))
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
func dataUpdate(context *gin.Context) {

	dataname, _ := context.GetPostForm("dataname")
	datapath, _ := context.GetPostForm("datapath")
	introduction, _ := context.GetPostForm("introduction")

	URLDataID := context.Param("dataID")

	objID, _ := primitive.ObjectIDFromHex(URLDataID)
	currentData,_ := data.SelectOne(bson.M{"_id": objID, "is_delete": false})

	if _, isSuccess := data.SelectOne(bson.M{"dataname": dataname}); isSuccess {
		context.JSON(500, utils.ErrorResponse(nil, "数据集名已被占用"))
		return
	}

	currentData.SetValue(dataname, datapath, introduction)

	if data.UpdateOne(currentData) {
		context.JSON(200, utils.SuccessResponse("数据集更新成功"))
	} else {
		context.JSON(500, utils.ErrorResponse(nil, "数据集更新失败"))
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
func dataDelete(context *gin.Context) {
	dataID := context.Param("dataID")
	objID, _ := primitive.ObjectIDFromHex(dataID)

	if _, isSuccess := data.SelectOne(bson.M{"_id": objID}); !isSuccess {
		context.JSON(500, utils.ErrorResponse(nil, "该数据集不存在"))
		return
	}

	if data.DeleteOne(bson.M{"_id": objID}) {
		context.JSON(200, utils.SuccessResponse("数据集删除成功"))
	} else {
		context.JSON(500, utils.ErrorResponse(nil, "数据集删除失败"))
	}
}

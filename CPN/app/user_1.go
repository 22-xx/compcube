//go:build ignore

package app

import (
	"LadderCompetitionPlatform/models/user"
	"LadderCompetitionPlatform/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strconv"
)

func UserAppInit(router *gin.Engine) {
	router.GET("/user", userList)
	router.GET("/user/:userID", userRetrieve)
	router.POST("/user", userCreate)
	router.PUT("/user/:userID", userUpdate)
	router.DELETE("/user/:userID", userDelete)
}

// @Summary 管理员查询所有用户信息
// @Tags User
// @Produce json
// @Param pageNum query int false "页数" default(1)
// @Param pageSize query int false "每页数据条数" default(10)
// @Success 200 {object} map[string]any "success"
// @Router /user [get]
func userList(context *gin.Context) {
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

	selectResultList := user.SelectAll(bson.M{"is_delete": false})
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
		"userList": resList,
	}))
}

// @Summary 查询当前登录用户的详细信息
// @Tags User
// @Produce json
// @Param userID path string true "用户ID(任意值)"
// @Success 200 {object} user.User "success"
// @Failure 500 {object} map[string]any "该用户不存在"
// @Failure 500 {object} map[string]any "当前用户权限不足"
// @Router /user/{userID} [get]
func userRetrieve(context *gin.Context) {
	currentUser := context.MustGet("currentUser").(user.User)
	//URLUserID := context.Param("userID")
	//
	//if URLUserID == "0" {
	context.JSON(200, utils.SuccessResponse(currentUser.GetShortInfo()))
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

// @Summary 管理员新建用户
// @Tags User
// @Accept multipart/form-data
// @Produce json
// @Param username formData string true "用户名"
// @Param password formData string true "密码"
// @Param role formData string false "用户角色" default(user)	Enums(admin, user)
// @Param email formData string true "邮箱"
// @Param school formData string true "学校"
// @Success 200 {object} map[string]any "用户添加成功"
// @Failure 500 {object} map[string]any "用户名已被占用"
// @Failure 500 {object} map[string]any "邮箱已被占用"
// @Failure 500 {object} map[string]any "用户添加失败"
// @Router /user [post]
func userCreate(context *gin.Context) {
	username, _ := context.GetPostForm("username")
	password, _ := context.GetPostForm("password")
	role, _ := context.GetPostForm("role")
	email, _ := context.GetPostForm("email")
	school, _ := context.GetPostForm("school")

	newUser := user.UserInit()
	newUser.SetValue(username, password, school, email, role)

	if _, isSuccess := user.SelectOne(bson.M{"username": username}); isSuccess {
		context.JSON(500, utils.ErrorResponse(nil, "用户名已被占用"))
		return
	}

	if _, isSuccess := user.SelectOne(bson.M{"email": email}); isSuccess {
		context.JSON(500, utils.ErrorResponse(nil, "邮箱已被占用"))
		return
	}

	if user.InsertOne(*newUser) {
		context.JSON(200, utils.SuccessResponse("用户添加成功"))
	} else {
		context.JSON(500, utils.ErrorResponse(nil, "用户添加失败"))
	}
}

// @Summary 修改个人信息
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
func userUpdate(context *gin.Context) {
	currentUser := context.MustGet("currentUser").(user.User)
	username, _ := context.GetPostForm("username")
	password, _ := context.GetPostForm("password")
	school, _ := context.GetPostForm("school")
	email, _ := context.GetPostForm("email")
	role, _ := context.GetPostForm("role")

	if role == "admin" && currentUser.Role != "admin" {
		role = "user"
	}

	if _, isSuccess := user.SelectOne(bson.M{"username": username}); isSuccess {
		context.JSON(500, utils.ErrorResponse(nil, "用户名已被占用"))
		return
	}
	if _, isSuccess := user.SelectOne(bson.M{"email": email}); isSuccess {
		context.JSON(500, utils.ErrorResponse(nil, "邮箱已被占用"))
		return
	}

	currentUser.SetValue(username, password, school, email, role)

	if user.UpdateOne(currentUser) {
		context.JSON(200, utils.SuccessResponse("用户更新成功"))
	} else {
		context.JSON(500, utils.ErrorResponse(nil, "用户更新失败"))
	}
}

// @Summary 管理员删除用户
// @Tags User
// @Produce json
// @Param userID path string true "用户ID(任意值)"
// @Success 200 {object} map[string]any "用户删除成功"
// @Failure 500 {object} map[string]any "该用户不存在"
// @Failure 500 {object} map[string]any "用户删除失败"
// @Router /user/{userID} [delete]
func userDelete(context *gin.Context) {
	userID := context.Param("userID")
	objID, _ := primitive.ObjectIDFromHex(userID)

	if _, isSuccess := user.SelectOne(bson.M{"_id": objID}); !isSuccess {
		context.JSON(500, utils.ErrorResponse(nil, "该用户不存在"))
		return
	}

	if user.DeleteOne(bson.M{"_id": objID}) {
		context.JSON(200, utils.SuccessResponse("用户删除成功"))
	} else {
		context.JSON(500, utils.ErrorResponse(nil, "用户删除失败"))
	}
}

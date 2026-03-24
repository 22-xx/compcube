package app

import (
	"LadderCompetitionPlatform/models/user"
	"LadderCompetitionPlatform/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func LoginAppInit(router *gin.Engine) {
	router.POST("/login", login)
	router.POST("/register", register)
	router.OPTIONS("/getInfo", getInfo)
}

// @Summary 登录
// @Tags Login
// @Accept multipart/form-data
// @Produce json
// @Param username formData string true "用户名"
// @Param password formData string true "密码"
// @Success 200 {object} map[string]any "success"
// @Failure 500 {object} map[string]any "用户名密码错误"
// @Router /login [post]
func login(context *gin.Context) {
	username, _ := context.GetPostForm("username")
	password, _ := context.GetPostForm("password")

	loginUser, isSuccess := user.SelectOne(bson.M{"username": username, "password": password, "is_delete": false})
	if isSuccess {
		cookie := utils.CookieEncoder(loginUser.GetStringID())
		context.SetCookie("LCP-Cookie", cookie, 3600, "/", "", false, true)
		context.JSON(200, utils.SuccessResponse(loginUser.GetShortInfo()))
	} else {
		context.JSON(500, utils.ErrorResponse(nil, "用户名密码错误"))
	}
}

// @Summary 注册
// @Tags Login
// @Accept multipart/form-data
// @Produce json
// @Param username formData string true "用户名"
// @Param password formData string true "密码"
// @Param email formData string true "邮箱"
// @Param school formData string true "学校"
// @Success 200 {object} map[string]any "注册成功"
// @Failure 500 {object} map[string]any "用户名已被占用"
// @Failure 500 {object} map[string]any "邮箱已被占用"
// @Failure 500 {object} map[string]any "注册失败"
// @Router /register [post]
func register(context *gin.Context) {
	username, _ := context.GetPostForm("username")
	password, _ := context.GetPostForm("password")
	email, _ := context.GetPostForm("email")
	school, _ := context.GetPostForm("school")

	newUser := user.UserInit()
	newUser.SetValue(username, password, school, email, "user")

	if _, isSuccess := user.SelectOne(bson.M{"username": username}); isSuccess {
		context.JSON(500, utils.ErrorResponse(nil, "用户名已被占用"))
		return
	}

	if _, isSuccess := user.SelectOne(bson.M{"email": email}); isSuccess {
		context.JSON(500, utils.ErrorResponse(nil, "邮箱已被占用"))
		return
	}

	if user.InsertOne(*newUser) {
		context.JSON(200, utils.SuccessResponse("注册成功"))
	} else {
		context.JSON(500, utils.ErrorResponse(nil, "注册失败"))
	}
}

// @Summary 获取登录信息
// @Tags Login
// @Accept multipart/form-data
// @Produce json
// @Success 200 {object} map[string]any "success"
// @Router /getInfo [OPTIONS]
func getInfo(context *gin.Context) {
	currentUser := context.MustGet("currentUser").(user.User)
	context.JSON(200, utils.SuccessResponse(currentUser.GetShortInfo()))
}

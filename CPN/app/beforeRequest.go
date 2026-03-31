package app

import (
	"LadderCompetitionPlatform/models/user"
	"LadderCompetitionPlatform/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/storyicon/grbac"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

func BeforeRequestAppInit(router *gin.Engine) {
	router.Use(globalBaseLog)
	router.Use(RouterAuthorization())

	router.GET("/", func(context *gin.Context) {
		context.JSON(200, utils.SuccessResponse("访问成功"))
	})
}

func globalBaseLog(context *gin.Context) {
	startTime := time.Now()
	context.Next()
	endTime := time.Now()
	latencyTime := endTime.Sub(startTime)

	statusCode := context.Writer.Status() // 状态码
	reqMethod := context.Request.Method   // 请求方式
	reqUri := context.Request.RequestURI  // 请求路由
	clientIP := context.ClientIP()        // 请求IP

	// 日志格式
	utils.Logger.Infof("| %3d | %7s | %13v | http://%s%s |",
		statusCode,
		reqMethod,
		latencyTime,
		clientIP,
		reqUri,
	)
}

func QueryRolesByHeaders(context *gin.Context) (roles []string, err error) {
	cookie, _ := context.Cookie("LCP-Cookie")
	if cookie != "" {
		userID := utils.CookieDecoder(cookie)
		if userID != "" {
			objID, _ := primitive.ObjectIDFromHex(userID)
			loginUser, hasValue := user.SelectOne(bson.M{"_id": objID})

			if hasValue {
				roles = append(roles, loginUser.Role)
				context.Set("currentUser", loginUser)
				return roles, err
			} else {
				return roles, errors.New("登录信息无效")
			}
		} else {
			return roles, errors.New("登录信息无效")
		}
	}
	return roles, err
}

func RouterAuthorization() gin.HandlerFunc {
	rbac, err := grbac.New(grbac.WithYAML("./config/routerConfig.yaml", time.Minute))
	if err != nil {
		panic(err)
	}
	return func(c *gin.Context) {
		roles, err := QueryRolesByHeaders(c)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		state, _ := rbac.IsRequestGranted(c.Request, roles)
		if !state.IsGranted() {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
	}
}

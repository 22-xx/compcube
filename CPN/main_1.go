//go:build ignore

package main

import (
	"LadderCompetitionPlatform/app"
	"LadderCompetitionPlatform/config"
	"LadderCompetitionPlatform/models"
	"LadderCompetitionPlatform/utils"
	"github.com/gin-gonic/gin"
)

func Init(router *gin.Engine) {
	config.ConfigInit()
	models.DBConnectionInit()
	utils.LoggerInit()

	app.BeforeRequestAppInit(router)
	app.LoginAppInit(router)
	app.UserAppInit(router)
	app.CompetitionAppInit(router)
	app.RecordAppInit(router)
	app.StaticAppInit(router)
	app.DocsInit(router)
}

// @title Ladder Competition Platform
// @version 1.0
// @description ↑ ↑ ↓ ↓ ← → ← → B A
// @termsOfService  http://swagger.io/terms/

// @contact.name Ruizhe Ma
// @contact.email ruizhe_ma@tju.edu.cn

// @host 127.0.0.1:8000
// @BasePath /
func main() {
	router := gin.Default()
	Init(router)
	router.Run(config.Config.App.Host + ":" + config.Config.App.Port) // 监听并在 0.0.0.0:8000 上启动服务
}

//go:build ignore

package app

import (
	_ "LadderCompetitionPlatform/docs" // 不要忘了导入把你上一步生成的docs
	"github.com/gin-gonic/gin"
	knife4goFiles "github.com/go-webtools/knife4go"
	knife4goGin "github.com/go-webtools/knife4go/gin"
	//swaggerFiles "github.com/swaggo/files"
	//ginSwagger "github.com/swaggo/gin-swagger"
)

func DocsInit(router *gin.Engine) {
	//router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/docs/*any", knife4goGin.WrapHandler(knife4goFiles.Handler))
}

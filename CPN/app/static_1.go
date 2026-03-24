//go:build ignore

package app

import (
	"LadderCompetitionPlatform/config"
	"LadderCompetitionPlatform/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path"
)

func StaticAppInit(router *gin.Engine) {
	router.GET("/files/:filename", fileDownload)
}

// @Summary 文件下载
// @Tags Files
// @Produce json
// @Param filename path string true "文件名"
// @Success 200 {object} []byte
// @Failure 500 {object} map[string]any "该文件不存在"
// @Failure 500 {object} map[string]any "下载失败"
// @Router /files/{filename} [get]
func fileDownload(context *gin.Context) {
	fileName := context.Param("filename")

	file, err := os.Open(path.Join(config.Config.Path.FilePath, fileName))
	if err != nil {
		if os.IsNotExist(err) {
			context.JSON(500, utils.ErrorResponse(nil, "该文件不存在"))
			return
		} else {
			utils.Logger.Errorf("FileName: %s\nError: %s", fileName, err)
			context.JSON(500, utils.ErrorResponse(nil, "下载失败"))
			return
		}
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		utils.Logger.Errorf("FileName: %s\nError: %s", fileName, err)
		context.JSON(500, utils.ErrorResponse(nil, "下载失败"))
		return
	}

	filesize := fileInfo.Size()
	buffer := make([]byte, filesize)

	file.Read(buffer)

	context.Writer.WriteHeader(http.StatusOK)
	context.Header("Content-Disposition", "attachment; filename="+fileName)
	context.Header("Content-Type", "application/text/plain")
	context.Header("Accept-Length", fmt.Sprintf("%d", len(buffer)))
	context.Writer.Write(buffer)
}

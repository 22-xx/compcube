package utils

import (
	"LadderCompetitionPlatform/config"
	"bytes"
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"strings"
	"time"
)

var Logger *logrus.Logger

type MyFormatter struct {
}

func (m *MyFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	logLevel := strings.ToUpper(entry.Level.String())
	var newLog string
	if logLevel == "ERROR" {
		//fName := filepath.Base(entry.Caller.File)
		newLog = fmt.Sprintf("\n%s [%s] [%s:%d]\nMethod: %s\n%s\n",
			logLevel, timestamp, entry.Caller.File, entry.Caller.Line, entry.Caller.Function, entry.Message)
	} else {
		newLog = fmt.Sprintf("%s  [%s] %s\n", logLevel, timestamp, entry.Message)
	}

	b.WriteString(newLog)
	return b.Bytes(), nil
}

func newLfsHook(logName string) logrus.Hook {
	writer, err := rotatelogs.New(
		logName+"-%Y%m%d.log",
		// WithLinkName为最新的日志建立软连接,以方便随着找到当前日志文件
		rotatelogs.WithLinkName(logName),
		// WithRotationTime设置日志分割的时间
		rotatelogs.WithRotationTime(time.Hour*24),
		// WithRotationCount设置文件清理前最多保存的个数.
		rotatelogs.WithRotationCount(365),
	)

	if err != nil {
		Logger.Errorf("config local file system for logger error: %v", err)
	}

	// 使用了lfshook软件包创建了一个新的日志钩子，该钩子将日志记录到指定的日志文件中。
	// lfshook.WriterMap指定了每个日志级别所使用的写入器（writer）。
	// 在这个函数中，所有的日志级别都使用同一个写入器writer。
	lfsHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: writer,
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	}, &MyFormatter{})

	return lfsHook
}

func LoggerInit() {

	logFilePath := config.Config.Path.LogPath
	logFileName := config.Config.Path.LogName

	_, err := os.Stat(logFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			errDir := os.MkdirAll(logFilePath, 0755)
			if errDir != nil {
				fmt.Println("创建文件夹出错：" + errDir.Error())
			}
		}
	}

	fileName := path.Join(logFilePath, logFileName)
	Logger = logrus.New()
	Logger.SetLevel(logrus.InfoLevel)
	Logger.SetReportCaller(true)
	Logger.AddHook(newLfsHook(fileName))
}

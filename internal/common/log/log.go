package log

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"time"

	"github.com/TwiN/go-color"
	"github.com/YiNNx/WeVote/internal/config"
	nested "github.com/antonfisher/nested-logrus-formatter"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

func InitLogger() {
	Logger = getLogger(config.C.Log.Path, config.C.Server.DebugMode)
	Logger.Info("logger started")
}

func getLogger(logPath string, debug bool) *logrus.Logger {
	logger := logrus.New()
	if debug {
		logger.SetReportCaller(true)
	} else {
		logger.SetReportCaller(false)
	}
	logger.SetFormatter(formatter())

	logger.SetLevel(logrus.DebugLevel)

	wd, _ := os.Getwd()
	logDir := path.Join(wd, logPath)
	os.MkdirAll(logDir, os.ModePerm)
	baseLogPath := path.Join(wd, logPath, "log")
	writer, _ := rotatelogs.New(
		baseLogPath+"-%Y-%m-%d",
		rotatelogs.WithLinkName(baseLogPath),
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	lfHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: writer,
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	}, &logrus.JSONFormatter{})

	logger.AddHook(lfHook)
	return logger
}

func formatter() *nested.Formatter {
	fmtter := &nested.Formatter{
		HideKeys:        true,
		TimestampFormat: "15:04:05",
		CallerFirst:     false,
		CustomCallerFormatter: func(frame *runtime.Frame) string {
			funcInfo := runtime.FuncForPC(frame.PC)
			if funcInfo == nil {
				return "error during runtime.FuncForPC"
			}
			fullPath, line := funcInfo.FileLine(frame.PC)
			return fmt.Sprintf(color.InBlue(" ⇨ %v (line%v)"), filepath.Base(fullPath), line)
		},
	}
	fmtter.NoColors = false
	return fmtter
}

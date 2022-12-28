package logs

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/partyhall/partyhall/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.SugaredLogger

func Init() {
	isDebug := strings.HasPrefix(os.Args[0], "/tmp/")
	if isDebug {
		isDebug = os.Getenv("ENV") != "prod"
	}

	// Because the RPi does not have hwclock we need to do this stupid thing
	// If it had we could just use the current date & symlink latest
	logpath := utils.GetPath("latest.log")
	if _, err := os.Stat(logpath); !os.IsNotExist(err) {
		date := time.Now().Format("2006-01-02_14-04-05")
		err := os.Rename(logpath, utils.GetPath(fmt.Sprintf("logs_before_%v.log", date)))
		if err != nil {
			panic(err)
		}
	}

	pe := zap.NewProductionEncoderConfig()

	fileEncoder := zapcore.NewJSONEncoder(pe)

	pe.EncodeTime = zapcore.ISO8601TimeEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(pe)

	level := zap.InfoLevel
	if isDebug {
		level = zap.DebugLevel
	}

	var core *zapcore.Core

	if !isDebug {
		file, err := os.OpenFile(logpath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			fmt.Println("FAILED TO OPEN LOGFILE, ONLY PRINTING TO THE CONSOLE: ", err)
			c := zapcore.NewTee(zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level))
			core = &c
		} else {
			c := zapcore.NewTee(
				zapcore.NewCore(fileEncoder, zapcore.AddSync(file), level),
				zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level),
			)

			core = &c
		}
	} else {
		fmt.Println("Not running in prod: no logfile")
		c := zapcore.NewTee(zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level))
		core = &c
	}

	l := zap.New(*core)

	logger = l.Sugar()
}

// #region Error

func Error(args ...interface{}) {
	logger.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	logger.Errorf(template, args...)
}

func Errorln(args ...interface{}) {
	logger.Errorln(args...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	logger.Errorw(msg, keysAndValues...)
}

//#endregion

// #region Warn

func Warn(args ...interface{}) {
	logger.Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	logger.Warnf(template, args...)
}

func Warnln(args ...interface{}) {
	logger.Warnln(args...)
}

func Warnw(msg string, keysAndValues ...interface{}) {
	logger.Warnw(msg, keysAndValues...)
}

//#endregion

// #region Info

func Info(args ...interface{}) {
	logger.Info(args...)
}

func Infof(template string, args ...interface{}) {
	logger.Infof(template, args...)
}

func Infoln(args ...interface{}) {
	logger.Infoln(args...)
}

func Infow(msg string, keysAndValues ...interface{}) {
	logger.Infow(msg, keysAndValues...)
}

//#endregion

// #region Debug

func Debug(args ...interface{}) {
	logger.Debug(args...)
}

func Debugf(template string, args ...interface{}) {
	logger.Debugf(template, args...)
}

func Debugln(args ...interface{}) {
	logger.Infoln(args...)
}

func Debugw(msg string, keysAndValues ...interface{}) {
	logger.Debugw(msg, keysAndValues...)
}

//#endregion

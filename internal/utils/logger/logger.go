package logger

import (
	"fmt"
	"os"
	"time"
)

const (
	RED    = "\x1b[31m"
	GREEN  = "\x1b[32m"
	YELLOW = "\x1b[33mv"
	RESET  = "\x1b[0m"
)

type Logger struct {
	LogDir     string
	DateFormat string
}

func NewLogger(logDir string, dateFormat string) *Logger {
	return &Logger{
		LogDir:     logDir,
		DateFormat: dateFormat,
	}
}

func (l *Logger) CreateLogger() {
	if _, err := os.Stat(l.LogDir); os.IsNotExist(err) {
		err := os.Mkdir(l.LogDir, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	actualDate := time.Now()
	fileName := actualDate.Format(l.DateFormat)
	file, err := os.OpenFile(l.LogDir+"/"+fileName+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o666)
	if err != nil {
		panic(err)
	}
	fmt.Println(GREEN + "[INFO: " + actualDate.Format(l.DateFormat) + "] Logger created successfully" + RESET)
	fileRes := fmt.Sprintf("date:%s,type:success,message:Successfully created a logger,data:%v", fileName, map[string]interface{}{})
	file.Write([]byte(fileRes))
}

func (l *Logger) Info(msg string, data ...any) {
	actualDate := time.Now()
	fileName := actualDate.Format(l.DateFormat)
	fmt.Print(GREEN + "[INFO: " + fileName + "] " + msg)
	if len(data) > 0 {
		fmt.Print(" ")
		for _, d := range data {
			fmt.Print(d, " ")
		}
	}
	fmt.Println(RESET)
	file, err := os.OpenFile(l.LogDir+"/"+fileName+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o666)
	if err != nil {
		panic(err)
	}
	fileRes := fmt.Sprintf("date:%s,type:info,message:%s,data:%v", fileName, msg, data)
	file.Write([]byte(fileRes))
}

func (l *Logger) Warn(msg string, data ...any) {
	actualDate := time.Now()
	fileName := actualDate.Format(l.DateFormat)
	fmt.Print(YELLOW + "[WARN: " + fileName + "] " + msg)
	if len(data) > 0 {
		fmt.Print(" ")
		for _, d := range data {
			fmt.Print(d, " ")
		}
	}
	fmt.Println(RESET)
	file, err := os.OpenFile(l.LogDir+"/"+fileName+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o666)
	if err != nil {
		panic(err)
	}
	fileRes := fmt.Sprintf("date:%s,type:warn,message:%s,data:%v", fileName, msg, data)
	file.Write([]byte(fileRes))
}

func (l *Logger) Error(msg string, data ...any) {
	actualDate := time.Now()
	fileName := actualDate.Format(l.DateFormat)
	fmt.Print(RED + "[ERROR: " + fileName + "] " + msg)
	if len(data) > 0 {
		fmt.Print(" ")
		for _, d := range data {
			fmt.Print(d, " ")
		}
	}
	fmt.Println(RESET)
	file, err := os.OpenFile(l.LogDir+"/"+fileName+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o666)
	if err != nil {
		panic(err)
	}
	fileRes := fmt.Sprintf("date:%s,type:error,message:%s,data:%v", fileName, msg, data)
	file.Write([]byte(fileRes))
}

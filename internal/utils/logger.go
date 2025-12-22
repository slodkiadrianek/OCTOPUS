package utils

import (
	"fmt"
	"os"
	"time"
)

type LoggerService interface {
	InitializeLogger()
	Info(msg string, data ...any)
	Warn(msg string, data ...any)
	Error(msg string, data ...any)
	Close() error
}

const (
	red    = "\x1b[31m"
	green  = "\x1b[32m"
	yellow = "\x1b[33m"
	reset  = "\x1b[0m"
)

type Logger struct {
	logDir     string
	dateFormat string
	file       *os.File
	startTime  string
}

func NewLogger(logDir string, dateFormat string) *Logger {
	return &Logger{
		logDir:     logDir,
		dateFormat: dateFormat,
	}
}

func (l *Logger) getActualDate() string {
	actualDate := time.Now()
	year := actualDate.Year()
	month := int(actualDate.Month())
	day := actualDate.Day()
	actualDateFormat := fmt.Sprintf("%d.%d.%d", day, month, year)
	return actualDateFormat
}

func (l *Logger) printDataToTheConsole(data ...any) {
	if len(data) > 0 {
		fmt.Print(" ")
		for _, d := range data {
			fmt.Print(d, " ")
		}
	}
}

func (l *Logger) InitializeLogger() {
	if _, err := os.Stat(l.logDir); os.IsNotExist(err) {
		if err := os.Mkdir(l.logDir, os.ModePerm); err != nil {
			panic(err)
		}
	}
	actualDate := time.Now()
	l.startTime = l.getActualDate()
	fileName := l.getActualDate()
	logTime := actualDate.Format("2006-01-02 15:04:05")

	file, err := os.OpenFile(l.logDir+"/"+fileName+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o666)
	if err != nil {
		panic(err)
	}

	l.file = file

	fmt.Println(green + "[INFO: " + actualDate.Format(l.dateFormat) + "] Logger created successfully" + reset)

	fileRes := fmt.Sprintf("date:%s,type:success,message:Successfully created a logger,data:%v\n", logTime,
		map[string]any{})

	_, err = l.file.Write([]byte(fileRes))
	if err != nil {
		fmt.Println("Something went wrong during writing to data to the file")
	}
}

func (l *Logger) Info(msg string, data ...any) {
	l.Validate()

	actualDate := time.Now()
	fileName := actualDate.Format(l.dateFormat)
	logTime := actualDate.Format("2006-01-02 15:04:05")

	fmt.Println(green + "[INFO: " + logTime + "] " + msg)

	l.printDataToTheConsole(data...)

	fmt.Println(reset)

	fileRes := fmt.Sprintf("date:%s,type:info,message:%s,data:%v\n", fileName, msg, data)

	_, err := l.file.WriteString(fileRes)
	if err != nil {
		fmt.Println("Something went wrong during writing to data to the file", err)
	}
}

func (l *Logger) Warn(msg string, data ...any) {
	l.Validate()

	actualDate := time.Now()
	fileName := actualDate.Format(l.dateFormat)
	logTime := actualDate.Format("2006-01-02 15:04:05")

	fmt.Print(yellow + "[WARN: " + logTime + "] " + msg)

	l.printDataToTheConsole(data...)

	fmt.Println(reset)

	fileRes := fmt.Sprintf("date:%s,type:warn,message:%s,data:%v\n", fileName, msg, data)

	_, err := l.file.Write([]byte(fileRes))
	if err != nil {
		fmt.Println("Something went wrong during writing to data to the file")
	}
}

func (l *Logger) Error(msg string, data ...any) {
	l.Validate()

	actualDate := time.Now()
	fileName := actualDate.Format(l.dateFormat)
	logTime := actualDate.Format("2006-01-02 15:04:05")

	fmt.Print(red + "[ERROR: " + logTime + "] " + msg)

	l.printDataToTheConsole(data...)

	fmt.Println(reset)

	fileRes := fmt.Sprintf("date:%s,type:error,message:%s,data:%v\n", fileName, msg, data)
	_, err := l.file.Write([]byte(fileRes))
	if err != nil {
		fmt.Println("Something went wrong during writing to data to the file")
	}
}

func (l *Logger) Validate() {
	actualDateFormat := l.getActualDate()
	if actualDateFormat != l.startTime {

		fmt.Println("Closing old file and creating the new one for new date")

		err := l.file.Close()
		if err != nil {
			fmt.Println("Something went wrong during writing to data to the file")
		}

		l.startTime = actualDateFormat
		fileName := actualDateFormat

		file, err := os.OpenFile(l.logDir+"/"+fileName+".log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o666)
		if err != nil {
			panic(err)
		}

		l.file = file
	}
}

func (l *Logger) Close() error {
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

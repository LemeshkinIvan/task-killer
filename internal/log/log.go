package log

import (
	"fmt"
	"os"
	"time"
)

type LogLevel string

const INFO LogLevel = "INFO"
const WARN LogLevel = "WARN"
const FATAL LogLevel = "FATAL"

type Color string

const BLUE Color = "\033[34m"
const YELLOW Color = "\033[33m"
const RED Color = "\033[31m"

const Close string = "\033[0m"

type Logger struct {
	timeMask         string
	DefaultLogFolder string
	DefaultNameFile  string
	IsDebug          bool
	EnableWriteLog   bool
}

type LoggerCfg struct {
	LogFolder      string
	NameFile       string
	IsDebug        bool
	EnableWriteLog bool
}

func NewLogger(cfg LoggerCfg) (*Logger, error) {
	if cfg.NameFile == "" {
		cfg.NameFile = "app.log"
	}

	if cfg.EnableWriteLog {
		// create folder and open file
		if cfg.LogFolder == "" {
			path, err := setDefPath()
			if err != nil {
				return nil, err
			}

			cfg.LogFolder = path
		}
	}

	return &Logger{
		timeMask:         "02-01-2006 15:04:05",
		DefaultLogFolder: cfg.LogFolder,
		DefaultNameFile:  cfg.NameFile,
		IsDebug:          cfg.IsDebug,
		EnableWriteLog:   cfg.EnableWriteLog,
	}, nil
}

var logChan = make(chan string, 100) // буфер!
var name = os.Getenv("USERNAME")

func (l *Logger) Log(msg string, level LogLevel) {
	if level == "" {
		level = INFO
	}

	var header string
	switch level {
	case INFO:
		header = string(BLUE) + string(INFO) + Close
	case WARN:
		header = string(YELLOW) + string(WARN) + Close
	case FATAL:
		header = string(RED) + string(FATAL) + Close
	default:
	}

	if l.IsDebug {
		fmt.Println(l.formatMsg(msg, header))
	}

	if l.EnableWriteLog {
		//fmt.Println(l.formatMsg(rawMsg))
	}
}

func (l *Logger) formatMsg(raw string, header string) string {
	var msg string
	msg = time.Now().Format(l.timeMask) + " " + header

	if msg[len(msg)-1] == '\n' {
		msg = msg[:len(msg)-1]
	}
	msg += " " + raw
	return msg
}

func (l *Logger) LogInFile(fileName string) {
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	for msg := range logChan {
		file.WriteString(msg + "\n")
	}
}

func (l *Logger) CreateLogFolder(path string) error {
	if len(path) == 0 {
		path = fmt.Sprintf(l.DefaultLogFolder, name)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.MkdirAll(path, 0755); err != nil {
			return err
		}
	}

	return nil
}

func (l *Logger) CreateLogFile() error {
	_, err := os.Create(l.DefaultLogFolder + "\\" + l.DefaultNameFile)
	if err != nil {
		return err
	}
	return nil
}

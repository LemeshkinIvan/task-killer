package log

import (
	"fmt"
	"io"
	"os"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
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

// terrible
var headers map[LogLevel]string = map[LogLevel]string{
	INFO:  string(BLUE) + string(INFO) + Close,
	WARN:  string(YELLOW) + string(WARN) + Close,
	FATAL: string(RED) + string(FATAL) + Close,
}

type Logger struct {
	timeMask       string
	LogFolder      string
	NameFile       string
	IsDebug        bool
	EnableWriteLog bool

	Writer io.Writer
	Ch     chan string
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

	l := &Logger{
		timeMask:       "02-01-2006 15:04:05",
		NameFile:       cfg.NameFile,
		IsDebug:        cfg.IsDebug,
		EnableWriteLog: cfg.EnableWriteLog,
	}

	if cfg.EnableWriteLog {
		if cfg.LogFolder == "" {
			path, err := setDefPath()
			if err != nil {
				return nil, err
			}

			cfg.LogFolder = path
		}

		l.LogFolder = cfg.LogFolder

		err := l.initLogFolder()
		if err != nil {
			return nil, err
		}

		path := l.LogFolder + "\\" + l.NameFile

		if cfg.IsDebug {
			l.writeToStdOut(path, string(INFO))
		}

		l.Writer = &lumberjack.Logger{
			Filename:   path,
			MaxSize:    10, // MB
			MaxBackups: 5,
			MaxAge:     7,    // дней
			Compress:   true, // gzip
		}

		l.Ch = make(chan string, 100)

		if cfg.EnableWriteLog && l.Writer == nil {
			return nil, fmt.Errorf("logger writer is nil")
		}

		go l.listen()
	}

	return l, nil
}

func (l *Logger) listen() {
	for msg := range l.Ch {
		if l.Writer == nil {
			continue
		}

		_, err := l.Writer.Write([]byte(msg + "\n"))
		if err != nil {
			l.writeToStdOut(err.Error(), string(WARN))
		}
	}
}

func (l *Logger) Close() {
	if l.EnableWriteLog {
		close(l.Ch)
	}
}

func (l *Logger) initLogFolder() error {
	if len(l.LogFolder) == 0 {
		return fmt.Errorf("cant create log folder: path is empty")
	}

	if _, err := os.Stat(l.LogFolder); os.IsNotExist(err) {
		if err := os.MkdirAll(l.LogFolder, 0755); err != nil {
			return err
		}
	}

	return nil
}

func (l *Logger) Log(msg string, level LogLevel) {
	if level == "" {
		level = INFO
	}

	header, ok := headers[level]
	if !ok {
		header = headers[INFO]
	}

	if l.IsDebug {
		l.writeToStdOut(msg, header)
	}

	if !l.EnableWriteLog || l.Ch == nil {
		return
	}

	select {
	case l.Ch <- string(level) + " " + msg:
	default:
		l.writeToStdOut("log dropped (channel full)", string(WARN))
	}
}

func (l *Logger) writeToStdOut(msg string, header string) {
	fmt.Println(l.formatMsg(msg, header))
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

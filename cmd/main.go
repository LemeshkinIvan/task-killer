package main

import (
	"errors"
	"fmt"
	"os"

	"task-killer/internal/cfg"
	"task-killer/internal/constants"
	"task-killer/internal/log"
	w "task-killer/internal/watcher"
	"time"
)

const (
	path    = "../cfg/cfg.json"
	devPath = "../cfg/cfg.json"

	defaultTimeIdle    = 10 * time.Second // сколько ждать если getConfig дал ошибку
	defaultTimeRequest = 2 * time.Second  // при истечении пойдет за файлом

	// switch true if you need stdout log
	isDebug        = false
	enableWriteLog = false
)

func main() {
	logger, err := log.NewLogger(log.LoggerCfg{
		IsDebug:        isDebug,
		EnableWriteLog: enableWriteLog,
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}

	logger.Log(constants.StartProgram, log.INFO)

	for {
		var config *cfg.ConfigDTO
		var err error

		logger.Log(constants.GetConfig, log.INFO)
		for config == nil {
			config, err = cfg.GetConfig(devPath)

			if err != nil {
				logger.Log(err.Error(), log.WARN)
			}

			time.Sleep(defaultTimeRequest)
		}

		logger.Log(constants.ConfigIsLoaded, log.INFO)

		sleepDur, err := time.ParseDuration(config.TimeSleep)
		if err != nil {
			sleepDur = defaultTimeIdle
			logger.Log(err.Error(), log.WARN)
			logger.Log(constants.SetDefaultSleepTime, log.WARN)
		}

		watcher, err := w.NewWin32Watcher(w.WatcherInit{
			Log:       logger,
			IsDebug:   isDebug,
			Blacklist: config.Blacklist,
		})
		if err != nil {
			logger.Log(err.Error(), log.FATAL)
			os.Exit(-1)
		}

		if err := watcher.StartWatcherWin32(); err != nil {
			if errors.Is(err, w.ErrBlacklistLen) {
				logger.Log(err.Error(), log.WARN)
			} else {
				logger.Log(err.Error(), log.FATAL)
				os.Exit(-1)
			}
		}

		logger.Log(constants.GetSleepingMsg(config.TimeSleep), log.INFO)
		time.Sleep(sleepDur)
	}
}

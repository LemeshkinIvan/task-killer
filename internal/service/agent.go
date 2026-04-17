package service

import (
	"encoding/json"
	"errors"
	"fmt"
	appmsg "task-killer/internal/app_msg"
	"task-killer/internal/cli"
	cfg "task-killer/internal/config"
	"task-killer/internal/domain/metrics"
	"task-killer/internal/log"
	w "task-killer/internal/watcher/win"
	"time"
)

const (
	defaultTimeIdle    = 10 * time.Second // сколько ждать при завершении интерации watcher
	defaultTimeRequest = 2 * time.Second  // при истечении пойдет за файлом
)

// orchestrator
type AgentService struct {
	cfg           *cfg.ConfigMutationDTO
	configManager *cfg.ConfigManager
	logger        *log.Logger
	watcher       *w.Win32Watcher
	args          *cli.CMDFlags
}

func NewAgentService(
	args *cli.CMDFlags,
	configManager *cfg.ConfigManager,
	logger *log.Logger,
	watcher *w.Win32Watcher,
) (*AgentService, error) {
	if args == nil {
		return nil, fmt.Errorf("cmd flags is nil")
	}

	if configManager == nil {
		return nil, fmt.Errorf("configManager is nil")
	}

	if logger == nil {
		return nil, fmt.Errorf("logger is nil")
	}

	if watcher == nil {
		return nil, fmt.Errorf("watcher is nil")
	}

	return &AgentService{
		args:          args,
		configManager: configManager,
		logger:        logger,
		watcher:       watcher,
	}, nil
}

func (s *AgentService) Run() error {
	s.logger.Log(appmsg.StartAgent, log.INFO)

	var i uint64

	for {
		i++

		watcherMetrics, err := s.runIterator()
		if err != nil {
			return err
		}

		iteratorMetrics := metrics.IteratorMetrics{
			Iteration: i,
			Watcher:   watcherMetrics,
		}

		s.logMetrics(iteratorMetrics)
		s.sleep()
	}
}

func (s *AgentService) logMetrics(metrics metrics.IteratorMetrics) {
	line := fmt.Sprintf(
		"iteration: %d | protected: %d | killed: %d | skipped: %d | failKilled: %d",
		metrics.Iteration,
		metrics.Watcher.Protected,
		metrics.Watcher.Killed,
		metrics.Watcher.Skipped,
		metrics.Watcher.FailedKilled,
	)

	b, _ := json.MarshalIndent(line, "", "  ")
	s.logger.Log(string(b[:]), log.INFO)
}

func (s *AgentService) runIterator() (metrics.WatcherMetrics, error) {
	config, err := s.getConfigWithRetry()
	if err != nil {
		return metrics.WatcherMetrics{}, err
	}

	if config == nil {
		return metrics.WatcherMetrics{}, fmt.Errorf("config is nil")
	}

	s.cfg = config

	res, err := s.runWatcher(config.Blacklist)
	if err != nil {
		if errors.Is(err, w.ErrBlacklistLen) {
			s.logger.Log(err.Error(), log.WARN)
			return metrics.WatcherMetrics{}, err
		} else {
			s.logger.Log(err.Error(), log.FATAL)
			return metrics.WatcherMetrics{}, err
		}
	}

	return res, nil
}

func (s *AgentService) sleep() {
	sleepDur, err := time.ParseDuration(s.cfg.TimeSleep)
	if err != nil {
		sleepDur = defaultTimeIdle
		s.logger.Log(err.Error(), log.WARN)
		s.logger.Log(appmsg.SetDefaultSleepTime, log.WARN)
	}

	//s.logger.Log(appmsg.GetSleepingMsg(s.cfg.TimeSleep), log.INFO)
	time.Sleep(sleepDur)

}

func (s *AgentService) getConfigWithRetry() (*cfg.ConfigMutationDTO, error) {
	var config *cfg.ConfigMutationDTO
	var err error

	for config == nil {
		config, err = s.configManager.GetConfig()

		if err != nil {
			s.logger.Log(err.Error(), log.WARN)
		}

		time.Sleep(defaultTimeRequest)
	}
	return config, nil
}

func (s *AgentService) runWatcher(blacklist []string) (metrics.WatcherMetrics, error) {
	return s.watcher.StartWatcher(blacklist)
}

func (s *AgentService) Shutdown() {
	if s.args.EnableLogFile {
		s.logger.Close()
	}
}

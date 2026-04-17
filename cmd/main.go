package main

import (
	"fmt"
	"os"

	"task-killer/internal/cli"
	cfg "task-killer/internal/config"
	"task-killer/internal/log"
	p "task-killer/internal/provider"
	"task-killer/internal/service"
	"task-killer/internal/watcher/win"
)

const exitFailure = -1

func main() {
	agent, err := initAgent()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(exitFailure)
	}

	defer agent.Shutdown()

	if err := agent.Run(); err != nil {
		fmt.Println(err.Error())
		os.Exit(exitFailure)
	}
}

func initAgent() (*service.AgentService, error) {
	args, err := cli.GetCMDFlags()
	if err != nil {
		return nil, err
	}

	logger, err := log.NewLogger(log.LoggerCfg{
		IsDebug:        args.IsDebug,
		EnableWriteLog: args.EnableLogFile,
	})
	if err != nil {
		return nil, err
	}

	configManager, err := initConfigManager(args)
	if err != nil {
		return nil, err
	}

	watcher, err := win.NewWin32Watcher(win.WatcherInit{
		Log:     logger,
		IsDebug: false,
	})
	if err != nil {
		return nil, err
	}

	service, err := service.NewAgentService(args, configManager, logger, watcher)
	if err != nil {
		return nil, err
	}

	return service, nil
}

func initConfigManager(args *cli.CMDFlags) (*cfg.ConfigManager, error) {
	if args == nil {
		return nil, fmt.Errorf("args is nil")
	}

	var provider p.Provider
	var err error

	switch args.Conn {
	case cli.Local:
		provider, err = p.NewLocalProvider(args.Path)
	case cli.SMB:
		provider, err = p.NewSMBManager(p.SMBInput{
			Addr:     args.Path,
			User:     "",
			Password: "",
			Domain:   "",
		})
	case cli.HTTP:
		provider, err = p.NewHTTPProvider(args.Path)
	default:
		return nil, fmt.Errorf("unknown provider type")
	}

	if err != nil {
		return nil, err
	}

	return cfg.NewConfigManager(provider), nil
}

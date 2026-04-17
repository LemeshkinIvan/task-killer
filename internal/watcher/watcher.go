package watcher

// as crossplatform idea
// not important
type Watcher interface {
	StartWatcher(blacklist []string) error
	GetProcessPath(pid uint32) (string, error)
	KillProcess(pid uint32) error
}

package metrics

type WatcherMetrics struct {
	Protected    int32
	Skipped      int32
	Killed       int32
	FailedKilled int32
}

type IteratorMetrics struct {
	Iteration uint64
	Watcher   WatcherMetrics
}

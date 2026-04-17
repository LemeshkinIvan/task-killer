package appmsg

const (
	StartAgent          = "start agent"
	GetConfig           = "start getConfig"
	ConfigIsLoaded      = "config is loaded successfully"
	SetDefaultSleepTime = "set sleep time"
	StartSleeping       = "Program is over cycle. Sleeping, dur:"
)

// func GetSleepingMsg(dur string) string {
// 	return StartAgent + " " + dur
// }

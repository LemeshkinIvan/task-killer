package constants

const (
	StartProgram        = "start program"
	GetConfig           = "start getConfig"
	ConfigIsLoaded      = "config is loaded successfully"
	SetDefaultSleepTime = "set sleep time"
	StartSleeping       = "Program is over cycle. Sleeping, dur:"
)

func GetSleepingMsg(dur string) string {
	return StartProgram + " " + dur
}

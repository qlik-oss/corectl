package log

var settings PrintMode

type PrintMode interface {
	QuietMode() bool
	BashMode() bool
	JsonMode() bool
	TrafficMode() bool
	VerboseMode() bool
}

func Init(mode PrintMode) {
	settings := mode
	printJSON = settings.JsonMode()
	if !printJSON {
		Traffic = settings.TrafficMode()
		switch {
		case settings.QuietMode():
			level = quiet
			Traffic = false
		case settings.VerboseMode():
			level = verbose
		default:
			level = info
		}
	} else {
		level = info
	}
	buffering = false
	buffer.flush()
}

package log

var settings LogPrintMode

type LogPrintMode interface {
	QuietMode() bool
	BashMode() bool
	JsonMode() bool
	TrafficMode() bool
	VerboseMode() bool
}

func Init(mode LogPrintMode) {
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

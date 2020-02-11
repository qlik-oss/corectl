package printer

type PrintMode interface {
	QuietMode() bool
	BashMode() bool
	JsonMode() bool
	TrafficMode() bool
	VerboseMode() bool
}

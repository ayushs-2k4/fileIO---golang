package logger

type Level int

const (
	Debug Level = iota
	Error
	Warn
	Info
)

func (l Level) String() string {
	switch l {
	case Error:
		return "ERROR"
	case Warn:
		return "WARN"
	case Debug:
		return "DEBUG"
	case Info:
		return "INFO"
	default:
		return "N/A"
	}
}

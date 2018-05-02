package hrlog

type Level uint32

// These are the different logging levels. You can set the logging Level to log
// on your instance of logger, obtained with `hrlog.New()`
const (
	//PanicLevel Level, highest Level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	LevelPanic Level = iota
	// FatalLevel Level. Logs and then calls `os.Exit(1)`. It will exit even if the
	// logging Level is set to Panic.
	LevelFatal
	// ErrorLevel Level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	LevelError
	// WarnLevel Level. Non-critical entries that deserve eyes.
	LevelWarn
	// InfoLevel Level. General operational entries about what's going on inside the
	// application.
	LevelInfo
	// LevelDebug Level. Usually only enabled when debugging. Very verbose logging.
	LevelDebug
)

const (
	TagDebug = "[DEBUG]: "
	TagInfo  = "[INFO]: "
	TagWarn  = "[WARN]: "
	TagError = "[ERROR]: "
	TagFatal = "[FATAL]: "
	TagPanic = "[PANIC]: "
)

func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "LevelDebug"
	case LevelInfo:
		return "LevelInfo"
	case LevelWarn:
		return "LevelWarn"
	case LevelError:
		return "LevelError"
	case LevelFatal:
		return "LevelFatal"
	case LevelPanic:
		return "LevelPanic"
	}
	return "Unknow"
}

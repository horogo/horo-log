package hrlog

type level uint32

// These are the different logging levels. You can set the logging level to log
// on your instance of logger, obtained with `hrlog.New()`
const (
	//PanicLevel level, highest level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	LevelPanic level = iota
	// FatalLevel level. Logs and then calls `os.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	LevelFatal
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	LevelError
	// WarnLevel level. Non-critical entries that deserve eyes.
	LevelWarn
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	LevelInfo
	// LevelDebug level. Usually only enabled when debugging. Very verbose logging.
	LevelDebug
)

const (
	TagDebug = "[DEBUG]: "
	TagInfo = "[INFO]: "
	TagWarn = "[WARN]: "
	TagError = "[ERROR]: "
	TagFatal = "[FATAL]: "
	TagPanic = "[PANIC]: "
)
package hrlog

import (
	"io"
	"sync"
	"os"
	"sync/atomic"
	"fmt"
	"time"
	"runtime"
)

// A Logger represents an active logging object. Multiple loggers can be used
// simultaneously even if they are using the same same writers
type Logger struct {
	out    io.Writer
	level  level
	flag   uint32
	prefix string
	buf    []byte
	mu     sync.Mutex
}

var std = New()

func New() *Logger {
	return &Logger{
		out:   os.Stderr,
		level: LevelDebug,
		flag:  LstdFlags,
	}
}

// Cheap integer to fixed-width decimal ASCII. Give a negative width to avoid zero-padding.
func itoa(buf *[]byte, i int, wid int) {
	// Assemble decimal in reverse order.
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	*buf = append(*buf, b[bp:]...)
}

// formatHeader writes log header to buf in following order:
//   * l.prefix (if it's not blank),
//   * date and/or time (if corresponding flags are provided),
//   * file and line number (if corresponding flags are provided).
func (l *Logger) formatHeader(buf *[]byte, t time.Time, file string, line int) {
	if l.flag&(Ldate|Ltime|Lmicroseconds) != 0 {
		if l.flag&LUTC != 0 {
			t = t.UTC()
		}
		if l.flag&Ldate != 0 {
			year, month, day := t.Date()
			itoa(buf, year, 4)
			*buf = append(*buf, '/')
			itoa(buf, int(month), 2)
			*buf = append(*buf, '/')
			itoa(buf, day, 2)
			*buf = append(*buf, ' ')
		}
		if l.flag&(Ltime|Lmicroseconds) != 0 {
			hour, min, sec := t.Clock()
			itoa(buf, hour, 2)
			*buf = append(*buf, ':')
			itoa(buf, min, 2)
			*buf = append(*buf, ':')
			itoa(buf, sec, 2)
			if l.flag&Lmicroseconds != 0 {
				*buf = append(*buf, '.')
				itoa(buf, t.Nanosecond()/1e3, 6)
			}
			*buf = append(*buf, ' ')
		}
	}
	if l.flag&(Lshortfile|Llongfile) != 0 {
		if l.flag&Lshortfile != 0 {
			short := file
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					short = file[i+1:]
					break
				}
			}
			file = short
		}
		*buf = append(*buf, file...)
		*buf = append(*buf, ':')
		itoa(buf, line, -1)
		*buf = append(*buf, ": "...)
	}
	*buf = append(*buf, l.prefix...)
}

func (l *Logger) write(depth int, s string) {
	now := time.Now()
	var file string
	var line int

	l.mu.Lock()
	defer l.mu.Unlock()

	if l.flag&(Lshortfile|Llongfile) != 0 {
		// Release lock while getting caller info - it's expensive.
		l.mu.Unlock()
		var ok bool
		_, file, line, ok = runtime.Caller(depth)
		if !ok {
			file = "???"
			line = 0
		}
		l.mu.Lock()
	}
	l.buf = l.buf[:0]
	l.formatHeader(&l.buf, now, file, line)
	l.buf = append(l.buf, s...)
	if len(s) == 0 || s[len(s)-1] != '\n' {
		l.buf = append(l.buf, '\n')
	}
	_, err := l.out.Write(l.buf)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write to log, %v\n", err)
	}
}

func (l *Logger) SetOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.out = w
}

func (l *Logger) getLevel() level {
	return level(atomic.LoadUint32((*uint32)(&l.level)))
}

func (l *Logger) SetLevel(level level) {
	atomic.StoreUint32((*uint32)(&l.level), uint32(level))
}

func (l *Logger) SetFlags(flag uint32) {
	atomic.StoreUint32((*uint32)(&l.flag), flag)
}

func (l *Logger) SetPrefix(prefix string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.prefix = prefix
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	if l.getLevel() >= LevelDebug {
		l.write(2, TagDebug+fmt.Sprintf(format, v...))
	}
}

func (l *Logger) Debug(v ...interface{}) {
	if l.getLevel() >= LevelDebug {
		l.write(2, TagDebug+fmt.Sprint(v...))
	}
}

func (l *Logger) Debugln(v ...interface{}) {
	if l.getLevel() >= LevelDebug {
		l.write(2, TagDebug+fmt.Sprintln(v...))
	}
}

func (l *Logger) Infof(format string, v ...interface{}) {
	if l.getLevel() >= LevelInfo {
		l.write(2, TagInfo+fmt.Sprintf(format, v...))
	}
}

func (l *Logger) Info(v ...interface{}) {
	if l.getLevel() >= LevelInfo {
		l.write(2, TagInfo+fmt.Sprint(v...))
	}
}

func (l *Logger) Infoln(v ...interface{}) {
	if l.getLevel() >= LevelInfo {
		l.write(2, TagInfo+fmt.Sprintln(v...))
	}
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	if l.getLevel() >= LevelWarn {
		l.write(2, TagWarn+fmt.Sprintf(format, v...))
	}
}

func (l *Logger) Warn(v ...interface{}) {
	if l.getLevel() >= LevelWarn {
		l.write(2, TagWarn+fmt.Sprint(v...))
	}
}

func (l *Logger) Warnln(v ...interface{}) {
	if l.getLevel() >= LevelWarn {
		l.write(2, TagWarn+fmt.Sprintln(v...))
	}
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	if l.getLevel() >= LevelError {
		l.write(2, TagError+fmt.Sprintf(format, v...))
	}
}

func (l *Logger) Error(v ...interface{}) {
	if l.getLevel() >= LevelError {
		l.write(2, TagError+fmt.Sprint(v...))
	}
}

func (l *Logger) Errorln(v ...interface{}) {
	if l.getLevel() >= LevelError {
		l.write(2, TagError+fmt.Sprintln(v...))
	}
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	if l.getLevel() >= LevelFatal {
		l.write(2, TagFatal+fmt.Sprintf(format, v...))
		os.Exit(1)
	}
}

func (l *Logger) Fatal(v ...interface{}) {
	if l.getLevel() >= LevelFatal {
		l.write(2, TagFatal+fmt.Sprint(v...))
		os.Exit(1)
	}
}

func (l *Logger) Fatalln(v ...interface{}) {
	if l.getLevel() >= LevelFatal {
		l.write(2, TagFatal+fmt.Sprintln(v...))
		os.Exit(1)
	}
}

func (l *Logger) Panicf(format string, v ...interface{}) {
	s := fmt.Sprintf(format, v...)
	l.write(2, TagPanic+s)
	panic(s)
}

func (l *Logger) Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	l.write(2, TagPanic+s)
	panic(s)
}

func (l *Logger) Panicln(v ...interface{}) {
	s := fmt.Sprintln(v...)
	l.write(2, TagPanic+s)
	panic(s)
}

func SetOutput(w io.Writer) {
	std.SetOutput(w)
}

func SetLevel(level level) {
	std.SetLevel(level)
}

func SetFlags(flag uint32) {
	std.SetFlags(flag)
}

func SetPrefix(prefix string) {
	std.SetPrefix(prefix)
}

func Debugf(format string, v ...interface{}) {
	std.Debugf(format, v...)
}

func Debug(v ...interface{}) {
	std.Debug(v...)
}

func Debugln(v ...interface{}) {
	std.Debugln(v...)
}

func Infof(format string, v ...interface{}) {
	std.Infof(format, v...)
}

func Info(v ...interface{}) {
	std.Info(v...)
}

func Infoln(v ...interface{}) {
	std.Infoln(v...)
}

func Warnf(format string, v ...interface{}) {
	std.Warnf(format, v...)
}

func Warn(v ...interface{}) {
	std.Warn(v...)
}

func Warnln(v ...interface{}) {
	std.Warnln(v...)
}

func Errorf(format string, v ...interface{}) {
	std.Errorf(format, v...)
}

func Error(v ...interface{}) {
	std.Error(v...)
}

func Errorln(v ...interface{}) {
	std.Errorln(v...)
}

func Fatalf(format string, v ...interface{}) {
	std.Fatalf(format, v...)
}

func Fatal(v ...interface{}) {
	std.Fatal(v...)
}

func Fatalln(v ...interface{}) {
	std.Fatalln(v...)
}

func Panicf(format string, v ...interface{}) {
	std.Panicf(format, v...)
}

func Panic(v ...interface{}) {
	std.Panic(v...)
}

func Panicln(v ...interface{}) {
	std.Panicln(v...)
}
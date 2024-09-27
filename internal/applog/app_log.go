package applog

import (
	"log"
)

const (
	Reset     = "\033[0m"
	Black     = "\033[30m"
	Red       = "\033[31m"
	Green     = "\033[32m"
	Yellow    = "\033[33m"
	Blue      = "\033[34m"
	Magenta   = "\033[35m"
	Cyan      = "\033[36m"
	White     = "\033[37m"
	BgBlack   = "\033[40m"
	BgRed     = "\033[41m"
	BgGreen   = "\033[42m"
	BgYellow  = "\033[43m"
	BgBlue    = "\033[44m"
	BgMagenta = "\033[45m"
	BgCyan    = "\033[46m"
	BgWhite   = "\033[47m"
)

type AppLogger struct {
	log *log.Logger
}

func New(log *log.Logger) *AppLogger {
	return &AppLogger{
		log: log,
	}
}

// TODO debug mechanism to disbale debug prints.
func (l AppLogger) Debug(msg string, params ...any) {
	l.log.Println(BgGreen+White+"DEBUG:"+Reset+" "+msg, params)
}

func (l AppLogger) Info(msg string) {
	l.log.Println(BgBlue + White + "INFO:" + Reset + " " + msg)
}

func (l AppLogger) Error(msg string, err error) {
	l.log.Println(BgRed+Magenta+"ERROR:"+Reset+" "+msg+" ", err)
}

func (l AppLogger) Fatal(msg string, err error) {
	l.log.Fatal(BgRed+White+"FATAL:"+Reset+" "+msg+" ", err)
}

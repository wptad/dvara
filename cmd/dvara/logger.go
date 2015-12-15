package main

import "log"

// stdLogger provides a logger backed by the standard library logger. This is a
// placeholder until we can open source our logger.
type stdLogger struct {
	LogPrefix string
	verbose   bool
}

func (l *stdLogger) prepend(level string, format string) string {
	return level + l.LogPrefix + format
}

func (l *stdLogger) Error(args ...interface{}) { log.Printf(l.prepend("ERROR:", "%s"), args...) }
func (l *stdLogger) Errorf(format string, args ...interface{}) {
	log.Printf(l.prepend("ERROR:", format), args...)
}
func (l *stdLogger) Warn(args ...interface{}) { log.Printf(l.prepend("WARN:", "%s"), args...) }
func (l *stdLogger) Warnf(format string, args ...interface{}) {
	log.Printf(l.prepend("WARN:", format), args...)
}
func (l *stdLogger) Info(args ...interface{}) { log.Printf(l.prepend("INFO:", "%s"), args...) }
func (l *stdLogger) Infof(format string, args ...interface{}) {
	log.Printf(l.prepend("INFO:", format), args...)
}
func (l *stdLogger) Debug(args ...interface{}) {
	if l.verbose {
		log.Printf(l.prepend("DEBUG:", "%s"), args...)
	}
}
func (l *stdLogger) Debugf(format string, args ...interface{}) {
	if l.verbose {
		log.Printf(l.prepend("DEBUG:", format), args...)
	}
}

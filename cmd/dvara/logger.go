package main

import "log"

// stdLogger provides a logger backed by the standard library logger. This is a
// placeholder until we can open source our logger.
type stdLogger struct{}

func (l *stdLogger) SetPrefix(prefix string) {
	log.SetPrefix(prefix)
}

func (l *stdLogger) Error(args ...interface{})                 { log.Printf("ERROR:%s", args...) }
func (l *stdLogger) Errorf(format string, args ...interface{}) { log.Printf("ERROR:"+format, args...) }
func (l *stdLogger) Warn(args ...interface{})                  { log.Printf("WARN:%s", args...) }
func (l *stdLogger) Warnf(format string, args ...interface{})  { log.Printf("WARN:"+format, args...) }
func (l *stdLogger) Info(args ...interface{})                  { log.Printf("INFO:%s", args...) }
func (l *stdLogger) Infof(format string, args ...interface{})  { log.Printf("INFO:"+format, args...) }
func (l *stdLogger) Debug(args ...interface{})                 { log.Printf("DEBUG:%s", args...) }
func (l *stdLogger) Debugf(format string, args ...interface{}) { log.Printf("DEBUG:"+format, args...) }

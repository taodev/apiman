package logger

import (
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/google/uuid"
)

var defaultLogger *Logger

type Logger struct {
	name     string
	fullpath string

	file      *os.File
	writeChan chan string

	wait sync.WaitGroup
}

func (l *Logger) open() (err error) {
	uuid := uuid.New()
	now := time.Now().Format("20060102150405")

	workDir, err := os.Getwd()
	if err != nil {
		return
	}

	dir := filepath.Join(workDir, "logs")
	if err = os.MkdirAll(dir, 0777); err != nil {
		return
	}

	var buf [32]byte
	hex.Encode(buf[:], uuid[:])
	l.fullpath = filepath.Join(dir, fmt.Sprintf("%s-%s-%s.log", l.name, now, buf))

	// 创建日志文件
	if l.file, err = os.Create(l.fullpath); err != nil {
		return
	}

	l.writeChan = make(chan string, 1024)

	l.wait.Add(1)
	go l.writeLoop()

	return
}

func (l *Logger) Close() {
	close(l.writeChan)

	l.file.Close()
}

func (l *Logger) writeLoop() {
	defer l.wait.Done()

	for {
		message := <-l.writeChan

		fmt.Fprintln(l.file, message)
	}
}

func (l *Logger) Write(message string) {
	if l.writeChan == nil {
		return
	}

	l.writeChan <- message
}

func DefaultClose() {
	defaultLogger.Close()
}

func NewLogger(name string) *Logger {
	l := &Logger{
		name: name,
	}

	if err := l.open(); err != nil {
		l = nil
		return l
	}

	return l
}

func Default(name string) *Logger {
	defaultLogger = NewLogger(name)
	return defaultLogger
}

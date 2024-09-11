package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var defaultLogger *Logger

type Options struct {
	// 名称
	Name string `yaml:"name,omitempty"`
	// 日期后缀
	DateSuffix bool `yaml:"date_suffix,omitempty"`
	// 后缀名
	Suffix string `yaml:"suffix,omitempty"`
	// 每天新建文件
	EveryDay bool `yaml:"every_day,omitempty"`
}

type Logger struct {
	name string

	options Options

	writeChan chan string

	wait sync.WaitGroup
}

func (l *Logger) open() (err error) {
	l.writeChan = make(chan string, 1024)

	l.wait.Add(1)
	go l.writeLoop()

	return
}

func (l *Logger) Close() {
	close(l.writeChan)
	l.wait.Wait()
}

func openFile(option Options, date string) (fp *os.File, err error) {
	var fileName string
	if option.DateSuffix {
		fileName = fmt.Sprintf("%s-%s.%s", option.Name, date, option.Suffix)
	} else {
		fileName = option.Name + "." + option.Suffix
	}

	if fp, err = os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err != nil {
		dir := filepath.Dir(fileName)
		if err = os.MkdirAll(dir, 0644); err != nil {
			return
		}

		if fp, err = os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err != nil {
			return
		}
	}

	return
}

func (l *Logger) writeLoop() {
	var fp *os.File

	defer func() {
		l.writeChan = nil

		if fp != nil {
			fp.Close()
		}

		l.wait.Done()
	}()

	opt := l.options

	var currentDate string
	var date string

	var err error

	for message := range l.writeChan {
		if opt.EveryDay {
			date = time.Now().Format("20060102")
			if date != currentDate {
				if fp != nil {
					fp.Close()
					fp = nil
				}

				if fp, err = openFile(opt, date); err != nil {
					log.Println(err)
					continue
				}

				currentDate = date
			}
		} else {
			if fp == nil {
				date = time.Now().Format("20060102")
				if fp, err = openFile(opt, date); err != nil {
					log.Println(err)
					continue
				}
			}
		}

		fmt.Fprintln(fp, message)
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

func NewLogger(opts Options) *Logger {
	if len(opts.Name) <= 0 {
		opts.Name = "./logs/log"
	}

	if len(opts.Suffix) <= 0 {
		opts.Suffix = ".yaml"
	}

	l := &Logger{
		options: opts,
	}

	if err := l.open(); err != nil {
		l = nil
		return l
	}

	return l
}

func Default(opts Options) *Logger {
	defaultLogger = NewLogger(opts)
	return defaultLogger
}

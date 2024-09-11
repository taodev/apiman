package logger

import "fmt"

func Logf(format string, args ...any) {
	defaultLogger.Write(fmt.Sprintf(format, args...))
}

func Log(args ...any) {
	defaultLogger.Write(fmt.Sprint(args...))
}

func Print(args ...any) {
	message := fmt.Sprint(args...)
	fmt.Println(message)

	defaultLogger.Write(message)
}

func Printf(format string, args ...any) {
	message := fmt.Sprintf(format, args...)
	fmt.Println(message)

	defaultLogger.Write(message)
}

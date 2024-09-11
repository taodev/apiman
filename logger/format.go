package logger

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

func Logf(format string, args ...any) {
	defaultLogger.Write(fmt.Sprintf(format, args...))
}

func Log(args ...any) {
	defaultLogger.Write(fmt.Sprint(args...))
}

func LogYaml(args ...any) {
	out, err := yaml.Marshal(args)
	if err != nil {
		return
	}

	defaultLogger.Write(string(out))
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

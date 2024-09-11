package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"sync"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/taodev/apiman/logger"
	"github.com/taodev/apiman/storage"
)

var (
	configPath string
	workDir    string
	// 是否打印详细信息
	verboseVar bool

	loggerDir        string
	loggerName       string
	loggerSuffix     string
	loggerEveryDay   bool
	loggerDateSuffix bool

	globalWait   sync.WaitGroup
	globalCtx    context.Context
	globalCancel context.CancelFunc

	runPass bool
)

var mainCommand = &cobra.Command{
	Use: "apiman",
}

func init() {
	mainCommand.PersistentFlags().StringVarP(&configPath, "config", "c", "", "config file (default is $HOME/.apiman.yaml)")
	mainCommand.PersistentFlags().StringVarP(&workDir, "work", "w", "", "work dir (default is $PWD)")

	mainCommand.PersistentFlags().BoolVarP(&verboseVar, "verbose", "v", false, "verbose")
}

func preRun(cmd *cobra.Command, args []string) {
	var err error

	// 配置工作目录
	if len(workDir) <= 0 {
		workDir, err = os.Getwd()
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		workDir, err = filepath.Abs(workDir)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	if err = os.Chdir(workDir); err != nil {
		fmt.Println(err)
		return
	}

	// 判断env配置文件是否存在
	envPath := filepath.Join(workDir, ".env.yaml")
	if err = storage.LoadEnv(envPath); err != nil {
		fmt.Println(err)
		return
	}

	// 初始化日志
	logger.Default(logger.Options{
		Name:       filepath.Join(loggerDir, loggerName),
		Suffix:     loggerSuffix,
		EveryDay:   loggerEveryDay,
		DateSuffix: loggerDateSuffix,
		NoPrint:    !verboseVar,
	})

	// 设置线程数
	runtime.GOMAXPROCS(runtime.NumCPU())

	globalCtx, globalCancel = context.WithCancel(context.Background())
	go func() {
		osSignals := make(chan os.Signal, 1)
		signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM)
		defer signal.Stop(osSignals)

		<-osSignals

		globalCancel()
	}()
}

func postRun(cmd *cobra.Command, args []string) {
	globalWait.Wait()
	logger.DefaultClose()

	if !runPass {
		os.Exit(-1)
	}
}

func main() {
	if err := mainCommand.Execute(); err != nil {
		panic(err)
	}
}

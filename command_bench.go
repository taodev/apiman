package main

import (
	"time"

	"github.com/spf13/cobra"
)

var (
	numBench  int
	numWorker int
	interval  int
)

var commandBench = &cobra.Command{
	Use:     "bench",
	Short:   "run bench",
	PreRun:  preRun,
	PostRun: postRun,
	Run: func(cmd *cobra.Command, args []string) {
		runBench(args)
	},
}

func runBench(args []string) {
	for nw := 0; nw < numWorker; nw++ {
		globalWait.Add(1)
		go func() {
			defer globalWait.Done()

			for nb := 0; nb < numBench; nb++ {
				select {
				case <-globalCtx.Done():
					return
				default:
					runCase(args)
				}

				if interval <= 0 {
					continue
				}

				if nb < numBench-1 {
					<-time.After(time.Duration(interval) * time.Millisecond)
				}
			}
		}()
	}
}

func init() {
	// 运行次数
	commandBench.Flags().IntVarP(&numBench, "num-bench", "", 1, "bench")
	// 线程数
	commandBench.Flags().IntVarP(&numWorker, "num-worker", "", 1, "numWorker")
	// 间隔时间参数(毫秒)
	commandBench.Flags().IntVarP(&interval, "interval", "", 0, "interval in millisecond")

	// 日志配置
	commandBench.Flags().StringVarP(&loggerDir, "logger-dir", "", "logs", "logger dir")
	commandBench.Flags().StringVarP(&loggerName, "logger-name", "", "api-request", "logger name")
	commandBench.Flags().StringVarP(&loggerSuffix, "logger-suffix", "", ".yaml", "logger suffix")
	commandBench.Flags().BoolVarP(&loggerEveryDay, "logger-everyday", "", false, "logger every day")
	commandBench.Flags().BoolVarP(&loggerDateSuffix, "logger-datesuffix", "", true, "logger date suffix")

	mainCommand.AddCommand(commandBench)
}

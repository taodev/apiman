package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/taodev/apiman/runner"
)

var (
	bench     int
	numWorker int
	interval  int
)

var commandCase = &cobra.Command{
	Use:     "case",
	Short:   "run case",
	PreRun:  preRun,
	PostRun: postRun,
	Run: func(cmd *cobra.Command, args []string) {
		runCase(args)
	},
}

func runCase(args []string) {
	var err error
	for nw := 0; nw < numWorker; nw++ {
		globalWait.Add(1)
		go func() {
			defer globalWait.Done()

			for nb := 0; nb < bench; nb++ {
				for i := 0; i < len(args); i++ {
					select {
					case <-globalCtx.Done():
						log.Println("safe exit")
						return
					default:
						runner := new(runner.Runner)
						if err = runner.Do(workDir, configPath, args[i]); err != nil {
							fmt.Println("runner:", err)
							os.Exit(-1)
							return
						}
					}
				}

				if nb < bench-1 {
					<-time.After(time.Duration(1000) * time.Millisecond)
				}
			}
		}()
	}
}

func init() {
	// 运行次数
	commandCase.Flags().IntVarP(&bench, "bench", "b", 1, "bench")
	// 线程数
	commandCase.Flags().IntVarP(&numWorker, "num-worker", "n", 1, "numWorker")
	// 间隔时间参数(毫秒)
	commandCase.Flags().IntVarP(&interval, "interval", "i", 1, "interval in millisecond")
	mainCommand.AddCommand(commandCase)
}

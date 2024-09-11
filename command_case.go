package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/spf13/cobra"
	"github.com/taodev/apiman/runner"
	"github.com/taodev/apiman/storage"
)

var (
	bench     int
	numWorker int
	interval  int
)

var commandCase = &cobra.Command{
	Use:   "case",
	Short: "run case",
	Run: func(cmd *cobra.Command, args []string) {
		var err error
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

		var waitGroup sync.WaitGroup

		// 设置线程数
		runtime.GOMAXPROCS(runtime.NumCPU())

		for nw := 0; nw < numWorker; nw++ {
			waitGroup.Add(1)
			go func() {
				defer waitGroup.Done()

				for nb := 0; nb < bench; nb++ {
					for i := 0; i < len(args); i++ {
						runner := new(runner.Runner)
						if err = runner.Do(workDir, configPath, args[i]); err != nil {
							fmt.Println("runner:", err)
							os.Exit(-1)
							return
						}
					}

					if nb < bench-1 {
						<-time.After(time.Duration(1000) * time.Millisecond)
					}
				}
			}()
		}

		waitGroup.Wait()
	},
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

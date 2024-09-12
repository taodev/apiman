package main

import (
	"fmt"
	"sync"
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

type benchStat struct {
	BeginTime time.Time
	Time      time.Duration
	Total     int
	NumPass   int
	NumFail   int
}

func runBench(args []string) {
	var stat benchStat
	stat.BeginTime = time.Now()

	var locker sync.RWMutex
	runStat := true

	globalWait.Add(1)
	go func() {
		defer globalWait.Done()

		for runStat {
			select {
			case <-globalCtx.Done():
				return
			case <-time.After(time.Second):
				locker.RLock()
				stat.Time = time.Since(stat.BeginTime)
				fmt.Println("time:", stat.Time, "total:", stat.Total, "pass:", stat.NumPass, "fail:", stat.NumFail)
				locker.RUnlock()
			}
		}
	}()

	var runWait sync.WaitGroup

	for nw := 0; nw < numWorker; nw++ {
		runWait.Add(1)
		go func() {
			defer runWait.Done()

			for nb := 0; nb < numBench; nb++ {
				select {
				case <-globalCtx.Done():
					return
				default:
					_, pass, _ := runCase(args)
					locker.Lock()
					if pass {
						stat.NumPass++
					} else {
						stat.NumFail++
					}
					stat.Total++
					locker.Unlock()
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

	runWait.Wait()
	runStat = false
}

func init() {
	// 运行次数
	commandBench.Flags().IntVarP(&numBench, "num-bench", "", 1, "numBench")
	// 线程数
	commandBench.Flags().IntVarP(&numWorker, "num-worker", "", 1, "numWorker")
	// 间隔时间参数(毫秒)
	commandBench.Flags().IntVarP(&interval, "interval", "", 0, "interval in millisecond")

	mainCommand.AddCommand(commandBench)
}

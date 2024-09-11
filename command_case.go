package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/taodev/apiman/runner"
)

var commandCase = &cobra.Command{
	Use:     "case",
	Short:   "run case",
	PreRun:  preRun,
	PostRun: postRun,
	Run: func(cmd *cobra.Command, args []string) {
		results, err := runCase(args)

		if err != nil {
			fmt.Println("case:", err)
			os.Exit(-1)
			return
		}

		for _, result := range results {
			if !verboseVar {
				fmt.Println(result.String())
			}

			if !result.Pass() {
				runPass = false
			}
		}
	},
}

func runCase(args []string) (results []*runner.CaseResult, err error) {
	for i := 0; i < len(args); i++ {
		r := runner.NewRunner(globalCtx)
		var result *runner.CaseResult
		if result, err = r.Do(workDir, configPath, args[i]); err != nil {
			fmt.Println("runner:", err)
			os.Exit(-1)
			return
		}

		results = append(results, result)
	}

	return
}

func init() {
	// 日志配置
	commandCase.Flags().StringVarP(&loggerDir, "logger-dir", "", "logs", "logger dir")
	commandCase.Flags().StringVarP(&loggerName, "logger-name", "", "api-request", "logger name")
	commandCase.Flags().StringVarP(&loggerSuffix, "logger-suffix", "", ".yaml", "logger suffix")
	commandCase.Flags().BoolVarP(&loggerEveryDay, "logger-everyday", "", false, "logger every day")
	commandCase.Flags().BoolVarP(&loggerDateSuffix, "logger-datesuffix", "", true, "logger date suffix")

	mainCommand.AddCommand(commandCase)
}

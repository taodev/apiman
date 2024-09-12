package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/taodev/apiman/runner"
)

var commandRun = &cobra.Command{
	Use:     "run",
	Short:   "run case",
	PreRun:  preRun,
	PostRun: postRun,
	Run: func(cmd *cobra.Command, args []string) {
		results, _, err := runCase(args)

		if err != nil {
			fmt.Println("run:", err)
			os.Exit(-1)
			return
		}

		for _, result := range results {
			if !verboseVar {
				fmt.Println(result.String())
			}
		}
	},
}

func runCase(args []string) (results []*runner.CaseResult, pass bool, err error) {
	pass = true

	r := runner.NewRunner(workDir, configPath, globalCtx)
	if err = r.Load(); err != nil {
		return
	}

	if len(args) <= 0 {
		args = r.GetAllCases()
	}

	for i := 0; i < len(args); i++ {
		var result *runner.CaseResult
		if result, err = r.Do(args[i]); err != nil {
			fmt.Println("runner:", err)
			pass = false
			os.Exit(-1)
			return
		}

		if !result.Pass() {
			pass = false
			runPass = false
		}

		results = append(results, result)
	}

	return
}

func init() {
	mainCommand.AddCommand(commandRun)
}

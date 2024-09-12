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
		results, _, err := runCase(args)

		if err != nil {
			fmt.Println("case:", err)
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

	for i := 0; i < len(args); i++ {
		r := runner.NewRunner(globalCtx)
		var result *runner.CaseResult
		if result, err = r.Do(workDir, configPath, args[i]); err != nil {
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
	mainCommand.AddCommand(commandCase)
}

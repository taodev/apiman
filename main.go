package main

import "github.com/spf13/cobra"

var (
	configPath string
)

var mainCommand = &cobra.Command{
	Use: "apiman",
}

func init() {
	mainCommand.PersistentFlags().StringVarP(&configPath, "config", "c", "", "config file (default is $HOME/.apiman.yaml)")
}

func main() {
	if err := mainCommand.Execute(); err != nil {
		panic(err)
	}
}

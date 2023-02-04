package cmd

import (
	"fmt"

	"github.com/EvertonTomalok/go-template/internal/app"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "server",
	Short: "Run http server",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()

		config := app.Configure(ctx)

		fmt.Printf("%+v\n", config)

	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}

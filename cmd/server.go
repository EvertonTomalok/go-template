package cmd

import (
	"github.com/EvertonTomalok/go-template/internal/app"
	"github.com/EvertonTomalok/go-template/internal/ui/rest"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run http server",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		defer app.CloseConnections(ctx)

		config := app.Configure(ctx)

		app.InitDB(ctx, config, true)
		app.InitMongoDB(ctx, config, true)
		rest.RunServer(ctx, config)
	},
}

func init() {
	serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.AddCommand(serverCmd)
}

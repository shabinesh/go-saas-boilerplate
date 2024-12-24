package cmd

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/shabinesh/app/infra/web"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the app server",
	Long:  `Start the app server`,
	Run: func(cmd *cobra.Command, args []string) {
		// Open new database connection
		db, err := pgx.Connect(context.Background(), "postgres://postgres:password@localhost:5432/myapp")
		if err != nil {
			log.Fatal(err)
		}

		web.StartServer(db)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}

package cmd

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/shabinesh/transcription/infra/web"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the transcription server",
	Long:  `Start the transcription server`,
	Run: func(cmd *cobra.Command, args []string) {
		// Open new database connection
		db, err := sql.Open("sqlite3", "test.db")
		if err != nil {
			log.Fatal(err)
		}
		web.StartServer(db)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}

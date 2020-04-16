package cmd

import (
	"fmt"
	"lingva/server"
	"log"
	"os"

	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts the actual GraphQL server.",
	Run: func(cmd *cobra.Command, args []string) {
		port := os.Getenv("PORT")
		if len(port) == 0 {
			log.Fatal("NO PORT SPECIFIED!")
		}
		fmt.Printf("Starting Lingva GraphQL server at :%s!\n", port)
		s, err := server.NewHTTPServer(port)
		if err != nil {
			log.Fatal(err)
		}
		s.ListenAndServe()

	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}

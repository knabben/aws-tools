package cmd

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/redis.v5"
)

var Redis string
var client *redis.Client

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "HTTP Stuff",
	Long:  `REST API to manage resources`,
	Run: func(cmd *cobra.Command, args []string) {
		// Subscribe for keyevents expired
		client = redis.NewClient(&redis.Options{Addr: Redis, Password: "", DB: 0})
		pubsub, err := client.Subscribe("__keyevent@0__:expired")
		if err != nil {
			panic(err)
		}
		defer pubsub.Close()

		go func() {
			for {
				msg, err := pubsub.ReceiveMessage()
				if err != nil {
					panic(err)
				}
				splitter := strings.Split(msg.Payload, "|")
				ip := splitter[0]
				machine := splitter[1]
				output := deleteIPFromSG(ip, machine)
				if output {
					fmt.Println("Removed " + ip + " from " + machine)
				}
			}
		}()

		http.Handle("/api/", CreateHTTPAPIHandler())
		log.Print(http.ListenAndServe(fmt.Sprintf(":%d", 8085), nil))
	},
}

func init() {
	RootCmd.AddCommand(serveCmd)

	RootCmd.Flags().StringVarP(&Redis, "redis", "r", "localhost:6379", "Redis URL")
}

package cmd

import "github.com/spf13/cobra"

var IP string
var Machine string

var Region = "us-east-1"

var sgCmd = &cobra.Command{
	Use:   "sg",
	Short: "SecurityGroup Management",
	Long:  `A SecurityGroup management command tool`,
	Run: func(cmd *cobra.Command, args []string) {
		InsertIPOnSg(IP, Machine)
	},
}

func init() {
	RootCmd.AddCommand(sgCmd)

	sgCmd.PersistentFlags().StringVarP(&IP, "ip", "i", "", "IPAddress")
	sgCmd.PersistentFlags().StringVarP(&Machine, "machine", "m", "", "Machine")
}

package cmd

import (
	"backend/internal/global"
	"fmt"
	"github.com/spf13/cobra"
)

var (
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number of server",
		Run: func(cmd *cobra.Command, args []string) {
			version()
		},
	}
)

func version() {
	fmt.Printf("ServerConfig version %s -- HEAD\n", global.Version)
}

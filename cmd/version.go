package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version = "version is set by build process"

func init() {
	rootCMD.AddCommand(versionCMD)
}

var versionCMD = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run:   versionFunc,
}

func versionFunc(cmd *cobra.Command, args []string) {
	fmt.Println("kubeswap " + Version)
}

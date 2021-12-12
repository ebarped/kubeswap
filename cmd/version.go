package cmd

import (
	"fmt"

	"github.com/ebarped/kubeswap/internal/version"
	"github.com/spf13/cobra"
)

func init() {
	rootCMD.AddCommand(versionCMD)
}

var versionCMD = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Run:   versionFunc,
}

func versionFunc(cmd *cobra.Command, args []string) {
	fmt.Println("kubeswap " + version.Version)
}

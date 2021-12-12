package cmd

// this command has to:
// 1º. check that the kubeconfig file exists & is well formed
// 2º. add a new key-value to the DB, key=name, value=kubeconfig_content

import (
	"github.com/spf13/cobra"
)

var (
	name       string // represents the key to store the kubeconfig in the db
	kubeconfig string // represents the path to the kubeconfig file
)

var addCMD = &cobra.Command{
	Use:   "add <name> -f <kubeconfig>",
	Short: "adds a new kubeconfig to the database",
	Run:   addFunc,
}

// init adds all flags of this command
func init() {
	rootCMD.AddCommand(addCMD)

	addCMD.Flags().StringVarP(&name, "name", "", "", "name of the kubeconfig")
	addCMD.Flags().StringVarP(&kubeconfig, "kubeconfig", "", "", "kubeconfig's path")
}

func addFunc(cmd *cobra.Command, args []string) {
	log.Info().Msgf("Executing add command with values %s:%s", name, kubeconfig)
}

package cmd

import (
	"os"

	"github.com/ebarped/kubeswap/pkg/kubeconfig"
	"github.com/ebarped/kubeswap/pkg/kv"
	"github.com/spf13/cobra"
)

var addCMD = &cobra.Command{
	Use:   "add -n <name> -f <kubeconfig>",
	Short: "Adds a new kubeconfig to the database",
	Run:   addFunc,
}

// init adds this command and his flags
func init() {
	rootCMD.AddCommand(addCMD)
	addCMD.Flags().StringP("name", "n", "", "name of the kubeconfig")
	addCMD.MarkFlagRequired("name")
	addCMD.Flags().String("kubeconfig", "", "kubeconfig's path")
	addCMD.MarkFlagRequired("kubeconfig")
}

func addFunc(cmd *cobra.Command, args []string) {
	name, _ := cmd.Flags().GetString("name")
	kubeconfigPath, _ := cmd.Flags().GetString("kubeconfig")

	retcode := 0
	defer func() { os.Exit(retcode) }()

	log.Debug().Str("command", "add").Str("name", name).Str("kubeconfig path", kubeconfigPath).Str("database", dbPath).Send()

	kc, err := kubeconfig.New(name, kubeconfigPath)
	if err != nil {
		log.Error().Str("error", err.Error()).Msg("error creating new Kubeconfig struct")
		retcode = 1
		return
	}

	db, err := kv.Open(dbPath)
	if err != nil {
		log.Error().Str("error", err.Error()).Str("db", dbPath).Msg("error opening kv database")
		retcode = 1
		return
	}
	defer db.Close()

	log.Debug().Str("action", "adding new kubeconfig to the database").Str("key", kc.Name).Str("value", kc.Content).Send()

	err = db.PutKubeconfig(kc.Name, []byte(kc.Content))
	if err != nil {
		log.Error().Str("error", err.Error()).Msg("error putting new key-value to the database")
		retcode = 1
		return
	}
	log.Debug().Str("command", "add").Str("key", kc.Name).Str("result", "successful").Send()
}

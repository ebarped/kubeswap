package cmd

import (
	"github.com/ebarped/kubeswap/pkg/kubeconfig"
	"github.com/ebarped/kubeswap/pkg/kv"
	"github.com/spf13/cobra"
)

var (
	name           string // key to store the kubeconfig in the db
	kubeconfigPath string // path to the kubeconfig
)

var addCMD = &cobra.Command{
	Use:   "add <name> -f <kubeconfig>",
	Short: "adds a new kubeconfig to the database",
	Run:   addFunc,
}

// init adds this command and his flags
func init() {
	rootCMD.AddCommand(addCMD)
	addCMD.Flags().StringVarP(&name, "name", "", "", "name of the kubeconfig")
	addCMD.MarkFlagRequired("name")
	addCMD.Flags().StringVar(&kubeconfigPath, "kubeconfig", "", "kubeconfig's path")
	addCMD.MarkFlagRequired("kubeconfig")
}

func addFunc(cmd *cobra.Command, args []string) {
	log.Debug().Str("command", "add").Str("name", name).Str("kubeconfig path", kubeconfigPath).Str("database", dbPath).Send()

	kc, err := kubeconfig.New(name, kubeconfigPath)
	if err != nil {
		log.Fatal().Str("error", err.Error()).Msg("error creating new Kubeconfig struct")
	}

	db, err := kv.Open(dbPath)
	if err != nil {
		log.Fatal().Str("error", err.Error()).Msg("error opening kv database")
	}
	defer db.CloseDB()

	log.Debug().Str("action", "adding new kubeconfig to the database").Str("key", kc.Name).Str("value", kc.Content).Send()

	err = db.PutKubeconfig(kc.Name, []byte(kc.Content))
	if err != nil {
		log.Error().Str("error", err.Error()).Msg("error putting new key-value to the database")
	}
}

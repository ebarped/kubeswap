package cmd

// this command has to:
// 1º. check that the kubeconfig file exists & is well formed
// 2º. add a new key-value to the DB, key=name, value=kubeconfig_content

import (
	"github.com/ebarped/kubeswap/pkg/kubeconfig"
	"github.com/ebarped/kubeswap/pkg/kv"
	"github.com/spf13/cobra"
)

var (
	name           string // key to store the kubeconfig in the db
	kubeconfigPath string // path to the kubeconfig file
	dbPath         string // path to the db file
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
	addCMD.PersistentFlags().StringVar(&dbPath, "db", dbPath, "db file path")
	addCMD.Flags().StringVarP(&kubeconfigPath, "kubeconfig", "", "", "kubeconfig's path")
}

func addFunc(cmd *cobra.Command, args []string) {
	log.Info().Str("command", "add").Str("name", name).Str("kubeconfig path", kubeconfigPath).Str("database", dbPath).Send()

	kc, err := kubeconfig.New(name, kubeconfigPath)
	if err != nil {
		log.Fatal().Str("error", err.Error()).Msg("error creating new Kubeconfig struct")
	}

	db, err := kv.Open(dbPath)
	if err != nil {
		log.Fatal().Str("error", err.Error()).Msg("error opening kv database")
	}
	defer db.CloseDB()

	kcContent, err := kc.Config()
	if err != nil {
		log.Fatal().Str("error", err.Error()).Msg("error obtaining the content of the kubeconfig")
	}

	err = db.Put(kc.Name(), kcContent)
	if err != nil {
		log.Fatal().Str("error", err.Error()).Msg("error putting new key-value to the database")
	}
	log.Info().Str("action", "adding new kubeconfig to the database").Str("key", kc.Name()).Str("value", string(kcContent)).Send()
}

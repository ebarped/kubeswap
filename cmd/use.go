package cmd

import (
	"os"

	"github.com/ebarped/kubeswap/pkg/kv"
	"github.com/spf13/cobra"
)

var useCMD = &cobra.Command{
	Use:   "use",
	Short: "Select kubeconfig to use",
	Run:   useFunc,
}

// init adds this command and his flags
func init() {
	rootCMD.AddCommand(useCMD)
	useCMD.Flags().StringP("name", "n", "", "name of the kubeconfig")
}

func useFunc(cmd *cobra.Command, args []string) {
	var name string

	retcode := 0
	defer func() { os.Exit(retcode) }()

	db, err := kv.Open(dbPath)
	if err != nil {
		log.Error().Str("error", err.Error()).Str("db", dbPath).Msg("error opening kv database")
		retcode = 1
		return
	}
	defer db.CloseDB()

	userHome, err := os.UserHomeDir()
	if err != nil {
		log.Fatal().Str("error", err.Error()).Msg("error getting the homeDir of the user")
	}
	kcPath := userHome + "/.kube/config"

	if cmd.Flags().Changed("name") {
		name, _ = cmd.Flags().GetString("name")
		err = useWithName(db, name, kcPath)
	} else {
		err = useWithoutName(db, kcPath)
	}
	if err != nil {
		log.Error().Str("error", err.Error()).Str("db", dbPath).Msg("error selecting kubeconfig to use")
		retcode = 1
		return
	}
	log.Debug().Str("command", "use").Str("key", name).Str("result", "successful").Send()
}

func useWithName(db *kv.DB, name, kubeconfigPath string) error {
	log.Debug().Str("command", "use").Str("with name", "true").Str("name", name).Str("database", dbPath).Send()

	kc, err := db.GetKubeconfig(name)
	if err != nil {
		return err
	}
	err = os.WriteFile(kubeconfigPath, []byte(kc.Content), 0o644)
	if err != nil {
		return err
	}
	return nil
}

func useWithoutName(db *kv.DB, kubeconfigPath string) error {
	log.Debug().Str("command", "use").Str("with name", "false").Str("database", dbPath).Send()
	return nil
}

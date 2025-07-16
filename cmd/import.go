package cmd

import (
	"fmt"
	"os"

	"github.com/ebarped/kubeswap/pkg/kubeconfig"
	"github.com/ebarped/kubeswap/pkg/kv"
	"github.com/spf13/cobra"
)

var importCMD = &cobra.Command{
	Use:   "import",
	Short: "Removes the current db file and recreates it importing all the $HOME/.kube kubeconfigs",
	Run:   importFunc,
}

// init adds this command and his flags
func init() {
	rootCMD.AddCommand(importCMD)
}

func importFunc(cmd *cobra.Command, args []string) {
	retcode := 0
	defer func() { os.Exit(retcode) }()

	log.Debug().Str("command", "import").Str("database", dbPath).Msg("")

	err := os.RemoveAll(dbPath)
	if err != nil {
		log.Error().Str("command", "import").Str("error", err.Error()).Msgf("error removing dabasase %q", dbPath)
		retcode = 1
		return
	}

	db, err := kv.Open(dbPath)
	if err != nil {
		log.Error().Str("command", "import").Str("error", err.Error()).Msg("error opening kv database")
		retcode = 1
		return
	}
	defer db.Close()

	// we dont have the filename arg, so we scan the $HOME/.kube/ directory
	files, err := os.ReadDir(defaultKubeconfigPath())
	if err != nil {
		log.Error().Str("command", "import").Msgf("error reading %s: %s", kcRootDir, err.Error())
	}

	for _, f := range files {
		if f.IsDir() || f.Name() == "config" {
			continue
		}
		kubeconfigPath := kcRootDir + f.Name()
		log.Debug().Str("command", "import").Str("file", f.Name()).Str("path", kubeconfigPath).Msg("loading kubeconfig...")
		kc, err := kubeconfig.New(f.Name(), kubeconfigPath)
		if err != nil {
			log.Debug().Str("command", "import").Str("file", f.Name()).Msg("not a valid kubeconfig")
			continue
		}

		log.Debug().Str("action", "adding new kubeconfig to the database").Str("key", kc.Name).Str("value", kc.Content).Send()

		err = db.PutKubeconfig(kc.Name, []byte(kc.Content))
		if err != nil {
			log.Error().Str("command", "import").Str("error", err.Error()).Msg("error putting new key-value to the database")
			retcode = 1
			return
		}
	}

	fmt.Println("Import successful! 🚀")
}
